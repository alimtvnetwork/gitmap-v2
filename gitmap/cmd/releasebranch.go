// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runReleaseBranch handles the 'release-branch' command.
func runReleaseBranch(args []string) {
	branch, assets, draft, verbose := parseReleaseBranchFlags(args)
	_, _ = verbose, assets

	if len(branch) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrReleaseBranchUsage)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	// TODO: implement release-branch workflow
	fmt.Printf("release-branch command not yet implemented (branch=%s, draft=%v)\n", branch, draft)
}

// parseReleaseBranchFlags parses flags for the release-branch command.
func parseReleaseBranchFlags(args []string) (branch, assets string, draft, verbose bool) {
	fs := flag.NewFlagSet(constants.CmdReleaseBranch, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	fs.Parse(args)

	branch = ""
	if fs.NArg() > 0 {
		branch = fs.Arg(0)
	}

	return branch, *assetsFlag, *draftFlag, *verboseFlag
}
