# 89 — Update PATH Sync

> Spec for the automatic binary sync step during `gitmap update`.

---

## Problem

On Windows, the compiled binary is deployed to a subdirectory (e.g., `E:\bin-run\gitmap\gitmap.exe`), but the system PATH may resolve to a different copy at the parent level (e.g., `E:\bin-run\gitmap.exe`). After a successful build-and-deploy, the active PATH binary remains stale, producing:

```
[FAIL] Active PATH version does not match deployed version.
```

---

## Solution

The update script now includes an **auto-sync** step between the deploy and verification phases. It compares the deployed binary against the active PATH binary and copies the newer version into the PATH location.

---

## 3-Step Sync Strategy

The sync runs inside the auto-generated PowerShell update script (`UpdatePSSync` constant).

### Step 1 — Detect Mismatch

1. Resolve the **active binary** via `Get-Command gitmap` (PATH lookup).
2. Resolve the **deployed binary** from the deploy target directory.
3. Compare their absolute paths — if identical, no sync is needed.
4. Run `gitmap version` on both — if versions match, no sync is needed.

### Step 2 — Copy-Item (Primary)

```powershell
Copy-Item -Path $resolvedDeployed -Destination $resolvedActive -Force
```

This is the default sync method. It overwrites the stale PATH binary with the freshly built one. Works in most cases where the file is not locked.

### Step 3 — Fallback Hint

If `Copy-Item` fails (file lock, permission denied, admin-mode conflict):

```
[WARN] Could not sync: <error>
[HINT] Run 'gitmap doctor --fix-path' manually.
```

The user is directed to `gitmap doctor --fix-path` for manual resolution.

> **Future enhancement**: Add a rename-then-copy fallback (`Move-Item` the locked file to `gitmap.exe.old`, then `Copy-Item` the new binary) and a kill-process fallback (`Stop-Process`) before resorting to the manual hint. These are not yet implemented.

---

## Script Execution Order

The update PowerShell script is assembled from template constants in this order:

| Order | Constant               | Purpose                                        |
|-------|------------------------|-------------------------------------------------|
| 1     | `UpdatePSHeader`       | Set working directory to source repo            |
| 2     | `UpdatePSDeployDetect` | Resolve deployed binary path from config/PATH   |
| 3     | `UpdatePSVersionBefore`| Capture pre-update version of active binary     |
| 4     | `UpdatePSRunUpdate`    | Execute `run.ps1 -Update` (build + deploy)      |
| 5     | **`UpdatePSSync`**     | **Auto-sync deployed → active PATH binary**     |
| 6     | `UpdatePSVersionAfter` | Capture post-update versions of both binaries   |
| 7     | `UpdatePSVerify`       | Compare versions and pass/fail the update       |
| 8     | `UpdatePSPostActions`  | Show changelog, run cleanup                     |

---

## Components

| Component                                | File                                  |
|------------------------------------------|---------------------------------------|
| Sync PowerShell block                    | `constants/constants_update.go`       |
| Script assembly (`buildUpdateScript`)    | `cmd/updatescript.go`                 |
| Script execution (`runUpdateScript`)     | `cmd/updatescript.go`                 |
| Deploy path detection                    | `constants/constants_update.go`       |
| Version verification                     | `constants/constants_update.go`       |

---

## Error Scenarios

| Scenario                         | Behavior                                              |
|----------------------------------|-------------------------------------------------------|
| Paths identical (same binary)    | Sync skipped silently                                 |
| Versions already match           | Sync skipped silently                                 |
| Copy succeeds                    | `[OK] Synced successfully.`                           |
| Copy fails (lock/permissions)    | `[WARN]` with hint to run `gitmap doctor --fix-path`  |
| Deployed binary not found        | Sync skipped (nothing to copy from)                   |
| Active binary not in PATH        | Sync skipped (no target to copy to)                   |

---

## See Also

- [spec/01-app/88-clone-direct-url.md](88-clone-direct-url.md) — Direct URL clone with auto-open
- [spec/09-pipeline/06-version-and-help.md](../09-pipeline/06-version-and-help.md) — Version display and update verification
- [gitmap/helptext/update.md](../../gitmap/helptext/update.md) — User-facing update documentation
