package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/user/gitmap/constants"
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
	jsonOutput := hasListFlag(args, "--json")

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

// runTempReleaseRemove handles "tr remove <version>|<v1> to <v2>|all".
func runTempReleaseRemove(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrTRRemoveUsage)
		os.Exit(1)
	}

	if args[0] == "all" {
		removeTempReleaseAll()

		return
	}

	if len(args) >= 3 && args[1] == "to" {
		removeTempReleaseRange(args[0], args[2])

		return
	}

	removeTempReleaseSingle(args[0])
}

// removeTempReleaseSingle removes one temp-release branch.
func removeTempReleaseSingle(version string) {
	branchName := resolveTRBranch(version)

	fmt.Printf(constants.MsgTRRemovePrompt, branchName)
	if !confirmAction() {
		return
	}

	removeBranches([]string{branchName})
	fmt.Printf(constants.MsgTRRemovedOne, branchName)
}

// removeTempReleaseRange removes branches from v1 to v2.
func removeTempReleaseRange(from, to string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	db.Migrate()

	releases, _ := db.ListTempReleases()
	fromBranch := resolveTRBranch(from)
	toBranch := resolveTRBranch(to)

	var targets []string

	inRange := false
	for _, r := range releases {
		if r.Branch == fromBranch {
			inRange = true
		}
		if inRange {
			targets = append(targets, r.Branch)
		}
		if r.Branch == toBranch {
			break
		}
	}

	if len(targets) == 0 {
		fmt.Print(constants.MsgTRNoneToRemove)

		return
	}

	fmt.Printf(constants.MsgTRRemoveRange, len(targets))
	for _, b := range targets {
		fmt.Printf(constants.MsgTRRemoveBranch, b)
	}

	fmt.Print(constants.MsgTRRemoveConfirm)
	if !confirmAction() {
		return
	}

	removeBranches(targets)
	cleanupTRFromDB(db, targets)
	fmt.Printf(constants.MsgTRRemoved, len(targets))
}

// removeTempReleaseAll removes all temp-release branches.
func removeTempReleaseAll() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	db.Migrate()

	releases, _ := db.ListTempReleases()
	if len(releases) == 0 {
		fmt.Print(constants.MsgTRNoneToRemove)

		return
	}

	var branches []string
	for _, r := range releases {
		branches = append(branches, r.Branch)
	}

	fmt.Printf(constants.MsgTRRemoveAll, len(branches))
	for _, b := range branches {
		fmt.Printf(constants.MsgTRRemoveBranch, b)
	}

	fmt.Print(constants.MsgTRRemoveConfirm)
	if !confirmAction() {
		return
	}

	removeBranches(branches)
	_ = db.DeleteAllTempReleases()
	fmt.Printf(constants.MsgTRRemoved, len(branches))
}

// resolveTRBranch adds the temp-release/ prefix if not present.
func resolveTRBranch(version string) string {
	if strings.HasPrefix(version, constants.TempReleaseBranchPrefix) {
		return version
	}

	return constants.TempReleaseBranchPrefix + version
}

// removeBranches deletes branches locally and from remote.
func removeBranches(branches []string) {
	for _, b := range branches {
		_ = release.DeleteLocalBranch(b)
	}

	if len(branches) > 0 {
		_ = release.DeleteRemoteBranches(branches)
	}
}

// cleanupTRFromDB removes temp-release records from the database.
func cleanupTRFromDB(db *store.DB, branches []string) {
	for _, b := range branches {
		_ = db.DeleteTempRelease(b)
	}
}

// confirmAction reads a y/N prompt from stdin.
func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}
