// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/user/gitmap/constants"
)

// CreateBranch creates a release branch from the given source ref.
func CreateBranch(branchName, sourceRef string) error {
	args := []string{constants.GitCheckout, constants.GitBranchFlag, branchName}
	if len(sourceRef) > 0 {
		args = append(args, sourceRef)
	}

	return runGitCmd(args...)
}

// CreateTag creates an annotated git tag.
func CreateTag(tag, message string) error {
	return runGitCmd(constants.GitTag, constants.GitTagAnnotateFlag, tag, constants.GitTagMessageFlag, message)
}

// PushBranchAndTag pushes the branch and tag to origin.
func PushBranchAndTag(branchName, tag string) error {
	err := runGitCmd(constants.GitPush, constants.GitOrigin, branchName)
	if err != nil {
		return fmt.Errorf("push branch: %w", err)
	}

	err = runGitCmd(constants.GitPush, constants.GitOrigin, tag)
	if err != nil {
		return fmt.Errorf("push tag: %w", err)
	}

	return nil
}

// CheckoutBranch checks out an existing branch.
func CheckoutBranch(branch string) error {
	return runGitCmd(constants.GitCheckout, branch)
}

// FetchBranch fetches the latest of a remote branch.
func FetchBranch(branch string) error {
	return runGitCmd(constants.GitFetch, constants.GitOrigin, branch)
}

// ResolveSourceRef returns the ref to use as the release base.
func ResolveSourceRef(commit, branch string) (string, string, error) {
	if len(commit) > 0 {
		return resolveFromCommit(commit)
	}
	if len(branch) > 0 {
		return resolveFromBranch(branch)
	}

	return resolveFromHead()
}

// resolveFromCommit validates and returns the commit ref.
func resolveFromCommit(commit string) (string, string, error) {
	if CommitExists(commit) {
		return commit, constants.GitCommitPrefix + commit, nil
	}

	return "", "", fmt.Errorf("commit %s not found", commit)
}

// resolveFromBranch fetches and returns the branch tip.
func resolveFromBranch(branch string) (string, string, error) {
	err := FetchBranch(branch)
	if err != nil {
		return "", "", fmt.Errorf("branch %s not found: %w", branch, err)
	}

	return constants.GitOriginPrefix + branch, branch, nil
}

// resolveFromHead returns HEAD as the source ref.
func resolveFromHead() (string, string, error) {
	branchName, err := CurrentBranchName()
	if err != nil {
		return constants.GitHEAD, constants.GitHEAD, nil
	}

	return constants.GitHEAD, branchName, nil
}

// runGitCmd executes a git command and prints stderr on failure.
func runGitCmd(args ...string) error {
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
