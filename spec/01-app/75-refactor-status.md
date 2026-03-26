# Refactor: cmd/status.go

## Problem
`status.go` is 219 lines with two responsibilities: command orchestration with scope resolution (flag parsing, alias/group/DB/JSON loading) and terminal display (banner, table headers, summary formatting with colored output).

## Target Layout

### status.go (~133 lines) — Command & Data Loading
Stays:
- `runStatus()`
- `parseStatusFlags()`
- `loadStatusByScope()`
- `loadRecordsByGroup()`
- `loadAllRecordsDB()`
- `loadRecordsJSONFallback()`
- `loadStatusRecords()`
- `type statusSummary`

### statusprint.go (~100 lines) — Display & Formatting
Moves:
- `printStatusBanner()`
- `printStatusTable()`
- `printStatusTableTracked()`
- `printStatusHeader()`
- `printStatusSummary()`
- `buildSummaryParts()`
- `appendSummaryPart()`

Imports: `fmt`, `strings`, `cloner`, `constants`, `model`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
