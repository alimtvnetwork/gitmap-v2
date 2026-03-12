package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/user/gitmap/constants"
)

// Install writes the completion script and adds a source line to the profile.
func Install(shell string) error {
	script, err := Generate(shell)
	if err != nil {
		return err
	}

	scriptPath, profilePath := resolvePaths(shell)

	return writeAndSource(script, scriptPath, profilePath, shell)
}

// resolvePaths returns the script file path and profile path for the shell.
func resolvePaths(shell string) (string, string) {
	switch shell {
	case constants.ShellPowerShell:
		return resolvePowerShellPaths()
	case constants.ShellBash:
		return resolveBashPaths()
	default:
		return resolveZshPaths()
	}
}

// resolvePowerShellPaths returns paths for PowerShell completion.
func resolvePowerShellPaths() (string, string) {
	appData := os.Getenv("APPDATA")
	if len(appData) == 0 {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, ".config")
	}

	scriptPath := filepath.Join(appData, constants.CompDirName, constants.CompFilePS)
	profile := os.Getenv("PROFILE")
	if len(profile) == 0 {
		profile = defaultPSProfile()
	}

	return scriptPath, profile
}

// defaultPSProfile returns the default PowerShell profile path.
func defaultPSProfile() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	}

	return filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1")
}

// resolveBashPaths returns paths for Bash completion.
func resolveBashPaths() (string, string) {
	home, _ := os.UserHomeDir()
	scriptPath := filepath.Join(home, ".local", "share", constants.CompDirName, constants.CompFileBash)
	profile := filepath.Join(home, ".bashrc")

	return scriptPath, profile
}

// resolveZshPaths returns paths for Zsh completion.
func resolveZshPaths() (string, string) {
	home, _ := os.UserHomeDir()
	scriptPath := filepath.Join(home, ".local", "share", constants.CompDirName, constants.CompFileZsh)
	profile := filepath.Join(home, ".zshrc")

	return scriptPath, profile
}

// writeAndSource writes the script file and adds a source line to the profile.
func writeAndSource(script, scriptPath, profilePath, shell string) error {
	if err := writeScriptFile(scriptPath, script); err != nil {
		return err
	}

	return addSourceLine(scriptPath, profilePath, shell)
}

// writeScriptFile creates directories and writes the completion script.
func writeScriptFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

// addSourceLine appends the source command to the profile if absent.
func addSourceLine(scriptPath, profilePath, shell string) error {
	sourceLine := buildSourceLine(scriptPath, shell)

	existing, err := os.ReadFile(profilePath)
	if err == nil && strings.Contains(string(existing), sourceLine) {
		fmt.Fprintf(os.Stderr, constants.MsgCompAlreadyDone, shell)

		return nil
	}

	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\n# gitmap shell completion\n%s\n", sourceLine)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCompProfileWrite, profilePath)

	return nil
}

// buildSourceLine returns the shell-appropriate source command.
func buildSourceLine(scriptPath, shell string) string {
	if shell == constants.ShellPowerShell {
		return fmt.Sprintf(". '%s'", scriptPath)
	}

	return fmt.Sprintf("source '%s'", scriptPath)
}
