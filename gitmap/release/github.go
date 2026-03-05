// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
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

// ghHTTPRelease creates a release via gh api as a fallback.
func ghHTTPRelease(tag, body string, draft bool, token string) error {
	draftStr := "false"
	if draft {
		draftStr = "true"
	}
	notesBody := body
	if len(notesBody) == 0 {
		notesBody = "Release " + tag
	}

	args := []string{
		"api", "repos/{owner}/{repo}/releases",
		"--method", "POST",
		"-f", "tag_name=" + tag,
		"-f", "name=" + tag,
		"-f", "body=" + notesBody,
		"-f", "draft=" + draftStr,
	}

	cmd := exec.Command(constants.GHBin, args...)
	cmd.Env = append(os.Environ(), "GITHUB_TOKEN="+token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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
