package constants

// Install CLI commands.
const (
	CmdInstall      = "install"
	CmdInstallAlias = "in"
)

// Install help text.
const HelpInstall = "  install (in) <tool> Install a developer tool by name"

// Supported tool names.
const (
	ToolVSCode         = "vscode"
	ToolNodeJS         = "node"
	ToolYarn           = "yarn"
	ToolBun            = "bun"
	ToolPnpm           = "pnpm"
	ToolPython         = "python"
	ToolGo             = "go"
	ToolGit            = "git"
	ToolGitLFS         = "git-lfs"
	ToolGHCLI          = "gh"
	ToolGitHubDesktop  = "github-desktop"
	ToolCPP            = "cpp"
	ToolPHP            = "php"
	ToolPowerShell     = "powershell"
)

// Package manager names.
const (
	PkgMgrChocolatey = "choco"
	PkgMgrWinget     = "winget"
	PkgMgrApt        = "apt"
	PkgMgrBrew       = "brew"
	PkgMgrDnf        = "dnf"
	PkgMgrPacman     = "pacman"
)

// Install flag names.
const (
	FlagInstallManager = "manager"
	FlagInstallVersion = "version"
	FlagInstallVerbose = "verbose"
	FlagInstallDryRun  = "dry-run"
	FlagInstallCheck   = "check"
	FlagInstallList    = "list"
)

// Install flag descriptions.
const (
	FlagDescInstallManager = "Force package manager (choco, winget, apt, brew)"
	FlagDescInstallVersion = "Install a specific version"
	FlagDescInstallVerbose = "Show full installer output"
	FlagDescInstallDryRun  = "Show install command without executing"
	FlagDescInstallCheck   = "Only check if tool is installed"
	FlagDescInstallList    = "List all supported tools"
)

// Chocolatey package IDs.
const (
	ChocoPkgVSCode        = "vscode"
	ChocoPkgNodeJS        = "nodejs"
	ChocoPkgYarn          = "yarn"
	ChocoPkgBun           = "bun"
	ChocoPkgPnpm          = "pnpm"
	ChocoPkgPython        = "python"
	ChocoPkgGo            = "golang"
	ChocoPkgGit           = "git"
	ChocoPkgGitLFS        = "git-lfs"
	ChocoPkgGHCLI         = "gh"
	ChocoPkgGitHubDesktop = "github-desktop"
	ChocoPkgCPP           = "mingw"
	ChocoPkgPHP           = "php"
)

// Winget package IDs.
const (
	WingetPkgVSCode    = "Microsoft.VisualStudioCode"
	WingetPkgPowerShell = "Microsoft.PowerShell"
)

// Install terminal messages.
const (
	MsgInstallChecking   = "Checking if %s is installed...\n"
	MsgInstallFound      = "%s is already installed (version: %s)\n"
	MsgInstallInstalling = "Installing %s...\n"
	MsgInstallSuccess    = "%s installed successfully.\n"
	MsgInstallDryCmd     = "[dry-run] Would run: %s\n"
	MsgInstallVerifying  = "Verifying %s installation...\n"
	MsgInstallListHeader = "Supported tools:\n"
	MsgInstallListRow    = "  %-20s %s\n"
)

// Install error messages.
const (
	ErrInstallToolRequired    = "Tool name is required. Use --list to see available tools."
	ErrInstallUnknownTool     = "Unknown tool: %s. Use --list to see available tools.\n"
	ErrInstallNoPkgMgr       = "No package manager found. Install Chocolatey or Winget first."
	ErrInstallFailed          = "Installation failed for %s: %v\n"
	ErrInstallVerifyFailed    = "Post-install verification failed for %s.\n"
	ErrInstallAdminRequired   = "%s requires administrator privileges to install.\n"
	ErrInstallNetworkRequired = "Network connection required for installation."
)

// Tool display names for --list output.
var InstallToolDescriptions = map[string]string{
	ToolVSCode:        "Visual Studio Code editor",
	ToolNodeJS:        "Node.js JavaScript runtime",
	ToolYarn:          "Yarn package manager",
	ToolBun:           "Bun JavaScript runtime",
	ToolPnpm:          "pnpm package manager",
	ToolPython:        "Python programming language",
	ToolGo:            "Go programming language",
	ToolGit:           "Git version control",
	ToolGitLFS:        "Git Large File Storage",
	ToolGHCLI:         "GitHub CLI",
	ToolGitHubDesktop: "GitHub Desktop application",
	ToolCPP:           "C++ compiler (MinGW/g++)",
	ToolPHP:           "PHP programming language",
	ToolPowerShell:    "PowerShell shell",
}