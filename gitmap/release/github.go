// Package release handles version parsing, release workflows,
// and release metadata management.
package release

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/verbose"
)

// CollectAssets gathers file paths for release attachment.
func CollectAssets(assetsPath string) []string {
	if len(assetsPath) == 0 {
		return nil
	}

	info, err := os.Stat(assetsPath)
	if err != nil {
		if verbose.IsEnabled() {
			verbose.Get().Log("assets: path not found: %s", assetsPath)
		}
		return nil
	}

	if info.IsDir() {
		files := collectDirFiles(assetsPath)
		if verbose.IsEnabled() {
			verbose.Get().Log("assets: collected %d file(s) from directory %s", len(files), assetsPath)
			for _, f := range files {
				verbose.Get().Log("assets: %s", filepath.Base(f))
			}
		}
		return files
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("assets: single file %s", assetsPath)
	}

	return []string{assetsPath}
}

// collectDirFiles returns all file paths in a directory.
func collectDirFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	files := make([]string, 0, len(entries))
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
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
