package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/user/gitmap/constants"
)

// isVSCodeAvailable checks if the VS Code CLI is on PATH.
func isVSCodeAvailable() bool {
	_, err := exec.LookPath(constants.VSCodeBin)

	return err == nil
}

// openInVSCode opens the given folder in VS Code.
// Tries multiple strategies to bypass "Another instance running as administrator":
// 1. code --reuse-window (standard)
// 2. code --new-window (bypasses some admin conflicts)
// 3. cmd /C start (launches in separate process context, avoids admin inheritance)
func openInVSCode(absPath string) {
	if !isVSCodeAvailable() {
		fmt.Fprintf(os.Stdout, constants.MsgVSCodeNotFound)

		return
	}

	fmt.Printf(constants.MsgVSCodeOpening, absPath)

	if tryVSCodeReuse(absPath) {
		fmt.Println(constants.MsgVSCodeOpened)

		return
	}

	if tryVSCodeNewWindow(absPath) {
		fmt.Println(constants.MsgVSCodeOpened)

		return
	}

	if tryVSCodeDetached(absPath) {
		fmt.Println(constants.MsgVSCodeOpened)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrVSCodeAdminLock)
}

// tryVSCodeReuse attempts to open with --reuse-window.
func tryVSCodeReuse(absPath string) bool {
	cmd := exec.Command(
		constants.VSCodeBin,
		constants.VSCodeFlagReuseWindow,
		absPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start() == nil
}

// tryVSCodeNewWindow attempts to open with --new-window.
func tryVSCodeNewWindow(absPath string) bool {
	cmd := exec.Command(
		constants.VSCodeBin,
		constants.VSCodeFlagNewWindow,
		absPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start() == nil
}

// tryVSCodeDetached launches VS Code via cmd /C start to avoid
// inheriting the admin elevation context of the current process.
func tryVSCodeDetached(absPath string) bool {
	cmd := exec.Command(
		"cmd", "/C", "start", "",
		constants.VSCodeBin,
		constants.VSCodeFlagNewWindow,
		absPath,
	)

	return cmd.Start() == nil
}
