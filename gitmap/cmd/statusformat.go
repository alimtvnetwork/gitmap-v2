package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/gitutil"
	"github.com/user/gitmap/model"
)

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

	fmt.Printf(constants.StatusRowFmt,
		truncate(rec.RepoName, 22),
		branchText, stateIcon, syncText, stashText, filesText)
}

// printMissingRepo prints a row for a repo not found on disk.
func printMissingRepo(name string, s *statusSummary) {
	fmt.Printf(constants.StatusMissingFmt,
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

	return constants.ColorDim + constants.StatusDash + constants.ColorReset
}

// buildFileCountParts assembles the file count display parts.
func buildFileCountParts(rs gitutil.RepoStatus) string {
	parts := make([]string, 0, 3)
	if rs.Staged > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.StatusStagedFmt+"%s", constants.ColorGreen, rs.Staged, constants.ColorReset))
	}
	if rs.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.StatusModifiedFmt+"%s", constants.ColorYellow, rs.Modified, constants.ColorReset))
	}
	if rs.Untracked > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.StatusUntrackedFmt+"%s", constants.ColorDim, rs.Untracked, constants.ColorReset))
	}

	return strings.Join(parts, " ")
}

// truncate shortens a string to maxLen, adding ellipsis if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-1] + "…"
}
