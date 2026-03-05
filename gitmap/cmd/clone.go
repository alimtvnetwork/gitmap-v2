package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/cloner"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/desktop"
	"github.com/user/gitmap/model"
)

// runClone handles the "clone" subcommand.
func runClone(args []string) {
	source, targetDir, safePull, ghDesktop := parseCloneFlags(args)
	if len(source) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrSourceRequired)
		fmt.Fprintln(os.Stderr, constants.ErrCloneUsage)
		os.Exit(1)
	}
	executeClone(source, targetDir, safePull, ghDesktop)
}

// executeClone runs the clone operation and prints the summary.
func executeClone(source, targetDir string, safePull, ghDesktop bool) {
	summary, err := cloner.CloneFromFile(source, targetDir, safePull)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFailed, err)
		os.Exit(1)
	}
	printSummary(summary)
	registerCloned(summary, targetDir, ghDesktop)
}

// registerCloned adds successfully cloned repos to GitHub Desktop.
func registerCloned(s model.CloneSummary, targetDir string, enabled bool) {
	if enabled {
		records := buildClonedRecords(s, targetDir)
		result := desktop.AddRepos(records)
		fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
	}
}

// buildClonedRecords creates records with absolute paths for cloned repos.
func buildClonedRecords(s model.CloneSummary, targetDir string) []model.ScanRecord {
	absTarget, _ := filepath.Abs(targetDir)
	records := make([]model.ScanRecord, 0, s.Succeeded)
	for _, r := range s.Cloned {
		r.Record.AbsolutePath = filepath.Join(absTarget, r.Record.RelativePath)
		records = append(records, r.Record)
	}

	return records
}

// printSummary displays clone results to the user.
func printSummary(s model.CloneSummary) {
	fmt.Printf(constants.MsgCloneComplete, s.Succeeded, s.Failed)
	if s.Failed > 0 {
		printFailures(s)
	}
}

// printFailures lists each failed clone operation.
func printFailures(s model.CloneSummary) {
	fmt.Println(constants.MsgFailedClones)
	for _, e := range s.Errors {
		fmt.Printf(constants.MsgFailedEntry,
			e.Record.RepoName, e.Record.RelativePath, e.Error)
	}
}
