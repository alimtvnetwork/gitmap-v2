package release

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/verbose"
)

// AutoCommitResult describes what happened during auto-commit.
type AutoCommitResult struct {
	Committed bool
	AllFiles  bool
	Message   string
}

// AutoCommit inspects working tree changes after returning to the original branch.
// If only .gitmap/release/ files changed, it commits and pushes silently.
// If other files also changed, it prompts the user. On decline, it commits only .gitmap/release/.
func AutoCommit(version string, dryRun bool) AutoCommitResult {
	fmt.Print(constants.MsgAutoCommitScanning)

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: starting for %s (dry-run=%v)", version, dryRun)
	}

	if dryRun {
		fmt.Print(constants.MsgAutoCommitDryRun)

		return AutoCommitResult{}
	}

	changed := listChangedFiles()
	if len(changed) == 0 {
		fmt.Print(constants.MsgAutoCommitNone)

		if verbose.IsEnabled() {
			verbose.Get().Log("autocommit: no changed files detected")
		}

		return AutoCommitResult{}
	}

	releaseFiles, otherFiles := classifyFiles(changed)

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: %d release file(s), %d other file(s)", len(releaseFiles), len(otherFiles))
	}

	commitMsg := fmt.Sprintf(constants.AutoCommitMsgFmt, version)

	if len(otherFiles) == 0 {
		return commitReleaseOnly(releaseFiles, commitMsg)
	}

	return promptAndCommit(releaseFiles, otherFiles, commitMsg)
}

// listChangedFiles returns all modified/untracked files in the working tree.
func listChangedFiles() []string {
	cmd := exec.Command(constants.GitBin, constants.GitStatus, constants.GitStatusShort)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return parsePorcelainOutput(string(out))
}

// parsePorcelainOutput extracts file paths from git status --porcelain output.
func parsePorcelainOutput(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var files []string

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		path := strings.TrimSpace(line[3:])
		if len(path) > 0 {
			files = append(files, path)
		}
	}

	return files
}

// classifyFiles separates .gitmap/release/ files (and legacy .release/ files) from everything else.
func classifyFiles(files []string) (releaseFiles, otherFiles []string) {
	for _, f := range files {
		if strings.HasPrefix(f, constants.DefaultReleaseDir+"/") || f == constants.DefaultReleaseDir ||
			strings.HasPrefix(f, constants.LegacyReleaseDir+"/") || f == constants.LegacyReleaseDir {
			releaseFiles = append(releaseFiles, f)
		} else {
			otherFiles = append(otherFiles, f)
		}
	}

	return releaseFiles, otherFiles
}

// commitReleaseOnly stages and commits only .gitmap/release/ files.
func commitReleaseOnly(files []string, msg string) AutoCommitResult {
	err := stageFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: staged %d release file(s)", len(files))
	}

	err = commitStaged(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: committed release-only: %s", msg)
	}

	fmt.Printf(constants.MsgAutoCommitReleaseOnly, msg)

	err = pushCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitPush, err)

		return AutoCommitResult{Committed: true, Message: msg}
	}

	branch, _ := CurrentBranchName()

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: pushed to %s", branch)
	}

	fmt.Printf(constants.MsgAutoCommitPushed, branch)

	return AutoCommitResult{Committed: true, Message: msg}
}

// promptAndCommit shows changed files and asks the user whether to commit all.
func promptAndCommit(releaseFiles, otherFiles []string, msg string) AutoCommitResult {
	fmt.Print(constants.MsgAutoCommitPrompt)

	for _, f := range otherFiles {
		fmt.Printf(constants.MsgAutoCommitFile, f)
	}

	fmt.Print(constants.MsgAutoCommitAsk)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return commitReleaseOnly(releaseFiles, msg)
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer == "y" || answer == "yes" {
		return commitAll(msg)
	}

	if len(releaseFiles) > 0 {
		result := commitReleaseOnly(releaseFiles, msg)
		fmt.Printf(constants.MsgAutoCommitPartial, msg)

		return result
	}

	return AutoCommitResult{}
}

// commitAll stages everything and commits.
func commitAll(msg string) AutoCommitResult {
	err := stageAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: staged all files")
	}

	err = commitStaged(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: committed all: %s", msg)
	}

	fmt.Printf(constants.MsgAutoCommitAll, msg)

	err = pushCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitPush, err)

		return AutoCommitResult{Committed: true, AllFiles: true, Message: msg}
	}

	branch, _ := CurrentBranchName()

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: pushed all to %s", branch)
	}

	fmt.Printf(constants.MsgAutoCommitPushed, branch)

	return AutoCommitResult{Committed: true, AllFiles: true, Message: msg}
}

// stageFiles runs git add on specific files.
func stageFiles(files []string) error {
	args := append([]string{constants.GitAdd}, files...)

	return runGitCmd(args...)
}

// stageAll runs git add -A.
func stageAll() error {
	return runGitCmd(constants.GitAdd, constants.GitAddAll)
}

// commitStaged runs git commit -m <msg>.
func commitStaged(msg string) error {
	return runGitCmd(constants.GitCommit, constants.GitCommitMsg, msg)
}

// pushCurrentBranch pushes the current branch to origin.
func pushCurrentBranch() error {
	branch, err := CurrentBranchName()
	if err != nil {
		return err
	}

	pushOutput, err := runGitCmdCombined(constants.GitPush, constants.GitOrigin, branch)
	if err == nil {
		return nil
	}

	if !isNonFastForwardPushError(pushOutput) {
		return formatGitCommandError(pushOutput, err)
	}

	return syncBranchAndRetryPush(branch, pushOutput)
}

func syncBranchAndRetryPush(branch, pushOutput string) error {
	if verbose.IsEnabled() {
		verbose.Get().Log(
			"autocommit: push rejected for %s, attempting rebase sync: %s",
			branch,
			singleLineGitOutput(pushOutput),
		)
	}

	fmt.Printf(constants.MsgAutoCommitSyncRetry, branch)

	pullOutput, err := runGitCmdCombined(
		constants.GitPull,
		constants.GitPullRebaseFlag,
		constants.GitOrigin,
		branch,
	)
	if err != nil {
		abortRebaseAfterFailure()

		return fmt.Errorf(
			"remote branch advanced; pull --rebase failed: %s",
			trimGitOutput(pullOutput),
		)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: rebase sync completed for %s", branch)
	}

	retryOutput, err := runGitCmdCombined(constants.GitPush, constants.GitOrigin, branch)
	if err != nil {
		return fmt.Errorf(
			"push retry after rebase failed: %s",
			trimGitOutput(retryOutput),
		)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: push retry succeeded for %s", branch)
	}

	return nil
}

func runGitCmdCombined(args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	out, err := cmd.CombinedOutput()

	return string(out), err
}

func isNonFastForwardPushError(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "fetch first") ||
		strings.Contains(lower, "non-fast-forward") ||
		strings.Contains(lower, "failed to push some refs")
}

func formatGitCommandError(output string, err error) error {
	trimmed := trimGitOutput(output)
	if len(trimmed) > 0 {
		return fmt.Errorf("%s", trimmed)
	}

	return err
}

func trimGitOutput(output string) string {
	trimmed := strings.TrimSpace(output)
	if len(trimmed) > 0 {
		return trimmed
	}

	return "unknown git error"
}

func singleLineGitOutput(output string) string {
	return strings.Join(strings.Fields(trimGitOutput(output)), " ")
}

func abortRebaseAfterFailure() {
	_, err := runGitCmdCombined(constants.GitRebase, constants.GitRebaseAbortFlag)
	if err != nil {
		return
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: aborted failed rebase")
	}
}
