export interface ChangelogEntry {
  version: string;
  items: string[];
}

export const changelog: ChangelogEntry[] = [
  {
    version: "v2.36.7",
    items: [
      "Added SkipMeta integration test: 6 test cases verifying SkipMeta prevents metadata and latest.json creation.",
      "Added release rollback integration test: 5 test cases verifying branch/tag cleanup on simulated push failure.",
      "Added end-to-end release test: full cycle from version bump through metadata commit on a temp repo.",
      "E2E edge-case coverage: dry-run, no-commit, skip-meta, and duplicate version blocking.",
      "Added edge-case test suite: pre-release parsing/comparison, bump resolution, version ordering, multi-release sequences, and rc-to-stable promotion.",
      "Added TUI Temp Releases view: 9th tab with flat list, detail panel, and grouped-by-prefix aggregation.",
      "Added --stop-on-fail flag to pull and exec: halts batch after first failure.",
      "Enhanced BatchProgress with per-item failure tracking, detailed failure reports, and exit code 3 on partial failures.",
    ],
  },
  {
    version: "v2.36.6",
    items: [
      "Split assets.go into assets.go + assetsbuild.go (build helpers).",
      "Split zipgroupops.go into zipgroupops.go + zipgroupshow.go (display logic).",
      "Split tui.go into tui.go + tuiview.go (rendering logic).",
      "Split aliasops.go into aliasops.go + aliassuggest.go (interactive suggestions).",
      "Split tempreleaseops.go into tempreleaseops.go + tempreleaselist.go (listing).",
      "Split listreleases.go into listreleases.go + listreleasesload.go (data loading).",
      "Split listversions.go into listversions.go + listversionsutil.go (collection utils).",
      "Split sshgen.go into sshgen.go + sshgenutil.go (validation utils).",
      "Split scanprojects.go into scanprojects.go + scanprojectsmeta.go (metadata).",
      "Split amendexec.go into amendexec.go + amendexecprint.go (output formatting).",
      "Split status.go into status.go + statusprint.go (table formatting).",
      "Split exec.go into exec.go + execprint.go (result formatting).",
      "Split logs.go into logs.go + logsview.go (view rendering).",
      "Split compress.go into compress.go + compresstar.go (tar logic).",
      "Added refactoring specs 65–78 for all 14 file splits.",
      "All source files comply with the 200-line limit; no functional changes.",
    ],
  },
  {
    version: "v2.36.5",
    items: [
      "Split ziparchive.go (362 lines) into three files: ziparchive.go, zipio.go, zipdryrun.go.",
      "Split autocommit.go (352 lines) into two files: autocommit.go, autocommitgit.go.",
      "Split seowriteloop.go (340 lines) into two files: seowriteloop.go, seowritegit.go.",
      "Split workflowbranch.go (310 lines) into two files: workflowbranch.go, workflowpending.go.",
      "Split workflow.go (291 lines) into two files: workflow.go, workflowvalidate.go.",
      "Added refactoring specs 60–64 for all five file splits.",
      "All release/ and cmd/ files comply with the 200-line limit; no functional changes.",
    ],
  },
  {
    version: "v2.36.4",
    items: [
      "Split workflowfinalize.go (498 lines) into four domain-specific files: workflowfinalize.go, workflowdryrun.go, workflowzip.go, workflowgithub.go.",
      "Split root.go (388 lines) into seven domain-specific dispatch files: root.go, rootcore.go, rootrelease.go, rootutility.go, rootdata.go, roottooling.go, rootprojectrepos.go.",
      "Eliminated dispatchMisc (166 lines); replaced by dispatchData + dispatchTooling.",
      "All cmd/ and release/ files now comply with the 200-line limit.",
      "Added refactoring specs for workflowfinalize.go and root.go dispatch modularization.",
    ],
  },
  {
    version: "v2.36.3",
    items: [
      "Bumped compiled version constant to v2.36.3.",
      "Post-release migration: re-runs legacy directory migration after returning to original branch.",
      "Prevents `.release/` from persisting when older branches restore tracked legacy files.",
      "Simplified doctor legacy directory check — confirms auto-migration succeeded instead of warning.",
      "Scan output now always writes to `.gitmap/output/` relative to scan root via `resolveOutputDir`.",
      "Updated all helptext output examples to use `.gitmap/output/` paths.",
      "Audited all specs, memory docs, and helptext for stale legacy path references.",
    ],
  },
  {
    version: "v2.36.1",
    items: [
      "Bumped compiled version constant to v2.36.1.",
      "Added automatic database migration from legacy UUID TEXT IDs to INTEGER AUTOINCREMENT IDs.",
      "Fixed FK constraint violation (787) during scan when legacy UUID IDs were present in the Repos table.",
    ],
  },
  {
    version: "v2.36.0",
    items: [
      "Bumped compiled version constant to v2.36.0.",
      "Added automatic legacy directory migration: gitmap-output/ → .gitmap/output/, .release/ → .gitmap/release/, .deployed/ → .gitmap/deployed/.",
      "Migration runs at CLI startup before any command dispatch; skips if target already exists.",
    ],
  },
  {
    version: "v2.35.1",
    items: [
      "Bumped compiled version constant to v2.35.1.",
      "Added legacy UUID data detection to all remaining DB query paths: group show, group list, stats, history, status, and export.",
      "All DB query errors from legacy string-based IDs now show a recovery prompt instead of raw SQL errors.",
    ],
  },
  {
    version: "v2.35.0",
    items: [
      "Bumped compiled version constant to v2.35.0.",
      "Consolidated `.release/` and `gitmap-output/` under unified `.gitmap/` directory (`release/`, `output/`).",
      "Centralized all path constants (`GitMapDir`, `DefaultReleaseDir`, `DefaultOutputDir`) for single-point configuration.",
      "Migrated all database primary keys from UUID strings to `INTEGER PRIMARY KEY AUTOINCREMENT` (`int64`).",
      "Removed `github.com/google/uuid` dependency.",
      "Added `doctor` check that warns if legacy `.release/` or `gitmap-output/` directories exist.",
      "Updated all helptext, spec documents, and docs site to reference `.gitmap/` paths.",
    ],
  },
  {
    version: "v2.34.0",
    items: [
      "Bumped compiled version constant to v2.34.0.",
      "Fixed list-releases to read .gitmap/release/v*.json from the current repo first, falling back to the database only when no local files exist.",
      "Added SourceRepo constant to release model for repo-sourced release records.",
    ],
  },
  {
    version: "v2.33.0",
    items: [
      "Bumped compiled version constant to v2.33.0.",
      "Fixed auto-commit push rejection when remote branch advances during release: added pull --rebase recovery with single retry.",
      "Added 16-stage summary table with anchor links to verbose logging spec.",
    ],
  },
  {
    version: "v2.32.0",
    items: [
      "Bumped compiled version constant to v2.32.0.",
      "Documented autocommit verbose logging as pipeline stage 16 in the verbose logging spec.",
    ],
  },
  {
    version: "v2.31.0",
    items: [
      "Bumped compiled version constant to v2.31.0.",
      "Added verbose logging to auto-commit step: logs version, file counts, staging, commit message, and push target.",
    ],
  },
  {
    version: "v2.30.0",
    items: [
      "Bumped compiled version constant to v2.30.0.",
      "Renamed TempReleases Commit column to CommitSha to avoid SQLite reserved keyword conflict.",
      "Added automatic database migration for existing TempReleases tables.",
      "Added JSON struct tags to model.TempRelease for backward-compatible serialization.",
    ],
  },
  {
    version: "v2.29.0",
    items: [
      "Bumped compiled version constant to v2.29.0.",
      "Fixed TempReleases SQL syntax error: quoted reserved keyword Commit in CREATE TABLE, INSERT, and SELECT statements.",
      "Documented metadata persistence and rollback log points in verbose logging spec (stages 14–15 of 15).",
    ],
  },
  {
    version: "v2.28.0",
    items: [
      "Bumped compiled version constant to v2.28.0.",
      "Added verbose logging to release pipeline: version resolution, source resolution, git operations, asset collection, staging, cross-compilation, compression, checksums, zip groups, ad-hoc zips, GitHub upload, retry, metadata persistence, and rollback.",
      "Updated verbose logging spec with all 15 pipeline stages documented.",
      "Added pull conflict handling to run.ps1 and run.sh with stash/discard/clean/quit prompt.",
      "Added --force-pull flag to both build scripts for non-interactive CI usage.",
      "Fixed set -e early exit bug in run.sh git pull error handling.",
      "Fixed parseCommitLines and hasListFlag redeclaration conflicts.",
    ],
  },
  {
    version: "v2.27.0",
    items: [
      "Bumped compiled version constant to v2.27.0.",
      "Added doctor validation checks for config.json, database migration, lock file, and network connectivity.",
      "Added TUI release trigger overlay with patch/minor/major/custom version bump selection.",
      "Integrated batch progress tracking into pull, exec, and status commands with success/fail/skip counters.",
      "Added BatchProgress tracker to cloner package with quiet mode for programmatic use.",
      "Added TUI interaction tests covering tab switching, browser navigation, fuzzy search, and release triggers.",
      "Added alias suggestion tests covering auto-suggestion, conflict detection, and idempotent re-runs.",
    ],
  },
  {
    version: "v2.24.0",
    items: [
      "Bumped compiled version constant to v2.24.0.",
      "Moved release metadata writing from the release branch to the original branch.",
      "Auto-commit now handles `.gitmap/release/` files after returning to the original branch.",
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
      "Persisted zip group metadata in `.gitmap/release/vX.Y.Z.json` via new `zipGroups` field.",
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
      "**Fixed**: Database now resolves to `<binary-location>/data/gitmap.db` instead of CWD-relative `.gitmap/output/data/`.",
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
      "Enhanced `gitmap scan` to import `.gitmap/release/v*.json` metadata files.",
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
      "Generic spec files in `spec/03-general/`.",
    ],
  },
];
