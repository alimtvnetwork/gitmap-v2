# Refactor: cmd/scanprojects.go

## Problem
`scanprojects.go` is 224 lines with two responsibilities: top-level project detection orchestration (detect, upsert, resolve IDs) and language-specific metadata persistence (Go runnables, C# project files, C# key files, stale cleanup).

## Target Layout

### scanprojects.go (~98 lines) — Orchestration
Stays:
- `detectAllProjects()`
- `upsertProjectsToDB()`
- `upsertProjectRecords()`
- `resolveDetectedProjectID()`
- `upsertProjectMetadata()`

### scanprojectsmeta.go (~130 lines) — Metadata Persistence
Moves:
- `upsertGoProjectMeta()`
- `upsertGoRunnables()`
- `upsertCSharpProjectMeta()`
- `upsertCSharpFiles()`
- `upsertCSharpKeyFiles()`
- `collectRepoIDs()`
- `cleanStaleProjects()`
- `collectKeepIDs()`

Imports: `fmt`, `os`, `constants`, `detector`, `model`, `store`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
