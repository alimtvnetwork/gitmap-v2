package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runGroupList handles "group list".
func runGroupList() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	groups, err := db.ListGroups()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}

	printGroupList(db, groups)
}

// printGroupList renders the group table to stdout.
func printGroupList(db *store.DB, groups []model.Group) {
	if len(groups) == 0 {
		fmt.Println(constants.MsgGroupEmpty)
		return
	}
	fmt.Println(constants.MsgGroupHeader)
	fmt.Println(constants.MsgListSeparator)
	for _, g := range groups {
		count, _ := db.CountGroupRepos(g.Name)
		fmt.Printf(constants.MsgGroupRowFmt, g.Name, count, g.Description)
	}
}
