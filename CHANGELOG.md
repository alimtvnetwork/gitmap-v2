# Changelog

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
