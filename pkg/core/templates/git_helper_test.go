package template

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTransport struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.handler(req)
}

// mockHTTPResponse tạo nhanh một response giả lập
func mockHTTPResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestDoGET_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	resp, err := doGET(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDoGET_SetsDefaultHeaders(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"))
		assert.Equal(t, "2022-11-28", r.Header.Get("X-GitHub-Api-Version"))
		assert.Empty(t, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp, err := doGET(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDoGET_WithToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "ghp_test_token_abc")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer ghp_test_token_abc", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp, err := doGET(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDoGET_InvalidURL(t *testing.T) {
	_, err := doGET("http://example.com/\x7f")
	require.Error(t, err)
}

func TestDoGET_ConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close() // đóng ngay

	_, err := doGET(url)
	require.Error(t, err)
}

func TestDoGET_Non2xx(t *testing.T) {
	tests := []struct {
		status int
		body   string
		errMsg string
	}{
		{http.StatusNotFound, `{"message":"Not Found"}`, "status 404"},
		{http.StatusForbidden, `{"message":"API rate limit exceeded"}`, "status 403"},
		{http.StatusInternalServerError, "  \n  server error \n\n  ", "status 500: server error"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			_, err := doGET(server.URL)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
			assert.Contains(t, err.Error(), strings.TrimSpace(tt.body))
		})
	}
}

func TestFetchFileViaAPI_Success(t *testing.T) {
	content := "package book\n"
	encoded := base64.StdEncoding.EncodeToString([]byte(content))
	// Giả lập GitHub API chèn \n và \r\n vào base64
	encodedWithNewlines := encoded[:5] + "\n" + encoded[5:] + "\r\n"

	oldTransport := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = oldTransport }()

	http.DefaultClient.Transport = &mockTransport{
		handler: func(req *http.Request) (*http.Response, error) {
			// SỬA Ở ĐÂY: Dùng json.Marshal thay vì fmt.Sprintf để tự động escape \n
			mockData := githubContent{
				Type:     "file",
				Encoding: "base64",
				Content:  encodedWithNewlines,
			}
			bodyBytes, err := json.Marshal(mockData)
			require.NoError(t, err)

			return mockHTTPResponse(http.StatusOK, string(bodyBytes)), nil
		},
	}

	result, err := fetchFileViaAPI("owner", "repo", "main", "internal/book/module.go")
	require.NoError(t, err)
	assert.Equal(t, content, result)
}

func TestFetchFileViaAPI_Errors(t *testing.T) {
	tests := []struct {
		name     string
		mockBody string
		status   int
		errMsg   string
	}{
		{
			name:     "API returns error status",
			status:   http.StatusNotFound,
			mockBody: `{"message":"Not Found"}`,
			errMsg:   "status 404",
		},
		{
			name:     "invalid JSON response",
			status:   http.StatusOK,
			mockBody: `{"type": "file", "encoding": `, // cố tình hỏng JSON
			errMsg:   "unexpected EOF",
		},
		{
			name:     "not a file",
			status:   http.StatusOK,
			mockBody: `{"type": "dir", "encoding": "base64", "content": ""}`,
			errMsg:   "not a file: internal/book (type=dir)",
		},
		{
			name:     "wrong encoding",
			status:   http.StatusOK,
			mockBody: `{"type": "file", "encoding": "utf-8", "content": ""}`,
			errMsg:   `unexpected encoding "utf-8"`,
		},
		{
			name:     "invalid base64",
			status:   http.StatusOK,
			mockBody: `{"type": "file", "encoding": "base64", "content": "invalid_base64!!!"}`,
			errMsg:   "decode base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTransport := http.DefaultClient.Transport
			defer func() { http.DefaultClient.Transport = oldTransport }()

			http.DefaultClient.Transport = &mockTransport{
				handler: func(req *http.Request) (*http.Response, error) {
					return mockHTTPResponse(tt.status, tt.mockBody), nil
				},
			}

			_, err := fetchFileViaAPI("owner", "repo", "main", "internal/book")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestParseGitHubRepo_Success(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		owner string
		repo  string
	}{
		{"https URL", "https://github.com/owner/repo", "owner", "repo"},
		{"with .git suffix", "https://github.com/owner/repo.git", "owner", "repo"},
		{"without protocol", "github.com/owner/repo", "owner", "repo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseGitHubRepo(tt.url)
			assert.NoError(t, err)
			assert.Equal(t, tt.owner, owner)
			assert.Equal(t, tt.repo, repo)
		})
	}
}

func TestParseGitHubRepo_Error(t *testing.T) {
	tests := []struct {
		name string
		url  string
		err  string
	}{
		{"invalid URL", "://invalid", "invalid repository URL"},
		{"not github", "https://gitlab.com/owner/repo", "must be hosted on github.com"},
		{"invalid path - single part", "https://github.com/owner", "invalid GitHub URL"},
		{"empty owner", "https://github.com//repo", "owner is empty"},
		{"empty repo", "https://github.com/owner/", "repo is empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseGitHubRepo(tt.url)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.err)
		})
	}
}

func TestFetchTemplateFilesFromGit_Success(t *testing.T) {
	fileContent := "package book"
	encodedContent := base64.StdEncoding.EncodeToString([]byte(fileContent))

	oldTransport := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = oldTransport }()

	http.DefaultClient.Transport = &mockTransport{
		handler: func(req *http.Request) (*http.Response, error) {
			path := req.URL.Path
			switch {
			case path == "/repos/owner/repo":
				return mockHTTPResponse(http.StatusOK, `{"default_branch": "main"}`), nil
			case path == "/repos/owner/repo/git/trees/main":
				treeBody := `{
					"tree": [
						{"path": "internal/book/module.go", "type": "blob"},
						{"path": "internal/book/handler/handler.go", "type": "blob"},
						{"path": "internal/shared/domain/book.go", "type": "blob"},
						{"path": "internal/other/file.go", "type": "blob"},
						{"path": "internal/book/empty_dir", "type": "tree"},
						{"path": "README.md", "type": "blob"}
					]
				}`
				return mockHTTPResponse(http.StatusOK, treeBody), nil
			case strings.HasPrefix(path, "/repos/owner/repo/contents/internal/book/"),
				strings.HasPrefix(path, "/repos/owner/repo/contents/internal/shared/domain/book.go"):
				body := fmt.Sprintf(`{"type": "file", "encoding": "base64", "content": "%s"}`, encodedContent)
				return mockHTTPResponse(http.StatusOK, body), nil
			}
			return mockHTTPResponse(http.StatusNotFound, `{"message": "Not Found"}`), nil
		},
	}

	files, err := fetchTemplateFilesFromGit("https://github.com/owner/repo", "book")
	require.NoError(t, err)

	// Phải có 3 file: module.go, handler.go và shared/domain/book.go
	assert.Len(t, files, 3)
	assert.Equal(t, fileContent, files["internal/book/module.go"])
	assert.Equal(t, fileContent, files["internal/book/handler/handler.go"])
	assert.Equal(t, fileContent, files["internal/shared/domain/book.go"])

	// Không lấy file của internal/other và README.md
	_, ok := files["internal/other/file.go"]
	assert.False(t, ok)
}

func TestFetchTemplateFilesFromGit_FetchFileError(t *testing.T) {
	oldTransport := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = oldTransport }()

	http.DefaultClient.Transport = &mockTransport{
		handler: func(req *http.Request) (*http.Response, error) {
			path := req.URL.Path
			switch {
			case path == "/repos/owner/repo":
				return mockHTTPResponse(http.StatusOK, `{"default_branch": "main"}`), nil
			case path == "/repos/owner/repo/git/trees/main":
				treeBody := `{"tree": [{"path": "internal/book/module.go", "type": "blob"}]}`
				return mockHTTPResponse(http.StatusOK, treeBody), nil
			default:
				// Trả về 500 cho request lấy content
				return mockHTTPResponse(http.StatusInternalServerError, `server error`), nil
			}
		},
	}

	_, err := fetchTemplateFilesFromGit("https://github.com/owner/repo", "book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fetch internal/book/module.go: GET")
	assert.Contains(t, err.Error(), "status 500")
}

func TestFetchTemplateFilesFromGit_InvalidRepo(t *testing.T) {
	_, err := fetchTemplateFilesFromGit("invalid", "book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be hosted on github.com")
}

func TestFetchTemplateFilesFromGit_InvalidURL(t *testing.T) {
	_, err := fetchTemplateFilesFromGit("://invalid-url", "book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid repository URL")
}

func TestFetchTemplateFilesFromGit_NonGitHub(t *testing.T) {
	_, err := fetchTemplateFilesFromGit("https://gitlab.com/owner/repo", "book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be hosted on github.com")
}

func TestFetchTemplateFilesFromGit_InvalidPath(t *testing.T) {
	_, err := fetchTemplateFilesFromGit("https://github.com/owner", "book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid GitHub URL")
}

func BenchmarkParseGitHubRepo(b *testing.B) {
	url := "https://github.com/owner/repo"
	b.ResetTimer()
	for b.Loop() {
		_, _, _ = parseGitHubRepo(url)
	}
}
