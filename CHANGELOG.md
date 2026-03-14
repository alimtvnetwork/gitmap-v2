# Changelog

## v2.10.0
- Version bump for next development cycle.

## v2.9.0
- Completed flags and examples for all 22 command entries on the documentation site.
- Added detailed flag tables and usage examples for `seo-write`, `doctor`, `update`, `pull`, `version`, `history-reset`, and `db-reset`.
- Filled in flags and examples for 15 commands missing both: `rescan`, `desktop-sync`, `status`, `latest-branch`, `release-branch`, `release-pending`, `changelog`, `group`, `list`, `diff-profiles`, `export`, `import`, `profile`, `bookmark`, and `stats`.

## v2.28.0
- Removed unused `detector` import from `cmd/scan.go` that caused build failure.
- Updated documentation site fonts: Ubuntu for headings, Poppins for body text, Ubuntu Mono for code blocks.

## v2.27.0
- Added `gitmap cd` (`go`) command: jump to any tracked repo by slug or partial name.
- Subcommands: `cd repos`, `cd set-default`, `cd clear-default`; supports `--group` and `--pick` flags.
- Added `gitmap watch` (`w`) command: live terminal dashboard monitoring repo status.
- Supports `--interval`, `--group`, `--no-fetch`, and `--json` snapshot mode.
- Added `gitmap diff-profiles` (`dp`) command: compare two profiles side-by-side.
- Supports `--all` and `--json` output flags.
- Added clone progress bars with retry logic and Windows long-path warnings.
- Built documentation site with interactive terminal preview for the watch command.
- Added `gitmap/Makefile` as a thin wrapper around `run.sh` for standard `make` workflows.
  - Targets: `build`, `run` (with `ARGS=`), `test`, `update`, `no-pull`, `no-deploy`, `clean`, `help`.
- Added Makefile documentation page to the docs site with target reference, examples, and argument-passing guide.
- Added `run.sh` cross-platform build script: Bash equivalent of `run.ps1` for Linux and macOS.
  - Full pipeline: pull, tidy, build, deploy with `-ldflags` version embedding.
  - Reads config from `powershell.json` via `jq` or `python3` fallback.
  - Supports `-t` (test with report), `-n` (no-pull), `-d` (no-deploy), and `-u` (update) flags.
- Added `gitmap gomod` (`gm`) command: rename Go module path across an entire repo with branch safety.
  - Replaces module directive in `go.mod` and all matching paths across **all files** by default.
  - Use `--ext "*.go,*.md,*.txt"` to restrict replacement to specific file extensions.
  - Creates `backup/before-replace-<slug>` and `feature/replace-<slug>` branches automatically.
  - Commits changes on the feature branch and merges back to the original branch.
  - Supports `--dry-run`, `--no-merge`, `--no-tidy`, `--verbose`, and `--ext` flags.

## v2.26.0
- Version bump to v2.26.0 following `gitmap profile` command addition.
- All profile subcommands (`create`, `list`, `switch`, `delete`, `show`) fully integrated and documented.

## v2.25.0
- Added `gitmap profile` (`pf`) command: manage multiple database profiles (work, personal, etc.).
- Subcommands: `create`, `list`, `switch`, `delete`, `show`.
- Each profile has its own SQLite database file (`gitmap-{name}.db`).
- Default profile uses existing `gitmap.db` for full backward compatibility.
- Profile config stored in `gitmap-output/data/profiles.json`.
- All commands automatically use the active profile's database.

## v2.24.0
- Added `gitmap import` (`im`) command: restore database from a `gitmap-export.json` backup file.
- Merge semantics: upserts repos/releases, INSERT OR IGNORE for history/bookmarks/groups.
- Group members re-linked by resolving `repoSlugs` against the Repos table.
- Requires `--confirm` flag to prevent accidental data changes.

## v2.23.0
- Added `gitmap export` (`ex`) command: export the full database as a portable JSON file.
- Exports all tables: repos, groups (with member repo slugs), releases, command history, and bookmarks.
- Default output: `gitmap-export.json`; accepts optional custom file path.
- Summary line shows counts for each exported section.

## v2.22.0
- Added `gitmap bookmark` (`bk`) command: save and replay frequently-used command+flag combinations.
- Subcommands: `save`, `list`, `run`, `delete` — full CRUD for saved bookmarks.
- `bookmark run <name>` replays the saved command through standard dispatch (appears in audit history).
- `bookmark list --json` outputs bookmarks as JSON.
- New `Bookmarks` SQLite table with unique name constraint.
- `db-reset --confirm` now also clears the Bookmarks table.

## v2.21.0
- Added `gitmap stats` (`ss`) command: aggregated usage statistics from command history.
- Shows most-used commands, success/fail counts, failure rates, and avg/min/max durations.
- Supports `--command <name>` filter and `--json` output.
- Summary row displays overall totals across all commands.

## v2.20.0
- Added `gitmap history` (`hi`) command: queryable audit trail of all CLI command executions.
- Three detail levels: `--detail basic` (command + timestamp), `--detail standard` (+ flags + duration), `--detail detailed` (+ args + repos + summary).
- Supports `--command <name>` filter, `--limit N`, and `--json` output.
- Added `gitmap history-reset` (`hr`) command: clears audit history (requires `--confirm`).
- New `CommandHistory` SQLite table auto-records every command with start/end timestamps, duration, exit code, and affected repo count.
- `db-reset --confirm` now also clears the CommandHistory table.

## v2.19.0
- Added `gitmap amend` (`am`) command: rewrite author name/email on existing commits with three modes (all, range, HEAD).
- Supports `--branch` flag to operate on a specific branch (auto-switches back to original branch after completion).
- SHA as first positional argument: `gitmap amend <sha> --name "Name"` rewrites from that commit to HEAD.
- `--dry-run` previews affected commits without modifying history or writing audit records.
- `--force-push` auto-runs `git push --force-with-lease` after amend.
- Audit trail: every amend operation writes a JSON log to `.gitmap/amendments/amend-<timestamp>.json` with full details.
- Database persistence: amendment records saved to `Amendments` SQLite table for queryable history.
- `db-reset --confirm` now also clears the `Amendments` table.
- Added `--author-name` and `--author-email` flags to `gitmap seo-write` (`sw`): set custom author on each commit.
- SEO-write dry-run now displays the author that would be used when author flags are set.

## v2.18.0
- Added `gitmap seo-write` (`sw`) command: automated SEO commit scheduler that stages, commits, and pushes files on a randomized interval.
- Supports CSV input mode (`--csv`) for user-provided title/description pairs.
- Supports template mode with placeholder substitution (`{service}`, `{area}`, `{url}`, `{company}`, `{phone}`, `{email}`, `{address}`).
- Pre-seeded `data/seo-templates.json` with 25 title and 20 description templates (500 unique combinations).
- Added `CommitTemplates` SQLite table for persistent template storage with auto-seeding on first run.
- Rotation mode: when pending files are exhausted, appends/reverts text in a target file to maintain commit activity.
- Configurable interval (`--interval min-max`), commit limit (`--max-commits`), file selection (`--files`), and dry-run preview.
- Added `--template <path>` flag to load templates from a custom JSON file at runtime.
- Added `--create-template` / `ct` shorthand to scaffold a sample `seo-templates.json` in the current directory.
- Graceful shutdown on Ctrl+C (finishes current commit before exiting).

## v2.17.0
- Added `Source` column to the `Releases` table: tracks whether each release was created via `gitmap release` (`release`) or imported from `.release/` files (`import`).
- Added `--source` flag to `gitmap list-releases` (`lr`): filter releases by origin (`--source release` or `--source import`).
- Added `--source` flag to `gitmap list-versions` (`lv`): cross-references git tags with the Releases DB to filter by source and display source metadata.
- Added `--source` flag to `gitmap changelog` (`cl`): filter changelog entries by release source.
- Terminal and JSON output for `list-releases` and `list-versions` now includes the Source field.

## v2.16.0
- Added `gitmap list-releases` (`lr`) command: queries the Releases DB table and displays stored releases with `--json` and `--limit N` support.
- Enhanced `gitmap scan` to import `.release/v*.json` metadata files into the Releases DB table automatically after each scan.

## v2.15.0
- Added `--limit N` flag to `gitmap list-versions` (`lv`): show only the top N versions (0 or omitted = all).

## v2.14.0
- Added `Releases` table to SQLite database: stores release metadata (version, tag, branch, commit, changelog, flags) persistently.
- Release workflow now auto-persists metadata to the database after successful releases.
- Converted all database table and column names from snake_case to PascalCase (`Repos`, `Groups`, `GroupRepos`, `Releases`).
- Added `store/release.go` with `UpsertRelease`, `ListReleases`, `FindReleaseByTag` methods.
- Added `model/release.go` with `ReleaseRecord` struct.
- Note: existing databases will need `gitmap db-reset --confirm` to adopt the new schema.

## v2.13.0
- Release metadata JSON (`.release/vX.Y.Z.json`) now includes a `changelog` field with notes from CHANGELOG.md (gracefully omitted if unreadable).
- `gitmap list-versions` (`lv`) now shows changelog notes as sub-points under each version in terminal output.
- `gitmap list-versions --json` includes changelog array per version in JSON output.

## v2.12.0
- Added `gitmap list-versions` (`lv`) command: lists all release tags sorted highest-first, with `--json` output support.
- Added `gitmap revert <version>` command: checks out a release tag and rebuilds/deploys via handoff (same mechanism as `update`).

## v2.11.0
- Added constants inventory audit section to compliance spec, documenting ~280 constants across 9 files and 17 categories.

## v2.10.0
- Full compliance audit (Wave 1 + Wave 2): all 75 source files pass code style rules.
  - Trimmed 4 oversized files: `workflow.go`, `terminal.go`, `safe_pull.go`, `setup.go` (all under 200 lines).
  - Fixed all negation and switch violations across `changelog.go`, `github.go`, `metadata.go`, `config.go`, `verbose.go`, `semver.go`.
  - Extracted missing constants to dedicated constants files.

## v2.9.0
- Full code style refactor of `latest-branch` command:
  - Split `cmd/latestbranch.go` into 3 files: handler, resolve, output (all under 200 lines).
  - Split `gitutil/latestbranch.go` into 2 files: core operations, resolve helpers.
  - All functions comply with 8-15 line limit. Positive logic throughout.
  - Blank line before every return. No magic strings. Chained if+return replaces switch.
  - Extracted git constants and display message constants.

## v2.8.0
- Added `--filter` flag to `latest-branch`: filter branches by glob pattern (e.g. `feature/*`) or substring match.

## v2.7.0
- Added `--sort` flag to `latest-branch`: supports `date` (default, descending) and `name` (alphabetical ascending).

## v2.6.0
- Centralized date display formatting: all dates now convert to local timezone and display as `DD-Mon-YYYY hh:mm AM/PM`.
- Added `gitutil/dateformat.go` with `FormatDisplayDate` and `FormatDisplayDateUTC` functions.
- Updated `latest-branch` terminal, JSON, and CSV output to use the new date format.

## v2.5.1
- Added `--no-fetch` flag to `latest-branch`: skips `git fetch --all --prune` when remote refs are already up to date.

## v2.5.0
- Added `--format` flag to `latest-branch`: supports `terminal` (default), `json`, and `csv` output formats.
  - CSV outputs a header row + data rows to stdout, suitable for piping and spreadsheets.
  - `--json` remains as shorthand for `--format json`.
- Refactored `latest-branch` output into dedicated functions per format.

## v2.4.1
- Added positional integer shorthand for `latest-branch`: `gitmap lb 3` is equivalent to `gitmap lb --top 3`.

## v2.4.0
- Added `gitmap latest-branch` (`lb`) command: finds the most recently updated remote branch by commit date and displays name, SHA, date, and subject.
  - Flags: `--remote`, `--all-remotes`, `--contains-fallback`, `--top N`, `--json`.
  - Positional integer shorthand: `gitmap lb 3` is equivalent to `gitmap lb --top 3`.

## v2.3.12
- Spec, issue post-mortems, and memory aligned to codify synchronous update handoff and rename-first PATH sync as permanent rules.
- Rename-first PATH sync in `-Update` mode: renames active binary to `.old` before copying, eliminating lock-retry loops.
- Parent `update` handoff uses `cmd.Start()` + `os.Exit(0)` to release file lock before worker runs.
- Handoff diagnostic log prints active exe and copy paths at update start.
- Spec consistency pass: all four update-flow specs now enforce identical rules.

## v2.3.10
- Fixed `Read-Host` error in non-interactive PowerShell sessions during update by removing trailing prompt.
- Parent `update` process now exits immediately (handoff copy runs synchronously via `update-runner`).
- Added diagnostic log at update start showing active exe path and handoff copy path.
- Update script now uses unique temp file names (`gitmap-update-*.ps1`) to avoid stale script collisions.

## v2.3.9
- Version bump for rebuild validation after update-runner handoff changes.

- Replaced `update --from-copy` with hidden `update-runner` command for cleaner handoff separation.
- Handoff copy now created in the same directory as the active binary (fallback to %TEMP% if locked).
- Added `-Update` flag to `run.ps1`: runs full update pipeline (pull, build, deploy, sync) with post-update validation and cleanup.
- Update script delegates entire pipeline to `run.ps1 -Update`.
- Before/after version output derived from actual executables, not static constants.
- Mandatory `update-cleanup` runs after successful update to remove handoff and `.old` artifacts.
- Cleanup now scans both `%TEMP%` and same-directory for leftover `gitmap-update-*.exe` files.

- Added `gitmap doctor --fix-path` flag: automatically syncs the active PATH binary from the deployed binary using retry (20×500ms), rename fallback, and stale-process termination, with clear confirmation output.
- Doctor diagnostics now suggest `--fix-path` when version mismatches are detected.

## v2.3.6
- Added stale-process fallback during PATH-binary sync (`update` + `run.ps1`): if copy+rename fail, it now stops stale `gitmap.exe` processes bound to the old path and retries once.
- Improved failure guidance to run the deployed binary directly when active PATH binary remains locked.

## v2.3.5
- Hardened `gitmap update` PATH sync with retry + rename fallback, and it now exits with failure if active PATH binary remains stale.
- Clarified update output labels to distinguish source version (`constants.go`) vs active executable version.
- Added same rename-fallback PATH sync behavior in `run.ps1`.

## v2.3.4
- Updated PATH-binary sync in `run.ps1` and `gitmap update` to use retry-on-lock behavior (20 attempts × 500ms), matching the self-update spec.
- Added explicit recovery guidance when active PATH binary is still locked, including an exact `Copy-Item` fix command.

## v2.3.3
- Added `gitmap doctor` command: reports PATH binary, deployed binary, version mismatches, git/go availability, and recommends exact fix commands.

## v2.3.2
- `gitmap update` now syncs the active PATH binary with the deployed binary, so commands like `release` are available immediately.
- `gitmap update` now prints changelog bullet points after update (or no-op update) for quick visibility.
- Added `gitmap changelog --open` and `gitmap changelog.md` to open `CHANGELOG.md` in the default app.

## v2.3.1
- Added `gitmap changelog` command for concise, CLI-friendly release notes.
- Improved `gitmap update` output to show deployed binary/version and warn if PATH points to another binary.
- `gitmap update` now prints latest changelog notes after a successful update.

## v2.3.0
- Added `gitmap release-pending` (`rp`) to release all `release/v*` branches missing tags.
- `gitmap release` and `gitmap release-branch` now switch back to the previous branch after completion.

## v2.2.3
- Fixed PowerShell parser-breaking characters in update/deploy output paths.
- Improved deployment rollback messaging in `run.ps1`.

## v2.2.2
- Added additional parser safety fixes for update script output.

## v2.2.1
- Patched PowerShell parsing edge cases affecting update flow.
