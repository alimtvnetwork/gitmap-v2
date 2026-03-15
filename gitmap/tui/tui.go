package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

const viewCount = 4

// view indices.
const (
	viewBrowser   = 0
	viewActions   = 1
	viewGroups    = 2
	viewDashboard = 3
)

// rootModel is the top-level Bubble Tea model.
type rootModel struct {
	db        *store.DB
	repos     []model.ScanRecord
	groups    []model.Group
	activeTab int
	width     int
	height    int
	browser   browserModel
	actions   actionsModel
	groupsMgr groupsModel
	dashboard dashboardModel
	quitting  bool
}

// Run launches the interactive TUI.
func Run(db *store.DB, cfg model.Config) error {
	repos, err := db.ListRepos()
	if err != nil {
		return fmt.Errorf(constants.ErrTUIDBOpen, err)
	}

	groups, _ := db.ListGroups()

	m := newRootModel(db, repos, groups, cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()

	return err
}

func newRootModel(db *store.DB, repos []model.ScanRecord, groups []model.Group, cfg model.Config) rootModel {
	return rootModel{
		db:        db,
		repos:     repos,
		groups:    groups,
		activeTab: viewBrowser,
		browser:   newBrowserModel(repos),
		actions:   newActionsModel(),
		groupsMgr: newGroupsModel(groups),
		dashboard: newDashboardModel(repos, cfg.DashboardRefresh),
	}
}

func (m rootModel) Init() tea.Cmd {
	return m.dashboard.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case refreshMsg, tickMsg:
		dm, cmd := m.dashboard.Update(msg)
		m.dashboard = dm

		return m, cmd
	case tea.KeyMsg:
		if keys.quit(msg) && !m.browser.searching {
			m.quitting = true

			return m, tea.Quit
		}
		if keys.tab(msg) && !m.browser.searching {
			m.activeTab = (m.activeTab + 1) % viewCount

			return m, nil
		}
	}

	return m.updateActiveView(msg)
}

func (m rootModel) updateActiveView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.activeTab {
	case viewBrowser:
		bm, cmd := m.browser.Update(msg)
		m.browser = bm

		return m, cmd
	case viewActions:
		am, cmd := m.actions.Update(msg, m.browser.selected())
		m.actions = am

		return m, cmd
	case viewGroups:
		gm, cmd := m.groupsMgr.Update(msg)
		m.groupsMgr = gm

		return m, cmd
	case viewDashboard:
		dm, cmd := m.dashboard.Update(msg)
		m.dashboard = dm

		return m, cmd
	}

	return m, nil
}

func (m rootModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(styleTitle.Render(constants.TUITitle))
	b.WriteString("\n")
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")
	b.WriteString(m.renderContent())
	b.WriteString("\n")
	b.WriteString(m.renderStatusBar())

	return b.String()
}

func (m rootModel) renderTabs() string {
	labels := []string{
		constants.TUIViewBrowser,
		constants.TUIViewActions,
		constants.TUIViewGroups,
		constants.TUIViewDashboard,
	}

	var tabs []string
	for i, label := range labels {
		if i == m.activeTab {
			tabs = append(tabs, styleActiveTab.Render(label))
		} else {
			tabs = append(tabs, styleTab.Render(label))
		}
	}

	return strings.Join(tabs, " ")
}

func (m rootModel) renderContent() string {
	switch m.activeTab {
	case viewBrowser:
		return m.browser.View()
	case viewActions:
		return m.actions.View()
	case viewGroups:
		return m.groupsMgr.View()
	case viewDashboard:
		return m.dashboard.View()
	}

	return ""
}

func (m rootModel) renderStatusBar() string {
	hints := []string{constants.TUIQuitHint, constants.TUITabHint}

	switch m.activeTab {
	case viewBrowser:
		hints = append(hints, constants.TUISelectHint)
	case viewActions:
		hints = append(hints, constants.TUIBatchHint)
	case viewGroups:
		hints = append(hints, constants.TUIGroupHint)
	case viewDashboard:
		hints = append(hints, constants.TUIDashHint)
	}

	return styleStatusBar.Render(strings.Join(hints, "  │  "))
}
