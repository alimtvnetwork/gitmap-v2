# Issue: `gitmap update` fails with "file is being used by another process"

## Root Cause

When `gitmap update` runs from `E:\bin-run\gitmap\gitmap.exe`, the process holds a file lock on the binary. The update triggers `run.ps1` which tries to `Copy-Item` over the same binary during the deploy step — but the original process hasn't exited yet, so Windows blocks the overwrite.

## Solution (Final)

Five-layer fix:

1. **Copy-and-handoff** (`gitmap/cmd/update.go`):
   - Parent copies itself to same directory as `gitmap-update-<pid>.exe` (fallback to `%TEMP%`)
   - Launches the copy with hidden `update-runner` command
   - Parent **exits immediately** via `cmd.Start()` + `os.Exit(0)` to release the file lock
   - **MUST use `cmd.Start()`, never `cmd.Run()`** — synchronous wait holds the lock

2. **Rename-first PATH sync** (`run.ps1` in `-Update` mode):
   - Renames the active binary to `.old` (Windows allows renaming a running exe)
   - Copies deployed binary to the active path
   - Falls back to copy-retry loop (20 x 500ms) only if rename fails

3. **Deploy with rollback** (`run.ps1`):
   - Backs up existing binary as `.old` before overwriting
   - `Copy-Item` wrapped in a retry loop (20 attempts, 500ms delay)
   - On failure after retries → restores `.old` backup
   - On success → leaves `.old` in place for cleanup command

4. **Auto-cleanup** (generated PowerShell script):
   - After successful update, runs `gitmap update-cleanup`
   - Removes `%TEMP%\gitmap-update-*.exe` temp copies
   - Removes `*.old` backup files from deploy directory
   - Also available as manual command for ad-hoc use

5. **Version comparison** (generated PowerShell script):
   - Compares old vs new version after rebuild
   - Warns if version unchanged (constant not bumped)

## Key Learnings

- **Windows file locks are held until the process fully terminates** — `cmd.Start()` + `os.Exit(0)` isn't always instant
- **Always add retry logic** for file operations on deployed binaries
- **A delay before rebuild** gives the parent process time to fully release handles
- **Don't assume `os.Exit(0)` releases locks immediately** — the OS may keep the handle briefly
- **Keep `.old` backups until explicitly cleaned** — serves as manual rollback if new version has issues
- **Auto-cleanup at end of update** — best of both worlds: cleanup happens but user can still roll back before it runs
- **Skip unnecessary rebuilds** — check `git pull` output before building
- **Use rename-first, not copy-first** — Windows blocks overwrite of running exe but allows rename
- **Never add `Read-Host` to generated scripts** — they run in non-interactive PowerShell

## What NOT to Repeat

- Don't use `cmd.Run()` in `runUpdate()` — this holds the lock during the entire pipeline
- Don't use copy-overwrite as the primary PATH sync strategy — use rename-first
- Don't run `run.ps1` (which overwrites the binary) from the same process that holds the lock
- Don't skip the deploy retry — even with the handoff, a small timing window exists
- Don't auto-delete backups on startup — use an explicit cleanup command
- Don't add `Read-Host` or interactive prompts to generated scripts
- Always bump the version so the user can confirm the update actually applied
