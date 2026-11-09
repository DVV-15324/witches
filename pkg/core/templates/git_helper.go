package template

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type githubContent struct {
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
	Content  string `json:"content"`
	Size     int    `json:"size"`
}

func doGET(rawURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"GET %s: status %d: %s",
			rawURL,
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}

	return resp, nil
}

func fetchFileViaAPI(owner, repo, branch, path string) (string, error) {
	apiURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		owner, repo, path, branch,
	)
	resp, err := doGET(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var c githubContent
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return "", err
	}
	if c.Type != "file" {
		return "", fmt.Errorf("not a file: %s (type=%s)", path, c.Type)
	}
	if c.Encoding != "base64" {
		return "", fmt.Errorf("unexpected encoding %q for %s", c.Encoding, path)
	}
	// Contents API chèn \n mỗi 60 ký tự base64, phải strip trước khi decode.
	clean := strings.NewReplacer("\n", "", "\r", "").Replace(c.Content)
	raw, err := base64.StdEncoding.DecodeString(clean)
	if err != nil {
		return "", fmt.Errorf("decode base64 for %s: %w", path, err)
	}
	return string(raw), nil
}
func parseGitHubRepo(repoURL string) (owner, repo string, err error) {
	repoURL = strings.TrimSpace(repoURL)

	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = "https://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid repository URL: %w", err)
	}

	if u.Host == "" || u.Host == ":" || u.Host == "://" {
		return "", "", fmt.Errorf("invalid repository URL: host is empty or malformed")
	}
	if u.Host != "github.com" {
		return "", "", fmt.Errorf("repository must be hosted on github.com, got: %s", u.Host)
	}

	path := strings.TrimPrefix(u.Path, "/")

	parts := strings.Split(path, "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf(
			"invalid GitHub URL: expected https://github.com/{owner}/{repo}, got: %s",
			repoURL,
		)
	}

	owner = parts[0]
	repo = strings.TrimSuffix(parts[1], ".git")

	if owner == "" {
		return "", "", fmt.Errorf("invalid GitHub repository: owner is empty")
	}

	if repo == "" {
		return "", "", fmt.Errorf("invalid GitHub repository: repo is empty")
	}

	return owner, repo, nil
}

type gitHubTree struct {
	SHA  string `json:"sha"`
	URL  string `json:"url"`
	Tree []struct {
		Path string `json:"path"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"tree"`
}

func fetchTemplateFilesFromGit(repoURL string, targetDomain string) (map[string]string, error) {
	owner, repo, err := parseGitHubRepo(repoURL)
	if err != nil {
		return nil, err
	}

	fmt.Printf("  Owner: %s, Repo: %s\n", owner, repo)

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	resp, err := doGET(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repoData struct {
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return nil, err
	}
	branch := repoData.DefaultBranch
	fmt.Printf("  Default branch: %s\n", branch)

	treeURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", owner, repo, branch)
	resp, err = doGET(treeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}
	defer resp.Body.Close()

	var treeData gitHubTree
	if err := json.NewDecoder(resp.Body).Decode(&treeData); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("decode tree error: %w, body: %s", err, string(body))
	}

	templateFiles := make(map[string]string)

	fmt.Println("  Scanning files in internal/...")
	for _, item := range treeData.Tree {
		if item.Type != "blob" {
			continue
		}

		if !strings.HasPrefix(item.Path, "internal/") && !strings.HasPrefix(item.Path, "pkg/") {
			continue
		}

		var shouldFetch bool
		var fileKey string

		internalPrefix := fmt.Sprintf("internal/%s/", targetDomain)
		if strings.HasPrefix(item.Path, internalPrefix) {
			shouldFetch = true
			fileKey = item.Path
			fmt.Printf("  Found file: %s\n", item.Path)
		}

		if !shouldFetch {
			sharedDomainPath := fmt.Sprintf("internal/shared/domain/%s.go", targetDomain)
			if item.Path == sharedDomainPath {
				shouldFetch = true
				fileKey = item.Path
				fmt.Printf("  Found shared domain: %s\n", item.Path)
			}
		}

		if !shouldFetch {
			continue
		}

		content, err := fetchFileViaAPI(owner, repo, branch, item.Path)
		if err != nil {
			return nil, fmt.Errorf("fetch %s: %w", item.Path, err)
		}

		templateFiles[fileKey] = content
		fmt.Printf("  Fetched: %s\n", item.Path)
	}

	sqlFiles := []string{}
	for path := range templateFiles {
		if strings.HasSuffix(path, ".sql") {
			sqlFiles = append(sqlFiles, path)
		}
	}

	if len(sqlFiles) > 0 {
		fmt.Printf("  Found %d SQL migration files:\n", len(sqlFiles))
		for _, f := range sqlFiles {
			fmt.Printf("     - %s\n", f)
		}
	} else {
		fmt.Println("No SQL files found in template repository")
	}

	return templateFiles, nil
}
