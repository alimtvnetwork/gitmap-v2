// Package gitutil extracts Git metadata by running git commands.
package gitutil

import (
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
)

// RemoteURL returns the origin remote URL for a repo at the given path.
func RemoteURL(repoPath string) (string, error) {
	out, err := runGit(repoPath,
		constants.GitConfigCmd, constants.GitGetFlag, constants.GitRemoteOrigin)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

// CurrentBranch returns the current branch name for a repo.
func CurrentBranch(repoPath string) (string, error) {
	out, err := runGit(repoPath,
		constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if err != nil {
		return constants.DefaultBranch, err
	}

	return strings.TrimSpace(out), nil
}

// runGit executes a git command in the given directory and returns stdout.
func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
