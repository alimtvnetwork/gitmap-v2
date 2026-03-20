export interface ChangelogEntry {
  version: string;
  items: string[];
}

export const changelog: ChangelogEntry[] = [
  {
    version: "v2.24.0",
    items: [
      "Bumped compiled version constant to v2.24.0.",
      "Moved release metadata writing from the release branch to the original branch.",
      "Auto-commit now handles `.release/` files after returning to the original branch.",
      "Removed `commitReleaseMeta` step from the release branch workflow.",
      "Simplified `pushAndFinalize` to complete without metadata writes.",
    ],
  },
  {
    version: "v2.23.0",
    items: [
      "Bumped compiled version constant to v2.23.0.",
      "Added `--notes` / `-N` flag to `release-branch` and `release-pending` commands.",
      "Updated docs site Release page with metadata-first workflow and release notes documentation.",
    ],
  },
  {
    version: "v2.22.0",
    items: [
      "Bumped compiled version constant to v2.22.0.",
      "Persisted zip group metadata in `.release/vX.Y.Z.json` via new `zipGroups` field.",
      "Documented `-A`/`--alias` flag in help text for `pull`, `exec`, `status`, and `cd`.",
      "Added shell completion support for `alias` and `zip-group` subcommands.",
      "Added `--list-aliases` and `--list-zip-groups` completion list flags.",
      "Added unit tests for `collectZipGroupNames`.",
    ],
  },
  {
    version: "v2.21.0",
    items: [
      "Bumped compiled version constant to v2.21.0.",
      "Refactored `assetsupload.go` into three focused files: `githubapi.go`, `assetsupload.go`, `remoteorigin.go`.",
      "Rebuilt Project Detection docs page with pipeline visualization, metadata deep-dive, DB schema, and package layout.",
      "Added detection docs link from Projects dashboard page.",
      "Added unit tests for `store/location.go` covering symlink resolution, fallback, and double-nesting prevention.",
      "Added unit tests for `remoteorigin.go` covering HTTPS, SSH, and invalid URL parsing.",
    ],
  },
  {
    version: "v2.20.0",
    items: [
      "**Fixed**: `OpenDefault()` double-nesting bug where profile config resolved to `<binary>/data/data/profiles.json`.",
      "Added `DefaultDBPath()` diagnostic helper to `store/location.go`.",
      "`gitmap ls` now prints resolved DB path when `--verbose` or zero repos found.",
      "Created path resolution contract spec for database diagnostics.",
    ],
  },
  {
    version: "v2.19.0",
    items: [
      "Bumped compiled version constant to v2.19.0.",
    ],
  },
  {
    version: "v2.18.0",
    items: [
      "Added batch status terminal demo to Batch Actions page showing dirty/clean state across repos.",
      "Fixed missing `os/exec` import in release asset upload.",
      "Resolved `deriveSlug` redeclaration conflict in project repos output.",
      "Removed unused `os` import from audit command.",
    ],
  },
  {
    version: "v2.17.0",
    items: [
      "Added 30-second auto-refresh timer to TUI dashboard via `tea.Tick`.",
      "Dashboard refresh interval configurable via `dashboardRefresh` in `config.json`.",
      "Added `--refresh` flag to `interactive` command for CLI-level override.",
      "Refresh interval validates with fallback to default 30s when missing or invalid.",
    ],
  },
  {
    version: "v2.16.0",
    items: [
      "Wired real `gitutil.Status()` into TUI dashboard for live dirty/clean indicators.",
      "Dashboard now shows ahead/behind counts and stash per repo.",
      "Async background refresh on TUI startup; manual refresh via `r` key.",
      "Summary bar with aggregate dirty/behind/stash counts and UTC timestamp.",
    ],
  },
  {
    version: "v2.15.1",
    items: [
      "**Fixed**: Database now resolves to `<binary-location>/data/gitmap.db` instead of CWD-relative `gitmap-output/data/`.",
      "Added `store.OpenDefault()` and `store.OpenDefaultProfile()` for binary-relative database access.",
      "Added `store/location.go` with `BinaryDataDir()` using `os.Executable()` + `filepath.EvalSymlinks()`.",
      "Updated all 13 database callers across the codebase to use binary-relative paths.",
      "Removed unused `resolveAuditOutputDir()` and `resolveDefaultOutputDir()` helpers.",
    ],
  },
  {
    version: "v2.15.0",
    items: [
      "Added cross-platform build support: `run.sh` (Linux/macOS) with full parity to `run.ps1`.",
      "Fixed Makefile flags to match `run.sh` argument format.",
      "Added GitHub Actions CI workflow: test on push, cross-compile 6 OS/arch targets.",
      "Added GitHub Actions Release workflow: auto-release on `v*` tags with compression and checksums.",
      "Added interactive TUI mode (`gitmap interactive` / `gitmap i`) built with Bubble Tea.",
      "TUI repo browser with fuzzy search, multi-select, and keyboard navigation.",
      "TUI batch actions: pull, exec, status across selected repos.",
      "TUI group management: browse, create, delete groups interactively.",
      "TUI status dashboard with live repo status view.",
      "Added Build System section to Architecture documentation page.",
      "Added spec documents: `42-cross-platform.md` and `43-interactive-tui.md`.",
    ],
  },
  {
    version: "v2.14.0",
    items: [
      "Added Go release assets: automatic cross-compilation for 6 OS/arch targets (windows/linux/darwin × amd64/arm64).",
      "Added GitHub Releases API integration for asset upload — no `gh` CLI needed.",
      "Added `--compress` flag to wrap release assets in `.zip` (Windows) or `.tar.gz` (Linux/macOS).",
      "Added `--checksums` flag to generate SHA256 `checksums.txt` for all release assets.",
      "Added `--no-assets` flag to skip automatic Go binary compilation.",
      "Added `--targets` flag for custom cross-compile target selection.",
      "Added `--list-targets` flag to print resolved target matrix and exit.",
      "Added config-driven release targets: `release.targets` in `config.json` overrides the default OS/arch matrix.",
      "Added config-driven `release.checksums` and `release.compress` booleans.",
      "Improved `gitmap ls <type>` output with labeled fields and inline `cd` examples.",
      "Added shell completion for `release`, `release-branch`, `group`, `multi-group`, and `list` commands.",
      "Fixed duplicate hints after `gitmap ls <type>` output.",
    ],
  },
  {
    version: "v2.13.0",
    items: [
      "Added group activation: `gitmap g <name>` sets a persistent active group for batch pull/status/exec.",
      "Added `multi-group` (mg) command for selecting and operating on multiple groups at once.",
      "Added `gitmap ls <type>` filtering: `gitmap ls go`, `gitmap ls node`, `gitmap ls groups`.",
      "Added contextual helper hints shown after command output to aid discoverability.",
      "Added Settings table for persistent key-value configuration in SQLite.",
      "Release metadata JSON now includes a `changelog` field.",
      "`gitmap list-versions` now shows changelog notes.",
    ],
  },
  {
    version: "v2.12.0",
    items: [
      "Added `gitmap list-versions` (`lv`) command: show all release tags sorted highest-first with changelog.",
      "Added `gitmap revert <version>` command: checkout tag and handoff rebuild.",
      "Added global ⌘K command palette searching across commands, flags, and pages.",
    ],
  },
  {
    version: "v2.11.0",
    items: [
      "Added constants inventory audit section documenting ~280 constants.",
      "Added Changelog page with timeline view and expand/collapse controls.",
      "Added Flag Reference page with sortable, searchable table of all flags.",
      "Added Interactive Examples page with animated terminal demos.",
    ],
  },
  {
    version: "v2.10.0",
    items: [
      "Full compliance audit (Wave 1 + Wave 2): all 75+ source files pass code style rules.",
      "Trimmed oversized files, fixed negation/switch violations, extracted constants.",
    ],
  },
  {
    version: "v2.9.0",
    items: [
      "Full code style refactor of `latest-branch` command.",
      "Split handler into 3 files, all under 200 lines.",
      "Added `--filter` flag to `latest-branch`.",
      "Added `--sort` flag to `latest-branch`.",
    ],
  },
  {
    version: "v2.8.0",
    items: [
      "Added `gitmap cd` (`go`) command: jump to any tracked repo by slug or partial name.",
      "Added `gitmap watch` (`w`) command: live terminal dashboard monitoring repo status.",
      "Added `gitmap diff-profiles` (`dp`) command: compare two profiles side-by-side.",
      "Added clone progress bars with retry logic and Windows long-path warnings.",
      "Added `gitmap gomod` (`gm`) command: rename Go module path across an entire repo.",
      "Added `gitmap/Makefile` as a thin wrapper around `run.sh` for Linux/macOS.",
      "Added `run.sh` cross-platform build script for Linux and macOS.",
    ],
  },
  {
    version: "v2.5.0",
    items: [
      "Added `gitmap profile` (`pf`) command: manage multiple database profiles.",
      "Added `gitmap export` (`ex`) command: export the full database as portable JSON.",
      "Added `gitmap import` (`im`) command: restore database from a backup file.",
      "Added `gitmap bookmark` (`bk`) command: save and replay command+flag combinations.",
      "Added `gitmap stats` (`ss`) command: aggregated usage statistics from command history.",
      "Added `gitmap history` (`hi`) command: queryable audit trail of CLI executions.",
      "Added `gitmap amend` (`am`) command: rewrite author name/email on existing commits.",
      "Added `gitmap seo-write` (`sw`) command: automated SEO commit scheduler.",
    ],
  },
  {
    version: "v2.4.0",
    items: [
      "Added `gitmap latest-branch` (`lb`) command.",
      "Positional integer shorthand: `gitmap lb 3` equals `gitmap lb --top 3`.",
      "Added `--format` flag to `latest-branch` (terminal, json, csv).",
      "Centralized date display formatting with local timezone conversion.",
    ],
  },
  {
    version: "v2.3.12",
    items: [
      "Added `gitmap list-releases` (`lr`) command.",
      "Enhanced `gitmap scan` to import `.release/v*.json` metadata files.",
      "Added `Source` column to the `Releases` table.",
      "Added `Releases` table to SQLite database.",
      "Converted all DB table/column names from snake_case to PascalCase.",
    ],
  },
  {
    version: "v2.3.10",
    items: [
      "Self-update hardening: rename-first strategy, stale-process fallback.",
      "Update enhancements: skip-if-current, version comparison, rollback safety.",
      "Added `update-cleanup` command with auto-run.",
    ],
  },
  {
    version: "v2.3.7",
    items: [
      "Release-pending, changelog, doctor commands.",
      "Database with repos and group management.",
      "Generic spec files in `spec/02-general/`.",
    ],
  },
];
