package constants

// Notes.
const (
	NoteNoRemote    = "no remote configured"
	UnknownRepoName = "unknown"
)

// GitHub Desktop.
const (
	GitHubDesktopBin   = "github"
	OSWindows          = "windows"
	MsgDesktopNotFound = "GitHub Desktop CLI not found — skipping."
	MsgDesktopAdded    = "  ✓ Added to GitHub Desktop: %s\n"
	MsgDesktopFailed   = "  ✗ Failed to add %s: %v\n"
	MsgDesktopSummary  = "GitHub Desktop: %d added, %d failed\n"
)

// Latest-branch display messages.
const (
	MsgLatestBranchFetching     = "  Fetching remotes..."
	MsgLatestBranchFetchWarning = "  Warning: fetch failed: %v\n"
	LBUnknownBranch             = "<unknown>"
)

// Generic error formats.
const (
	ErrGenericFmt = "Error: %v\n"
	ErrBareFmt    = "%v\n"
)

// OS platform constants.
const OSDarwin = "darwin"

// OS file-explorer commands.
const (
	CmdExplorer     = "explorer"
	CmdOpen         = "open"
	CmdXdgOpen      = "xdg-open"
	CmdWindowsShell = "cmd"
	CmdArgSlashC    = "/c"
	CmdArgStart     = "start"
	CmdArgEmpty     = ""
)

// Desktop sync error messages.
const (
	ErrDesktopReadFailed  = "Error reading %s: %v\n"
	ErrDesktopParseFailed = "Error parsing JSON from %s: %v\n"
	ErrNoAbsPath          = "no absolute path"
)

// Command dispatch errors.
const (
	ErrUnknownCommand  = "Unknown command: %s\n"
	ErrUnknownGroupSub = "Unknown group subcommand: %s\n"
)

// Docs command.
const (
	DocsURL        = "https://gitmap.dev/docs"
	MsgDocsOpened  = "  ✓ Opened %s\n"
	ErrDocsOpen    = "  ✗ Failed to open browser: %v\n"
)

// Version display.
const MsgVersionFmt = "gitmap v%s\n"

// CLI messages.
const (
	MsgFoundRepos         = "Found %d repositories.\n"
	MsgCSVWritten         = "CSV written to %s\n"
	MsgJSONWritten        = "JSON written to %s\n"
	MsgTextWritten        = "Text clone list written to %s\n"
	MsgStructureWritten   = "Folder structure written to %s\n"
	MsgCloneScript        = "Clone script written to %s\n"
	MsgDirectClone        = "Direct clone script written to %s\n"
	MsgDirectCloneSSH     = "Direct SSH clone script written to %s\n"
	MsgDesktopScript      = "Desktop registration script written to %s\n"
	MsgCloneComplete      = "\nClone complete: %d succeeded, %d failed\n"
	MsgAutoSafePull       = "Existing repos detected — safe-pull enabled automatically.\n"
	MsgOpenedFolder       = "Opened output folder: %s\n"
	MsgVerboseLogFile     = "Verbose log: %s\n"
	MsgDesktopSyncStart   = "\n  Syncing repos to GitHub Desktop from %s...\n"
	MsgDesktopSyncSkipped = "  ⊘ Skipped (already exists): %s\n"
	MsgDesktopSyncAdded   = "  ✓ Added to GitHub Desktop: %s\n"
	MsgDesktopSyncFailed  = "  ✗ Failed: %s — %v\n"
	MsgDesktopSyncDone    = "\n  GitHub Desktop sync: %d added, %d skipped, %d failed\n"
	MsgNoOutputDir        = "Error: gitmap-output/ not found in current directory.\nRun 'gitmap scan' first to generate output files."
	MsgNoJSONFile         = "Error: %s not found.\nRun 'gitmap scan' first to generate the JSON output."
	MsgFailedClones       = "\nFailed clones:"
	MsgFailedEntry        = "  - %s (%s): %s\n"
	MsgPullStarting       = "\n  Pulling %s (%s)...\n"
	MsgPullSuccess        = "  ✓ %s is up to date.\n"
	MsgPullFailed         = "  ✗ Pull failed for %s: %s\n"
	MsgPullAvailable      = "\nAvailable repos:"
	MsgPullListEntry      = "  - %s\n"
	WarnVerboseLogFailed  = "Warning: could not create verbose log: %v\n"
	MsgRescanReplay       = "\n  Rescanning with cached flags (dir: %s)...\n"
	MsgScanCacheSaved     = "Scan cache written to %s\n"
	MsgDBUpsertDone       = "Database updated: %d repos upserted\n"
	MsgDBUpsertFailed     = "Warning: database upsert failed: %v\n"
	MsgUpdateStarting     = "\n  Updating gitmap from source repo...\n"
	MsgUpdateRepoPath     = "  → Repo path: %s\n"
	MsgUpdateVersion      = "\n  ✓ Updated to gitmap v%s\n"
)

// List and group messages.
const (
	MsgListHeader       = "SLUG                 REPO NAME"
	MsgListSeparator    = "──────────────────────────────────────────"
	MsgListRowFmt       = "%-20s %s\n"
	MsgListVerboseFmt   = "%-20s %-20s %s\n"
	MsgListEmpty        = "No repos tracked. Run 'gitmap scan' first."
	MsgListDBPath       = "  → Database: %s\n"
	MsgGroupCreated     = "Group created: %s\n"
	MsgGroupDeleted     = "Group deleted: %s\n"
	MsgGroupAdded       = "Added %s to group %s\n"
	MsgGroupRemoved     = "Removed %s from group %s\n"
	MsgGroupHeader      = "GROUP           REPOS   DESCRIPTION"
	MsgGroupRowFmt      = "%-15s %-7d %s\n"
	MsgGroupShowHeader  = "Group: %s (%d repos)\n"
	MsgGroupShowRowFmt  = "  %-16s %s\n"
	MsgGroupEmpty       = "No groups defined. Use 'gitmap group create <name>' to create one."
	MsgGroupActivated   = "Active group set: %s\n"
	MsgGroupActiveShow  = "Active group: %s\n"
	MsgGroupNoActive    = "No active group. Use 'gitmap g <name>' to set one."
	MsgGroupCleared     = "Active group cleared."
	ErrGroupNameReq     = "Error: group name is required"
	ErrGroupUsage       = "Usage: gitmap group <create|add|remove|list|show|delete> [args]"
	ErrGroupSlugReq     = "Error: at least one slug is required"
	ErrListDBFailed     = "Error: could not open database: %v\nRun 'gitmap scan' first.\n"
	ErrNoDatabase       = "No database found. Run 'gitmap scan' first."
	MsgDBResetDone      = "Database reset: all repos and groups cleared.\n"
	ErrDBResetFailed    = "Error: database reset failed: %v\n"
	ErrDBResetNoConfirm = "Error: this will delete all tracked repos and groups.\nRun with --confirm to proceed: gitmap db-reset --confirm"
)

// Latest-branch error messages.
const (
	ErrLatestBranchNotRepo   = "Error: not inside a Git repository."
	ErrLatestBranchNoRefs    = "Error: no remote-tracking branches found for remote '%s'.\n"
	ErrLatestBranchNoRefsAll = "Error: no remote-tracking branches found on any remote."
	ErrLatestBranchNoCommits = "Error: could not read commit info for remote branches."
	ErrLatestBranchNoMatch   = "Error: no branches matching filter '%s'.\n"
)

// CLI error messages.
const (
	ErrSourceRequired         = "Error: source file is required"
	ErrCloneUsage             = "Usage: gitmap clone <source|json|csv|text> [--target-dir <dir>] [--safe-pull]"
	ErrShorthandNotFound      = "Error: %s not found.\nRun 'gitmap scan' first to generate output files.\n"
	ErrConfigLoad             = "Error loading config: %v\n"
	ErrScanFailed             = "Scan error: %v\n"
	ErrCloneFailed            = "Clone error: %v\n"
	ErrOutputFailed           = "Output error: %v\n"
	ErrCreateDir              = "Cannot create directory: %v\n"
	ErrCreateFile             = "Cannot create file: %v\n"
	ErrNoRepoPath             = "Error: repo path not embedded. Binary was not built with run.ps1."
	ErrUpdateFailed           = "Update error: %v\n"
	ErrPullSlugRequired       = "Error: repo name is required"
	ErrPullUsage              = "Usage: gitmap pull <repo-name> [--verbose]"
	ErrPullLoadFailed         = "Error: could not load gitmap.json: %v\n"
	ErrPullNotFound           = "Error: no repo found matching '%s'\n"
	ErrPullNotRepo            = "Error: %s is not a git repository\n"
	ErrRescanNoCache          = "Error: no previous scan found. Run 'gitmap scan' first.\n%v\n"
	ErrSetupLoadFailed        = "Error: could not load git-setup.json: %v\n"
	ErrStatusLoadFailed       = "Error: could not load gitmap.json for status: %v\nRun 'gitmap scan' first.\n"
	ErrExecUsage              = "Usage: gitmap exec <git-args...>\nExample: gitmap exec fetch --prune"
	ErrExecLoadFailed         = "Error: could not load gitmap.json: %v\nRun 'gitmap scan' first.\n"
	ErrReleaseVersionRequired = "Error: version is required.\nProvide a version argument, use --bump, or create a version.json file."
	ErrReleaseUsage           = "Usage: gitmap release [version] [--assets <path>] [--commit <sha>] [--branch <name>] [--bump major|minor|patch] [--draft] [--dry-run]"
	ErrReleaseBranchUsage     = "Usage: gitmap release-branch <release/vX.Y.Z> [--assets <path>] [--draft]"
	ErrReleaseAlreadyExists   = "Error: version %s is already released. See .release/%s.json for details.\n"
	ErrReleaseTagExists       = "Error: tag %s already exists.\n"
	ErrReleaseBranchNotFound  = "Error: branch %s does not exist.\n"
	ErrReleaseCommitNotFound  = "Error: commit %s not found.\n"
	ErrReleaseInvalidVersion  = "Error: '%s' is not a valid version.\n"
	ErrReleaseBumpNoLatest    = "Error: no previous release found. Create an initial release before using --bump.\n"
	ErrReleaseBumpConflict    = "Error: --bump cannot be used with an explicit version argument.\n"
	ErrReleaseCommitBranch    = "Error: --commit and --branch are mutually exclusive.\n"
	ErrReleasePushFailed      = "Error: failed to push to remote: %v\n"
	ErrReleaseVersionLoad     = "Error: could not read version.json: %v\n"
	ErrReleaseMetaWrite       = "Error: could not write release metadata: %v\n"
	ErrChangelogRead            = "Error: could not read CHANGELOG.md: %v\n"
	ErrChangelogVersionNotFound = "Error: version %s not found in CHANGELOG.md\n"
	ErrChangelogOpen            = "Error: could not open CHANGELOG.md: %v\n"
)

// List-versions error messages.
const (
	ErrListVersionsNoTags = "Error: no version tags found. Create a release first."
)

// List-releases messages.
const (
	MsgListReleasesEmpty     = "No releases found."
	MsgListReleasesHeader    = "Releases (%d found)\n"
	MsgListReleasesSeparator = "────────────────────────────────────────────────────────────────────────"
	MsgListReleasesColumns   = "  VERSION    TAG          BRANCH              DRAFT  LATEST  SOURCE   DATE"
	MsgListReleasesRowFmt    = "  %-10s %-12s %-19s %-6s %-7s %-8s %s\n"
	ErrListReleasesFailed    = "Error: failed to load releases: %v\n"
	MsgYes                   = "yes"
	MsgNo                    = "no"
)

// Release import messages.
const (
	MsgReleasesImported   = "Releases imported: %d from .release/\n"
	WarnReleaseImportSkip = "Warning: skipping %s: %v\n"
	ReleaseGlob           = "v*.json"
)

// Pending metadata discovery messages.
const (
	MsgPendingMetaFound    = "  → Found %d unreleased version(s) from .release/ metadata\n"
	MsgPendingMetaRelease  = "  → Creating release from metadata: %s (commit: %s)\n"
	WarnPendingMetaNoCommit = "  ⚠ Skipping %s: commit %s not found in repository\n"
	WarnPendingMetaNoSHA   = "  ⚠ Skipping %s: no commit SHA in metadata\n"
)

// Clear release JSON messages.
const (
	MsgClearReleaseDone      = "  ✓ Removed .release/%s.json\n"
	MsgClearReleaseDryRun    = "  [dry-run] Would remove %s\n"
	ErrClearReleaseUsage     = "Usage: gitmap clear-release-json <version> [--dry-run]\nExample: gitmap clear-release-json v2.20.0"
	ErrClearReleaseNotFound  = "Error: no release file found for %s\n"
	ErrClearReleaseFailed    = "Error: could not remove release file: %v\n"
)

// Revert messages.
const (
	MsgRevertCheckout       = "  → Checking out %s...\n"
	MsgRevertStarting       = "\n  Building reverted version...\n"
	MsgRevertDone           = "\n  ✓ Revert complete.\n"
	ErrRevertUsage          = "Usage: gitmap revert <version>\nExample: gitmap revert v2.9.0"
	ErrRevertTagNotFound    = "Error: tag %s not found locally. Run 'git fetch --tags' first.\n"
	ErrRevertCheckoutFailed = "Error: git checkout failed: %v\n"
	ErrRevertFailed         = "Revert error: %v\n"
	RevertScriptLogExec     = "executing revert script: %s"
	RevertScriptLogExit     = "revert script exited: err=%v"
)
