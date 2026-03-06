package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runGroupAdd handles "group add <group> <slug...>".
func runGroupAdd(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrGroupSlugReq)
		os.Exit(1)
	}
	groupName := args[0]
	slugs := args[1:]
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	for _, slug := range slugs {
		addSlugToGroup(db, groupName, slug)
	}
}

// addSlugToGroup resolves a slug and adds matching repos to the group.
func addSlugToGroup(db interface {
	FindBySlug(string) ([]model.ScanRecord, error)
	AddRepoToGroup(string, string) error
}, groupName, slug string) {
	repos, err := db.(*store.DB).FindBySlug(slug)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBNoMatch+"\n", slug)
		return
	}
	if len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrDBNoMatch+"\n", slug)
		return
	}
	for _, r := range repos {
		err := db.(*store.DB).AddRepoToGroup(groupName, r.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		fmt.Printf(constants.MsgGroupAdded, r.Slug, groupName)
	}
}
