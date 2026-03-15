package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
)

// runClearReleaseJSON handles the "clear-release-json" subcommand.
func runClearReleaseJSON(args []string) {
	checkHelp("clear-release-json", args)

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrClearReleaseUsage)
		os.Exit(1)
	}

	version := args[0]
	v, err := release.Parse(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReleaseInvalidVersion, version)
		os.Exit(1)
	}

	filename := v.String() + constants.ExtJSON
	path := filepath.Join(constants.DefaultReleaseDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.ErrClearReleaseNotFound, v.String())
		os.Exit(1)
	}

	err = os.Remove(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrClearReleaseFailed, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgClearReleaseDone, v.String())
}
