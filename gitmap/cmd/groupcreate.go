package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runGroupCreate handles "group create <name>".
func runGroupCreate(args []string) {
	name, desc, color := parseGroupCreateFlags(args)
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrGroupNameReq)
		os.Exit(1)
	}
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.CreateGroup(name, desc, color)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgGroupCreated, name)
}

// parseGroupCreateFlags parses flags for group create.
func parseGroupCreateFlags(args []string) (name, desc, color string) {
	fs := flag.NewFlagSet(constants.CmdGroupCreate, flag.ExitOnError)
	descFlag := fs.String("description", "", constants.FlagDescGroupDesc)
	colorFlag := fs.String("color", "", constants.FlagDescGroupColor)
	fs.Parse(args)

	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	return name, *descFlag, *colorFlag
}
