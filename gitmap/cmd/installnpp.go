package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/user/gitmap/constants"
)

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

	err := os.MkdirAll(target, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create settings directory: %v\n", err)

		return
	}

	syncNppSettingsFiles(target)
}

// nppSettingsTarget returns the Notepad++ AppData settings path.
func nppSettingsTarget() string {
	if runtime.GOOS != "windows" {
		fmt.Fprintln(os.Stderr, "Notepad++ settings sync is only supported on Windows.")

		return ""
	}

	appData := os.Getenv("APPDATA")
	if appData == "" {
		fmt.Fprintln(os.Stderr, "APPDATA environment variable not set.")

		return ""
	}

	return filepath.Join(appData, "Notepad++")
}

// syncNppSettingsFiles copies settings from the data directory.
func syncNppSettingsFiles(target string) {
	source := filepath.Join("data", "npp-settings")

	entries, err := os.ReadDir(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Settings source not found: %s\n", source)

		return
	}

	for _, entry := range entries {
		copySettingsFile(source, target, entry.Name())
	}

	fmt.Printf("Settings synced to %s\n", target)
}

// copySettingsFile copies a single settings file to the target.
func copySettingsFile(source, target, name string) {
	src := filepath.Join(source, name)
	dst := filepath.Join(target, name)

	data, err := os.ReadFile(src)
	if err != nil {
		return
	}

	_ = os.WriteFile(dst, data, 0644)
}
