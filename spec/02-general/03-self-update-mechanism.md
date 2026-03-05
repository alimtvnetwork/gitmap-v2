# Self-Update Mechanism

## Overview

This document describes a reusable pattern for CLI tools that can
update themselves from source, handling the Windows file-lock problem
where a running binary cannot overwrite itself.

## The Problem

On Windows, a running `.exe` holds a file lock. If the update process
tries to overwrite the binary while it's still running, the OS blocks
the operation with "file is being used by another process."

## Solution: Copy-and-Handoff

A three-layer approach that reliably bypasses file locks:

### Layer 1 — Copy and Re-launch

1. The parent binary copies itself to `%TEMP%\toolname-update-<pid>.exe`.
2. It launches the temp copy with `update --from-copy`.
3. The parent **exits immediately** (`os.Exit(0)`) to release the lock.

```go
func runUpdate() {
    if isFromCopy() {
        executeUpdate(constants.RepoPath)
        return
    }
    tempPath := copyToTemp()
    launchCopy(tempPath)
    os.Exit(0)
}
```

### Layer 2 — Delayed Rebuild

The temp copy generates and runs a PowerShell script that:

1. Waits 1–2 seconds for the parent to fully terminate.
2. Runs the build script (`run.ps1`) which pulls, builds, and deploys.
3. Reports the updated version.

```go
func buildUpdateScript(repoPath, runPS1 string) string {
    return fmt.Sprintf(`
Start-Sleep -Seconds 1.2
Push-Location '%s'
& '%s'
Pop-Location
`, repoPath, runPS1)
}
```

### Layer 3 — Deploy Retry

The build script's deploy step retries file copy operations
(20 attempts × 500ms delay) to handle the race condition where
the parent hasn't fully released its handle yet.

## Flow Diagram

```
User runs: gitmap update
   │
   ├─ Parent copies self → %TEMP%\gitmap-update-1234.exe
   ├─ Parent launches copy with: update --from-copy
   ├─ Parent exits (releases file lock)
   │
   └─ Temp copy starts
      ├─ Waits 1.2 seconds
      ├─ Generates temp PowerShell script
      ├─ Runs: run.ps1 (pull → build → deploy)
      │    └─ Deploy retries if binary still locked
      ├─ Prints updated version
      └─ Cleans up temp script
```

## Prerequisites

- The source repo path must be embedded at build time via `-ldflags`.
- The build script (`run.ps1`) must exist at the embedded repo path.
- If the binary was not built with `run.ps1`, the update command
  prints an error and exits.

## Error Handling

| Scenario | Behavior |
|----------|----------|
| No embedded repo path | Print error, exit 1 |
| Repo path doesn't exist | Print error, exit 1 |
| Build fails | PowerShell script exits with error |
| Deploy locked after 20 retries | Throws, script fails |

## Cleanup

- The temp copy binary persists in `%TEMP%` (OS cleans up eventually).
- The generated PowerShell script is deleted after execution.

## Key Learnings

1. **`os.Exit(0)` doesn't release locks instantly** — always add a
   delay before the rebuild attempts to overwrite.
2. **Always add retry logic** for file operations on deployed binaries.
3. **Never run the build script from the same process** that holds
   the file lock on the target binary.
4. **Bump the version on every change** so the user can confirm
   the update actually applied.
