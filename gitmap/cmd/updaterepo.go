package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// pathExists checks if a directory exists on disk.
func pathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// promptRepoPath asks the user to enter the source repo path interactively.
func promptRepoPath() string {
	fmt.Fprint(os.Stderr, constants.MsgUpdatePathMissing)
	fmt.Fprint(os.Stderr, constants.MsgUpdatePathPrompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	path := strings.TrimSpace(input)
	if len(path) == 0 {
		return ""
	}

	if !pathExists(path) {
		fmt.Fprintf(os.Stderr, constants.ErrUpdatePathInvalid, path)

		return ""
	}

	return path
}

// saveRepoPathToDB persists the source repo path in the Settings table.
func saveRepoPathToDB(path string) {
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()

	_ = db.SetSetting(constants.SettingSourceRepoPath, path)
}

// loadRepoPathFromDB reads the source repo path from the Settings table.
func loadRepoPathFromDB() string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	return db.GetSetting(constants.SettingSourceRepoPath)
}