// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// GitHubRelease creates a GitHub release using gh CLI or HTTP fallback.
func GitHubRelease(tag, body string, assets []string, draft bool) error {
	if ghAvailable() {
		return ghRelease(tag, body, assets, draft)
	}

	fmt.Fprint(os.Stderr, constants.ErrReleaseGHNotFound)

	if tryInstallGH() {
		fmt.Print(constants.MsgReleaseGHInstalled)
		return ghRelease(tag, body, assets, draft)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if len(token) == 0 {
		return fmt.Errorf(constants.ErrReleaseGHTokenMissing)
	}

	return ghHTTPRelease(tag, body, draft, token)
}

// ghAvailable checks if the gh CLI is installed.
func ghAvailable() bool {
	_, err := exec.LookPath(constants.GHBin)

	return err == nil
}

// ghRelease creates a release using the gh CLI.
func ghRelease(tag, body string, assets []string, draft bool) error {
	args := buildGHArgs(tag, body, draft)
	args = appendAssetArgs(args, assets)

	cmd := exec.Command(constants.GHBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// buildGHArgs constructs the base gh release create arguments.
func buildGHArgs(tag, body string, draft bool) []string {
	args := []string{"release", "create", tag, "--title", tag}
	if draft {
		args = append(args, "--draft")
	}
	if len(body) > 0 {
		args = append(args, "--notes", body)
	} else {
		args = append(args, "--generate-notes")
	}

	return args
}

// appendAssetArgs adds file attachment arguments for gh CLI.
func appendAssetArgs(args []string, assets []string) []string {
	for _, asset := range assets {
		args = append(args, asset)
	}

	return args
}

// ghHTTPRelease creates a release via the GitHub REST API using net/http.
func ghHTTPRelease(tag, body string, draft bool, token string) error {
	repoSlug, err := detectRepoSlug()
	if err != nil {
		return fmt.Errorf("could not detect GitHub repo: %w", err)
	}

	notesBody := body
	if len(notesBody) == 0 {
		notesBody = "Release " + tag
	}

	payload := map[string]interface{}{
		"tag_name": tag,
		"name":     tag,
		"body":     notesBody,
		"draft":    draft,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := "https://api.github.com/repos/" + repoSlug + "/releases"
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// detectRepoSlug extracts "owner/repo" from the git remote origin URL.
func detectRepoSlug() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitConfigCmd, constants.GitGetFlag, constants.GitRemoteOrigin)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("no remote origin configured")
	}

	remote := strings.TrimSpace(string(out))

	// Handle HTTPS: https://github.com/owner/repo.git
	if strings.HasPrefix(remote, "https://") {
		remote = strings.TrimPrefix(remote, "https://github.com/")
		remote = strings.TrimSuffix(remote, ".git")
		return remote, nil
	}

	// Handle SSH: git@github.com:owner/repo.git
	if strings.Contains(remote, ":") {
		parts := strings.SplitN(remote, ":", 2)
		slug := strings.TrimSuffix(parts[1], ".git")
		return slug, nil
	}

	return "", fmt.Errorf("unrecognized remote format: %s", remote)
}

// CollectAssets gathers file paths for release attachment.
func CollectAssets(assetsPath string) []string {
	if len(assetsPath) == 0 {
		return nil
	}

	info, err := os.Stat(assetsPath)
	if err != nil {
		return nil
	}

	if info.IsDir() == false {
		return []string{assetsPath}
	}

	return collectDirFiles(assetsPath)
}

// collectDirFiles returns all file paths in a directory.
func collectDirFiles(dir string) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() == false {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return files
}

// DetectChangelog returns the content of CHANGELOG.md if it exists.
func DetectChangelog() string {
	data, err := os.ReadFile(constants.ChangelogFile)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// DetectReadme returns the path to README.md if it exists.
func DetectReadme() string {
	_, err := os.Stat(constants.ReadmeFile)
	if err != nil {
		return ""
	}

	return constants.ReadmeFile
}
