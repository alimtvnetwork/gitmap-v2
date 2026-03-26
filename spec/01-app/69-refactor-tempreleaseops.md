# Refactor: cmd/tempreleaseops.go

## Problem
`tempreleaseops.go` is 234 lines with two responsibilities: branch creation workflow (create, pattern parsing, sequence management, dry-run) and listing/display (list, print, flag detection).

## Target Layout

### tempreleaseops.go (~170 lines) — Create Workflow
Stays:
- `runTempReleaseCreate()`
- `executeTRCreate()`
- `parseVersionPattern()`
- `resolveAutoStart()`
- `validateSequenceRange()`
- `formatSeq()`
- `createTempBranches()`
- `pushTempBranches()`
- `printTRDryRun()`

### tempreleaselist.go (~80 lines) — List & Display
Moves:
- `runTempReleaseList()`
- `printTRList()`
- `hasTRListFlag()`

Imports: `encoding/json`, `fmt`, `os`, `constants`, `model`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
