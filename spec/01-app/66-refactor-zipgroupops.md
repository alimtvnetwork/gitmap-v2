# Refactor: cmd/zipgroupops.go

## Problem
`zipgroupops.go` is 252 lines with two responsibilities: mutation commands (remove, delete, rename, sync) and read-only display commands (list, show, folder expansion).

## Target Layout

### zipgroupops.go (~133 lines) — Mutation Commands
Stays:
- `runZipGroupRemove()`
- `runZipGroupDelete()`
- `runZipGroupRename()`
- `parseZipGroupRenameFlags()`
- `executeZipGroupRename()`
- `syncZipGroupJSON()`

### zipgroupshow.go (~140 lines) — Display Commands
Moves:
- `runZipGroupList()`
- `printZipGroupList()`
- `runZipGroupShow()`
- `executeZipGroupShow()`
- `printZipGroupShow()`
- `expandFolder()`

Imports: `fmt`, `os`, `path/filepath`, `constants`, `model`, `store`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
