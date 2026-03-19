package release

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
)

// AutoCommitResult describes what happened during auto-commit.
type AutoCommitResult struct {
	Committed bool
	AllFiles  bool
	Message   string
}

// AutoCommit inspects working tree changes after returning to the original branch.
// If only .release/ files changed, it commits and pushes silently.
// If other files also changed, it prompts the user. On decline, it commits only .release/.
func AutoCommit(version string, dryRun bool) AutoCommitResult {
	fmt.Print(constants.MsgAutoCommitScanning)

	if dryRun {
		fmt.Print(constants.MsgAutoCommitDryRun)

		return AutoCommitResult{}
	}

	changed := listChangedFiles()
	if len(changed) == 0 {
		fmt.Print(constants.MsgAutoCommitNone)

		return AutoCommitResult{}
	}

	releaseFiles, otherFiles := classifyFiles(changed)

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

// classifyFiles separates .release/ files from everything else.
func classifyFiles(files []string) (releaseFiles, otherFiles []string) {
	for _, f := range files {
		if strings.HasPrefix(f, constants.DefaultReleaseDir+"/") || f == constants.DefaultReleaseDir {
			releaseFiles = append(releaseFiles, f)
		} else {
			otherFiles = append(otherFiles, f)
		}
	}

	return releaseFiles, otherFiles
}

// commitReleaseOnly stages and commits only .release/ files.
func commitReleaseOnly(files []string, msg string) AutoCommitResult {
	err := stageFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	err = commitStaged(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	fmt.Printf(constants.MsgAutoCommitReleaseOnly, msg)

	err = pushCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitPush, err)

		return AutoCommitResult{Committed: true, Message: msg}
	}

	branch, _ := CurrentBranchName()
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

	err = commitStaged(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	fmt.Printf(constants.MsgAutoCommitAll, msg)

	err = pushCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitPush, err)

		return AutoCommitResult{Committed: true, AllFiles: true, Message: msg}
	}

	branch, _ := CurrentBranchName()
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

	return runGitCmd(constants.GitPush, constants.GitOrigin, branch)
}
