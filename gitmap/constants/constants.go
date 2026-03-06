// Package constants defines all shared constant values for gitmap.
// No magic strings — all literals used for comparison, defaults,
// formats, and file extensions live here.
package constants

// Version.
const Version = "2.3.7"

// RepoPath is set at build time via -ldflags.
var RepoPath = ""

// Clone modes.
const (
	ModeHTTPS = "https"
	ModeSSH   = "ssh"
)

// Output formats.
const (
	OutputTerminal = "terminal"
	OutputCSV      = "csv"
	OutputJSON     = "json"
)

// URL prefixes.
const (
	PrefixHTTPS = "https://"
	PrefixSSH   = "git@"
)

// File extensions.
const (
	ExtCSV  = ".csv"
	ExtJSON = ".json"
	ExtTXT  = ".txt"
	ExtGit  = ".git"
)

// Default file names.
const (
	DefaultCSVFile       = "gitmap.csv"
	DefaultJSONFile      = "gitmap.json"
	DefaultTextFile      = "gitmap.txt"
	DefaultVerboseLogDir = "gitmap-output"
	DefaultStructureFile = "folder-structure.md"
	DefaultCloneScript          = "clone.ps1"
	DefaultDirectCloneScript    = "direct-clone.ps1"
	DefaultDirectCloneSSHScript = "direct-clone-ssh.ps1"
	DefaultDesktopScript        = "register-desktop.ps1"
	DefaultScanCacheFile        = "last-scan.json"
	DefaultConfigPath           = "./data/config.json"
	DefaultSetupConfigPath      = "./data/git-setup.json"
	DefaultOutputDir     = "./gitmap-output"
	DefaultOutputFolder  = "gitmap-output"
	DefaultBranch        = "main"
	DefaultDir           = "."
	DefaultVersionFile   = "version.json"
)

// DefaultReleaseDir is a var so tests can override it.
var DefaultReleaseDir = ".release"

const (
	DefaultLatestFile    = "latest.json"
)

// Git commands and arguments.
const (
	GitBin          = "git"
	GitClone        = "clone"
	GitPull         = "pull"
	GitBranchFlag   = "-b"
	GitDirFlag      = "-C"
	GitFFOnlyFlag   = "--ff-only"
	GitConfigCmd    = "config"
	GitGetFlag      = "--get"
	GitRemoteOrigin = "remote.origin.url"
	GitRevParse     = "rev-parse"
	GitAbbrevRef    = "--abbrev-ref"
	GitHEAD         = "HEAD"
	GitTag          = "tag"
	GitCheckout     = "checkout"
	GitPush         = "push"
	GitLsRemote     = "ls-remote"
	GitLsRemoteTags = "--tags"
)

// Clone instruction format.
const (
	CloneInstructionFmt = "git clone -b %s %s %s"
	HTTPSFromSSHFmt     = "https://%s/%s"
	SSHFromHTTPSFmt     = "git@%s:%s"
)

// Notes.
const (
	NoteNoRemote    = "no remote configured"
	UnknownRepoName = "unknown"
)

// CLI commands.
const (
	CmdScan             = "scan"
	CmdScanAlias        = "s"
	CmdClone            = "clone"
	CmdCloneAlias       = "c"
	CmdUpdate           = "update"
	CmdUpdateRunner     = "update-runner"
	CmdUpdateCleanup    = "update-cleanup"
	CmdVersion          = "version"
	CmdVersionAlias     = "v"
	CmdHelp             = "help"
	CmdDesktopSync      = "desktop-sync"
	CmdDesktopSyncAlias = "ds"
	CmdPull             = "pull"
	CmdPullAlias        = "p"
	CmdRescan           = "rescan"
	CmdRescanAlias      = "rs"
	CmdSetup            = "setup"
	CmdStatus           = "status"
	CmdStatusAlias      = "st"
	CmdExec             = "exec"
	CmdExecAlias        = "x"
	CmdRelease          = "release"
	CmdReleaseAlias     = "r"
	CmdReleaseBranch      = "release-branch"
	CmdReleaseBranchAlias = "rb"
	CmdReleasePending      = "release-pending"
	CmdReleasePendingAlias = "rp"
	CmdChangelog          = "changelog"
	CmdChangelogMD        = "changelog.md"
	CmdDoctor             = "doctor"
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

// ANSI color codes.
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[97m"
	ColorDim    = "\033[90m"
)

// Terminal output sections.
const (
	TermBannerTop    = "  ╔══════════════════════════════════════╗"
	TermBannerTitle  = "  ║            gitmap v%s               ║"
	TermBannerBottom = "  ╚══════════════════════════════════════╝"
	TermFoundFmt     = "  ✓ Found %d repositories"
	TermReposHeader  = "  ■ Repositories"
	TermTreeHeader   = "  ■ Folder Structure"
	TermCloneHeader  = "  ■ How to Clone on Another Machine"
	TermSeparator    = "  ──────────────────────────────────────────"
)

// Terminal repo entry formats.
const (
	TermRepoIcon  = "  📦 %s\n"
	TermPathLine  = "     Path:  %s\n"
	TermCloneLine = "     Clone: %s\n"
)

// Terminal clone help text.
const (
	TermCloneStep1     = "  1. Copy the output files to the target machine:"
	TermCloneCmd1      = "     gitmap-output/gitmap.json  (or .csv / .txt)"
	TermCloneStep2     = "  2. Clone via JSON (shorthand):"
	TermCloneCmd2      = "     gitmap clone json --target-dir ./projects"
	TermCloneCmd2Alt   = "     gitmap c json               # alias"
	TermCloneStep3     = "  3. Clone via CSV (shorthand):"
	TermCloneCmd3      = "     gitmap clone csv --target-dir ./projects"
	TermCloneCmd3Alt   = "     gitmap c csv                # alias"
	TermCloneStep3t    = "  4. Clone via text (shorthand):"
	TermCloneCmd3t     = "     gitmap clone text --target-dir ./projects"
	TermCloneCmd3tAlt  = "     gitmap c text               # alias"
	TermCloneStep3b    = "  5. Or specify a file path directly:"
	TermCloneCmd3b     = "     gitmap clone ./gitmap-output/gitmap.json --target-dir ./projects"
	TermCloneStep4     = "  6. Or run the PowerShell script directly:"
	TermCloneCmd4HTTPS = "     .\\direct-clone.ps1       # HTTPS clone commands"
	TermCloneCmd4SSH   = "     .\\direct-clone-ssh.ps1   # SSH clone commands"
	TermCloneStep5     = "  7. Full clone script with progress & error handling:"
	TermCloneCmd5      = "     .\\clone.ps1 -TargetDir .\\projects"
	TermCloneStep6     = "  8. Sync repos to GitHub Desktop:"
	TermCloneCmd6      = "     gitmap desktop-sync         # or: gitmap ds"
	TermCloneNote      = "  Note: safe-pull is auto-enabled when existing repos are detected."
)

// JSON formatting.
const JSONIndent = "  "

// CLI messages.
const (
	MsgFoundRepos       = "Found %d repositories.\n"
	MsgCSVWritten       = "CSV written to %s\n"
	MsgJSONWritten      = "JSON written to %s\n"
	MsgTextWritten      = "Text clone list written to %s\n"
	MsgStructureWritten = "Folder structure written to %s\n"
	MsgCloneScript      = "Clone script written to %s\n"
	MsgDirectClone      = "Direct clone script written to %s\n"
	MsgDirectCloneSSH   = "Direct SSH clone script written to %s\n"
	MsgDesktopScript    = "Desktop registration script written to %s\n"
	MsgCloneComplete    = "\nClone complete: %d succeeded, %d failed\n"
	MsgAutoSafePull     = "Existing repos detected — safe-pull enabled automatically.\n"
	MsgOpenedFolder     = "Opened output folder: %s\n"
	MsgVerboseLogFile   = "Verbose log: %s\n"
	MsgDesktopSyncStart   = "\n  Syncing repos to GitHub Desktop from %s...\n"
	MsgDesktopSyncSkipped = "  ⊘ Skipped (already exists): %s\n"
	MsgDesktopSyncAdded   = "  ✓ Added to GitHub Desktop: %s\n"
	MsgDesktopSyncFailed  = "  ✗ Failed: %s — %v\n"
	MsgDesktopSyncDone    = "\n  GitHub Desktop sync: %d added, %d skipped, %d failed\n"
	MsgNoOutputDir        = "Error: gitmap-output/ not found in current directory.\nRun 'gitmap scan' first to generate output files."
	MsgNoJSONFile         = "Error: %s not found.\nRun 'gitmap scan' first to generate the JSON output."
	MsgFailedClones     = "\nFailed clones:"
	MsgFailedEntry      = "  - %s (%s): %s\n"
	MsgPullStarting     = "\n  Pulling %s (%s)...\n"
	MsgPullSuccess      = "  ✓ %s is up to date.\n"
	MsgPullFailed       = "  ✗ Pull failed for %s: %s\n"
	MsgPullAvailable    = "\nAvailable repos:"
	MsgRescanReplay     = "\n  Rescanning with cached flags (dir: %s)...\n"
	MsgScanCacheSaved   = "Scan cache written to %s\n"
	MsgUpdateStarting   = "\n  Updating gitmap from source repo...\n"
	MsgUpdateRepoPath   = "  → Repo path: %s\n"
	MsgUpdateVersion    = "\n  ✓ Updated to gitmap v%s\n"
)

// Folder structure Markdown.
const (
	StructureTitle       = "# Folder Structure"
	StructureDescription = "Git repositories discovered by gitmap."
	StructureRepoFmt     = "📦 **%s** (`%s`) — %s"
	TreeBranch           = "├──"
	TreeCorner           = "└──"
	TreePipe             = "│   "
	TreeSpace            = "    "
)

// Clone shorthands.
const (
	ShorthandJSON = "json"
	ShorthandCSV  = "csv"
	ShorthandText = "text"
)

// CLI error messages.
const (
	ErrSourceRequired    = "Error: source file is required"
	ErrCloneUsage        = "Usage: gitmap clone <source|json|csv|text> [--target-dir <dir>] [--safe-pull]"
	ErrShorthandNotFound = "Error: %s not found.\nRun 'gitmap scan' first to generate output files.\n"
	ErrConfigLoad     = "Error loading config: %v\n"
	ErrScanFailed     = "Scan error: %v\n"
	ErrCloneFailed    = "Clone error: %v\n"
	ErrOutputFailed   = "Output error: %v\n"
	ErrCreateDir      = "Cannot create directory: %v\n"
	ErrCreateFile     = "Cannot create file: %v\n"
	ErrNoRepoPath       = "Error: repo path not embedded. Binary was not built with run.ps1."
	ErrUpdateFailed     = "Update error: %v\n"
	ErrPullSlugRequired = "Error: repo name is required"
	ErrPullUsage        = "Usage: gitmap pull <repo-name> [--verbose]"
	ErrPullLoadFailed   = "Error: could not load gitmap.json: %v\n"
	ErrPullNotFound     = "Error: no repo found matching '%s'\n"
	ErrPullNotRepo      = "Error: %s is not a git repository\n"
	ErrRescanNoCache    = "Error: no previous scan found. Run 'gitmap scan' first.\n%v\n"
	ErrSetupLoadFailed  = "Error: could not load git-setup.json: %v\n"
	ErrStatusLoadFailed = "Error: could not load gitmap.json for status: %v\nRun 'gitmap scan' first.\n"
	ErrExecUsage        = "Usage: gitmap exec <git-args...>\nExample: gitmap exec fetch --prune"
	ErrExecLoadFailed   = "Error: could not load gitmap.json: %v\nRun 'gitmap scan' first.\n"
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

// CLI help text.
const (
	HelpUsage      = "Usage: gitmap <command> [flags]"
	HelpCommands   = "Commands:"
	HelpScan       = "  scan (s) [dir]      Scan directory for Git repos"
	HelpClone      = "  clone (c) <source|json|csv|text>  Re-clone from file (shorthands auto-resolve)"
	HelpUpdate     = "  update              Self-update from source repo"
	HelpUpdateCleanup = "  update-cleanup      Remove leftover update temp files and .old backups"
	HelpVersion    = "  version (v)         Show version number"
	HelpDesktopSync = "  desktop-sync (ds)   Sync repos to GitHub Desktop from output"
	HelpPull        = "  pull (p) <name>     Pull a specific repo by its name"
	HelpRescan      = "  rescan (rs)         Re-run last scan with cached flags"
	HelpSetup       = "  setup               Configure Git diff/merge tool, aliases & core settings"
	HelpStatus      = "  status (st)         Show dirty/clean, ahead/behind, stash for all repos"
	HelpExec        = "  exec (x) <args...>  Run any git command across all repos"
	HelpRelease     = "  release (r) [ver]   Create release branch, tag, and push"
	HelpReleaseBr   = "  release-branch (rb) Complete release from existing release branch"
	HelpReleasePend = "  release-pending (rp) Release all pending branches without tags"
	HelpChangelog   = "  changelog (cl) [ver] Show concise release notes (use --open or changelog.md)"
	HelpDoctor      = "  doctor [--fix-path] Diagnose PATH, deploy, and version issues"
	HelpHelp        = "  help                Show this help message"
	HelpScanFlags  = "Scan flags:"
	HelpConfig     = "  --config <path>     Config file (default: ./data/config.json)"
	HelpMode       = "  --mode ssh|https    Clone URL style (default: https)"
	HelpOutput     = "  --output csv|json|terminal  Output format (default: terminal)"
	HelpOutputPath = "  --output-path <dir> Output directory (default: ./gitmap-output)"
	HelpOutFile        = "  --out-file <path>   Exact output file path"
	HelpGitHubDesktop  = "  --github-desktop    Add repos to GitHub Desktop"
	HelpOpen           = "  --open              Open output folder after scan"
	HelpQuiet          = "  --quiet             Suppress clone help section (for CI/scripted use)"
	HelpCloneFlags    = "Clone flags:"
	HelpTargetDir     = "  --target-dir <dir>  Base directory for clones (default: .)"
	HelpSafePull      = "  --safe-pull         Pull existing repos with retry + unlock diagnostics (auto-enabled)"
	HelpVerbose       = "  --verbose           Write detailed debug log to a timestamped file"
	HelpReleaseFlags  = "Release flags:"
	HelpAssets        = "  --assets <path>     Directory or file to attach to the release"
	HelpCommit        = "  --commit <sha>      Create release from a specific commit"
	HelpRelBranch     = "  --branch <name>     Create release from latest commit of a branch"
	HelpBump          = "  --bump major|minor|patch  Auto-increment from latest released version"
	HelpDraft         = "  --draft             Create an unpublished draft release"
	HelpDryRun        = "  --dry-run           Preview release steps without executing"
)

// Flag descriptions.
const (
	FlagDescConfig     = "Path to config file"
	FlagDescMode       = "Clone URL style: https or ssh"
	FlagDescOutput     = "Output format: terminal, csv, json"
	FlagDescOutFile    = "Exact output file path"
	FlagDescOutputPath = "Output directory for CSV/JSON"
	FlagDescTargetDir  = "Base directory for cloned repos"
	FlagDescSafePull   = "If repo exists, run safe git pull with retries and unlock diagnostics"
	FlagDescGHDesktop  = "Add discovered repos to GitHub Desktop"
	FlagDescOpen       = "Open output folder after scan completes"
	FlagDescQuiet      = "Suppress terminal clone help section"
	FlagDescVerbose    = "Write detailed stdout/stderr debug log to a timestamped file"
	FlagDescSetupConfig = "Path to git-setup.json config file"
	FlagDescDryRun     = "Preview changes without applying them"
	FlagDescAssets     = "Directory or file to attach to the release"
	FlagDescCommit     = "Create release from a specific commit"
	FlagDescRelBranch  = "Create release from latest commit of a branch"
	FlagDescBump       = "Auto-increment version: major, minor, or patch"
	FlagDescDraft      = "Create an unpublished draft release"
	FlagDescLatest        = "Show only the latest changelog entry"
	FlagDescLimit         = "Number of changelog versions to show"
	FlagDescOpenChangelog = "Open CHANGELOG.md with the default system app"
)

// Clone and Desktop scripts are now generated from Go templates
// embedded in formatter/templates/. See clone.ps1.tmpl and desktop.ps1.tmpl.

// Directory permissions.
const DirPermission = 0o755

// Safe-pull defaults.
const (
	SafePullRetryAttempts    = 4
	SafePullRetryDelayMS     = 600
	WindowsPathWarnThreshold = 240
)

// Verbose log file.
const VerboseLogFileFmt = "gitmap-verbose-%s.log"

// Setup section headers.
const (
	SetupSectionDiff  = "Diff Tool"
	SetupSectionMerge = "Merge Tool"
	SetupSectionAlias = "Aliases"
	SetupSectionCred  = "Credential Helper"
	SetupSectionCore  = "Core Settings"
)

// Release messages.
const (
	MsgReleaseStart       = "\n  Creating release %s...\n"
	MsgReleaseBranch      = "  ✓ Created branch %s\n"
	MsgReleaseTag         = "  ✓ Created tag %s\n"
	MsgReleasePushed      = "  ✓ Pushed branch and tag to origin\n"
	MsgReleaseMeta        = "  ✓ Release metadata written to %s\n"
	MsgReleaseLatest      = "  ✓ Marked %s as latest release\n"
	MsgReleaseAttach      = "  ✓ Attached %s\n"
	MsgReleaseChangelog   = "  ✓ Using CHANGELOG.md as release body\n"
	MsgReleaseReadme      = "  ✓ Attached README.md\n"
	MsgReleaseDryRun      = "  [dry-run] %s\n"
	MsgReleaseComplete    = "\n  Release %s complete.\n"
	MsgReleaseBranchStart  = "\n  Completing release from %s...\n"
	MsgReleaseVersionRead  = "  → Version from %s: %s\n"
	MsgReleaseBumpResult   = "  → Bumped %s → %s\n"
	MsgReleaseSwitchedBack = "  ✓ Switched back to %s\n"
	MsgReleasePendingNone  = "  No pending release branches found.\n"
	MsgReleasePendingFound = "\n  Found %d pending release branch(es).\n"
	MsgReleasePendingFailed = "  ✗ Failed to release %s: %v\n"
	ReleaseBranchPrefix    = "release/"
	ChangelogFile          = "CHANGELOG.md"
	ReadmeFile             = "README.md"
)
