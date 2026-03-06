package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// runStatus handles the "status" subcommand.
func runStatus() {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadStatusRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrStatusLoadFailed, err)
		os.Exit(1)
	}

	printStatusBanner(len(records))
	summary := printStatusTable(records)
	printStatusSummary(summary)
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
	fmt.Printf("  %s\n\n", strings.Join(parts, constants.SummaryJoinSep))
}

// buildSummaryParts assembles summary line segments.
func buildSummaryParts(s statusSummary) []string {
	parts := []string{
		fmt.Sprintf(constants.SummaryReposFmt, s.Total),
	}
	if s.Clean > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.SummaryCleanFmt+"%s", constants.ColorGreen, s.Clean, constants.ColorReset))
	}
	if s.Dirty > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.SummaryDirtyFmt+"%s", constants.ColorYellow, s.Dirty, constants.ColorReset))
	}
	if s.Ahead > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.SummaryAheadFmt+"%s", constants.ColorCyan, s.Ahead, constants.ColorReset))
	}
	if s.Behind > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.SummaryBehindFmt+"%s", constants.ColorYellow, s.Behind, constants.ColorReset))
	}
	if s.Stashed > 0 {
		parts = append(parts, fmt.Sprintf(constants.SummaryStashedFmt, s.Stashed))
	}
	if s.Missing > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.SummaryMissingFmt+"%s", constants.ColorYellow, s.Missing, constants.ColorReset))
	}

	return parts
}
