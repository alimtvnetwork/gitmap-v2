// Package cmd — amendexec.go handles git operations for the amend command.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// listCommitsForAmend returns commits that will be rewritten.
func listCommitsForAmend(f amendFlags) []model.CommitEntry {
	var args []string

	if f.commitHash == "" {
		args = []string{"log", "--format=%H %s", "--reverse"}
	} else if f.commitHash == "HEAD" {
		args = []string{"log", "--format=%H %s", "-1"}
	} else {
		args = []string{"log", "--format=%H %s", "--reverse", f.commitHash + "^..HEAD"}
	}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendListCommits, err)

		return nil
	}

	return parseCommitLines(string(out))
}

// parseCommitLines splits git log output into CommitEntry slices.
func parseCommitLines(output string) []model.CommitEntry {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var entries []model.CommitEntry

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		msg := ""
		if len(parts) > 1 {
			msg = parts[1]
		}

		entries = append(entries, model.CommitEntry{
			SHA:     parts[0],
			Message: msg,
		})
	}

	return entries
}

// detectPreviousAuthor reads the author of the first commit in the range.
func detectPreviousAuthor(commits []model.CommitEntry) (string, string) {
	if len(commits) == 0 {
		return "", ""
	}

	sha := commits[0].SHA
	nameOut, _ := exec.Command("git", "log", "-1", "--format=%an", sha).Output()
	emailOut, _ := exec.Command("git", "log", "-1", "--format=%ae", sha).Output()

	return strings.TrimSpace(string(nameOut)), strings.TrimSpace(string(emailOut))
}

// getCurrentBranch returns the current Git branch name.
func getCurrentBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "main"
	}

	return strings.TrimSpace(string(out))
}

// switchBranch checks out the specified branch.
func switchBranch(branch string) {
	fmt.Printf(constants.MsgAmendCheckout, branch)

	cmd := exec.Command("git", "checkout", branch)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendCheckout, branch, err)
		os.Exit(1)
	}
}

// runFilterBranch executes the git filter-branch command.
func runFilterBranch(f amendFlags, commits []model.CommitEntry) {
	if f.commitHash == "HEAD" {
		runAmendHead(f)

		return
	}

	envFilter := buildEnvFilter(f)
	var args []string

	if f.commitHash == "" {
		args = []string{"filter-branch", "-f", "--env-filter", envFilter, "--", "HEAD"}
	} else {
		args = []string{"filter-branch", "-f", "--env-filter", envFilter, f.commitHash + "^..HEAD"}
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendFilter, err)
		os.Exit(1)
	}
}

// runAmendHead uses git commit --amend for single HEAD commit.
func runAmendHead(f amendFlags) {
	author := buildAuthorString(f)
	args := []string{"commit", "--amend", "--no-edit", "--author", author}

	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendCommitAmend, err)
		os.Exit(1)
	}
}

// buildEnvFilter constructs the env-filter script for filter-branch.
func buildEnvFilter(f amendFlags) string {
	var lines []string

	if f.name != "" {
		lines = append(lines, "export GIT_AUTHOR_NAME='"+f.name+"'")
		lines = append(lines, "export GIT_COMMITTER_NAME='"+f.name+"'")
	}

	if f.email != "" {
		lines = append(lines, "export GIT_AUTHOR_EMAIL='"+f.email+"'")
		lines = append(lines, "export GIT_COMMITTER_EMAIL='"+f.email+"'")
	}

	return strings.Join(lines, "\n")
}

// buildAuthorString creates the --author flag value.
func buildAuthorString(f amendFlags) string {
	name := f.name
	email := f.email

	if name == "" {
		out, _ := exec.Command("git", "config", "user.name").Output()
		name = strings.TrimSpace(string(out))
	}

	if email == "" {
		out, _ := exec.Command("git", "config", "user.email").Output()
		email = strings.TrimSpace(string(out))
	}

	return name + " <" + email + ">"
}

// runForcePush executes git push --force-with-lease.
func runForcePush() {
	cmd := exec.Command("git", "push", "--force-with-lease")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendForcePush, err)

		return
	}

	fmt.Print(constants.MsgAmendForcePush)
}

// printAmendHeader outputs the operation header.
func printAmendHeader(f amendFlags, commits []model.CommitEntry, branch, prevName, prevEmail string) {
	if f.commitHash == "" {
		fmt.Printf(constants.MsgAmendHeaderAll, len(commits), branch)
	} else {
		fmt.Printf(constants.MsgAmendHeader, len(commits), commits[0].SHA[:7], commits[len(commits)-1].SHA[:7], branch)
	}

	oldAuthor := prevName + " <" + prevEmail + ">"
	newAuthor := buildAuthorString(f)
	fmt.Printf(constants.MsgAmendAuthor, oldAuthor, newAuthor)
}

// printAmendProgress outputs per-commit progress lines.
func printAmendProgress(commits []model.CommitEntry) {
	for i, c := range commits {
		sha := c.SHA
		if len(sha) > 7 {
			sha = sha[:7]
		}

		fmt.Printf(constants.MsgAmendProgress, i+1, len(commits), sha, c.Message)
	}
}

// printAmendDryRun outputs dry-run preview.
func printAmendDryRun(commits []model.CommitEntry, f amendFlags) {
	fmt.Printf(constants.MsgAmendDryHeader, len(commits))

	for i, c := range commits {
		sha := c.SHA
		if len(sha) > 7 {
			sha = sha[:7]
		}

		pn, pe := detectPreviousAuthor([]model.CommitEntry{c})
		fmt.Printf(constants.MsgAmendDryLine, i+1, sha, c.Message, pn, pe)
	}

	fmt.Print(constants.MsgAmendDrySkip)
}
