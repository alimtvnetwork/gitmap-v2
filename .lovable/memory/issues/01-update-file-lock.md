# Issue: `gitmap update` fails with "file is being used by another process"

## Root Cause

When `gitmap update` runs from `E:\bin-run\gitmap\gitmap.exe`, the process holds a file lock on the binary. The update triggers `run.ps1` which tries to `Copy-Item` over the same binary during the deploy step — but the original process hasn't exited yet, so Windows blocks the overwrite.

## Solution (Final)

Four-layer fix:

1. **Copy-and-handoff** (`gitmap/cmd/update.go`):
   - Parent copies itself to `%TEMP%\gitmap-update-<pid>.exe`
   - Launches the copy with `update --from-copy`
   - Parent **exits immediately** (`os.Exit(0)`) to release the file lock

2. **Skip-if-current** (generated PowerShell script):
   - Captures current deployed version before pulling
   - Runs `git pull` and checks output
   - If "Already up to date" → exits early (no rebuild)

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

## What NOT to Repeat

- Don't run `run.ps1` (which overwrites the binary) from the same process that holds the lock
- Don't skip the deploy retry — even with the handoff, a small timing window exists
- Don't auto-delete backups on startup — use an explicit cleanup command
- Always bump the version so the user can confirm the update actually applied
