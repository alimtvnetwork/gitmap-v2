package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

// runZipGroupCreate handles "zip-group create <name> [--archive <name>]".
func runZipGroupCreate(args []string) {
	name, archiveName := parseZipGroupCreateFlags(args)
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}
	executeZipGroupCreate(name, archiveName)
}

// parseZipGroupCreateFlags parses flags for zip-group create.
func parseZipGroupCreateFlags(args []string) (name, archive string) {
	fs := flag.NewFlagSet(constants.SubCmdZGCreate, flag.ExitOnError)
	archiveFlag := fs.String("archive", "", constants.FlagDescZGArchive)
	fs.Parse(args)

	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	return name, *archiveFlag
}

// executeZipGroupCreate opens the DB and creates the zip group.
func executeZipGroupCreate(name, archiveName string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.CreateZipGroup(name, archiveName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgZGCreated, name)
	printHints(zipGroupCreateHints())
}

// runZipGroupAdd handles "zip-group add <group> <path...>".
func runZipGroupAdd(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	groupName := args[0]
	paths := args[1:]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	for _, p := range paths {
		addOneZipGroupItem(db, groupName, p)
	}
}

// addOneZipGroupItem adds a single path to a zip group.
func addOneZipGroupItem(db *store.DB, groupName, path string) {
	isFolder := isDirectory(path)

	err := db.AddZipGroupItem(groupName, path, isFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return
	}

	fmt.Printf(constants.MsgZGItemAdded, path, groupName)
}

// isDirectory checks if a path is a directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// runZipGroupRemove handles "zip-group remove <group> <path>".
func runZipGroupRemove(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	groupName := args[0]
	path := args[1]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.RemoveZipGroupItem(groupName, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgZGItemRemoved, path, groupName)
}

// runZipGroupList handles "zip-group list".
func runZipGroupList() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	groups, err := db.ListZipGroupsWithCount()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	printZipGroupList(groups)
	printHints(zipGroupListHints())
}

// printZipGroupList renders the zip group table to stdout.
func printZipGroupList(groups []store.ZipGroupWithCount) {
	if len(groups) == 0 {
		fmt.Println("  No zip groups defined.")

		return
	}

	fmt.Printf(constants.MsgZGListHeader, len(groups))

	for _, g := range groups {
		archive := g.ArchiveName
		if len(archive) == 0 {
			archive = g.Name + ".zip"
		}

		fmt.Printf(constants.MsgZGListRow, g.Name, g.ItemCount, archive)
	}
}

// runZipGroupShow handles "zip-group show <name>".
func runZipGroupShow(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	name := args[0]
	executeZipGroupShow(name)
}

// executeZipGroupShow opens the DB and displays group items.
func executeZipGroupShow(name string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	items, err := db.ListZipGroupItems(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	group, _ := db.FindZipGroupByName(name)
	printZipGroupShow(group, items)
}

// printZipGroupShow renders items in a zip group.
func printZipGroupShow(group model.ZipGroup, items []model.ZipGroupItem) {
	fmt.Printf(constants.MsgZGShowHeader, group.Name, len(items))

	for _, item := range items {
		if item.IsFolder {
			fmt.Printf(constants.MsgZGShowFolder, item.Path)
		} else {
			fmt.Printf(constants.MsgZGShowFile, item.Path)
		}
	}

	if len(group.ArchiveName) > 0 {
		fmt.Printf(constants.MsgZGShowArchive, group.ArchiveName)
	}

	printHints(zipGroupShowHints())
}

// runZipGroupDelete handles "zip-group delete <name>".
func runZipGroupDelete(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	name := args[0]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.DeleteZipGroup(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgZGDeleted, name)
}

// runZipGroupRename handles "zip-group rename <group> --archive <name>".
func runZipGroupRename(args []string) {
	name, archiveName := parseZipGroupRenameFlags(args)
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}
	if len(archiveName) == 0 {
		fmt.Fprintln(os.Stderr, constants.FlagDescZGArchive)
		os.Exit(1)
	}
	executeZipGroupRename(name, archiveName)
}

// parseZipGroupRenameFlags parses flags for zip-group rename.
func parseZipGroupRenameFlags(args []string) (name, archive string) {
	fs := flag.NewFlagSet(constants.SubCmdZGRename, flag.ExitOnError)
	archiveFlag := fs.String("archive", "", constants.FlagDescZGArchive)
	fs.Parse(args)

	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	return name, *archiveFlag
}

// executeZipGroupRename sets a custom archive name for a group.
func executeZipGroupRename(name, archiveName string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.UpdateZipGroupArchive(name, archiveName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgZGArchiveSet, archiveName, name)
}

