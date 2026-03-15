package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/gitutil"
	"github.com/user/gitmap/model"
)

// statusEntry holds computed git status for one repo.
type statusEntry struct {
	Slug   string
	Branch string
	Status string
	Ahead  int
	Behind int
	Stash  int
}

// refreshMsg carries freshly computed statuses.
type refreshMsg struct {
	entries []statusEntry
}

type dashboardModel struct {
	repos      []model.ScanRecord
	entries    []statusEntry
	cursor     int
	loading    bool
}

func newDashboardModel(repos []model.ScanRecord) dashboardModel {
	return dashboardModel{repos: repos, loading: true}
}

// refreshStatuses collects live git status for all repos.
func refreshStatuses(repos []model.ScanRecord) tea.Cmd {
	return func() tea.Msg {
		entries := make([]statusEntry, 0, len(repos))

		for _, r := range repos {
			rs := gitutil.Status(r.AbsolutePath)
			entries = append(entries, statusEntry{
				Slug:   r.Slug,
				Branch: rs.Branch,
				Status: statusLabel(rs.Dirty, rs.Unreachable),
				Ahead:  rs.Ahead,
				Behind: rs.Behind,
				Stash:  rs.StashCount,
			})
		}

		return refreshMsg{entries: entries}
	}
}

func statusLabel(dirty, unreachable bool) string {
	if unreachable {
		return "error"
	}
	if dirty {
		return "dirty"
	}

	return "clean"
}

func (m dashboardModel) Init() tea.Cmd {
	return refreshStatuses(m.repos)
}

func (m dashboardModel) Update(msg tea.Msg) (dashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		m.entries = msg.entries
		m.loading = false

		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m dashboardModel) handleKey(msg tea.KeyMsg) (dashboardModel, tea.Cmd) {
	max := len(m.entries) - 1
	if max < 0 {
		max = len(m.repos) - 1
	}

	switch {
	case keys.down(msg):
		if m.cursor < max {
			m.cursor++
		}
	case keys.up(msg):
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.refresh(msg):
		m.loading = true

		return m, refreshStatuses(m.repos)
	}

	return m, nil
}

func (m dashboardModel) View() string {
	if len(m.repos) == 0 {
		return styleHint.Render(constants.TUINoRepos)
	}
	if m.loading {
		return styleHint.Render(constants.TUIRefreshing)
	}

	var b strings.Builder

	header := fmt.Sprintf("  %-4s %-20s %-12s %-8s %-6s %-6s %-6s",
		"", constants.TUIColSlug, constants.TUIColBranch,
		constants.TUIColStatus, constants.TUIColAhead,
		constants.TUIColBehind, constants.TUIColStash)
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, e := range m.entries {
		line := m.formatRow(e)
		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(m.summaryLine()))

	return b.String()
}

func (m dashboardModel) formatRow(e statusEntry) string {
	styledStatus := formatStatus(e.Status)

	return fmt.Sprintf("%-20s %-12s %-8s %-6s %-6s %-6s",
		e.Slug, e.Branch, styledStatus,
		formatCount(e.Ahead), formatCount(e.Behind), formatCount(e.Stash))
}

func formatStatus(status string) string {
	switch status {
	case "dirty":
		return styleDirty.Render("dirty")
	case "error":
		return styleDirty.Render("error")
	default:
		return styleClean.Render("clean")
	}
}

func formatCount(n int) string {
	if n == 0 {
		return "-"
	}

	return fmt.Sprintf("%d", n)
}

func (m dashboardModel) summaryLine() string {
	dirty, behind, stash := 0, 0, 0
	for _, e := range m.entries {
		if e.Status == "dirty" {
			dirty++
		}
		if e.Behind > 0 {
			behind++
		}
		if e.Stash > 0 {
			stash++
		}
	}

	ts := time.Now().UTC().Format("15:04:05")

	return fmt.Sprintf("  %d repos  •  %d dirty  •  %d behind  •  %d stash  •  %s UTC  •  r: refresh",
		len(m.entries), dirty, behind, stash, ts)
}
