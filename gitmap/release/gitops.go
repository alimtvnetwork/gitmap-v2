// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	return runGitCmd(constants.GitTag, "-a", tag, "-m", message)
}

// PushBranchAndTag pushes the branch and tag to origin.
func PushBranchAndTag(branchName, tag string) error {
	err := runGitCmd(constants.GitPush, "origin", branchName)
	if err != nil {
		return fmt.Errorf("push branch: %w", err)
	}

	err = runGitCmd(constants.GitPush, "origin", tag)
	if err != nil {
		return fmt.Errorf("push tag: %w", err)
	}

	return nil
}

// TagExistsLocally checks if a git tag exists in the local repo.
func TagExistsLocally(tag string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitTag, "-l", tag)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// TagExistsRemote checks if a git tag exists on the remote.
func TagExistsRemote(tag string) bool {
	cmd := exec.Command(constants.GitBin,
		constants.GitLsRemote, constants.GitLsRemoteTags, "origin", tag)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// BranchExists checks if a local branch exists.
func BranchExists(branch string) bool {
	cmd := exec.Command(constants.GitBin, "branch", "--list", branch)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// CheckoutBranch checks out an existing branch.
func CheckoutBranch(branch string) error {
	return runGitCmd(constants.GitCheckout, branch)
}

// CurrentCommitSHA returns the full SHA of HEAD.
func CurrentCommitSHA() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, constants.GitHEAD)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// CurrentBranchName returns the current branch name.
func CurrentBranchName() (string, error) {
	cmd := exec.Command(constants.GitBin,
		constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// CommitExists checks if a commit SHA is valid.
func CommitExists(sha string) bool {
	cmd := exec.Command(constants.GitBin, "cat-file", "-t", sha)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "commit"
}

// FetchBranch fetches the latest of a remote branch.
func FetchBranch(branch string) error {
	return runGitCmd("fetch", "origin", branch)
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
	if CommitExists(commit) == false {
		return "", "", fmt.Errorf("commit %s not found", commit)
	}

	return commit, "commit:" + commit, nil
}

// resolveFromBranch fetches and returns the branch tip.
func resolveFromBranch(branch string) (string, string, error) {
	err := FetchBranch(branch)
	if err != nil {
		return "", "", fmt.Errorf("branch %s not found: %w", branch, err)
	}

	return "origin/" + branch, branch, nil
}

// resolveFromHead returns HEAD as the source ref.
func resolveFromHead() (string, string, error) {
	branchName, err := CurrentBranchName()
	if err != nil {
		return constants.GitHEAD, "HEAD", nil
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
