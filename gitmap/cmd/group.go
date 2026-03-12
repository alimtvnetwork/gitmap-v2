package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runGroup handles the "group" subcommand and routes to sub-handlers.
func runGroup(args []string) {
	checkHelp("group", args)
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrGroupUsage)
		os.Exit(1)
	}
	dispatchGroup(args[0], args[1:])
}

// dispatchGroup routes group subcommands to their handlers.
func dispatchGroup(sub string, args []string) {
	if sub == constants.CmdGroupCreate {
		runGroupCreate(args)

		return
	}
	if sub == constants.CmdGroupAdd {
		runGroupAdd(args)

		return
	}
	if sub == constants.CmdGroupRemove {
		runGroupRemove(args)

		return
	}
	if sub == constants.CmdGroupList {
		runGroupList()

		return
	}
	if sub == constants.CmdGroupShow {
		runGroupShow(args)

		return
	}
	if sub == constants.CmdGroupDelete {
		runGroupDelete(args)

		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrUnknownGroupSub, sub)
	fmt.Fprintln(os.Stderr, constants.ErrGroupUsage)
	os.Exit(1)
}
