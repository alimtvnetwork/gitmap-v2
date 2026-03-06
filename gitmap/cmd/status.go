package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/gitutil"
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
	Total    int
	Clean    int
	Dirty    int
	Ahead    int
	Behind   int
	Stashed  int
	Missing  int
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

	// Print header row.
	fmt.Printf("  %s%-22s %-12s %-8s %-10s %-8s %-6s%s\n",
		constants.ColorWhite,
		constants.StatusTableColumns[0], constants.StatusTableColumns[1],
		constants.StatusTableColumns[2], constants.StatusTableColumns[3],
		constants.StatusTableColumns[4], constants.StatusTableColumns[5],
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim,
		constants.TermTableRule, constants.ColorReset)

	for _, rec := range records {
		printOneStatus(rec, &s)
	}

	return s
}

// printOneStatus prints a single repo's status row.
func printOneStatus(rec model.ScanRecord, s *statusSummary) {
	if _, err := os.Stat(rec.AbsolutePath); os.IsNotExist(err) {
		printMissingRepo(rec.RepoName, s)

		return
	}
	rs := gitutil.Status(rec.AbsolutePath)
	stateIcon := formatStateIcon(rs.Dirty, s)
	syncText := formatSyncText(rs.Ahead, rs.Behind, s)
	stashText := formatStashText(rs.StashCount, s)
	filesText := formatFileCounts(rs)
	branchText := fmt.Sprintf("%s%s%s", constants.ColorCyan, truncate(rs.Branch, 12), constants.ColorReset)

	fmt.Printf("  %-22s %s  %s  %s  %s  %s\n",
		truncate(rec.RepoName, 22),
		branchText, stateIcon, syncText, stashText, filesText)
}

// printMissingRepo prints a row for a repo not found on disk.
func printMissingRepo(name string, s *statusSummary) {
	fmt.Printf("  %s%-22s %s⊘ not found%s\n",
		constants.ColorDim, truncate(name, 22),
		constants.ColorYellow, constants.ColorReset)
	s.Missing++
}

// formatStateIcon returns the clean/dirty indicator and updates summary.
func formatStateIcon(dirty bool, s *statusSummary) string {
	if dirty {
		s.Dirty++

		return constants.ColorYellow + constants.StatusIconDirty + constants.ColorReset
	}
	s.Clean++

	return constants.ColorGreen + constants.StatusIconClean + constants.ColorReset
}

// formatSyncText returns the ahead/behind indicator and updates summary.
func formatSyncText(ahead, behind int, s *statusSummary) string {
	if ahead > 0 && behind > 0 {
		s.Ahead++
		s.Behind++

		return fmt.Sprintf("%s"+constants.StatusSyncBothFmt+"%s", constants.ColorYellow, ahead, behind, constants.ColorReset)
	}
	if ahead > 0 {
		s.Ahead++

		return fmt.Sprintf("%s"+constants.StatusSyncUpFmt+"%s", constants.ColorCyan, ahead, constants.ColorReset)
	}
	if behind > 0 {
		s.Behind++

		return fmt.Sprintf("%s"+constants.StatusSyncDownFmt+"%s", constants.ColorYellow, behind, constants.ColorReset)
	}

	return constants.ColorDim + constants.StatusSyncDash + constants.ColorReset
}

// formatStashText returns the stash indicator and updates summary.
func formatStashText(stashCount int, s *statusSummary) string {
	if stashCount > 0 {
		s.Stashed++

		return fmt.Sprintf("%s"+constants.StatusStashFmt+"%s", constants.ColorCyan, stashCount, constants.ColorReset)
	}

	return constants.ColorDim + constants.StatusDash + constants.ColorReset
}

// formatFileCounts returns staged/modified/untracked counts.
func formatFileCounts(rs gitutil.RepoStatus) string {
	if rs.Dirty {

		return buildFileCountParts(rs)
	}

	return constants.ColorDim + "—" + constants.ColorReset
}

// buildFileCountParts assembles the file count display parts.
func buildFileCountParts(rs gitutil.RepoStatus) string {
	parts := make([]string, 0, 3)
	if rs.Staged > 0 {
		parts = append(parts, fmt.Sprintf("%s+%d%s", constants.ColorGreen, rs.Staged, constants.ColorReset))
	}
	if rs.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%s~%d%s", constants.ColorYellow, rs.Modified, constants.ColorReset))
	}
	if rs.Untracked > 0 {
		parts = append(parts, fmt.Sprintf("%s?%d%s", constants.ColorDim, rs.Untracked, constants.ColorReset))
	}

	return strings.Join(parts, " ")
}

// printStatusSummary shows the final totals.
func printStatusSummary(s statusSummary) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 70), constants.ColorReset)

	parts := []string{
		fmt.Sprintf("%d repos", s.Total),
	}
	if s.Clean > 0 {
		parts = append(parts, fmt.Sprintf("%s%d clean%s", constants.ColorGreen, s.Clean, constants.ColorReset))
	}
	if s.Dirty > 0 {
		parts = append(parts, fmt.Sprintf("%s%d dirty%s", constants.ColorYellow, s.Dirty, constants.ColorReset))
	}
	if s.Ahead > 0 {
		parts = append(parts, fmt.Sprintf("%s%d ahead%s", constants.ColorCyan, s.Ahead, constants.ColorReset))
	}
	if s.Behind > 0 {
		parts = append(parts, fmt.Sprintf("%s%d behind%s", constants.ColorYellow, s.Behind, constants.ColorReset))
	}
	if s.Stashed > 0 {
		parts = append(parts, fmt.Sprintf("%d stashed", s.Stashed))
	}
	if s.Missing > 0 {
		parts = append(parts, fmt.Sprintf("%s%d missing%s", constants.ColorYellow, s.Missing, constants.ColorReset))
	}

	fmt.Printf("  %s\n\n", strings.Join(parts, " · "))
}

// truncate shortens a string to maxLen, adding ellipsis if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-1] + "…"
}
