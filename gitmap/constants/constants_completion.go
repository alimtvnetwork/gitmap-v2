package constants

// Completion shells.

// Completion shells.
const (
	ShellPowerShell = "powershell"
	ShellBash       = "bash"
	ShellZsh        = "zsh"
)

// Completion list flags.
const (
	CompListRepos    = "--list-repos"
	CompListGroups   = "--list-groups"
	CompListCommands = "--list-commands"
)

// Completion file names.
const (
	CompFilePS   = "completions.ps1"
	CompFileBash = "completions.bash"
	CompFileZsh  = "completions.zsh"
	CompDirName  = "gitmap"
)

// Completion help text.
const HelpCompletion = "  completion (cmp)    Generate or install shell tab-completion scripts"

// Completion messages.
const (
	MsgCompInstalled    = "Shell completion installed for %s\n"
	MsgCompAlreadyDone  = "Shell completion already configured for %s\n"
	MsgCompProfileWrite = "Added source line to %s\n"
	ErrCompUsage        = "usage: gitmap completion <powershell|bash|zsh> [--list-repos|--list-groups|--list-commands]\n"
	ErrCompUnknownShell = "unknown shell: %s (use powershell, bash, or zsh)\n"
	ErrCompProfileWrite = "failed to update profile %s: %v\n"
)

// Completion flag descriptions.
const (
	FlagDescCompListRepos    = "Print repo slugs one per line"
	FlagDescCompListGroups   = "Print group names one per line"
	FlagDescCompListCommands = "Print all command names one per line"
)
