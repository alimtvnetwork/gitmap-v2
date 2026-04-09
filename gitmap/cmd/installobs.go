package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// runOBSSettingsOnly syncs OBS Studio settings from the bundled settings folder.
func runOBSSettingsOnly() {
	fmt.Println("  Syncing OBS Studio settings...")
	syncOBSSettings()
}

// syncOBSSettings copies OBS Studio settings from the bundled folder.
func syncOBSSettings() {
	target := obsSettingsTarget()
	if target == "" {
		return
	}

	fmt.Printf("  -> OBS settings target: %s\n", target)

	sourcePath := resolveSettingsPath(
		filepath.Join("settings", "03 - obs"),
		filepath.Join("data", "obs-settings"),
	)

	info, err := os.Stat(sourcePath)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "  Error: OBS settings source not found at %s: %v\n", sourcePath, err)

		return
	}

	copied, copyErr := copyDirRecursive(sourcePath, target)
	if copyErr != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to copy OBS settings: %v\n", copyErr)

		return
	}

	fmt.Printf("  -> Synced %d file(s) to %s\n", copied, target)
}

// obsSettingsTarget returns the OBS Studio config directory.
func obsSettingsTarget() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			fmt.Fprint(os.Stderr, "  Error: APPDATA environment variable not set\n")

			return ""
		}

		return filepath.Join(appData, "obs-studio")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: could not resolve home directory: %v\n", err)

			return ""
		}

		return filepath.Join(home, "Library", "Application Support", "obs-studio")
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: could not resolve home directory: %v\n", err)

			return ""
		}

		return filepath.Join(home, ".config", "obs-studio")
	}
}

// copyDirRecursive copies all files from src to dst recursively.
func copyDirRecursive(src, dst string) (int, error) {
	copied := 0

	entries, err := os.ReadDir(src)
	if err != nil {
		return 0, err
	}

	err = os.MkdirAll(dst, 0o755)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// Skip readme files.
		name := strings.ToLower(entry.Name())
		if name == "readme.txt" || name == "readme.md" {
			continue
		}

		if entry.IsDir() {
			n, dirErr := copyDirRecursive(srcPath, dstPath)
			if dirErr != nil {
				fmt.Fprintf(os.Stderr, "  ! Failed to copy directory %s: %v\n", entry.Name(), dirErr)

				continue
			}

			copied += n

			continue
		}

		copyErr := copyFile(srcPath, dstPath)
		if copyErr != nil {
			fmt.Fprintf(os.Stderr, "  ! Failed to copy %s: %v\n", entry.Name(), copyErr)

			continue
		}

		copied++
	}

	return copied, nil
}

// copyFile copies a single file from src to dst.
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
