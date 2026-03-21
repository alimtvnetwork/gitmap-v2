package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

type logsModel struct {
	db      *store.DB
	entries []model.CommandHistoryRecord
	cursor  int
	detail  bool
}

func newLogsModel(db *store.DB) logsModel {
	entries, _ := db.ListHistory()

	return logsModel{
		db:      db,
		entries: entries,
	}
}

func (m logsModel) Update(msg tea.Msg) (logsModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m.handleKey(keyMsg), nil
}

func (m logsModel) handleKey(msg tea.KeyMsg) logsModel {
	max := len(m.entries) - 1
	if max < 0 {
		return m
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
	case keys.enter(msg):
		m.detail = !m.detail
	case keys.refresh(msg):
		m.entries, _ = m.db.ListHistory()
		m.cursor = 0
	}

	return m
}

func (m logsModel) View() string {
	if len(m.entries) == 0 {
		return styleHint.Render(constants.TUILogEmpty)
	}

	if m.detail && m.cursor < len(m.entries) {
		return m.viewDetail()
	}

	return m.viewList()
}

func (m logsModel) viewList() string {
	var b strings.Builder

	header := fmt.Sprintf("  %-4s %-16s %-10s %-30s %-10s %-6s %s",
		"", constants.TUIColCommand, constants.TUIColAlias,
		constants.TUIColArgs, constants.TUIColDuration,
		constants.TUIColExit, constants.TUIColDate)
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, e := range m.entries {
		line := formatLogRow(e)
		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(fmt.Sprintf("  %d log(s)  •  enter: detail  •  r: refresh", len(m.entries))))

	return b.String()
}

func (m logsModel) viewDetail() string {
	e := m.entries[m.cursor]

	var b strings.Builder
	b.WriteString(styleGroupName.Render(fmt.Sprintf("  Command: %s", e.Command)))
	b.WriteString("\n\n")

	writeField(&b, "Alias", e.Alias)
	writeField(&b, "Args", e.Args)
	writeField(&b, "Flags", e.Flags)
	writeField(&b, "Started", e.StartedAt)
	writeField(&b, "Finished", e.FinishedAt)
	writeField(&b, "Duration", formatDurationMs(e.DurationMs))
	writeField(&b, "Exit Code", fmt.Sprintf("%d", e.ExitCode))
	writeField(&b, "Repo Count", fmt.Sprintf("%d", e.RepoCount))

	if len(e.Summary) > 0 {
		b.WriteString("\n")
		writeField(&b, "Summary", e.Summary)
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render("  enter: back to list"))

	return b.String()
}
