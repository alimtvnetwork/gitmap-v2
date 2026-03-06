// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
)

// runChangelog handles the 'changelog' command.
func runChangelog(args []string) {
	version, latest, limit, openFile := parseChangelogFlags(args)
	if strings.EqualFold(version, constants.ChangelogFile) || strings.EqualFold(version, constants.CmdChangelogMD) {
		openFile = true
		version = ""
	}

	if openFile {
		err := openChangelogFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrChangelogOpen, err)
			os.Exit(1)
		}
		if latest == false && len(version) == 0 {
			return
		}
	}

	entries, err := release.ReadChangelog()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChangelogRead, err)
		os.Exit(1)
	}

	if latest {
		printChangelogEntries(entries, 1)
		return
	}

	if len(version) > 0 {
		entry, found := release.FindChangelogEntry(entries, version)
		if found == false {
			fmt.Fprintf(os.Stderr, constants.ErrChangelogVersionNotFound, release.NormalizeVersion(version))
			os.Exit(1)
		}
		printChangelogEntry(entry)
		return
	}

	printChangelogEntries(entries, limit)
}

// parseChangelogFlags parses flags for the changelog command.
func parseChangelogFlags(args []string) (version string, latest bool, limit int, openFile bool) {
	fs := flag.NewFlagSet(constants.CmdChangelog, flag.ExitOnError)
	latestFlag := fs.Bool("latest", false, constants.FlagDescLatest)
	limitFlag := fs.Int("limit", 5, constants.FlagDescLimit)
	openFlag := fs.Bool("open", false, constants.FlagDescOpenChangelog)
	fs.Parse(args)

	version = ""
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}

	if *limitFlag < 1 {
		*limitFlag = 1
	}

	return version, *latestFlag, *limitFlag, *openFlag
}

// printChangelogEntries prints the newest N changelog entries.
func printChangelogEntries(entries []release.ChangelogEntry, limit int) {
	if limit > len(entries) {
		limit = len(entries)
	}
	for i := 0; i < limit; i++ {
		printChangelogEntry(entries[i])
	}
}

// printChangelogEntry prints a single changelog entry.
func printChangelogEntry(entry release.ChangelogEntry) {
	fmt.Printf("\n%s\n", entry.Version)
	for _, note := range entry.Notes {
		fmt.Printf("  - %s\n", note)
	}
}

// openChangelogFile opens CHANGELOG.md with the default OS app.
func openChangelogFile() error {
	absPath, err := filepath.Abs(constants.ChangelogFile)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "start", "", absPath)
		return cmd.Run()
	}
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", absPath)
		return cmd.Run()
	}

	cmd := exec.Command("xdg-open", absPath)
	return cmd.Run()
}
