
# Plan: `gitmap-updater` — Release-Based Update Tool

## Problem
When gitmap is installed from a GitHub release (not built from source), `gitmap update` fails because there's no embedded repo path. Users need a way to update without cloning the source.

## Proposed Solution
A standalone `gitmap-updater` binary that:
1. Checks GitHub releases API for the latest version tag
2. Compares with the locally installed `gitmap version`
3. Creates a handoff copy of itself (to avoid Windows file locks)
4. Downloads and runs the release's `install.ps1` to update gitmap
5. Waits for completion and verifies

## Architecture & Flow

```
gitmap update (no repo path)
  → Detects gitmap-updater on PATH
  → Launches: gitmap-updater run
  
gitmap-updater run:
  1. GET github.com/api/releases/latest → extract tag
  2. Compare tag vs `gitmap version` output
  3. If same → "Already up to date", exit
  4. Create handoff copy: gitmap-updater-tmp-{pid}.exe
  5. Launch handoff copy with: gitmap-updater-tmp update-worker
  6. Worker downloads install.ps1 from release assets
  7. Worker runs install.ps1 (pinned to detected version)
  8. Worker verifies `gitmap version` matches expected
  9. Cleanup temp copy
```

## Risk Analysis

| Risk | Severity | Mitigation |
|------|----------|------------|
| Windows file lock on running binary | High | Handoff copy pattern (already proven in `gitmap update`) |
| GitHub API rate limit (60/hr unauthenticated) | Medium | Cache last-check timestamp; only check once per hour |
| Network failure mid-download | Medium | Atomic install via install.ps1 (already has rename-first) |
| Updater updating itself creates circular dependency | High | Updater does NOT self-update in v1; only updates gitmap |
| install.ps1 not present in older releases | Low | Fall back to direct binary download from assets |
| Platform detection (Windows vs Linux vs macOS) | Medium | Reuse architecture detection logic; Windows-first (PS1) |
| Updater binary not on PATH | Low | `gitmap update` prints install instructions if not found |

## Key Decision: Skip Updater Self-Update in v1
The user mentioned the updater should update itself first. This creates a complex chicken-and-egg problem. **Recommendation**: In v1, the updater only updates gitmap. Self-update can be added in v2 once the core flow is proven. This dramatically reduces risk.

## Tasks

### 1. Create `gitmap-updater` Go module structure
- New directory: `gitmap-updater/` (separate binary, same repo)
- `main.go`, `cmd/root.go`, `cmd/check.go`, `cmd/run.go`
- Minimal dependency footprint

### 2. Implement GitHub release version check
- HTTP GET to `https://api.github.com/repos/{owner}/{repo}/releases/latest`
- Parse `tag_name` from JSON response
- Compare with `gitmap version` output (exec + capture)

### 3. Implement handoff copy mechanism
- Copy self to temp path (`gitmap-updater-tmp-{pid}.exe`)
- Launch copy with `update-worker` subcommand
- Worker performs the actual download + install
- Cleanup temp on exit

### 4. Implement install.ps1 download & execution
- Download `install.ps1` from release assets URL
- Write to temp file with UTF-8 BOM
- Execute via PowerShell with `-ExecutionPolicy Bypass`
- Capture exit code and verify

### 5. Integrate with `gitmap update` fallback
- In `gitmap/cmd/update.go`: if `RepoPath` is empty AND `--repo-path` not given, look for `gitmap-updater` on PATH
- If found, exec `gitmap-updater run` and exit
- If not found, show current friendly error (with added note about installing gitmap-updater)

### 6. Add CI build for gitmap-updater
- Add to release workflow: build `gitmap-updater` for same 6 targets
- Include in release assets alongside gitmap binaries

### 7. Documentation & helptext
- `helptext/update.md` — document the fallback to gitmap-updater
- New `helptext/updater.md` or README section
- CHANGELOG entry

## Files to Create/Modify

**New files (~6):**
- `gitmap-updater/main.go`
- `gitmap-updater/cmd/root.go`
- `gitmap-updater/cmd/check.go`
- `gitmap-updater/cmd/run.go`
- `gitmap-updater/cmd/worker.go`
- `gitmap-updater/go.mod`

**Modified files (~4):**
- `gitmap/cmd/update.go` — add updater fallback
- `gitmap/constants/constants_update.go` — add updater-related constants
- `.github/workflows/release.yml` — build updater binaries
- `CHANGELOG.md`

## Estimated Scope
~10 files, medium-large feature. Core logic is straightforward since it reuses proven patterns (handoff copy, install.ps1, rename-first).
