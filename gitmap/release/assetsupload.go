// Package release — assetsupload.go uploads release assets via GitHub API.
package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// GitHubRelease represents the response from creating a GitHub release.
type GitHubRelease struct {
	ID        int    `json:"id"`
	UploadURL string `json:"upload_url"`
}

// CreateGitHubRelease creates a release via the GitHub API and returns the release ID.
func CreateGitHubRelease(owner, repo, tag, body, token string, draft bool) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)

	payload := map[string]interface{}{
		"tag_name": tag,
		"name":     constants.ReleaseTagPrefix + tag,
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

// UploadAsset uploads a single file to a GitHub release.
func UploadAsset(owner, repo string, releaseID int, filePath, token string) error {
	filename := filepath.Base(filePath)
	url := fmt.Sprintf("https://uploads.github.com/repos/%s/%s/releases/%d/assets?name=%s",
		owner, repo, releaseID, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open asset: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodPost, url, file)
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("upload error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// UploadAllAssets uploads all assets to a GitHub release with retry.
func UploadAllAssets(owner, repo string, releaseID int, assets []string, token string) {
	for _, asset := range assets {
		err := UploadAsset(owner, repo, releaseID, asset, token)
		if err != nil {
			// Retry once on failure.
			fmt.Fprintf(os.Stderr, constants.ErrAssetUploadRetry, filepath.Base(asset))

			retryErr := UploadAsset(owner, repo, releaseID, asset, token)
			if retryErr != nil {
				fmt.Fprintf(os.Stderr, constants.ErrAssetUploadFailed, filepath.Base(asset), retryErr)

				continue
			}
		}

		fmt.Printf(constants.MsgAssetUploaded, filepath.Base(asset))
	}
}

// ParseRemoteOrigin extracts owner/repo from the git remote origin URL.
func ParseRemoteOrigin() (string, string, error) {
	url := getRemoteURL()
	if len(url) == 0 {
		return "", "", fmt.Errorf("no remote origin URL found")
	}

	return parseGitURL(url)
}

// getRemoteURL reads the origin remote URL via git config.
func getRemoteURL() string {
	out, err := gitOutput("config", "--get", "remote.origin.url")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(out)
}

// parseGitURL extracts owner/repo from HTTPS or SSH git URLs.
func parseGitURL(url string) (string, string, error) {
	// HTTPS: https://github.com/owner/repo.git
	if strings.HasPrefix(url, "https://") {
		return parseHTTPSURL(url)
	}

	// SSH: git@github.com:owner/repo.git
	if strings.Contains(url, "@") {
		return parseSSHURL(url)
	}

	return "", "", fmt.Errorf("unrecognized remote URL format: %s", url)
}

// parseHTTPSURL parses https://github.com/owner/repo.git
func parseHTTPSURL(url string) (string, string, error) {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return "", "", fmt.Errorf("invalid HTTPS remote: %s", url)
	}

	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	return owner, repo, nil
}

// parseSSHURL parses git@github.com:owner/repo.git
func parseSSHURL(url string) (string, string, error) {
	url = strings.TrimSuffix(url, ".git")
	colonIdx := strings.LastIndex(url, ":")

	if colonIdx < 0 {
		return "", "", fmt.Errorf("invalid SSH remote: %s", url)
	}

	path := url[colonIdx+1:]
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid SSH remote path: %s", url)
	}

	return parts[0], parts[1], nil
}

// setGitHubHeaders sets common headers for GitHub API requests.
func setGitHubHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
}

// gitOutput runs a git command and returns stdout.
func gitOutput(args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	out, err := cmd.Output()

	return string(out), err
}
