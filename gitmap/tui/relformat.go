package tui

import "fmt"

// formatRelRow renders a single release list row.
func formatRelRow(r interface {
	GetVersion() string
	GetTag() string
	GetBranch() string
	IsDraft() bool
	IsLatestRelease() bool
	GetSource() string
	GetDate() string
}) string {
	return fmt.Sprintf("%-12s %-14s %-20s %-8s %-8s %-8s %s",
		r.GetVersion(), r.GetTag(), truncateStr(r.GetBranch(), 20),
		boolLabel(r.IsDraft()), boolLabel(r.IsLatestRelease()),
		r.GetSource(), r.GetDate())
}

// writeField writes a labeled field to the builder.
func writeField(b interface{ WriteString(string) (int, error) }, label, value string) {
	if len(value) == 0 {
		value = "-"
	}

	b.WriteString(fmt.Sprintf("  %-16s %s\n", label+":", value))
}

// shortSHA truncates a commit SHA to 8 characters.
func shortSHA(sha string) string {
	if len(sha) > 8 {
		return sha[:8]
	}

	return sha
}

// boolLabel returns yes/no for a boolean.
func boolLabel(v bool) string {
	if v {
		return "yes"
	}

	return "no"
}

// truncateStr truncates a string with ellipsis.
func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-1] + "…"
}
