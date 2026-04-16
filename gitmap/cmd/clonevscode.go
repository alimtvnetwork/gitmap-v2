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
// It uses --reuse-window to avoid admin-mode conflicts
// and --new-window as a fallback.
func openInVSCode(absPath string) {
	if !isVSCodeAvailable() {
		fmt.Fprintf(os.Stdout, constants.MsgVSCodeNotFound)

		return
	}

	fmt.Printf(constants.MsgVSCodeOpening, absPath)

	cmd := exec.Command(
		constants.VSCodeBin,
		constants.VSCodeFlagReuseWindow,
		absPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		// Fallback: try --new-window to bypass admin conflicts.
		fallback := exec.Command(
			constants.VSCodeBin,
			constants.VSCodeFlagNewWindow,
			absPath,
		)
		fallback.Stdout = os.Stdout
		fallback.Stderr = os.Stderr
		fbErr := fallback.Start()
		if fbErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrVSCodeOpenFailed, fbErr)

			return
		}
	}

	fmt.Println(constants.MsgVSCodeOpened)
}
