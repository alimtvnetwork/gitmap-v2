# Changelog

## v2.3.8
- Replaced `update --from-copy` with hidden `update-runner` command for cleaner handoff separation.
- Handoff copy now created in the same directory as the active binary (fallback to %TEMP% if locked).
- Added `-Update` flag to `run.ps1`: skips `git pull` (delegated to update script), runs build+deploy+sync only.
- Update script now calls `run.ps1 -Update` instead of managing git pull separately.
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
