package constants

// Clone-next command messages.
const (
	MsgCloneNextCloning      = "Cloning %s into %s...\n"
	MsgCloneNextDone         = "✓ Cloned %s\n"
	MsgCloneNextDesktop      = "✓ Registered %s with GitHub Desktop\n"
	MsgCloneNextRemovePrompt = "Remove current folder %s? [y/N] "
	MsgCloneNextRemoved      = "✓ Removed %s\n"
)

// Clone-next error and warning messages.
const (
	ErrCloneNextUsage       = "Usage: gitmap clone-next <v++|vN> [flags]"
	ErrCloneNextCwd         = "Error: cannot determine current directory: %v\n"
	ErrCloneNextNoRemote    = "Error: not a git repo or no remote origin: %v\n"
	ErrCloneNextBadVersion  = "Error: %v\n"
	ErrCloneNextExists      = "Error: target directory already exists: %s\nUse 'cd' to switch to it.\n"
	ErrCloneNextFailed      = "Error: clone failed for %s\n"
	WarnCloneNextRemoveFailed = "Warning: could not remove %s: %v\n"
)

// Clone-next flag descriptions.
const (
	FlagDescCloneNextDelete    = "Auto-remove current folder after clone"
	FlagDescCloneNextKeep      = "Keep current folder without prompting"
	FlagDescCloneNextNoDesktop = "Skip GitHub Desktop registration"
)

// Clone-next help strings for usage output.
const (
	HelpCloneNextFlags   = "Clone-Next Flags:"
	HelpCNDelete         = "  --delete            Auto-remove current version folder after clone"
	HelpCNKeep           = "  --keep              Keep current folder without prompting for removal"
	HelpCNNoDesktop      = "  --no-desktop        Skip GitHub Desktop registration"
	HelpCNSSHKey         = "  --ssh-key, -K       SSH key name to use for clone"
	HelpCNVerbose        = "  --verbose           Show detailed clone-next output"
)
