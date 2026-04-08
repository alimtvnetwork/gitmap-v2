package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/user/gitmap/constants"
)

const maxNppFileSize = 10 * 1024 * 1024 // 10 MB limit per extracted file

// resolveNppInstallName maps install-npp to the npp binary for install.
func resolveNppInstallName(tool string) string {
	if tool == constants.ToolNppInstall {
		return constants.ToolNpp
	}

	return tool
}

// runNppSettingsOnly syncs Notepad++ settings without installing the binary.
func runNppSettingsOnly() {
	fmt.Print(constants.MsgInstallNppSkipBin)
	runNppSettings()
}

// runNppSettings syncs Notepad++ settings to the AppData directory.
func runNppSettings() {
	fmt.Print(constants.MsgInstallNppSettings)

	target := nppSettingsTarget()
	if target == "" {
		return
	}

	err := os.MkdirAll(target, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppDirCreate, target, err)

		return
	}

	extractNppSettingsZip(target)
}

// nppSettingsTarget returns the Notepad++ AppData settings path.
func nppSettingsTarget() string {
	if runtime.GOOS != "windows" {
		fmt.Fprintf(os.Stderr, constants.ErrNppWindowsOnly, runtime.GOOS)

		return ""
	}

	appData := os.Getenv("APPDATA")
	if appData == "" {
		fmt.Fprint(os.Stderr, constants.ErrNppNoAppData)

		return ""
	}

	return filepath.Join(appData, "Notepad++")
}

// resolveNppDataPath resolves the npp-settings data path relative to the binary.
func resolveNppDataPath(subpath string) string {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve executable path: %v\n", err)

		return filepath.Join("data", subpath)
	}

	// Resolve symlinks to find the real binary location.
	realExe, err := filepath.EvalSymlinks(exe)
	if err != nil {
		realExe = exe
	}

	binDir := filepath.Dir(realExe)

	// Search order: binary-relative, then CWD-relative.
	candidates := []string{
		filepath.Join(binDir, "data", subpath),
		filepath.Join("data", subpath),
	}

	for _, candidate := range candidates {
		abs, _ := filepath.Abs(candidate)
		if _, statErr := os.Stat(candidate); statErr == nil {
			fmt.Printf("  → Resolved data path: %s\n", abs)

			return candidate
		}

		fmt.Printf("  → Searched: %s (not found)\n", abs)
	}

	// Return the binary-relative path as default (will produce a clear error).
	return candidates[0]
}
