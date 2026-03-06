# Issue: PATH sync lock loop during `gitmap update` (v2.3.9 → v2.3.10)

## Summary

`gitmap update` repeatedly failed to sync the active PATH binary because:
1. The parent process held a lock on the active binary while waiting synchronously for the handoff copy.
2. `run.ps1` used copy-first (overwrite) for PATH sync instead of rename-first, causing 20 retries then failure.
3. A generated PowerShell script contained `Read-Host` which failed in non-interactive mode.

## Observed Symptoms

- `run.ps1` deploy succeeded but PATH sync looped 20 retries then fell through.
- `Read-Host` threw `PSInvalidOperationException` in non-interactive PowerShell sessions.
- Build failures from missing closing brace after `os.Exit(0)` when switching between `cmd.Run()` and `cmd.Start()`.

## Root Causes

### RC1: Parent held lock during sync (synchronous wait)
- `runUpdate()` used `cmd.Run()` to wait for the handoff copy to finish.
- This means the parent process (the active PATH binary) was still alive during the entire update pipeline, including the PATH sync step.
- Windows blocked `Copy-Item` because the parent held the file lock.

### RC2: Copy-first instead of rename-first
- `run.ps1` PATH sync tried `Copy-Item` (overwrite) as the primary strategy.
- On Windows, you cannot overwrite a running `.exe`, but you CAN rename it.
- The rename fallback existed but only triggered after exhausting all 20 copy retries.

### RC3: `Read-Host` in non-interactive session
- The generated update script ended with `Read-Host` ("Press Enter to continue...").
- When PowerShell runs via `exec.Command` from Go, it operates in non-interactive mode.
- `Read-Host` is not available in non-interactive mode, causing a hard error.

### RC4: Syntax error from incomplete refactoring
- When switching from `cmd.Run()` back to `cmd.Start()` + `os.Exit(0)`, the closing brace `}` for `runUpdate()` was lost.
- This caused `go build` to fail with "unexpected name runUpdateRunner".

## Solutions Applied

1. **Handoff uses `cmd.Start()` + `os.Exit(0)`** — parent exits immediately, releasing the file lock before deploy/sync runs.
2. **Rename-first PATH sync in `-Update` mode** — `run.ps1` now renames the active binary to `.old` first, then copies the new one, avoiding the lock entirely.
3. **Removed `Read-Host`** — update script exits cleanly without waiting for user input.
4. **Diagnostic log** — prints active exe path and handoff copy path at update start.
5. **Unique temp script names** — `gitmap-update-*.ps1` instead of fixed `gitmap-update.ps1` to avoid stale collisions.

## Acceptance Criteria

- Parent process exits before `run.ps1` starts (verified by zero lock retries during PATH sync).
- PATH sync uses rename-first in `-Update` mode; falls back to copy-retry loop otherwise.
- No `Read-Host` or interactive prompts in generated update scripts.
- Build compiles cleanly with no syntax errors.

## Prevention Rules

1. **Never use `cmd.Run()` in `runUpdate()`** — the parent MUST exit before the worker touches any binaries.
2. **Always use rename-first for PATH sync during update** — copy-overwrite is unreliable on Windows when any process may hold the binary.
3. **Never add interactive prompts to generated scripts** — they run in non-interactive PowerShell sessions.
4. **After switching between `cmd.Run()` and `cmd.Start()`, verify the function has a closing brace** — this is a mechanical error that breaks the build.
5. **Any update-flow change must update ALL of:**
   - `gitmap/cmd/update.go`
   - `run.ps1`
   - `spec/01-app/09-build-deploy.md`
   - `spec/02-general/02-powershell-build-deploy.md`
   - `spec/02-general/03-self-update-mechanism.md`
   - `.lovable/memory/issues/` (if new failure mode)
