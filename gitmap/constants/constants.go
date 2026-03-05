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
	DefaultCSVFile    = "gitmap.csv"
	DefaultJSONFile   = "gitmap.json"
	DefaultConfigPath = "./data/config.json"
	DefaultOutputDir  = "./gitmap-output"
	DefaultBranch     = "main"
	DefaultDir        = "."
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

// Directory permissions.
const DirPermission = 0o755
