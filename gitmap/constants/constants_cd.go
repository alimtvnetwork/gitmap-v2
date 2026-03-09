package constants

// CD CLI commands.
const (
	CmdCD      = "cd"
	CmdCDAlias = "go"
)

// CD subcommands.
const (
	CmdCDRepos      = "repos"
	CmdCDSetDefault = "set-default"
	CmdCDClearDefault = "clear-default"
)

// CD help text.
const HelpCD = "  cd (go) <name>      Navigate to a tracked repo directory"

// CD file.
const CDDefaultsFile = "cd-defaults.json"

// CD messages.
const (
	MsgCDMultipleHeader  = "Multiple locations found for \"%s\":\n"
	MsgCDMultipleRowFmt  = "  %d  %s\n"
	MsgCDPickPrompt      = "\nPick [1-%d]: "
	MsgCDReposHeader     = "TRACKED REPOS\n"
	MsgCDReposRowFmt     = "  %d  %s\n"
	MsgCDDefaultSet      = "Default set for %s: %s\n"
	MsgCDDefaultCleared  = "Default cleared for %s\n"
	ErrCDUsage           = "usage: gitmap cd <repo-name|repos> [--group <name>] [--pick]\n"
	ErrCDNotFound        = "no repo found matching '%s'\n"
	ErrCDInvalidPick     = "invalid selection\n"
	ErrCDSetDefaultUsage = "usage: gitmap cd set-default <name> <path>\n"
	ErrCDClearDefaultUsage = "usage: gitmap cd clear-default <name>\n"
	ErrCDDefaultNotFound = "no default set for '%s'\n"
)

// CD flag descriptions.
const (
	FlagDescCDGroup = "Filter repos list by group"
	FlagDescCDPick  = "Force interactive picker even if a default is set"
)
