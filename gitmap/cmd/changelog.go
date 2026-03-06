// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
)

// runChangelog handles the 'changelog' command.
func runChangelog(args []string) {
	version, latest, limit := parseChangelogFlags(args)

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
func parseChangelogFlags(args []string) (version string, latest bool, limit int) {
	fs := flag.NewFlagSet(constants.CmdChangelog, flag.ExitOnError)
	latestFlag := fs.Bool("latest", false, constants.FlagDescLatest)
	limitFlag := fs.Int("limit", 5, constants.FlagDescLimit)
	fs.Parse(args)

	version = ""
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}

	if *limitFlag < 1 {
		*limitFlag = 1
	}

	return version, *latestFlag, *limitFlag
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
