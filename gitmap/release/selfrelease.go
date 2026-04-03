package release

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// ExecuteSelf resolves the gitmap source repo from the running binary,
// switches to that directory, runs Execute, then returns to the original dir.
func ExecuteSelf(opts Options) error {
	srcRoot, err := resolveSourceRepo()
	if err != nil {
		return err
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not determine current directory: %w", err)
	}

	fmt.Printf(constants.MsgSelfReleaseSwitch, srcRoot)

	err = os.Chdir(srcRoot)
	if err != nil {
		return fmt.Errorf("could not switch to source repo: %w", err)
	}

	releaseErr := Execute(opts)

	// Always attempt to return to original directory.
	if cdErr := os.Chdir(originalDir); cdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not return to %s: %v\n", originalDir, cdErr)
	} else {
		fmt.Printf(constants.MsgSelfReleaseReturn, originalDir)
	}

	return releaseErr
}

// IsInsideGitRepo checks if the current directory is inside a Git repository.
func IsInsideGitRepo() bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}

	return findGitRoot(dir) != ""
}

// resolveSourceRepo finds the git root of the gitmap source from the executable path.
func resolveSourceRepo() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf(constants.ErrSelfReleaseExec, err)
	}

	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf(constants.ErrSelfReleaseExec, err)
	}

	root := findGitRoot(filepath.Dir(exe))
	if root == "" {
		return "", fmt.Errorf("%s", constants.ErrSelfReleaseNoRepo)
	}

	return root, nil
}

// findGitRoot walks up from dir looking for a .git directory.
func findGitRoot(dir string) string {
	for {
		if info, err := os.Stat(filepath.Join(dir, ".git")); err == nil && info.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}

		dir = parent
	}
}
