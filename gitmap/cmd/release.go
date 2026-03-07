// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
	"github.com/user/gitmap/store"
)

// runRelease handles the 'release' command.
func runRelease(args []string) {
	version, assets, commit, branch, bump, draft, dryRun, verbose := parseReleaseFlags(args)
	_ = verbose
	validateReleaseFlags(version, bump, commit, branch)
	executeRelease(version, assets, commit, branch, bump, draft, dryRun, verbose)
}

// executeRelease builds options and runs the release workflow.
func executeRelease(version, assets, commit, branch, bump string, draft, dryRun, verbose bool) {
	opts := release.Options{
		Version: version, Assets: assets,
		Commit: commit, Branch: branch,
		Bump: bump, Draft: draft,
		DryRun: dryRun, Verbose: verbose,
	}
	err := release.Execute(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	persistReleaseToDB()
}

// validateReleaseFlags checks for mutually exclusive flags.
func validateReleaseFlags(version, bump, commit, branch string) {
	if len(bump) > 0 && len(version) > 0 {
		fmt.Fprint(os.Stderr, constants.ErrReleaseBumpConflict)
		os.Exit(1)
	}
	if len(commit) > 0 && len(branch) > 0 {
		fmt.Fprint(os.Stderr, constants.ErrReleaseCommitBranch)
		os.Exit(1)
	}
}

// parseReleaseFlags parses flags for the release command.
func parseReleaseFlags(args []string) (version, assets, commit, branch, bump string, draft, dryRun, verbose bool) {
	fs := flag.NewFlagSet(constants.CmdRelease, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	commitFlag := fs.String("commit", "", constants.FlagDescCommit)
	branchFlag := fs.String("branch", "", constants.FlagDescRelBranch)
	bumpFlag := fs.String("bump", "", constants.FlagDescBump)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	fs.Parse(args)

	version = ""
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}

	return version, *assetsFlag, *commitFlag, *branchFlag, *bumpFlag, *draftFlag, *dryRunFlag, *verboseFlag
}

// persistReleaseToDB saves the release metadata to SQLite if available.
func persistReleaseToDB() {
	meta := release.LastMeta
	if meta == nil {
		return
	}

	db, err := store.Open(constants.DefaultOutputFolder)
	if err != nil {
		return
	}
	defer db.Close()

	_ = db.Migrate()
	_ = db.UpsertRelease(*meta)
}
