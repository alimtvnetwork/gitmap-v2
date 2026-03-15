package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
	"github.com/user/gitmap/tui"
)

// runInteractive launches the full-screen TUI.
func runInteractive() {
	checkHelp("interactive", os.Args[2:])

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTUIDBOpen, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := tui.Run(db); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
