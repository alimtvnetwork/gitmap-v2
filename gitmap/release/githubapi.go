// Package release — githubapi.go provides GitHub API types and helpers.
package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GitHubRelease represents the response from creating a GitHub release.
type GitHubRelease struct {
	ID        int    `json:"id"`
	UploadURL string `json:"upload_url"`
}

// CreateGitHubRelease creates a release via the GitHub API and returns the release ID.
func CreateGitHubRelease(owner, repo, tag, name, body, token string, draft bool) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)

	payload := map[string]interface{}{
		"tag_name": tag,
		"name":     name,
		"body":     body,
		"draft":    draft,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal release payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	setGitHubHeaders(req, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(respBody))
	}

	var release GitHubRelease
	err = json.NewDecoder(resp.Body).Decode(&release)

	return &release, err
}

// setGitHubHeaders sets common headers for GitHub API requests.
func setGitHubHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
}
