package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// runExec handles the "exec" subcommand.
func runExec(args []string) {
	checkHelp("exec", args)
	groupName, all, gitArgs := parseExecFlags(args)
	if len(gitArgs) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrExecUsage)
		os.Exit(1)
	}

	records := loadExecByScope(groupName, all)
	printExecBanner(gitArgs, len(records))
	succeeded, failed, missing := execAllRepos(records, gitArgs)
	printExecSummary(succeeded, failed, missing, len(records))
}

// execAllRepos runs a git command across all repos and returns totals.
func execAllRepos(records []model.ScanRecord, gitArgs []string) (int, int, int) {
	var succeeded, failed, missing int
	for _, rec := range records {
		s, f, m := execOneRepo(rec, gitArgs)
		succeeded += s
		failed += f
		missing += m
	}

	return succeeded, failed, missing
}

// execOneRepo runs a git command in one repo, returning increment counts.
func execOneRepo(rec model.ScanRecord, gitArgs []string) (int, int, int) {
	_, err := os.Stat(rec.AbsolutePath)
	if err == nil && execInRepo(rec, gitArgs) {
		return 1, 0, 0
	}
	if err == nil {
		return 0, 1, 0
	}

	fmt.Printf(constants.ExecMissingFmt,
		constants.ColorDim, truncate(rec.RepoName, 22),
		constants.ColorYellow, constants.ColorReset)

	return 0, 0, 1
}

// parseExecFlags parses --group and --all flags, returning remaining args as git args.
func parseExecFlags(args []string) (groupName string, all bool, gitArgs []string) {
	fs := flag.NewFlagSet(constants.CmdExec, flag.ExitOnError)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	fs.Parse(args)

	return *gFlag, *aFlag, fs.Args()
}

// loadExecByScope returns records filtered by group, all DB repos, or JSON fallback.
func loadExecByScope(groupName string, all bool) []model.ScanRecord {
	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}
	if all {
		return loadAllRecordsDB()
	}

	return loadExecRecordsJSON()
}

// loadExecRecordsJSON reads ScanRecords from gitmap.json.
func loadExecRecordsJSON() []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadExecRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrExecLoadFailed, err)
		os.Exit(1)
	}

	return records
}

// loadExecRecords reads ScanRecords from a JSON file.
func loadExecRecords(path string) ([]model.ScanRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var records []model.ScanRecord
	err = json.Unmarshal(data, &records)

	return records, err
}

// execInRepo runs a git command inside a single repo directory.
func execInRepo(rec model.ScanRecord, gitArgs []string) bool {
	cmd := exec.Command(constants.GitBin, gitArgs...)
	cmd.Dir = rec.AbsolutePath
	cmd.Stdout = nil
	cmd.Stderr = nil

	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	printExecResult(rec.RepoName, output, err)

	return err == nil
}

// printExecResult prints the success or failure line for one repo.
func printExecResult(name, output string, err error) {
	if err == nil {
		fmt.Printf(constants.ExecSuccessFmt, constants.ColorGreen, truncate(name, 22), constants.ColorReset)
	} else {
		fmt.Printf(constants.ExecFailFmt, constants.ColorYellow, truncate(name, 22), constants.ColorReset)
	}

	printExecOutput(output)
}

// printExecOutput prints indented command output lines.
func printExecOutput(output string) {
	if len(output) == 0 {
		return
	}
	for _, line := range strings.Split(output, "\n") {
		fmt.Printf(constants.ExecOutputLineFmt, constants.ColorDim, line, constants.ColorReset)
	}
}

// printExecBanner shows the command header.
func printExecBanner(gitArgs []string, count int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.ExecBannerTop, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.ExecBannerTitle, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.ExecBannerBottom, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %s"+constants.ExecCommandFmt+"%s\n", constants.ColorWhite, strings.Join(gitArgs, " "), constants.ColorReset)
	fmt.Printf("  %s"+constants.ExecRepoCountFmt+"%s\n", constants.ColorDim, count, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermSeparator, constants.ColorReset)
	fmt.Println()
}

// printExecSummary shows final totals.
func printExecSummary(succeeded, failed, missing, total int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.ExecSummaryRule, constants.ColorReset)
	parts := buildExecSummaryParts(succeeded, failed, missing, total)
	line := strings.Join(parts, constants.SummaryJoinSep)
	fmt.Printf("  %s\n\n", line)
}

// buildExecSummaryParts assembles exec summary line segments.
func buildExecSummaryParts(succeeded, failed, missing, total int) []string {
	parts := []string{fmt.Sprintf(constants.SummaryReposFmt, total)}
	parts = appendSummaryPart(parts, succeeded, constants.ColorGreen, constants.SummarySucceededFmt)
	parts = appendSummaryPart(parts, failed, constants.ColorYellow, constants.SummaryFailedFmt)
	parts = appendSummaryPart(parts, missing, constants.ColorYellow, constants.SummaryMissingFmt)

	return parts
}
