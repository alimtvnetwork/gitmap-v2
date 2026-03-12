// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
)

// runReleaseBranch handles the 'release-branch' command.
func runReleaseBranch(args []string) {
	checkHelp("release-branch", args)
	branch, assets, draft, dryRun, verbose := parseReleaseBranchFlags(args)
	_ = verbose

	if len(branch) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrReleaseBranchUsage)
		os.Exit(1)
	}

	err := release.ExecuteFromBranch(branch, assets, draft, dryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}
}

// parseReleaseBranchFlags parses flags for the release-branch command.
func parseReleaseBranchFlags(args []string) (branch, assets string, draft, dryRun, verbose bool) {
	fs := flag.NewFlagSet(constants.CmdReleaseBranch, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	fs.Parse(args)

	branch = ""
	if fs.NArg() > 0 {
		branch = fs.Arg(0)
	}

	return branch, *assetsFlag, *draftFlag, *dryRunFlag, *verboseFlag
}
