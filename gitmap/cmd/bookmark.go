package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// runBookmark handles the "bookmark" subcommand routing.
func runBookmark(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)
		os.Exit(1)
	}

	sub := args[0]
	rest := args[1:]

	routeBookmarkSub(sub, rest)
}

// routeBookmarkSub routes to the appropriate bookmark subcommand.
func routeBookmarkSub(sub string, args []string) {
	if sub == constants.CmdBookmarkSave {
		runBookmarkSave(args)

		return
	}
	if sub == constants.CmdBookmarkList {
		runBookmarkList(args)

		return
	}
	if sub == constants.CmdBookmarkRun {
		runBookmarkRun(args)

		return
	}
	if sub == constants.CmdBookmarkDelete {
		runBookmarkDelete(args)

		return
	}

	fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)
	os.Exit(1)
}
