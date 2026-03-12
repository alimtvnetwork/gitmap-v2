export interface CommandDef {
  name: string;
  alias?: string;
  description: string;
  usage?: string;
  flags?: { flag: string; description: string }[];
  examples?: { command: string; description?: string }[];
  category: string;
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
  },
  {
    category: "scanning",
    name: "rescan", alias: "rs", description: "Re-scan previously scanned directories",
    usage: "gitmap rescan [--output csv|json|terminal]",
  },
  {
    category: "scanning",
    name: "desktop-sync", alias: "ds", description: "Sync tracked repos with GitHub Desktop",
    usage: "gitmap desktop-sync",
  },

  // --- Monitoring & Status ---
  {
    category: "monitoring",
    name: "status", alias: "st", description: "Show repo status dashboard",
    usage: "gitmap status [--group <name>] [--all]",
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
  },
  {
    category: "monitoring",
    name: "exec", alias: "x", description: "Run git command across all repos",
    usage: "gitmap exec <git-args...>",
    examples: [{ command: "gitmap exec fetch --prune" }],
  },
  {
    category: "monitoring",
    name: "latest-branch", alias: "lb", description: "Find most recently updated remote branch",
    usage: "gitmap latest-branch [--top N] [--format json|csv|terminal]",
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
  },
  {
    category: "release",
    name: "release-branch", alias: "rb", description: "Create a release branch without tagging",
    usage: "gitmap release-branch [version] [--bump major|minor|patch]",
  },
  {
    category: "release",
    name: "release-pending", alias: "rp", description: "Show unreleased commits since last tag",
    usage: "gitmap release-pending [--json]",
  },
  {
    category: "release",
    name: "changelog", alias: "cl", description: "Show release notes",
    usage: "gitmap changelog [version] [--open] [--source]",
  },
  {
    category: "release",
    name: "list-versions", alias: "lv", description: "List all available Git release tags",
    usage: "gitmap list-versions [--json] [--limit N]",
    flags: [
      { flag: "--json", description: "Output as structured JSON" },
      { flag: "--limit N", description: "Show only the top N versions (0 = all)" },
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
  },
  {
    category: "release",
    name: "revert", alias: undefined, description: "Revert to a specific release version",
    usage: "gitmap revert <version>",
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
  },
  {
    category: "navigation",
    name: "list", alias: "ls", description: "Show all tracked repos with slugs",
    usage: "gitmap list [--group <name>] [--verbose]",
  },
  {
    category: "navigation",
    name: "group", alias: "g", description: "Manage repo groups",
    usage: "gitmap group <create|add|remove|list|show|delete> [args]",
  },
  {
    category: "navigation",
    name: "diff-profiles", alias: "dp", description: "Compare repos across two profiles",
    usage: "gitmap diff-profiles <profileA> <profileB> [--all] [--json]",
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
  },
  {
    category: "history",
    name: "history-reset", alias: "hr", description: "Clear command execution history",
    usage: "gitmap history-reset [--confirm]",
  },
  {
    category: "history",
    name: "stats", alias: "ss", description: "Show aggregated usage and performance metrics",
    usage: "gitmap stats [--json]",
  },
  {
    category: "history",
    name: "amend", alias: "am", description: "Rewrite commit author info",
    usage: "gitmap amend [commit-hash] --name <name> --email <email>",
  },
  {
    category: "history",
    name: "amend-list", alias: "al", description: "List previous author amendments",
    usage: "gitmap amend-list [--json]",
  },

  // --- Project Detection ---
  {
    category: "detection",
    name: "go-repos", alias: "gr", description: "List detected Go projects",
    usage: "gitmap go-repos [--json]",
  },
  {
    category: "detection",
    name: "node-repos", alias: "nr", description: "List detected Node.js projects",
    usage: "gitmap node-repos [--json]",
  },
  {
    category: "detection",
    name: "react-repos", alias: "rr", description: "List detected React projects",
    usage: "gitmap react-repos [--json]",
  },
  {
    category: "detection",
    name: "cpp-repos", alias: "cr", description: "List detected C++ projects",
    usage: "gitmap cpp-repos [--json]",
  },
  {
    category: "detection",
    name: "csharp-repos", alias: "csr", description: "List detected C# projects",
    usage: "gitmap csharp-repos [--json]",
  },

  // --- Data & Profiles ---
  {
    category: "data",
    name: "export", alias: "ex", description: "Export database to file",
    usage: "gitmap export [--json]",
  },
  {
    category: "data",
    name: "import", alias: "im", description: "Import repos from file",
    usage: "gitmap import <file>",
  },
  {
    category: "data",
    name: "profile", alias: "pf", description: "Manage database profiles",
    usage: "gitmap profile <create|list|switch|delete|show> [name]",
  },
  {
    category: "data",
    name: "bookmark", alias: "bk", description: "Save and run bookmarked commands",
    usage: "gitmap bookmark <save|list|run|delete> [args]",
  },
  {
    category: "data",
    name: "db-reset", alias: undefined, description: "Reset the local SQLite database",
    usage: "gitmap db-reset [--confirm]",
  },

  // --- Utilities ---
  {
    category: "utilities",
    name: "setup", alias: undefined, description: "Interactive first-time configuration wizard",
    usage: "gitmap setup",
  },
  {
    category: "utilities",
    name: "doctor", alias: undefined, description: "Diagnose PATH, deploy, and version issues",
    usage: "gitmap doctor [--fix-path]",
  },
  {
    category: "utilities",
    name: "update", alias: undefined, description: "Self-update from source repo",
    usage: "gitmap update",
  },
  {
    category: "utilities",
    name: "version", alias: "v", description: "Show version number",
    usage: "gitmap version",
  },
  {
    category: "utilities",
    name: "seo-write", alias: "sw", description: "Auto-commit SEO messages",
    usage: "gitmap seo-write [flags]",
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
  },
];
