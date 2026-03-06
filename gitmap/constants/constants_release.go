package constants

// Setup section headers.
const (
	SetupSectionDiff  = "Diff Tool"
	SetupSectionMerge = "Merge Tool"
	SetupSectionAlias = "Aliases"
	SetupSectionCred  = "Credential Helper"
	SetupSectionCore  = "Core Settings"
	SetupGlobalFlag   = "--global"
)

// Release messages.
const (
	MsgReleaseStart         = "\n  Creating release %s...\n"
	MsgReleaseBranch        = "  ✓ Created branch %s\n"
	MsgReleaseTag           = "  ✓ Created tag %s\n"
	MsgReleasePushed        = "  ✓ Pushed branch and tag to origin\n"
	MsgReleaseMeta          = "  ✓ Release metadata written to %s\n"
	MsgReleaseLatest        = "  ✓ Marked %s as latest release\n"
	MsgReleaseAttach        = "  ✓ Attached %s\n"
	MsgReleaseChangelog     = "  ✓ Using CHANGELOG.md as release body\n"
	MsgReleaseReadme        = "  ✓ Attached README.md\n"
	MsgReleaseDryRun        = "  [dry-run] %s\n"
	MsgReleaseComplete      = "\n  Release %s complete.\n"
	MsgReleaseBranchStart   = "\n  Completing release from %s...\n"
	MsgReleaseVersionRead   = "  → Version from %s: %s\n"
	MsgReleaseBumpResult    = "  → Bumped %s → %s\n"
	MsgReleaseSwitchedBack  = "  ✓ Switched back to %s\n"
	MsgReleasePendingNone   = "  No pending release branches found.\n"
	MsgReleasePendingFound  = "\n  Found %d pending release branch(es).\n"
	MsgReleasePendingFailed = "  ✗ Failed to release %s: %v\n"
	ReleaseBranchPrefix     = "release/"
	ChangelogFile           = "CHANGELOG.md"
	ReadmeFile              = "README.md"
	ReleaseTagPrefix        = "Release "
)
