package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/gitutil"
)

// runGoMod is the entry point for the gomod command.
func runGoMod(args []string) {
	newPath, dryRun, noMerge, noTidy, verbose := parseGoModFlags(args)

	if len(newPath) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrGoModUsage)
		os.Exit(1)
	}

	oldPath := readModulePath()
	validateGoModPreconditions(oldPath, newPath)

	if dryRun {
		runGoModDryRun(oldPath, newPath)

		return
	}

	originalBranch := goModCurrentBranch()
	slug := deriveSlug(newPath)
	backupBranch, featureBranch := createGoModBranches(slug)

	fileCount := replaceModulePath(oldPath, newPath, verbose)
	runGoModTidy(noTidy)
	commitGoModChanges(oldPath, newPath, fileCount)

	if noMerge {
		printGoModSummaryNoMerge(oldPath, newPath, fileCount, backupBranch, featureBranch)

		return
	}

	mergeGoModBranch(originalBranch, featureBranch, newPath)
	printGoModSummary(oldPath, newPath, fileCount, backupBranch, featureBranch, originalBranch)
}

// parseGoModFlags parses flags for the gomod command.
func parseGoModFlags(args []string) (string, bool, bool, bool, bool) {
	fs := flag.NewFlagSet(constants.CmdGoMod, flag.ExitOnError)
	dryRun := fs.Bool(constants.FlagGoModDryRun, false, constants.FlagDescGoModDryRun)
	noMerge := fs.Bool(constants.FlagGoModNoMerge, false, constants.FlagDescGoModNoMerge)
	noTidy := fs.Bool(constants.FlagGoModNoTidy, false, constants.FlagDescGoModNoTidy)
	verbose := fs.Bool("verbose", false, constants.FlagDescVerbose)
	fs.Parse(args)

	newPath := ""
	if fs.NArg() > 0 {
		newPath = fs.Arg(0)
	}

	return newPath, *dryRun, *noMerge, *noTidy, *verbose
}

// validateGoModPreconditions checks all prerequisites before starting.
func validateGoModPreconditions(oldPath, newPath string) {
	if oldPath == newPath {
		fmt.Printf(constants.MsgGoModNothingRename, oldPath)
		os.Exit(0)
	}

	requireInsideWorkTree()

	if isWorkTreeDirty() {
		fmt.Fprint(os.Stderr, constants.ErrGoModDirtyTree)
		os.Exit(1)
	}
}

// runGoModDryRun previews changes without modifying files.
func runGoModDryRun(oldPath, newPath string) {
	files := findGoFilesWithPath(oldPath)

	fmt.Print(constants.MsgGoModDryHeader)
	fmt.Printf(constants.MsgGoModDryOld, oldPath)
	fmt.Printf(constants.MsgGoModDryNew, newPath)
	fmt.Printf(constants.MsgGoModDryFiles, len(files))

	for _, f := range files {
		fmt.Printf(constants.MsgGoModDryFile, f)
	}
}

// printGoModSummary prints the final summary after merge.
func printGoModSummary(oldPath, newPath string, fileCount int, backup, feature, merged string) {
	fmt.Print(constants.MsgGoModSummary)
	fmt.Printf(constants.MsgGoModOld, oldPath)
	fmt.Printf(constants.MsgGoModNew, newPath)
	fmt.Printf(constants.MsgGoModFiles, fileCount)
	fmt.Printf(constants.MsgGoModBackupBranch, backup)
	fmt.Printf(constants.MsgGoModFeatureBranch, feature)
	fmt.Printf(constants.MsgGoModMergedInto, merged)
}

// printGoModSummaryNoMerge prints the summary when --no-merge is used.
func printGoModSummaryNoMerge(oldPath, newPath string, fileCount int, backup, feature string) {
	fmt.Print(constants.MsgGoModSummary)
	fmt.Printf(constants.MsgGoModOld, oldPath)
	fmt.Printf(constants.MsgGoModNew, newPath)
	fmt.Printf(constants.MsgGoModFiles, fileCount)
	fmt.Printf(constants.MsgGoModBackupBranch, backup)
	fmt.Printf(constants.MsgGoModLeftOn, feature)
}

// runGoModTidy runs go mod tidy unless --no-tidy is set.
func runGoModTidy(noTidy bool) {
	if noTidy {
		return
	}

	err := goModTidy()
	if err != nil {
		fmt.Printf(constants.MsgGoModTidyWarn, err)
	}
}
