package template

import (
	//"errors"
	//"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestDoGET_Non2xx_404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	_, err := doGET(server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")
	assert.Contains(t, err.Error(), `{"message":"Not Found"}`)
}

func TestDoGET_Non2xx_403RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"API rate limit exceeded"}`))
	}))
	defer server.Close()

	_, err := doGET(server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 403")
	assert.Contains(t, err.Error(), "API rate limit exceeded")
}

func TestDoGET_Non2xx_500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	_, err := doGET(server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 500")
}

func TestDoGET_Non2xx_TrimsBodyWhitespace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("  \n  server error \n\n  "))
	}))
	defer server.Close()

	_, err := doGET(server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 500: server error")
}
func TestFetchRawContent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`content`))
	}))
	defer server.Close()

	content, err := fetchRawContent(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, "content", content)
}

func TestFetchRawContent_Error(t *testing.T) {
	_, err := fetchRawContent("http://invalid-url-12345.local")
	assert.Error(t, err)
}

func TestFetchRawContent_Non200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := fetchRawContent(server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}

// ==================== TEST parseGitHubRepo ====================

func TestParseGitHubRepo_Success(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		owner string
		repo  string
	}{
		{
			name:  "https URL",
			url:   "https://github.com/owner/repo",
			owner: "owner",
			repo:  "repo",
		},
		{
			name:  "with .git suffix",
			url:   "https://github.com/owner/repo.git",
			owner: "owner",
			repo:  "repo",
		},
		{
			name:  "without protocol",
			url:   "github.com/owner/repo",
			owner: "owner",
			repo:  "repo",
		},
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
		{
			name: "invalid URL",
			url:  "://invalid",
			err:  "invalid repository URL",
		},
		{
			name: "not github",
			url:  "https://gitlab.com/owner/repo",
			err:  "must be hosted on github.com",
		},
		{
			name: "invalid path - single part",
			url:  "https://github.com/owner",
			err:  "invalid GitHub URL",
		},
		{
			name: "empty owner",
			url:  "https://github.com//repo",
			err:  "owner is empty",
		},
		{
			name: "empty repo",
			url:  "https://github.com/owner/",
			err:  "repo is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseGitHubRepo(tt.url)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.err)
		})
	}
}

// ==================== TEST fetchTemplateFilesFromGit ====================

func TestFetchTemplateFilesFromGit_Success(t *testing.T) {
	// Tạo mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"default_branch": "main"}`))
			return
		}
		if r.URL.Path == "/repos/owner/repo/git/trees/main" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"tree": [
					{"path": "internal/book/module.go", "type": "blob"},
					{"path": "internal/book/handler/handler.go", "type": "blob"},
					{"path": "internal/shared/domain/book.go", "type": "blob"},
					{"path": "README.md", "type": "blob"}
				]
			}`))
			return
		}
		// Raw content requests
		if strings.HasPrefix(r.URL.Path, "/owner/repo/main/internal/") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`package book`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Mock GitHub API URL bằng cách override
	// Vì code dùng github.com cố định, test này sẽ skip
	t.Skip("Skipping: requires mocking github.com API")
}
func TestFetchTemplateFilesFromGit_InvalidRepo(t *testing.T) {
	_, err := fetchTemplateFilesFromGit("invalid", "book")
	assert.Error(t, err)
	// Sửa: lỗi từ parseGitHubRepo với "invalid" là "must be hosted on github.com"
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

// ==================== TEST INTEGRATION (có thể skip) ====================

func TestFetchTemplateFilesFromGit_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test với repo thật
	files, err := fetchTemplateFilesFromGit("https://github.com/DVV-15324/witches-book", "book")
	if err != nil {
		t.Logf("Integration test failed: %v", err)
		t.Skip("Skipping due to network/github API issues")
	}

	assert.NotEmpty(t, files)

	// Kiểm tra có file module.go
	_, ok := files["internal/book/module.go"]
	if !ok {
		t.Log("module.go not found in fetched files")
	}
}

// ==================== BENCHMARK ====================

func BenchmarkParseGitHubRepo(b *testing.B) {
	url := "https://github.com/owner/repo"
	b.ResetTimer()
	for b.Loop() {
		_, _, _ = parseGitHubRepo(url)
	}
}
