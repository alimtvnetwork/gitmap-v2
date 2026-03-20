package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// InstallCDFunction writes the gcd shell wrapper to the user's profile.
func InstallCDFunction(shell string) error {
	snippet := cdSnippet(shell)
	if len(snippet) == 0 {
		return fmt.Errorf(constants.ErrCompUnknownShell, shell)
	}

	profilePath := cdProfilePath(shell)

	return appendCDFunction(snippet, profilePath)
}

// cdSnippet returns the gcd function body for the given shell.
func cdSnippet(shell string) string {
	switch shell {
	case constants.ShellPowerShell:
		return constants.CDFuncPowerShell
	case constants.ShellBash:
		return constants.CDFuncBash
	case constants.ShellZsh:
		return constants.CDFuncZsh
	default:
		return ""
	}
}

// cdProfilePath returns the profile path to write the cd function to.
func cdProfilePath(shell string) string {
	switch shell {
	case constants.ShellPowerShell:
		_, profile := resolvePowerShellPaths()
		return profile
	case constants.ShellBash:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".bashrc")
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".zshrc")
	}
}

// appendCDFunction appends the gcd function to the profile if not present.
func appendCDFunction(snippet, profilePath string) error {
	existing, err := os.ReadFile(profilePath)
	if err == nil && strings.Contains(string(existing), constants.CDFuncMarker) {
		fmt.Fprintf(os.Stderr, constants.MsgCDFuncAlready)

		return nil
	}

	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\n%s\n%s\n", constants.CDFuncMarker, snippet)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCDFuncInstalled)

	return nil
}
