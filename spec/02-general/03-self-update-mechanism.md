# Self-Update Mechanism

## Overview

A reusable pattern for CLI tools that update themselves from source,
solving the Windows file-lock problem where a running binary cannot
overwrite itself. This guide is framework-agnostic and applies to
any compiled CLI tool (Go, Rust, C#, etc.) deployed on Windows.

## The Problem

On Windows, a running `.exe` holds a file lock on its own binary.
If the update process tries to overwrite the file while the original
process is still running, the OS blocks the operation:

> "The process cannot access the file because it is being used by
> another process."

This does not occur on Linux/macOS, where a running binary can be
replaced on disk (the OS keeps the old inode until the process exits).

## Solution: Copy-and-Handoff

A three-layer approach that reliably bypasses file locks:

### Layer 1 — Copy and Re-launch

1. The running binary copies itself to a temp location:
   `%TEMP%\<toolname>-update-<pid>.exe`
2. It launches the temp copy with a flag (e.g. `update --from-copy`)
   to indicate it's the delegated updater.
3. The parent **exits immediately** to release the file lock.

```
# Pseudocode — applies to any compiled language
func runUpdate():
    if isRunningFromCopy():
        executeUpdate(repoPath)
        return
    
    tempPath = copyBinaryToTemp()
    launchProcess(tempPath, ["update", "--from-copy"])
    exit(0)  # Release file lock immediately
```

### Layer 2 — Delayed Rebuild

The temp copy orchestrates the actual update. It generates and runs
a script (PowerShell, bash, etc.) that:

1. **Waits 1–2 seconds** for the parent process to fully terminate
   and release all OS-level file handles.
2. **Runs the build pipeline** (pull source, resolve deps, build,
   deploy).
3. **Reports the updated version** so the user can confirm success.

```
# Example: generated PowerShell script
Start-Sleep -Seconds 1.2
Push-Location '<repo-path>'
& '<build-script>'
Pop-Location
```

### Layer 3 — Deploy Retry

The build pipeline's deploy/copy step wraps file operations in a
retry loop to handle the race condition where the parent process
hasn't fully released its handle yet.

Recommended defaults:
- **Max attempts:** 15–20
- **Delay between attempts:** 300–500ms
- **Total timeout:** ~10 seconds

```
# Pseudocode — retry loop for file copy
attempts = 0
while attempts < maxAttempts:
    try:
        copyFile(source, destination)
        break
    catch fileLocked:
        log("Target binary in use, retrying...")
        sleep(500ms)
        attempts++
if attempts >= maxAttempts:
    fail("Could not overwrite binary after retries")
```

## Flow Diagram

```
User runs: <tool> update
   │
   ├─ Parent copies self → %TEMP%\<tool>-update-<pid>.exe
   ├─ Parent launches copy with: update --from-copy
   ├─ Parent exits (releases file lock)
   │
   └─ Temp copy starts
      ├─ Waits 1–2 seconds
      ├─ Generates temp build script
      ├─ Runs: build pipeline (pull → build → deploy)
      │    └─ Deploy step retries if binary still locked
      ├─ Prints updated version
      └─ Cleans up temp script
```

## Prerequisites

- The **source repo path** must be available to the binary at runtime.
  Common approaches:
  - Embedded at build time via linker flags (e.g. Go `-ldflags`)
  - Stored in a config file next to the binary
  - Resolved from an environment variable
- A **build script** must exist at the known repo path.
- If the repo path is missing or invalid, the update command should
  print a clear error and exit.

## Recommended Enhancements

Beyond the core three-layer pattern, consider these improvements:

### Rollback Safety

Before overwriting the deployed binary, **back it up**:

```
# Before deploy
backup = destination + ".bak"
rename(destination, backup)

# After successful deploy
delete(backup)

# On failure
rename(backup, destination)  # Restore working binary
```

This ensures a failed build never leaves the user without a working
binary.

### Skip If Already Up-to-Date

Compare the current version against the latest available version
(from a metadata file, git tag, or version endpoint) and skip the
update if they match:

```
if currentVersion == latestVersion:
    print("Already up to date")
    exit(0)
```

### Checksum Verification

After building the new binary, verify its integrity:

```
expected = readHashFile(buildOutput + ".sha256")
actual = sha256(newBinary)
if expected != actual:
    fail("Binary checksum mismatch — build may be corrupted")
```

### Proactive Temp Cleanup

On startup (not just after update), scan `%TEMP%` for leftover
update copies from previous runs and delete them:

```
# On tool startup
for file in glob("%TEMP%/<tool>-update-*.exe"):
    tryDelete(file)
```

### Exit Code Propagation

The temp copy should propagate the build script's exit code so
the calling process (or CI) can detect failures:

```
result = runBuildScript(scriptPath)
exit(result.exitCode)
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| No repo path configured | Print error, exit 1 |
| Repo path doesn't exist | Print error, exit 1 |
| Build/compile fails | Script exits with error |
| Deploy locked after retries | Throw/fail with clear message |
| Temp copy fails to launch | Print error, exit 1 |
| Version unchanged after update | Warn user (possible build issue) |

## Cleanup

- The temp copy binary persists in `%TEMP%` until proactively cleaned
  or the OS purges temp files.
- The generated build script is deleted immediately after execution.

## Platform Considerations

| Platform | File Lock Behavior | Self-Update Approach |
|----------|--------------------|----------------------|
| Windows | Binary locked while running | Copy-and-handoff (this pattern) |
| Linux | Binary replaceable on disk | Direct overwrite (simpler) |
| macOS | Binary replaceable on disk | Direct overwrite (simpler) |

On Linux/macOS, the copy-and-handoff pattern still works but is
unnecessary — a simple in-place replace suffices. Consider
platform-detecting which path to take if cross-platform support
is needed.

## Key Learnings

1. **`exit()` doesn't release locks instantly** — the OS may hold
   file handles briefly after process termination. Always add a
   delay before attempting to overwrite.
2. **Always add retry logic** for file operations on deployed
   binaries — even with the handoff, a small timing window exists.
3. **Never run the build script from the same process** that holds
   the file lock on the target binary.
4. **Bump the version on every change** so the user can confirm
   the update actually applied.
5. **Always provide a rollback path** — if the update fails
   mid-deploy, the user should still have a working binary.
6. **Log verbosely during update** — self-update failures are hard
   to debug without detailed logs of each step.
