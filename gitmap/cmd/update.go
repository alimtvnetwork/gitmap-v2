package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
	"github.com/user/gitmap/verbose"
)

// runUpdate handles the "update" subcommand.
// It creates a handoff copy and runs a hidden worker command from that copy.
func runUpdate() {
	requireOnline()
	repoPath := resolveRepoPath()

	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateExecFind, err)
		os.Exit(1)
	}

	copyPath := createHandoffCopy(selfPath)
	fmt.Printf(constants.MsgUpdateActive, selfPath, copyPath)
	launchHandoff(copyPath, repoPath)
}

// resolveRepoPath returns the repo path from --repo-path flag or embedded constant.
// If neither is available, it attempts to delegate to gitmap-updater.
func resolveRepoPath() string {
	if flagVal := getFlagValue(constants.FlagRepoPath); len(flagVal) > 0 {
		saveRepoPathToDB(flagVal)

		return flagVal
	}

	if len(constants.RepoPath) > 0 && pathExists(constants.RepoPath) {
		return constants.RepoPath
	}

	if saved := loadRepoPathFromDB(); len(saved) > 0 && pathExists(saved) {
		return saved
	}

	if prompted := promptRepoPath(); len(prompted) > 0 {
		saveRepoPathToDB(prompted)

		return prompted
	}

	// Try to fall back to gitmap-updater for release-based update
	if tryUpdaterFallback() {
		os.Exit(0)
	}

	fmt.Fprint(os.Stderr, constants.ErrNoRepoPath)
	os.Exit(1)

	return ""
}

// tryUpdaterFallback looks for gitmap-updater on PATH and launches it.
func tryUpdaterFallback() bool {
	updaterPath, err := exec.LookPath(constants.UpdaterBin)
	if err != nil {
		return false
	}

	fmt.Printf(constants.MsgUpdaterFallback, updaterPath)
	cmd := exec.Command(updaterPath, "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		return false
	}

	return true
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
func launchHandoff(copyPath, repoPath string) {
	copyArgs := []string{constants.CmdUpdateRunner}
	if hasFlag(constants.FlagVerbose) {
		copyArgs = append(copyArgs, constants.FlagVerbose)
	}

	copyArgs = append(copyArgs, constants.FlagRepoPath, repoPath)

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
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}

	fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
	os.Exit(1)
}

// runUpdateRunner is a hidden command that performs the real update work.
func runUpdateRunner() {
	repoPath := resolveRepoPath()

	initRunnerVerbose()
	fmt.Printf(constants.MsgUpdateStarting)
	fmt.Printf(constants.MsgUpdateRepoPath, repoPath)
	executeUpdate(repoPath)
}

// getFlagValue returns the value following a flag like --repo-path <value>.
func getFlagValue(name string) string {
	args := os.Args[2:]
	for i, arg := range args {
		if arg == name && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
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

// pathExists checks if a directory exists on disk.
func pathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// promptRepoPath asks the user to enter the source repo path interactively.
func promptRepoPath() string {
	fmt.Fprint(os.Stderr, constants.MsgUpdatePathMissing)
	fmt.Fprint(os.Stderr, constants.MsgUpdatePathPrompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	path := strings.TrimSpace(input)
	if len(path) == 0 {
		return ""
	}

	if !pathExists(path) {
		fmt.Fprintf(os.Stderr, constants.ErrUpdatePathInvalid, path)

		return ""
	}

	return path
}

// saveRepoPathToDB persists the source repo path in the Settings table.
func saveRepoPathToDB(path string) {
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()

	_ = db.SetSetting(constants.SettingSourceRepoPath, path)
}

// loadRepoPathFromDB reads the source repo path from the Settings table.
func loadRepoPathFromDB() string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	return db.GetSetting(constants.SettingSourceRepoPath)
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
