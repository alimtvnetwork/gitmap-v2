// Package constants defines all shared constant values for gitmap.
// No magic strings — all literals used for comparison, defaults,
// formats, and file extensions live here.
package constants

// Version.
const Version = "1.3.4"

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
	DefaultStructureFile = "folder-structure.md"
	DefaultCloneScript          = "clone.ps1"
	DefaultDirectCloneScript    = "direct-clone.ps1"
	DefaultDirectCloneSSHScript = "direct-clone-ssh.ps1"
	DefaultDesktopScript        = "register-desktop.ps1"
	DefaultConfigPath    = "./data/config.json"
	DefaultOutputDir     = "./gitmap-output"
	DefaultOutputFolder  = "gitmap-output"
	DefaultBranch        = "main"
	DefaultDir           = "."
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
	CmdVersion          = "version"
	CmdVersionAlias     = "v"
	CmdHelp             = "help"
	CmdDesktopSync      = "desktop-sync"
	CmdDesktopSyncAlias = "ds"
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
	TermCloneCmd1      = "     gitmap-output/gitmap.json  (or gitmap.csv)"
	TermCloneStep2     = "  2. Clone via JSON (shorthand):"
	TermCloneCmd2      = "     gitmap clone json --target-dir ./projects"
	TermCloneCmd2Alt   = "     gitmap c json               # alias"
	TermCloneStep3     = "  3. Clone via CSV (shorthand):"
	TermCloneCmd3      = "     gitmap clone csv --target-dir ./projects"
	TermCloneCmd3Alt   = "     gitmap c csv                # alias"
	TermCloneStep3b    = "  4. Or specify a file path directly:"
	TermCloneCmd3b     = "     gitmap clone ./gitmap-output/gitmap.json --target-dir ./projects"
	TermCloneStep4     = "  5. Or run the PowerShell script directly:"
	TermCloneCmd4HTTPS = "     .\\direct-clone.ps1       # HTTPS clone commands"
	TermCloneCmd4SSH   = "     .\\direct-clone-ssh.ps1   # SSH clone commands"
	TermCloneStep5     = "  6. Full clone script with progress & error handling:"
	TermCloneCmd5      = "     .\\clone.ps1 -TargetDir .\\projects"
	TermCloneStep6     = "  7. Sync repos to GitHub Desktop:"
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
	MsgDesktopSyncStart   = "\n  Syncing repos to GitHub Desktop from %s...\n"
	MsgDesktopSyncSkipped = "  ⊘ Skipped (already exists): %s\n"
	MsgDesktopSyncAdded   = "  ✓ Added to GitHub Desktop: %s\n"
	MsgDesktopSyncFailed  = "  ✗ Failed: %s — %v\n"
	MsgDesktopSyncDone    = "\n  GitHub Desktop sync: %d added, %d skipped, %d failed\n"
	MsgNoOutputDir        = "Error: gitmap-output/ not found in current directory.\nRun 'gitmap scan' first to generate output files."
	MsgNoJSONFile         = "Error: %s not found.\nRun 'gitmap scan' first to generate the JSON output."
	MsgFailedClones     = "\nFailed clones:"
	MsgFailedEntry      = "  - %s (%s): %s\n"
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
	ErrNoRepoPath     = "Error: repo path not embedded. Binary was not built with run.ps1."
	ErrUpdateFailed   = "Update error: %v\n"
)

// CLI help text.
const (
	HelpUsage      = "Usage: gitmap <command> [flags]"
	HelpCommands   = "Commands:"
	HelpScan       = "  scan (s) [dir]      Scan directory for Git repos"
	HelpClone      = "  clone (c) <source|json|csv|text>  Re-clone from file (shorthands auto-resolve)"
	HelpUpdate     = "  update              Self-update from source repo"
	HelpVersion    = "  version (v)         Show version number"
	HelpDesktopSync = "  desktop-sync (ds)   Sync repos to GitHub Desktop from output"
	HelpHelp       = "  help                Show this help message"
	HelpScanFlags  = "Scan flags:"
	HelpConfig     = "  --config <path>     Config file (default: ./data/config.json)"
	HelpMode       = "  --mode ssh|https    Clone URL style (default: https)"
	HelpOutput     = "  --output csv|json|terminal  Output format (default: terminal)"
	HelpOutputPath = "  --output-path <dir> Output directory (default: ./gitmap-output)"
	HelpOutFile        = "  --out-file <path>   Exact output file path"
	HelpGitHubDesktop  = "  --github-desktop    Add repos to GitHub Desktop"
	HelpOpen           = "  --open              Open output folder after scan"
	HelpCloneFlags    = "Clone flags:"
	HelpTargetDir     = "  --target-dir <dir>  Base directory for clones (default: .)"
	HelpSafePull      = "  --safe-pull         Pull existing repos with retry + unlock diagnostics (auto-enabled)"
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
)

// Clone and Desktop scripts are now generated from Go templates
// embedded in formatter/templates/. See clone.ps1.tmpl and desktop.ps1.tmpl.

// Directory permissions.
const DirPermission = 0o755

// Safe-pull defaults.
const (
	SafePullRetryAttempts   = 4
	SafePullRetryDelayMS    = 600
	WindowsPathWarnThreshold = 240
)
