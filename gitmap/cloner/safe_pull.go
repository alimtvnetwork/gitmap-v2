// Package cloner re-clones repos from structured files.
package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/verbose"
)

var (
	unlinkOldRegex    = regexp.MustCompile(`(?i)unable to unlink old '([^']+)'`)
	unlinkPromptRegex = regexp.MustCompile(`(?i)unlink of file '([^']+)' failed`)
)

func cloneOrPullOne(rec model.ScanRecord, targetDir string, safePull bool) model.CloneResult {
	dest := filepath.Join(targetDir, rec.RelativePath)
	if safePull && isGitRepo(dest) {
		return safePullRepo(rec, dest)
	}

	return cloneOne(rec, targetDir)
}

func isGitRepo(path string) bool {
	return IsGitRepo(path)
}

// IsGitRepo checks whether the given path contains a .git directory.
func IsGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))

	return err == nil
}

// SafePullOne runs safe-pull on a single repo. Exported for use by the pull command.
func SafePullOne(rec model.ScanRecord, repoDir string) model.CloneResult {
	return safePullRepo(rec, repoDir)
}

func safePullRepo(rec model.ScanRecord, repoDir string) model.CloneResult {
	log := verbose.Get()
	if log != nil {
		log.Log("safe-pull starting: %s → %s", rec.RepoName, repoDir)
	}

	var lastError string
	for attempt := 1; attempt <= constants.SafePullRetryAttempts; attempt++ {
		output, err := runGitPull(repoDir)
		if log != nil {
			log.Log("pull attempt %d/%d for %s: exit=%v output=%s",
				attempt, constants.SafePullRetryAttempts, rec.RepoName, err, trimOutput(output))
		}
		if err == nil {
			if log != nil {
				log.Log("safe-pull succeeded: %s (attempt %d)", rec.RepoName, attempt)
			}
			return model.CloneResult{Record: rec, Success: true}
		}

		cleared := clearReadOnlyAttrs(repoDir, output)
		if log != nil && cleared {
			log.Log("cleared read-only attributes for blocked files in %s", repoDir)
		}
		diagnosis := buildPullDiagnosis(repoDir, output)
		if log != nil {
			log.Log("diagnosis for %s: %s", rec.RepoName, diagnosis)
		}
		lastError = fmt.Sprintf(
			"safe-pull failed (attempt %d/%d): %v\n%s\nDiagnosis: %s",
			attempt,
			constants.SafePullRetryAttempts,
			err,
			trimOutput(output),
			diagnosis,
		)

		if attempt < constants.SafePullRetryAttempts {
			time.Sleep(time.Duration(constants.SafePullRetryDelayMS) * time.Millisecond)
		}
	}

	if log != nil {
		log.Log("safe-pull FAILED after all retries: %s — %s", rec.RepoName, lastError)
	}

	return model.CloneResult{Record: rec, Success: false, Error: lastError}
}

func runGitPull(repoDir string) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, repoDir, constants.GitPull, constants.GitFFOnlyFlag)
	out, err := cmd.CombinedOutput()

	return string(out), err
}

func clearReadOnlyAttrs(repoDir, output string) bool {
	if runtime.GOOS != constants.OSWindows {
		return false
	}

	paths := extractUnlinkPaths(output)
	if len(paths) == 0 {
		return false
	}

	cleared := false
	for _, relativePath := range paths {
		fullPath := filepath.Join(repoDir, filepath.FromSlash(relativePath))
		if clearReadOnly(fullPath) {
			cleared = true
		}
	}

	return cleared
}

func clearReadOnly(path string) bool {
	cmd := exec.Command("attrib", "-R", path)
	if err := cmd.Run(); err == nil {
		return true
	}

	return os.Chmod(path, 0o666) == nil
}

func buildPullDiagnosis(repoDir, output string) string {
	hints := make([]string, 0, 3)
	if hasUnlinkFailure(output) {
		hints = append(hints, "file lock/read-only attribute blocked replacing old files")
	}
	if hasPathLengthRisk(repoDir, output) {
		hints = append(hints, "Windows path length risk detected; use a shorter base path like C:\\src")
	}
	if strings.Contains(strings.ToLower(repoDir), "onedrive") {
		hints = append(hints, "repo is under a synced folder (OneDrive), which often locks files")
	}
	if len(hints) == 0 {
		hints = append(hints, "non-unlink git pull failure (check auth/merge or run pull manually for full output)")
	}

	return strings.Join(hints, "; ")
}

func hasUnlinkFailure(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "unable to unlink old") || strings.Contains(lower, "unlink of file")
}

func hasPathLengthRisk(repoDir, output string) bool {
	if runtime.GOOS != constants.OSWindows {
		return false
	}
	for _, relativePath := range extractUnlinkPaths(output) {
		fullPath := filepath.Join(repoDir, filepath.FromSlash(relativePath))
		if len(fullPath) >= constants.WindowsPathWarnThreshold {
			return true
		}
	}

	return false
}

func extractUnlinkPaths(output string) []string {
	matches := make([]string, 0, 2)
	for _, m := range unlinkOldRegex.FindAllStringSubmatch(output, -1) {
		if len(m) > 1 {
			matches = append(matches, m[1])
		}
	}
	for _, m := range unlinkPromptRegex.FindAllStringSubmatch(output, -1) {
		if len(m) > 1 {
			matches = append(matches, m[1])
		}
	}

	seen := map[string]struct{}{}
	unique := make([]string, 0, len(matches))
	for _, match := range matches {
		if _, ok := seen[match]; ok {
			continue
		}
		seen[match] = struct{}{}
		unique = append(unique, match)
	}

	return unique
}

func trimOutput(output string) string {
	trimmed := strings.TrimSpace(output)
	if len(trimmed) <= 1200 {
		return trimmed
	}

	return trimmed[:1200] + "..."
}
