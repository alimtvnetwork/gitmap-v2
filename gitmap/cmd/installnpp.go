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
	// Try relative to the binary first.
	exe, err := os.Executable()
	if err == nil {
		binDir := filepath.Dir(exe)
		candidate := filepath.Join(binDir, "data", subpath)

		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate
		}
	}

	// Fallback to CWD-relative path.
	return filepath.Join("data", subpath)
}

// extractNppSettingsZip extracts the bundled settings zip to the target.
func extractNppSettingsZip(target string) {
	zipPath := resolveNppDataPath(filepath.Join("npp-settings", "npp-settings.zip"))

	fmt.Printf(constants.MsgInstallNppExtract, target)
	fmt.Printf("  → Settings zip: %s\n", zipPath)

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppZipNotFound, zipPath, err)
		syncNppSettingsFallback(target)

		return
	}
	defer reader.Close()

	extracted := 0

	for _, file := range reader.File {
		extractZipEntry(target, file)
		extracted++
	}

	fmt.Printf("  ✓ Extracted %d files\n", extracted)
	fmt.Printf(constants.MsgNppSettingsSynced, target)
}

// extractZipEntry writes a single zip entry to the target directory.
func extractZipEntry(target string, file *zip.File) {
	cleanName := filepath.FromSlash(file.Name)
	destPath := filepath.Join(target, cleanName)

	absTarget, _ := filepath.Abs(target)
	absDest, _ := filepath.Abs(destPath)

	if !strings.HasPrefix(absDest, absTarget+string(os.PathSeparator)) {
		fmt.Fprintf(os.Stderr, constants.ErrNppExtractEntry, file.Name, destPath, fmt.Errorf("path traversal detected"))

		return
	}

	if file.FileInfo().IsDir() {
		err := os.MkdirAll(destPath, 0o755)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrNppDirCreate, destPath, err)
		}

		return
	}

	err := os.MkdirAll(filepath.Dir(destPath), 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppDirCreate, filepath.Dir(destPath), err)

		return
	}

	src, err := file.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppExtractEntry, file.Name, destPath, err)

		return
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileCreate, destPath, err)

		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, io.LimitReader(src, maxNppFileSize))
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileCopy, file.Name, destPath, err)
	}
}

// syncNppSettingsFallback copies loose settings files as a fallback.
func syncNppSettingsFallback(target string) {
	source := resolveNppDataPath("npp-settings")

	entries, err := os.ReadDir(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppSourceDir, source, err)

		return
	}

	for _, entry := range entries {
		if entry.Name() == "npp-settings.zip" {
			continue
		}
		copySettingsFile(source, target, entry.Name())
	}

	fmt.Printf(constants.MsgNppSettingsFallback, target)
}

// copySettingsFile copies a single settings file to the target.
func copySettingsFile(source, target, name string) {
	src := filepath.Join(source, name)
	dst := filepath.Join(target, name)

	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileRead, src, err)

		return
	}

	err = os.WriteFile(dst, data, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppFileWrite, dst, err)
	}
}
