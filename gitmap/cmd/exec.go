package cmd

import (
	"encoding/json"
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
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrExecUsage)
		os.Exit(1)
	}

	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadExecRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrExecLoadFailed, err)
		os.Exit(1)
	}

	gitArgs := args
	printExecBanner(gitArgs, len(records))

	var succeeded, failed, missing int
	for _, rec := range records {
		if _, err := os.Stat(rec.AbsolutePath); os.IsNotExist(err) {
			fmt.Printf("  %s⊘ %-22s %snot found%s\n",
				constants.ColorDim, truncateExec(rec.RepoName, 22),
				constants.ColorYellow, constants.ColorReset)
			missing++
			continue
		}

		ok := execInRepo(rec, gitArgs)
		if ok {
			succeeded++
		} else {
			failed++
		}
	}

	printExecSummary(succeeded, failed, missing, len(records))
}

// loadExecRecords reads ScanRecords from gitmap.json.
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

	if err != nil {
		fmt.Printf("  %s✗ %-22s%s\n", constants.ColorYellow, truncateExec(rec.RepoName, 22), constants.ColorReset)
		if len(output) > 0 {
			for _, line := range strings.Split(output, "\n") {
				fmt.Printf("    %s%s%s\n", constants.ColorDim, line, constants.ColorReset)
			}
		}

		return false
	}

	fmt.Printf("  %s✓ %-22s%s\n", constants.ColorGreen, truncateExec(rec.RepoName, 22), constants.ColorReset)
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			fmt.Printf("    %s%s%s\n", constants.ColorDim, line, constants.ColorReset)
		}
	}

	return true
}

// printExecBanner shows the command header.
func printExecBanner(gitArgs []string, count int) {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║           gitmap exec                ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %sCommand: git %s%s\n", constants.ColorWhite, strings.Join(gitArgs, " "), constants.ColorReset)
	fmt.Printf("  %s%d repos from gitmap-output/gitmap.json%s\n", constants.ColorDim, count, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermSeparator, constants.ColorReset)
	fmt.Println()
}

// printExecSummary shows final totals.
func printExecSummary(succeeded, failed, missing, total int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 50), constants.ColorReset)

	parts := []string{fmt.Sprintf("%d repos", total)}
	if succeeded > 0 {
		parts = append(parts, fmt.Sprintf("%s%d succeeded%s", constants.ColorGreen, succeeded, constants.ColorReset))
	}
	if failed > 0 {
		parts = append(parts, fmt.Sprintf("%s%d failed%s", constants.ColorYellow, failed, constants.ColorReset))
	}
	if missing > 0 {
		parts = append(parts, fmt.Sprintf("%s%d missing%s", constants.ColorYellow, missing, constants.ColorReset))
	}

	fmt.Printf("  %s\n\n", strings.Join(parts, " · "))
}

// truncateExec shortens a string to maxLen with ellipsis.
func truncateExec(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-1] + "…"
}
