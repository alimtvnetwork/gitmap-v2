package constants

// TUI defaults.
const DefaultDashboardRefresh = 30

// TUI command.
const (
	CmdInteractive      = "interactive"
	CmdInteractiveAlias = "i"
)

// TUI help text.
const (
	HelpInteractive = "  interactive (i)     Launch interactive TUI with repo browser and batch actions"
)

// TUI view labels.
const (
	TUIViewBrowser   = "Repos"
	TUIViewActions   = "Actions"
	TUIViewGroups    = "Groups"
	TUIViewDashboard = "Status"
)

// TUI status messages.
const (
	TUITitle          = "gitmap interactive"
	TUISearchPrompt   = "Search: "
	TUINoRepos        = "No repositories found. Run 'gitmap scan' first."
	TUINoGroups       = "No groups found. Press 'c' to create one."
	TUINoSelection    = "No repos selected. Use Space to select."
	TUIConfirmDelete  = "Delete group '%s'? (y/n)"
	TUIGroupCreated   = "Group '%s' created"
	TUIGroupDeleted   = "Group '%s' deleted"
	TUIActionPull     = "Pulling %d repo(s)..."
	TUIActionExec     = "Executing across %d repo(s)..."
	TUIActionStatus   = "Checking status of %d repo(s)..."
	TUIActionComplete = "Action complete: %d success, %d failed"
	TUIRefreshing     = "Refreshing..."
	TUIQuitHint       = "q/esc: quit"
	TUITabHint        = "tab: switch view"
	TUISelectHint     = "space: select  enter: detail  /: search"
	TUIBatchHint      = "p: pull  x: exec  s: status  g: add to group"
	TUIGroupHint      = "c: create  d: delete  enter: show members"
	TUIDashHint       = "r: refresh"
)

// TUI column headers.
const (
	TUIColSlug    = "Slug"
	TUIColBranch  = "Branch"
	TUIColPath    = "Path"
	TUIColType    = "Type"
	TUIColStatus  = "Status"
	TUIColAhead   = "Ahead"
	TUIColBehind  = "Behind"
	TUIColStash   = "Stash"
	TUIColGroup   = "Group"
	TUIColMembers = "Members"
)

// TUI errors.
const (
	ErrTUINoTerminal = "interactive mode requires a terminal — use standard commands instead"
	ErrTUIDBOpen     = "failed to open database for interactive mode: %v"
)
