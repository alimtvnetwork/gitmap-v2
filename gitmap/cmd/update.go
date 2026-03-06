package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/verbose"
)

// runUpdate handles the "update" subcommand.
// It creates a handoff copy and runs a hidden worker command from that copy.
func runUpdate() {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrNoRepoPath)
		os.Exit(1)
	}

	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateExecFind, err)
		os.Exit(1)
	}

	copyPath := createHandoffCopy(selfPath)
	fmt.Printf(constants.MsgUpdateActive, selfPath, copyPath)
	launchHandoff(copyPath)
}

// createHandoffCopy creates a temporary copy of the binary for handoff.
func createHandoffCopy(selfPath string) string {
	name := fmt.Sprintf(constants.UpdateCopyFmt, os.Getpid())
	copyPath := filepath.Join(filepath.Dir(selfPath), name)
	if copyFile(selfPath, copyPath) == nil {
		return copyPath
	}

	fallbackPath := filepath.Join(os.TempDir(), name)
	if err := copyFile(selfPath, fallbackPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateCopyFail, err)
		os.Exit(1)
	}

	return fallbackPath
}

// launchHandoff runs the handoff binary with update-runner command.
func launchHandoff(copyPath string) {
	copyArgs := []string{constants.CmdUpdateRunner}
	if hasFlag(constants.FlagVerbose) {
		copyArgs = append(copyArgs, constants.FlagVerbose)
	}

	cmd := exec.Command(copyPath, copyArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		handleHandoffError(err)
	}
}

// handleHandoffError exits with the handoff process exit code if available.
func handleHandoffError(err error) {
	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}

	fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
	os.Exit(1)
}

// runUpdateRunner is a hidden command that performs the real update work.
func runUpdateRunner() {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrNoRepoPath)
		os.Exit(1)
	}

	initRunnerVerbose()
	fmt.Printf(constants.MsgUpdateStarting)
	fmt.Printf(constants.MsgUpdateRepoPath, repoPath)
	executeUpdate(repoPath)
}

// initRunnerVerbose initializes verbose logging if --verbose flag is present.
func initRunnerVerbose() {
	if hasFlag(constants.FlagVerbose) {
		log, err := verbose.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)
		} else {
			defer log.Close()
			log.Log(constants.UpdateRunnerLogStart, constants.RepoPath)
		}
	}
}

// hasFlag checks if a flag is present in os.Args[2:].
func hasFlag(name string) bool {
	for _, arg := range os.Args[2:] {
		if arg == name {
			return true
		}
	}

	return false
}

// copyFile copies src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}
