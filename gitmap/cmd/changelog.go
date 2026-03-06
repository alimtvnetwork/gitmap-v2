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
	version, openFile = resolveChangelogAlias(version, openFile)
	if openFile {
		handleChangelogOpen(latest, version)
	}
	if latest == false && len(version) == 0 && openFile {
		return
	}

	dispatchChangelogOutput(version, latest, limit)
}

// resolveChangelogAlias detects if the version arg is actually a file-open alias.
func resolveChangelogAlias(version string, openFile bool) (string, bool) {
	if strings.EqualFold(version, constants.ChangelogFile) || strings.EqualFold(version, constants.CmdChangelogMD) {
		return "", true
	}

	return version, openFile
}

// handleChangelogOpen opens the changelog file and exits on error.
func handleChangelogOpen(latest bool, version string) {
	err := openChangelogFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChangelogOpen, err)
		os.Exit(1)
	}
	if latest == false && len(version) == 0 {
		os.Exit(0)
	}
}

// dispatchChangelogOutput prints the appropriate changelog entries.
func dispatchChangelogOutput(version string, latest bool, limit int) {
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
		printSingleVersion(entries, version)

		return
	}
	printChangelogEntries(entries, limit)
}

// printSingleVersion finds and prints one version's changelog.
func printSingleVersion(entries []release.ChangelogEntry, version string) {
	entry, found := release.FindChangelogEntry(entries, version)
	if found == false {
		fmt.Fprintf(os.Stderr, constants.ErrChangelogVersionNotFound, release.NormalizeVersion(version))
		os.Exit(1)
	}
	printChangelogEntry(entry)
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
	fmt.Printf(constants.ChangelogVersionFmt+"\n", entry.Version)
	for _, note := range entry.Notes {
		fmt.Printf(constants.ChangelogNoteFmt+"\n", note)
	}
}

// openChangelogFile opens CHANGELOG.md with the default OS app.
func openChangelogFile() error {
	absPath, err := filepath.Abs(constants.ChangelogFile)
	if err != nil {
		return err
	}

	return runOpenCommand(absPath)
}

// runOpenCommand executes the platform-specific open command.
func runOpenCommand(path string) error {
	if runtime.GOOS == constants.OSWindows {
		cmd := exec.Command("cmd", "/c", "start", "", path)

		return cmd.Run()
	}
	if runtime.GOOS == constants.OSDarwin {
		cmd := exec.Command("open", path)

		return cmd.Run()
	}
	cmd := exec.Command("xdg-open", path)

	return cmd.Run()
}
