package constants

// Hint header.
const MsgHintHeader = "\nHints:\n"

// Hint format.
const MsgHintRowFmt = "  → %-35s %s\n"

// Hint messages for project-repos commands (go-repos, node-repos, etc.).
const (
	HintCDRepo     = "gitmap cd <repo-name>"
	HintCDRepoDesc = "Navigate to a repo"

	HintGroupCreate     = "gitmap g create <name>"
	HintGroupCreateDesc = "Create a group"

	HintLsType     = "gitmap ls go"
	HintLsTypeDesc = "List only Go projects"

	HintGroupAdd     = "gitmap g add <group> <slug>"
	HintGroupAddDesc = "Add repos to a group"

	HintPullGroup     = "gitmap g pull"
	HintPullGroupDesc = "Pull repos in active group"

	HintGroupShow     = "gitmap g show <name>"
	HintGroupShowDesc = "Show repos in a group"

	HintGroupDelete     = "gitmap g delete <name>"
	HintGroupDeleteDesc = "Delete a group"

	HintLsGroups     = "gitmap ls groups"
	HintLsGroupsDesc = "List all groups"

	HintGPull     = "gitmap g pull"
	HintGPullDesc = "Pull active group repos"

	HintGStatus     = "gitmap g status"
	HintGStatusDesc = "Show active group status"

	HintGExec     = "gitmap g exec <cmd>"
	HintGExecDesc = "Run git across active group"

	HintGClear     = "gitmap g clear"
	HintGClearDesc = "Clear active group"

	HintCDSetDefault     = "gitmap cd set-default <name> <path>"
	HintCDSetDefaultDesc = "Set a default repo path"

	HintCDRepos     = "gitmap cd repos"
	HintCDReposDesc = "Browse all repos interactively"

	HintMGUsage     = "gitmap mg g1,g2"
	HintMGUsageDesc = "Select multiple groups"
)
