package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// runStatus handles the "status" subcommand.
func runStatus(args []string) {
	groupName, all := parseStatusFlags(args)
	records := loadStatusByScope(groupName, all)

	printStatusBanner(len(records))
	summary := printStatusTable(records)
	printStatusSummary(summary)
}

// parseStatusFlags parses --group and --all flags.
func parseStatusFlags(args []string) (groupName string, all bool) {
	fs := flag.NewFlagSet(constants.CmdStatus, flag.ExitOnError)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	fs.Parse(args)

	return *gFlag, *aFlag
}

// loadStatusByScope returns records filtered by group, all DB repos, or JSON fallback.
func loadStatusByScope(groupName string, all bool) []model.ScanRecord {
	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}
	if all {
		return loadAllRecordsDB()
	}

	return loadRecordsJSONFallback()
}

// loadRecordsByGroup loads repos from a specific group in the database.
func loadRecordsByGroup(groupName string) []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	records, err := db.ShowGroup(groupName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGenericFmt, err)
		os.Exit(1)
	}

	return records
}

// loadAllRecordsDB loads all repos from the database.
func loadAllRecordsDB() []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	records, err := db.ListRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGenericFmt, err)
		os.Exit(1)
	}

	return records
}

// loadRecordsJSONFallback loads records from gitmap.json.
func loadRecordsJSONFallback() []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadStatusRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrStatusLoadFailed, err)
		os.Exit(1)
	}

	return records
}

// loadStatusRecords reads ScanRecords from gitmap.json.
func loadStatusRecords(path string) ([]model.ScanRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var records []model.ScanRecord
	err = json.Unmarshal(data, &records)

	return records, err
}

// statusSummary aggregates counts across all repos.
type statusSummary struct {
	Total   int
	Clean   int
	Dirty   int
	Ahead   int
	Behind  int
	Stashed int
	Missing int
}

// printStatusBanner shows the dashboard header.
func printStatusBanner(count int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.StatusBannerTop, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.StatusBannerTitle, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.StatusBannerBottom, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %s"+constants.StatusRepoCountFmt+"%s\n", constants.ColorDim, count, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermSeparator, constants.ColorReset)
	fmt.Println()
}

// printStatusTable prints each repo's status and returns a summary.
func printStatusTable(records []model.ScanRecord) statusSummary {
	s := statusSummary{Total: len(records)}
	printStatusHeader()

	for _, rec := range records {
		printOneStatus(rec, &s)
	}

	return s
}

// printStatusHeader prints the table column header row.
func printStatusHeader() {
	fmt.Printf(constants.StatusHeaderFmt,
		constants.ColorWhite,
		constants.StatusTableColumns[0], constants.StatusTableColumns[1],
		constants.StatusTableColumns[2], constants.StatusTableColumns[3],
		constants.StatusTableColumns[4], constants.StatusTableColumns[5],
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim,
		constants.TermTableRule, constants.ColorReset)
}

// printStatusSummary shows the final totals.
func printStatusSummary(s statusSummary) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	parts := buildSummaryParts(s)
	line := strings.Join(parts, constants.SummaryJoinSep)
	fmt.Printf("  %s\n\n", line)
}

// buildSummaryParts assembles summary line segments.
func buildSummaryParts(s statusSummary) []string {
	parts := []string{fmt.Sprintf(constants.SummaryReposFmt, s.Total)}
	parts = appendSummaryPart(parts, s.Clean, constants.ColorGreen, constants.SummaryCleanFmt)
	parts = appendSummaryPart(parts, s.Dirty, constants.ColorYellow, constants.SummaryDirtyFmt)
	parts = appendSummaryPart(parts, s.Ahead, constants.ColorCyan, constants.SummaryAheadFmt)
	parts = appendSummaryPart(parts, s.Behind, constants.ColorYellow, constants.SummaryBehindFmt)
	parts = appendSummaryPart(parts, s.Stashed, "", constants.SummaryStashedFmt)
	parts = appendSummaryPart(parts, s.Missing, constants.ColorYellow, constants.SummaryMissingFmt)

	return parts
}

// appendSummaryPart conditionally appends a colored summary segment.
func appendSummaryPart(parts []string, count int, color, format string) []string {
	if count == 0 {
		return parts
	}
	if len(color) > 0 {
		colored := fmt.Sprintf("%s"+format+"%s", color, count, constants.ColorReset)

		return append(parts, colored)
	}

	return append(parts, fmt.Sprintf(format, count))
}
