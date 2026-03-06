package constants

// ANSI color codes.
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[97m"
	ColorDim    = "\033[90m"
)

// Status banner box-drawing.
const (
	StatusBannerTop    = "  ╔══════════════════════════════════════╗"
	StatusBannerTitle  = "  ║         gitmap status                ║"
	StatusBannerBottom = "  ╚══════════════════════════════════════╝"
	StatusRepoCountFmt = "  %d repos from gitmap-output/gitmap.json"
	StatusMissingFmt   = "  %-22s ⊘ not found"
)

// Status indicator strings.
const (
	StatusIconClean = "✓ clean"
	StatusIconDirty = "● dirty"
	StatusDash      = "—"
	StatusSyncDash  = "  —"
	StatusStashFmt  = "📦 %d"
	StatusSyncUpFmt   = "↑%d"
	StatusSyncDownFmt = "↓%d"
	StatusSyncBothFmt = "↑%d ↓%d"
	StatusStagedFmt   = "+%d"
	StatusModifiedFmt = "~%d"
	StatusUntrackedFmt = "?%d"
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
	TermTableRule    = "──────────────────────────────────────────────────────────────────────"
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

// CSV headers.
var ScanCSVHeaders = []string{
	"repoName", "httpsUrl", "sshUrl", "branch",
	"relativePath", "absolutePath", "cloneInstruction", "notes",
}

var LatestBranchCSVHeaders = []string{
	"branch", "remote", "sha", "commitDate", "subject", "ref",
}

// Status terminal table header columns.
var StatusTableColumns = []string{
	"REPO", "BRANCH", "STATE", "SYNC", "STASH", "FILES",
}

// Latest-branch terminal table header columns.
var LatestBranchTableColumns = []string{
	"DATE", "BRANCH", "SHA", "SUBJECT",
}
