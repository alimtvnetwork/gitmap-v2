package constants

// Install CLI commands.
const (
	CmdInstall        = "install"
	CmdInstallAlias   = "in"
	CmdUninstall      = "uninstall"
	CmdUninstallAlias = "un"
)

// Install help text.
const (
	HelpInstall   = "  install (in) <tool> Install a developer tool by name"
	HelpUninstall = "  uninstall (un) <tool> Remove a previously installed tool"
)

// Supported tool names — Core.
const (
	ToolVSCode        = "vscode"
	ToolNodeJS        = "node"
	ToolYarn          = "yarn"
	ToolBun           = "bun"
	ToolPnpm          = "pnpm"
	ToolPython        = "python"
	ToolGo            = "go"
	ToolGit           = "git"
	ToolGitLFS        = "git-lfs"
	ToolGHCLI         = "gh"
	ToolGitHubDesktop = "github-desktop"
	ToolCPP           = "cpp"
	ToolPHP           = "php"
	ToolPowerShell    = "powershell"
	ToolChocolatey    = "chocolatey"
	ToolWinget        = "winget"
	ToolNpp           = "npp"
	ToolNppSettings   = "npp-settings"
)

// Supported tool names — Databases.
const (
	ToolMySQL         = "mysql"
	ToolMariaDB       = "mariadb"
	ToolPostgreSQL    = "postgresql"
	ToolSQLite        = "sqlite"
	ToolMongoDB       = "mongodb"
	ToolCouchDB       = "couchdb"
	ToolRedis         = "redis"
	ToolCassandra     = "cassandra"
	ToolNeo4j         = "neo4j"
	ToolElasticsearch = "elasticsearch"
	ToolDuckDB        = "duckdb"
)

// Package manager names.
const (
	PkgMgrChocolatey = "choco"
	PkgMgrWinget     = "winget"
	PkgMgrApt        = "apt"
	PkgMgrBrew       = "brew"
	PkgMgrSnap       = "snap"
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
	FlagInstallStatus  = "status"
	FlagInstallUpgrade = "upgrade"
)

// Install flag descriptions.
const (
	FlagDescInstallManager = "Force package manager (choco, winget, apt, brew, snap)"
	FlagDescInstallVersion = "Install a specific version"
	FlagDescInstallVerbose = "Show full installer output"
	FlagDescInstallDryRun  = "Show install command without executing"
	FlagDescInstallCheck   = "Only check if tool is installed"
	FlagDescInstallList    = "List all supported tools"
	FlagDescInstallStatus  = "Show installed tools from database"
	FlagDescInstallUpgrade = "Upgrade an already-installed tool"
)

// Uninstall flag names.
const (
	FlagUninstallDryRun = "dry-run"
	FlagUninstallForce  = "force"
	FlagUninstallPurge  = "purge"
)

// Uninstall flag descriptions.
const (
	FlagDescUninstallDryRun = "Show uninstall command without executing"
	FlagDescUninstallForce  = "Skip confirmation prompt"
	FlagDescUninstallPurge  = "Remove config files too"
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
	ChocoPkgMySQL         = "mysql"
	ChocoPkgMariaDB       = "mariadb"
	ChocoPkgPostgreSQL    = "postgresql"
	ChocoPkgSQLite        = "sqlite"
	ChocoPkgMongoDB       = "mongodb"
	ChocoPkgCouchDB       = "couchdb"
	ChocoPkgRedis         = "redis-64"
	ChocoPkgNeo4j         = "neo4j-community"
	ChocoPkgElasticsearch = "elasticsearch"
	ChocoPkgDuckDB        = "duckdb"
	ChocoPkgNpp           = "notepadplusplus"
)

// Winget package IDs.
const (
	WingetPkgVSCode     = "Microsoft.VisualStudioCode"
	WingetPkgPowerShell = "Microsoft.PowerShell"
)

// Apt package IDs.
const (
	AptPkgNodeJS        = "nodejs"
	AptPkgPython        = "python3"
	AptPkgGo            = "golang"
	AptPkgGit           = "git"
	AptPkgGitLFS        = "git-lfs"
	AptPkgCPP           = "g++"
	AptPkgPHP           = "php"
	AptPkgMySQL         = "mysql-server"
	AptPkgMariaDB       = "mariadb-server"
	AptPkgPostgreSQL    = "postgresql"
	AptPkgSQLite        = "sqlite3"
	AptPkgMongoDB       = "mongod"
	AptPkgCouchDB       = "couchdb"
	AptPkgRedis         = "redis-server"
	AptPkgCassandra     = "cassandra"
	AptPkgElasticsearch = "elasticsearch"
)

// Brew package IDs.
const (
	BrewPkgNodeJS        = "node"
	BrewPkgPython        = "python"
	BrewPkgGo            = "go"
	BrewPkgGit           = "git"
	BrewPkgGitLFS        = "git-lfs"
	BrewPkgGHCLI         = "gh"
	BrewPkgCPP           = "gcc"
	BrewPkgPHP           = "php"
	BrewPkgMySQL         = "mysql"
	BrewPkgMariaDB       = "mariadb"
	BrewPkgPostgreSQL    = "postgresql"
	BrewPkgSQLite        = "sqlite"
	BrewPkgMongoDB       = "mongodb-community"
	BrewPkgCouchDB       = "couchdb"
	BrewPkgRedis         = "redis"
	BrewPkgNeo4j         = "neo4j"
	BrewPkgElasticsearch = "elasticsearch"
	BrewPkgDuckDB        = "duckdb"
)

// Snap package IDs.
const (
	SnapPkgCouchDB = "couchdb"
	SnapPkgRedis   = "redis"
)

// Install terminal messages.
const (
	MsgInstallChecking   = "Checking if %s is installed...\n"
	MsgInstallFound      = "%s is already installed (version: %s)\n"
	MsgInstallNotFound   = "%s is not installed.\n"
	MsgInstallInstalling = "Installing %s...\n"
	MsgInstallSuccess    = "%s installed successfully.\n"
	MsgInstallDryCmd     = "[dry-run] Would run: %s\n"
	MsgInstallVerifying  = "Verifying %s installation...\n"
	MsgInstallListHeader = "Supported tools:\n\n"
	MsgInstallListRow    = "  %-20s %s\n"
	MsgInstallRecorded   = "Recorded %s v%s in database.\n"
	MsgInstallStatusHdr  = "Installed tools:\n\n"
	MsgInstallStatusRow  = "  %-20s %-12s %-8s %s\n"
	MsgInstallExeVerify   = "Verifying %s binary at: %s\n"
	MsgInstallExeFound    = "Binary confirmed: %s\n"
	MsgInstallNppSettings = "Syncing Notepad++ settings...\n"
	MsgInstallNppSkipBin  = "Skipping Notepad++ installation (settings-only mode)\n"
)

// Install error messages.
const (
	ErrInstallToolRequired    = "Tool name is required. Use --list to see available tools."
	ErrInstallUnknownTool     = "Unknown tool: %s. Use --list to see available tools.\n"
	ErrInstallNoPkgMgr        = "No package manager found. Install Chocolatey or Winget first."
	ErrInstallFailed          = "Installation failed for %s: %v\n"
	ErrInstallVerifyFailed    = "Post-install verification failed for %s.\n"
	ErrInstallAdminRequired   = "%s requires administrator privileges to install.\n"
	ErrInstallNetworkRequired = "Network connection required for installation."
	ErrInstallExeNotFound = "Post-install verification failed: binary not found at %s\n"
)

// Uninstall messages.
const (
	MsgUninstallRemoving = "Removing %s...\n"
	MsgUninstallSuccess  = "%s uninstalled successfully.\n"
	MsgUninstallDryCmd   = "[dry-run] Would run: %s\n"
	MsgUninstallConfirm  = "Uninstall %s? (y/N): "
	ErrUninstallFailed   = "Uninstall failed for %s: %v\n"
	ErrUninstallNotFound = "%s is not tracked in the database. Use --force to try anyway.\n"
	ErrUninstallDBRemove = "Warning: could not remove %s from database: %v\n"
)

// Tool categories.
const (
	ToolCategoryCore     = "Core Tools"
	ToolCategoryDatabase = "Databases"
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
	ToolChocolatey:    "Chocolatey package manager",
	ToolWinget:        "Winget package manager",
	ToolMySQL:         "MySQL relational database",
	ToolMariaDB:       "MariaDB (MySQL-compatible fork)",
	ToolPostgreSQL:    "PostgreSQL relational database",
	ToolSQLite:        "SQLite embedded database",
	ToolMongoDB:       "MongoDB document database",
	ToolCouchDB:       "CouchDB document database (REST API)",
	ToolRedis:         "Redis in-memory key-value store",
	ToolCassandra:     "Apache Cassandra wide-column NoSQL",
	ToolNeo4j:         "Neo4j graph database",
	ToolElasticsearch: "Elasticsearch search and analytics",
	ToolDuckDB:        "DuckDB analytical columnar database",
	ToolNpp:           "Notepad++ text editor",
	ToolNppSettings:   "Notepad++ settings sync (settings only)",
}

// InstallToolCategories groups tools by category for display.
var InstallToolCategories = map[string][]string{
	ToolCategoryCore: {
		ToolVSCode, ToolNodeJS, ToolYarn, ToolBun, ToolPnpm,
		ToolPython, ToolGo, ToolGit, ToolGitLFS, ToolGHCLI,
		ToolGitHubDesktop, ToolCPP, ToolPHP, ToolPowerShell,
		ToolChocolatey, ToolWinget, ToolNpp, ToolNppSettings,
	},
	ToolCategoryDatabase: {
		ToolMySQL, ToolMariaDB, ToolPostgreSQL, ToolSQLite,
		ToolMongoDB, ToolCouchDB, ToolRedis, ToolCassandra,
		ToolNeo4j, ToolElasticsearch, ToolDuckDB,
	},
}
