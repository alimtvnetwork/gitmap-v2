// Package localdirs handles migration of legacy repo-local directories.
package localdirs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// MigrateLegacyDirs moves old directories into .gitmap/ if found.
func MigrateLegacyDirs() {
	migrations := []struct{ old, sub string }{
		{constants.LegacyOutputDir, constants.OutputDirName},
		{constants.LegacyReleaseDir, constants.ReleaseDirName},
		{constants.LegacyDeployedDir, constants.DeployedDirName},
	}

	for _, m := range migrations {
		migrateSingleDir(m.old, m.sub)
	}
}

// migrateSingleDir moves a legacy directory to .gitmap/<sub> if it exists.
// If the target already exists, it merges files from the legacy directory
// (without overwriting) and removes the legacy directory.
func migrateSingleDir(oldDir, subDir string) {
	if !dirExists(oldDir) {
		return
	}

	target := filepath.Join(constants.GitMapDir, subDir)
	if dirExists(target) {
		mergeAndRemoveLegacy(oldDir, target)

		return
	}

	ensureGitMapDir()
	if err := os.Rename(oldDir, target); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrMigrationFailed, oldDir, err)

		return
	}

	fmt.Printf(constants.MsgMigrated, oldDir, target)
}

// mergeAndRemoveLegacy copies files from oldDir into target (skipping files
// that already exist in target), then removes the legacy directory.
func mergeAndRemoveLegacy(oldDir, target string) {
	var merged, skipped int

	_ = filepath.WalkDir(oldDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(oldDir, path)
		if rel == "." {
			return nil
		}

		dest := filepath.Join(target, rel)

		if d.IsDir() {
			_ = os.MkdirAll(dest, constants.DirPermission)

			return nil
		}

		if fileExists(dest) {
			skipped++

			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}

		if writeErr := os.WriteFile(dest, data, constants.FilePermission); writeErr != nil {
			return nil
		}

		merged++

		return nil
	})

	if err := os.RemoveAll(oldDir); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrMigrationFailed, oldDir, err)

		return
	}

	fmt.Printf(constants.MsgMergedAndRemoved, oldDir, target, merged, skipped)
}

// dirExists checks if a directory exists at the given path.
func dirExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}

// ensureGitMapDir creates the .gitmap/ directory if it does not exist.
func ensureGitMapDir() {
	_ = os.MkdirAll(constants.GitMapDir, constants.DirPermission)
}
