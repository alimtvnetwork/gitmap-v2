# Refactor: release/assets.go

## Problem
`assets.go` is 261 lines with two responsibilities: public cross-compilation orchestration (types, detection, compile loop, staging) and low-level build helpers (single-target compilation, output naming, environment setup, file existence checks).

## Target Layout

### assets.go (~147 lines) — Orchestration & Public API
Stays:
- `type BuildTarget`
- `type CrossCompileResult`
- `DetectGoProject()`
- `ReadModuleName()`
- `BinaryName()`
- `FindMainPackages()`
- `CrossCompile()`
- `resolveBinName()`
- `CollectSuccessfulBuilds()`
- `EnsureStagingDir()`
- `CleanupStagingDir()`

### assetsbuild.go (~90 lines) — Build Helpers
Moves:
- `buildSingleTarget()`
- `formatOutputName()`
- `buildEnv()`
- `setEnv()`
- `fileExists()`

Imports: `fmt`, `os`, `os/exec`, `path/filepath`, `strings`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `release`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
