// Package cmd — migrate.go handles automatic migration of legacy directories.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// migrateLegacyDirs moves old directories into .gitmap/ if found.
func migrateLegacyDirs() {
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
func migrateSingleDir(oldDir, subDir string) {
	if !dirExists(oldDir) {
		return
	}

	target := filepath.Join(constants.GitMapDir, subDir)
	if dirExists(target) {
		fmt.Printf(constants.WarnMigrationSkipped, oldDir, target)

		return
	}

	ensureGitMapDir()
	if err := os.Rename(oldDir, target); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrMigrationFailed, oldDir, err)

		return
	}

	fmt.Printf(constants.MsgMigrated, oldDir, target)
}

// dirExists checks if a directory exists at the given path.
func dirExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

// ensureGitMapDir creates the .gitmap/ directory if it does not exist.
func ensureGitMapDir() {
	_ = os.MkdirAll(constants.GitMapDir, constants.DirPermission)
}
