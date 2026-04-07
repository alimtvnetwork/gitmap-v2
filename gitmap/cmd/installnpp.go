package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/user/gitmap/constants"
)

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

	err := os.MkdirAll(target, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create settings directory: %v\n", err)

		return
	}

	extractNppSettingsZip(target)
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

// extractNppSettingsZip extracts the bundled settings zip to the target.
func extractNppSettingsZip(target string) {
	zipPath := filepath.Join("data", "npp-settings", "npp-settings.zip")

	fmt.Printf(constants.MsgInstallNppExtract, target)

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Settings zip not found: %s\n", zipPath)
		syncNppSettingsFallback(target)

		return
	}
	defer reader.Close()

	for _, file := range reader.File {
		extractZipEntry(target, file)
	}

	fmt.Printf("Settings synced to %s\n", target)
}

// extractZipEntry writes a single zip entry to the target directory.
func extractZipEntry(target string, file *zip.File) {
	destPath := filepath.Join(target, file.Name)

	if file.FileInfo().IsDir() {
		_ = os.MkdirAll(destPath, 0755)

		return
	}

	_ = os.MkdirAll(filepath.Dir(destPath), 0755)

	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return
	}
	defer dst.Close()

	_, _ = io.Copy(dst, src)
}

// syncNppSettingsFallback copies loose settings files as a fallback.
func syncNppSettingsFallback(target string) {
	source := filepath.Join("data", "npp-settings")

	entries, err := os.ReadDir(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Settings source not found: %s\n", source)

		return
	}

	for _, entry := range entries {
		if entry.Name() == "npp-settings.zip" {
			continue
		}
		copySettingsFile(source, target, entry.Name())
	}

	fmt.Printf("Settings synced to %s (fallback)\n", target)
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
