package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runDBReset handles the "db-reset" subcommand.
func runDBReset() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBResetFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Reset()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBResetFailed, err)
		os.Exit(1)
	}
	fmt.Print(constants.MsgDBResetDone)
}
