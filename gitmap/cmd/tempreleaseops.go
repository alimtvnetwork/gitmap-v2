// Package cmd implements CLI command handlers for gitmap.
package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/release"
	"github.com/user/gitmap/store"
)

// runTempReleaseCreate creates temp-release branches from recent commits.
func runTempReleaseCreate(args []string) {
	count, pattern, start, dryRun, _ := parseTempReleaseCreateFlags(args)
	prefix, digitCount := parseVersionPattern(pattern)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	db.Migrate()

	if start == 0 {
		start = resolveAutoStart(db, prefix)
		fmt.Printf(constants.MsgTRSeqAuto, start)
	} else {
		fmt.Printf(constants.MsgTRSeqStart, start)
	}

	validateSequenceRange(start, count, digitCount)
	executeTRCreate(db, count, prefix, digitCount, start, dryRun)
}

// executeTRCreate fetches commits and creates branches or prints dry-run.
func executeTRCreate(db *store.DB, count int, prefix string, digitCount, start int, dryRun bool) {
	commits, err := release.ListRecentCommits(count)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	if len(commits) < count {
		fmt.Fprintf(os.Stderr, constants.ErrTRNotEnough, len(commits), count)
		count = len(commits)
	}

	if dryRun {
		printTRDryRun(commits, prefix, start, digitCount)

		return
	}

	branches := createTempBranches(commits, prefix, start, digitCount, db)
	pushTempBranches(branches)
	fmt.Print(constants.MsgTRComplete)
}

// parseVersionPattern extracts the prefix and digit count from a pattern.
func parseVersionPattern(pattern string) (string, int) {
	idx := strings.Index(pattern, "$")
	if idx < 0 {
		fmt.Fprintln(os.Stderr, constants.ErrTRNoPlaceholder)
		os.Exit(1)
	}

	prefix := pattern[:idx]
	dollarCount := 0

	for i := idx; i < len(pattern) && pattern[i] == '$'; i++ {
		dollarCount++
	}

	return prefix, dollarCount
}

// resolveAutoStart determines the next sequence number from DB.
func resolveAutoStart(db *store.DB, prefix string) int {
	max, err := db.MaxTempReleaseSeq(prefix)
	if err != nil || max == 0 {
		return 1
	}

	return max + 1
}

// validateSequenceRange checks that all sequences fit within the digit format.
func validateSequenceRange(start, count, digits int) {
	maxVal := int(math.Pow(10, float64(digits))) - 1
	endSeq := start + count - 1

	if endSeq > maxVal {
		fmt.Fprintf(os.Stderr, constants.ErrTROverflow+"\n", endSeq, digits, maxVal)
		os.Exit(1)
	}
}

// formatSeq zero-pads a sequence number to the given digit count.
func formatSeq(seq, digits int) string {
	return fmt.Sprintf("%0*d", digits, seq)
}

// createTempBranches creates branches from commits and records them in DB.
func createTempBranches(commits []release.TempReleaseCommit, prefix string, start, digits int, db *store.DB) []string {
	fmt.Printf(constants.MsgTRCreating, len(commits))

	var created []string

	for i, c := range commits {
		seq := start + i
		version := prefix + formatSeq(seq, digits)
		branchName := constants.TempReleaseBranchPrefix + version

		if release.BranchExists(branchName) {
			fmt.Printf(constants.MsgTRSkipExists, branchName)

			continue
		}

		err := release.CreateBranchFromSHA(branchName, c.SHA)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Failed to create %s: %v\n", branchName, err)

			continue
		}

		fmt.Printf(constants.MsgTRCreated, branchName, c.Short)
		_ = db.InsertTempRelease(branchName, prefix, seq, c.SHA, c.Message)
		created = append(created, branchName)
	}

	return created
}

// pushTempBranches pushes all created branches to origin.
func pushTempBranches(branches []string) {
	if len(branches) == 0 {
		return
	}

	fmt.Printf(constants.MsgTRPushing, len(branches))

	err := release.PushBranches(branches)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Push failed: %v\n", err)

		return
	}

	fmt.Printf(constants.MsgTRPushed, len(branches))
}

// printTRDryRun shows what would be created without executing.
func printTRDryRun(commits []release.TempReleaseCommit, prefix string, start, digits int) {
	fmt.Printf(constants.MsgTRDryRunHeader, len(commits))

	for i, c := range commits {
		seq := start + i
		version := prefix + formatSeq(seq, digits)
		branchName := constants.TempReleaseBranchPrefix + version

		fmt.Printf(constants.MsgTRDryRunEntry, branchName, c.Short, c.Message)
	}
}

// runTempReleaseList lists all temp-release branches.
func runTempReleaseList(args []string) {
	jsonOutput := hasTRListFlag(args, "--json")

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	db.Migrate()

	releases, err := db.ListTempReleases()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(releases, "", constants.JSONIndent)
		fmt.Println(string(data))

		return
	}

	printTRList(releases)
}

// printTRList prints temp-release records in terminal format.
func printTRList(releases []model.TempRelease) {
	if len(releases) == 0 {
		fmt.Print(constants.MsgTRListEmpty)

		return
	}

	fmt.Printf(constants.MsgTRListHeader, len(releases))

	for _, r := range releases {
		short := r.Commit
		if len(short) > constants.ShaDisplayLength {
			short = short[:constants.ShaDisplayLength]
		}

		msg := r.CommitMessage
		if len(msg) > 50 {
			msg = msg[:50]
		}

		fmt.Printf(constants.MsgTRListRow, r.Branch, short, msg, r.CreatedAt)
	}
}

// hasListFlag checks if a flag is present in the args.
func hasListFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}

	return false
}
