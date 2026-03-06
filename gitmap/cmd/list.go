package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// runList handles the "list" subcommand.
func runList(args []string) {
	groupFilter, verbose := parseListFlags(args)
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	records, err := loadListRecords(db, groupFilter)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}

	printListOutput(records, verbose)
}

// parseListFlags parses flags for the list command.
func parseListFlags(args []string) (group string, verbose bool) {
	fs := flag.NewFlagSet(constants.CmdList, flag.ExitOnError)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	vFlag := fs.Bool("verbose", false, constants.FlagDescListVerbose)
	fs.Parse(args)

	return *gFlag, *vFlag
}

// loadListRecords loads repos, optionally filtered by group.
func loadListRecords(db *store.DB, group string) ([]model.ScanRecord, error) {
	if len(group) > 0 {
		return db.ShowGroup(group)
	}

	return db.ListRepos()
}

// printListOutput renders the list table to stdout.
func printListOutput(records []model.ScanRecord, verbose bool) {
	if len(records) == 0 {
		fmt.Println(constants.MsgListEmpty)
		return
	}
	fmt.Println(constants.MsgListHeader)
	fmt.Println(constants.MsgListSeparator)
	for _, r := range records {
		printListRow(r, verbose)
	}
}

// printListRow prints a single row in list output.
func printListRow(r model.ScanRecord, verbose bool) {
	if verbose {
		fmt.Printf(constants.MsgListVerboseFmt, r.Slug, r.RepoName, r.AbsolutePath)
		return
	}
	fmt.Printf(constants.MsgListRowFmt, r.Slug, r.RepoName)
}

// openDB opens the gitmap database from the default output folder.
func openDB() (*store.DB, error) {
	db, err := store.Open(constants.DefaultOutputFolder)
	if err != nil {
		return nil, err
	}

	return db, db.Migrate()
}
