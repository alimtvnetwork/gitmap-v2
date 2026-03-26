package release

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// checkDuplicate verifies the version hasn't been released.
// If a release JSON exists but no tag or branch, prompts to remove it.
func checkDuplicate(v Version) error {
	if ReleaseExists(v) {
		tagExists := TagExistsLocally(v.String()) || TagExistsRemote(v.String())
		branchName := constants.ReleaseBranchPrefix + v.String()
		branchExists := BranchExists(branchName)

		if !tagExists && !branchExists {
			return handleOrphanedMeta(v)
		}

		return fmt.Errorf(constants.ErrReleaseAlreadyExists, v.String(), v.String())
	}
	if TagExistsLocally(v.String()) || TagExistsRemote(v.String()) {
		return fmt.Errorf(constants.ErrReleaseTagExists, v.String())
	}

	return nil
}

// handleOrphanedMeta detects a release JSON with no matching tag/branch
// and prompts the user to remove it before proceeding.
func handleOrphanedMeta(v Version) error {
	fmt.Printf(constants.MsgReleaseOrphanedMeta, v.String())
	fmt.Print(constants.MsgReleaseOrphanedPrompt)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return fmt.Errorf(constants.ErrReleaseAlreadyExists, v.String(), v.String())
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer != "y" && answer != "yes" {
		return fmt.Errorf(constants.ErrReleaseAborted)
	}

	filename := v.String() + constants.ExtJSON
	path := filepath.Join(constants.DefaultReleaseDir, filename)

	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf(constants.ErrReleaseOrphanedRemove, err)
	}

	fmt.Printf(constants.MsgReleaseOrphanedRemoved, v.String())

	return nil
}
