export interface CommandSeeAlso {
  name: string;
  description: string;
}

export interface CommandDef {
  name: string;
  alias?: string;
  description: string;
  usage?: string;
  flags?: { flag: string; description: string }[];
  examples?: { command: string; description?: string }[];
  category: string;
  seeAlso?: CommandSeeAlso[];
}

export const categories = [
  { key: "scanning", label: "Scanning & Cloning", description: "Discover, clone, and re-scan repositories" },
  { key: "monitoring", label: "Monitoring & Status", description: "Track repo state in real time" },
  { key: "release", label: "Release & Versioning", description: "Tag, branch, and publish releases" },
  { key: "navigation", label: "Navigation & Organization", description: "Move between repos and manage groups" },
  { key: "history", label: "History & Stats", description: "Audit trail and usage analytics" },
  { key: "detection", label: "Project Detection", description: "Query detected project types" },
  { key: "data", label: "Data & Profiles", description: "Import, export, and manage profiles" },
  { key: "utilities", label: "Utilities", description: "Setup, diagnostics, and maintenance" },
];

export const commands: CommandDef[] = [
  // --- Scanning & Cloning ---
  {
    category: "scanning",
    name: "scan", alias: "s", description: "Scan directory for Git repos",
    usage: "gitmap scan [dir] [--output csv|json|terminal] [--mode ssh|https]",
    flags: [
      { flag: "--config <path>", description: "Config file path" },
      { flag: "--mode ssh|https", description: "Clone URL style (default: https)" },
      { flag: "--output csv|json|terminal", description: "Output format (default: terminal)" },
      { flag: "--output-path <dir>", description: "Output directory" },
      { flag: "--github-desktop", description: "Add repos to GitHub Desktop" },
      { flag: "--open", description: "Open output folder after scan" },
      { flag: "--quiet", description: "Suppress clone help section" },
    ],
    examples: [
      { command: "gitmap scan ~/projects", description: "Scan a directory" },
      { command: "gitmap s --output json --mode ssh", description: "JSON output with SSH URLs" },
    ],
    seeAlso: [
      { name: "rescan", description: "Re-scan using cached parameters" },
      { name: "clone", description: "Clone repos from scan output" },
      { name: "status", description: "View repo statuses after scanning" },
      { name: "desktop-sync", description: "Sync scanned repos to GitHub Desktop" },
      { name: "export", description: "Export scanned data" },
    ],
  },
  {
    category: "scanning",
    name: "clone", alias: "c", description: "Re-clone repos from structured file",
    usage: "gitmap clone <source|json|csv|text> [--target-dir <dir>] [--safe-pull]",
    flags: [
      { flag: "--target-dir <dir>", description: "Base directory for clones" },
      { flag: "--safe-pull", description: "Pull existing repos with retry + diagnostics" },
      { flag: "--verbose", description: "Write detailed debug log" },
    ],
    examples: [
      { command: "gitmap clone json --target-dir ./projects" },
      { command: "gitmap c csv" },
    ],
    seeAlso: [
      { name: "scan", description: "Scan directories to create clone source" },
      { name: "pull", description: "Pull updates for existing repos" },
      { name: "desktop-sync", description: "Sync cloned repos to GitHub Desktop" },
    ],
  },
  {
    category: "scanning",
    name: "pull", alias: "p", description: "Pull a specific repo by name",
    usage: "gitmap pull <repo-name> [--group <name>] [--all] [--verbose]",
    flags: [
      { flag: "--group <name>", description: "Pull all repos in a group" },
      { flag: "--all", description: "Pull all tracked repos" },
      { flag: "--verbose", description: "Enable verbose logging" },
    ],
    seeAlso: [
      { name: "scan", description: "Scan directories to populate the database" },
      { name: "clone", description: "Clone repos from structured file" },
      { name: "status", description: "View repo statuses" },
      { name: "group", description: "Manage repo groups for targeted pulls" },
    ],
  },
  {
    category: "scanning",
    name: "rescan", alias: "rs", description: "Re-scan previously scanned directories",
    usage: "gitmap rescan [--output csv|json|terminal]",
    seeAlso: [
      { name: "scan", description: "Initial directory scan" },
      { name: "status", description: "View repo statuses" },
      { name: "clone", description: "Clone from scan output" },
    ],
  },
  {
    category: "scanning",
    name: "desktop-sync", alias: "ds", description: "Sync tracked repos with GitHub Desktop",
    usage: "gitmap desktop-sync",
    seeAlso: [
      { name: "scan", description: "Scan directories first" },
      { name: "clone", description: "Clone repos from scan output" },
      { name: "list", description: "List tracked repos" },
    ],
  },

  // --- Monitoring & Status ---
  {
    category: "monitoring",
    name: "status", alias: "st", description: "Show repo status dashboard",
    usage: "gitmap status [--group <name>] [--all]",
    seeAlso: [
      { name: "watch", description: "Live-refresh status dashboard" },
      { name: "scan", description: "Scan directories to populate data" },
      { name: "exec", description: "Run git commands across repos" },
      { name: "group", description: "Filter status by group" },
    ],
  },
  {
    category: "monitoring",
    name: "watch", alias: "w", description: "Live-refresh repo status dashboard",
    usage: "gitmap watch [--interval <seconds>] [--group <name>] [--no-fetch] [--json]",
    flags: [
      { flag: "--interval <seconds>", description: "Refresh interval (default: 30, min: 5)" },
      { flag: "--group <name>", description: "Monitor only repos in a group" },
      { flag: "--no-fetch", description: "Skip git fetch" },
      { flag: "--json", description: "Output single snapshot as JSON" },
    ],
    seeAlso: [
      { name: "status", description: "One-time status snapshot" },
      { name: "exec", description: "Run git commands across repos" },
      { name: "group", description: "Filter by group" },
    ],
  },
  {
    category: "monitoring",
    name: "exec", alias: "x", description: "Run git command across all repos",
    usage: "gitmap exec <git-args...>",
    examples: [{ command: "gitmap exec fetch --prune" }],
    seeAlso: [
      { name: "scan", description: "Scan directories to populate the database" },
      { name: "pull", description: "Pull repos (built-in alternative)" },
      { name: "status", description: "View repo statuses" },
    ],
  },
  {
    category: "monitoring",
    name: "latest-branch", alias: "lb", description: "Find most recently updated remote branch",
    usage: "gitmap latest-branch [--top N] [--format json|csv|terminal]",
    seeAlso: [
      { name: "status", description: "View repo statuses" },
      { name: "watch", description: "Live-refresh dashboard" },
      { name: "release-branch", description: "Create a release branch" },
    ],
  },

  // --- Release & Versioning ---
  {
    category: "release",
    name: "release", alias: "r", description: "Create release branch, tag, and push",
    usage: "gitmap release [version] [--bump major|minor|patch] [--draft] [--dry-run]",
    flags: [
      { flag: "--assets <path>", description: "Attach files to release" },
      { flag: "--commit <sha>", description: "Release from specific commit" },
      { flag: "--branch <name>", description: "Release from branch" },
      { flag: "--bump major|minor|patch", description: "Auto-increment version" },
      { flag: "--draft", description: "Create unpublished draft" },
      { flag: "--dry-run", description: "Preview without executing" },
    ],
    seeAlso: [
      { name: "release-branch", description: "Create branch without tagging" },
      { name: "release-pending", description: "Show unreleased commits" },
      { name: "changelog", description: "View release notes" },
      { name: "list-versions", description: "List available tags" },
    ],
  },
  {
    category: "release",
    name: "release-branch", alias: "rb", description: "Create a release branch without tagging",
    usage: "gitmap release-branch [version] [--bump major|minor|patch]",
    seeAlso: [
      { name: "release", description: "Full release with tag and push" },
      { name: "release-pending", description: "Show unreleased commits" },
      { name: "latest-branch", description: "Find most recent branch" },
    ],
  },
  {
    category: "release",
    name: "release-pending", alias: "rp", description: "Show unreleased commits since last tag",
    usage: "gitmap release-pending [--json]",
    seeAlso: [
      { name: "release", description: "Create a release from pending commits" },
      { name: "release-branch", description: "Create a release branch" },
      { name: "changelog", description: "View release notes" },
    ],
  },
  {
    category: "release",
    name: "changelog", alias: "cl", description: "Show release notes",
    usage: "gitmap changelog [version] [--open] [--source]",
    seeAlso: [
      { name: "release", description: "Create a release" },
      { name: "list-versions", description: "List available tags" },
      { name: "list-releases", description: "List stored release metadata" },
    ],
  },
  {
    category: "release",
    name: "list-versions", alias: "lv", description: "List all available Git release tags",
    usage: "gitmap list-versions [--json] [--limit N]",
    flags: [
      { flag: "--json", description: "Output as structured JSON" },
      { flag: "--limit N", description: "Show only the top N versions (0 = all)" },
    ],
    seeAlso: [
      { name: "list-releases", description: "List stored release metadata" },
      { name: "changelog", description: "View release notes" },
      { name: "release", description: "Create a release" },
      { name: "revert", description: "Revert to a specific version" },
    ],
  },
  {
    category: "release",
    name: "list-releases", alias: "lr", description: "List release metadata from the database",
    usage: "gitmap list-releases [--json] [--source manual|scan]",
    flags: [
      { flag: "--json", description: "Output as structured JSON" },
      { flag: "--source <type>", description: "Filter by release source" },
    ],
    seeAlso: [
      { name: "list-versions", description: "List Git tags" },
      { name: "release", description: "Create a release" },
      { name: "changelog", description: "View release notes" },
    ],
  },
  {
    category: "release",
    name: "revert", alias: undefined, description: "Revert to a specific release version",
    usage: "gitmap revert <version>",
    seeAlso: [
      { name: "list-versions", description: "List available versions" },
      { name: "release", description: "Create a new release" },
      { name: "changelog", description: "View release notes" },
    ],
  },

  // --- Navigation & Organization ---
  {
    category: "navigation",
    name: "cd", alias: "go", description: "Navigate to a tracked repo directory",
    usage: "gitmap cd <repo-name|repos> [--group <name>] [--pick]",
    examples: [
      { command: "gitmap cd myrepo", description: "Navigate to repo" },
      { command: "gitmap cd repos", description: "Interactive repo picker" },
      { command: "gitmap cd repos --group work", description: "Pick from group" },
    ],
    seeAlso: [
      { name: "list", description: "List all tracked repos with slugs" },
      { name: "scan", description: "Scan directories to populate database" },
      { name: "group", description: "Manage repo groups" },
      { name: "bookmark", description: "Save commands for re-execution" },
    ],
  },
  {
    category: "navigation",
    name: "list", alias: "ls", description: "Show all tracked repos with slugs",
    usage: "gitmap list [--group <name>] [--verbose]",
    seeAlso: [
      { name: "cd", description: "Navigate to a tracked repo" },
      { name: "scan", description: "Scan directories to populate data" },
      { name: "group", description: "Manage repo groups" },
      { name: "status", description: "View repo statuses" },
    ],
  },
  {
    category: "navigation",
    name: "group", alias: "g", description: "Manage repo groups",
    usage: "gitmap group <create|add|remove|list|show|delete> [args]",
    seeAlso: [
      { name: "list", description: "List all tracked repos" },
      { name: "cd", description: "Navigate to repos" },
      { name: "pull", description: "Pull repos by group" },
      { name: "status", description: "Filter status by group" },
    ],
  },
  {
    category: "navigation",
    name: "diff-profiles", alias: "dp", description: "Compare repos across two profiles",
    usage: "gitmap diff-profiles <profileA> <profileB> [--all] [--json]",
    seeAlso: [
      { name: "profile", description: "Manage database profiles" },
      { name: "list", description: "List tracked repos" },
      { name: "export", description: "Export database" },
    ],
  },

  // --- History & Stats ---
  {
    category: "history",
    name: "history", alias: "hi", description: "Show CLI command execution history",
    usage: "gitmap history [--limit N] [--json]",
    flags: [
      { flag: "--limit N", description: "Number of entries to show" },
      { flag: "--json", description: "Output as structured JSON" },
    ],
    seeAlso: [
      { name: "history-reset", description: "Clear command history" },
      { name: "stats", description: "View aggregated usage metrics" },
      { name: "bookmark", description: "Save commands for re-execution" },
    ],
  },
  {
    category: "history",
    name: "history-reset", alias: "hr", description: "Clear command execution history",
    usage: "gitmap history-reset [--confirm]",
    seeAlso: [
      { name: "history", description: "View command history" },
      { name: "db-reset", description: "Reset entire database" },
    ],
  },
  {
    category: "history",
    name: "stats", alias: "ss", description: "Show aggregated usage and performance metrics",
    usage: "gitmap stats [--json]",
    seeAlso: [
      { name: "history", description: "View command history" },
      { name: "scan", description: "Scan directories" },
      { name: "status", description: "View repo statuses" },
    ],
  },
  {
    category: "history",
    name: "amend", alias: "am", description: "Rewrite commit author info",
    usage: "gitmap amend [commit-hash] --name <name> --email <email>",
    seeAlso: [
      { name: "amend-list", description: "List previous amendments" },
      { name: "history", description: "View command history" },
    ],
  },
  {
    category: "history",
    name: "amend-list", alias: "al", description: "List previous author amendments",
    usage: "gitmap amend-list [--json]",
    seeAlso: [
      { name: "amend", description: "Rewrite commit author info" },
      { name: "history", description: "View command history" },
    ],
  },

  // --- Project Detection ---
  {
    category: "detection",
    name: "go-repos", alias: "gr", description: "List detected Go projects",
    usage: "gitmap go-repos [--json]",
    seeAlso: [
      { name: "node-repos", description: "List Node.js projects" },
      { name: "react-repos", description: "List React projects" },
      { name: "scan", description: "Scan directories first" },
      { name: "gomod", description: "Rename Go module paths" },
    ],
  },
  {
    category: "detection",
    name: "node-repos", alias: "nr", description: "List detected Node.js projects",
    usage: "gitmap node-repos [--json]",
    seeAlso: [
      { name: "react-repos", description: "List React projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },
  {
    category: "detection",
    name: "react-repos", alias: "rr", description: "List detected React projects",
    usage: "gitmap react-repos [--json]",
    seeAlso: [
      { name: "node-repos", description: "List Node.js projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },
  {
    category: "detection",
    name: "cpp-repos", alias: "cr", description: "List detected C++ projects",
    usage: "gitmap cpp-repos [--json]",
    seeAlso: [
      { name: "csharp-repos", description: "List C# projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },
  {
    category: "detection",
    name: "csharp-repos", alias: "csr", description: "List detected C# projects",
    usage: "gitmap csharp-repos [--json]",
    seeAlso: [
      { name: "cpp-repos", description: "List C++ projects" },
      { name: "go-repos", description: "List Go projects" },
      { name: "scan", description: "Scan directories first" },
    ],
  },

  // --- Data & Profiles ---
  {
    category: "data",
    name: "export", alias: "ex", description: "Export database to file",
    usage: "gitmap export [--json]",
    seeAlso: [
      { name: "import", description: "Import repos from file" },
      { name: "profile", description: "Manage database profiles" },
      { name: "scan", description: "Scan directories to populate data" },
    ],
  },
  {
    category: "data",
    name: "import", alias: "im", description: "Import repos from file",
    usage: "gitmap import <file>",
    seeAlso: [
      { name: "export", description: "Export database to file" },
      { name: "scan", description: "Scan directories" },
      { name: "profile", description: "Manage database profiles" },
    ],
  },
  {
    category: "data",
    name: "profile", alias: "pf", description: "Manage database profiles",
    usage: "gitmap profile <create|list|switch|delete|show> [name]",
    seeAlso: [
      { name: "diff-profiles", description: "Compare repos across profiles" },
      { name: "export", description: "Export database" },
      { name: "import", description: "Import repos" },
    ],
  },
  {
    category: "data",
    name: "bookmark", alias: "bk", description: "Save and run bookmarked commands",
    usage: "gitmap bookmark <save|list|run|delete> [args]",
    seeAlso: [
      { name: "history", description: "View command execution history" },
      { name: "scan", description: "Scan directories (common bookmark target)" },
      { name: "pull", description: "Pull repos (common bookmark target)" },
    ],
  },
  {
    category: "data",
    name: "db-reset", alias: undefined, description: "Reset the local SQLite database",
    usage: "gitmap db-reset [--confirm]",
    seeAlso: [
      { name: "history-reset", description: "Clear command history only" },
      { name: "scan", description: "Re-scan after reset" },
      { name: "setup", description: "Re-run setup wizard" },
    ],
  },

  // --- Utilities ---
  {
    category: "utilities",
    name: "setup", alias: undefined, description: "Configure Git settings and install shell completions",
    usage: "gitmap setup [--config <path>] [--dry-run]",
    flags: [
      { flag: "--config <path>", description: "Path to git-setup.json config file" },
      { flag: "--dry-run", description: "Preview changes without applying" },
    ],
    seeAlso: [
      { name: "completion", description: "Generate completion scripts manually" },
      { name: "doctor", description: "Diagnose issues" },
      { name: "scan", description: "Scan directories after setup" },
      { name: "update", description: "Self-update to latest version" },
    ],
  },
  {
    category: "utilities",
    name: "doctor", alias: undefined, description: "Diagnose PATH, deploy, and version issues",
    usage: "gitmap doctor [--fix-path]",
    seeAlso: [
      { name: "setup", description: "Re-run setup wizard" },
      { name: "update", description: "Self-update to latest version" },
      { name: "version", description: "Show current version" },
    ],
  },
  {
    category: "utilities",
    name: "update", alias: undefined, description: "Self-update from source repo",
    usage: "gitmap update",
    seeAlso: [
      { name: "version", description: "Show current version" },
      { name: "doctor", description: "Diagnose issues after update" },
    ],
  },
  {
    category: "utilities",
    name: "update-cleanup", alias: undefined, description: "Remove leftover update artifacts",
    usage: "gitmap update-cleanup",
    examples: [
      {
        command: "gitmap update-cleanup",
        description: "Remove temp binaries and .old backups from previous updates",
      },
    ],
    seeAlso: [
      { name: "update", description: "Self-update to latest version" },
      { name: "revert", description: "Revert to a previous version" },
    ],
  },
  {
    category: "utilities",
    name: "version", alias: "v", description: "Show version number",
    usage: "gitmap version",
    seeAlso: [
      { name: "update", description: "Self-update to latest version" },
      { name: "doctor", description: "Diagnose version issues" },
    ],
  },
  {
    category: "utilities",
    name: "seo-write", alias: "sw", description: "Auto-commit SEO messages",
    usage: "gitmap seo-write [flags]",
    seeAlso: [
      { name: "scan", description: "Scan directories" },
      { name: "history", description: "View command history" },
    ],
  },
  {
    category: "utilities",
    name: "gomod", alias: "gm", description: "Rename Go module path across repo with branch safety",
    usage: "gitmap gomod <new-module-path> [--ext *.go,*.md] [--dry-run] [--no-merge] [--no-tidy] [--verbose]",
    flags: [
      { flag: "--ext <exts>", description: "Comma-separated extensions to filter (e.g. *.go,*.md); default: all files" },
      { flag: "--dry-run", description: "Preview changes without modifying files or branches" },
      { flag: "--no-merge", description: "Commit on feature branch but do not merge back" },
      { flag: "--no-tidy", description: "Skip go mod tidy after replacement" },
      { flag: "--verbose", description: "Print each file path as it is modified" },
    ],
    examples: [
      { command: 'gitmap gomod "x/y"', description: "Rename module path in all files" },
      { command: 'gitmap gomod "x/y" --ext "*.go,*.md"', description: "Only replace in .go and .md files" },
      { command: 'gitmap gomod "github.com/new/name" --dry-run', description: "Preview what would change" },
      { command: 'gitmap gomod "github.com/new/name" --no-merge', description: "Replace but stay on feature branch" },
    ],
    seeAlso: [
      { name: "go-repos", description: "List detected Go projects" },
      { name: "scan", description: "Scan directories" },
    ],
  },
  {
    category: "utilities",
    name: "completion", alias: "cmp", description: "Generate or install shell tab-completion scripts",
    usage: "gitmap completion <powershell|bash|zsh> [--list-repos] [--list-groups] [--list-commands]",
    flags: [
      { flag: "--list-repos", description: "Print repo slugs, one per line (for script use)" },
      { flag: "--list-groups", description: "Print group names, one per line (for script use)" },
      { flag: "--list-commands", description: "Print all command names, one per line (for script use)" },
    ],
    examples: [
      { command: "gitmap completion powershell", description: "Print PowerShell completion script" },
      { command: "gitmap completion bash", description: "Print Bash completion script" },
      { command: "gitmap completion --list-repos", description: "List repo slugs for scripting" },
    ],
    seeAlso: [
      { name: "setup", description: "Auto-installs completions during setup" },
      { name: "cd", description: "Navigate to repos using tab-completed slugs" },
      { name: "group", description: "Group names are also tab-completed" },
    ],
  },
];
