package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

// runZipGroupRemove handles "zip-group remove <group> <path>".
func runZipGroupRemove(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	groupName := args[0]
	rawPath := args[1]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	// Resolve to full path for matching.
	_, _, fullPath, _, _ := resolveZipPath(rawPath)

	err = db.RemoveZipGroupItem(groupName, fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgZGItemRemoved, rawPath, groupName)
	syncZipGroupJSON(db)
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

// executeZipGroupShow opens the DB and displays group items with dynamic expansion.
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

// printZipGroupShow renders items in a zip group with dynamic folder expansion.
func printZipGroupShow(group model.ZipGroup, items []model.ZipGroupItem) {
	fmt.Printf(constants.MsgZGShowHeader, group.Name, len(items))

	for _, item := range items {
		if item.IsFolder {
			fmt.Printf(constants.MsgZGShowFolder, item.RelativePath)
			fmt.Printf(constants.MsgZGShowPaths, item.RepoPath, item.RelativePath, item.FullPath)

			// Dynamically expand folder contents.
			files := expandFolder(item.FullPath)
			if len(files) > 0 {
				fmt.Printf(constants.MsgZGShowExpanded, len(files))
				for _, f := range files {
					fmt.Printf(constants.MsgZGShowExpFile, f)
				}
			}
		} else {
			fmt.Printf(constants.MsgZGShowFile, item.RelativePath)
			fmt.Printf(constants.MsgZGShowPaths, item.RepoPath, item.RelativePath, item.FullPath)
		}
	}

	if len(group.ArchiveName) > 0 {
		fmt.Printf(constants.MsgZGShowArchive, group.ArchiveName)
	}

	printHints(zipGroupShowHints())
}

// expandFolder returns relative file paths inside a folder for display.
func expandFolder(folderPath string) []string {
	var files []string

	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		rel, relErr := filepath.Rel(folderPath, path)
		if relErr != nil {
			rel = path
		}

		files = append(files, rel)

		return nil
	})

	return files
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
	syncZipGroupJSON(db)
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
	syncZipGroupJSON(db)
}

// syncZipGroupJSON writes zip group data to .gitmap/zip-groups.json.
func syncZipGroupJSON(db *store.DB) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	err = db.WriteZipGroupsJSON(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrZGJSONWrite+"\n", err)
	}
}
