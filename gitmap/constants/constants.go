// Package constants defines all shared constant values for gitmap.
// No magic strings — all literals used for comparison, defaults,
// formats, and file extensions live here.
package constants

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
	ExtGit  = ".git"
)

// Default file names.
const (
	DefaultCSVFile       = "gitmap.csv"
	DefaultJSONFile      = "gitmap.json"
	DefaultStructureFile = "folder-structure.md"
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
	GitBranchFlag   = "-b"
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
	CmdScan  = "scan"
	CmdClone = "clone"
	CmdHelp  = "help"
)

// Terminal table headers.
const (
	TerminalHeader    = "REPO\tBRANCH\tPATH\tCLONE INSTRUCTION"
	TerminalSeparator = "----\t------\t----\t-----------------"
	TerminalRowFmt    = "%s\t%s\t%s\t%s\n"
)

// JSON formatting.
const JSONIndent = "  "

// CLI messages.
const (
	MsgFoundRepos       = "Found %d repositories.\n"
	MsgCSVWritten       = "CSV written to %s\n"
	MsgJSONWritten      = "JSON written to %s\n"
	MsgStructureWritten = "Folder structure written to %s\n"
	MsgCloneComplete    = "\nClone complete: %d succeeded, %d failed\n"
	MsgFailedClones     = "\nFailed clones:"
	MsgFailedEntry      = "  - %s (%s): %s\n"
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

// CLI error messages.
const (
	ErrSourceRequired = "Error: source file is required"
	ErrCloneUsage     = "Usage: gitmap clone <source> [--target-dir <dir>]"
	ErrConfigLoad     = "Error loading config: %v\n"
	ErrScanFailed     = "Scan error: %v\n"
	ErrCloneFailed    = "Clone error: %v\n"
	ErrOutputFailed   = "Output error: %v\n"
	ErrCreateDir      = "Cannot create directory: %v\n"
	ErrCreateFile     = "Cannot create file: %v\n"
)

// CLI help text.
const (
	HelpUsage      = "Usage: gitmap <command> [flags]"
	HelpCommands   = "Commands:"
	HelpScan       = "  scan [dir]          Scan directory for Git repos"
	HelpClone      = "  clone <source>      Re-clone from CSV/JSON/text file"
	HelpHelp       = "  help                Show this help message"
	HelpScanFlags  = "Scan flags:"
	HelpConfig     = "  --config <path>     Config file (default: ./data/config.json)"
	HelpMode       = "  --mode ssh|https    Clone URL style (default: https)"
	HelpOutput     = "  --output csv|json|terminal  Output format (default: terminal)"
	HelpOutputPath = "  --output-path <dir> Output directory (default: ./gitmap-output)"
	HelpOutFile    = "  --out-file <path>   Exact output file path"
	HelpCloneFlags = "Clone flags:"
	HelpTargetDir  = "  --target-dir <dir>  Base directory for clones (default: .)"
)

// Flag descriptions.
const (
	FlagDescConfig     = "Path to config file"
	FlagDescMode       = "Clone URL style: https or ssh"
	FlagDescOutput     = "Output format: terminal, csv, json"
	FlagDescOutFile    = "Exact output file path"
	FlagDescOutputPath = "Output directory for CSV/JSON"
	FlagDescTargetDir  = "Base directory for cloned repos"
)

// Directory permissions.
const DirPermission = 0o755
