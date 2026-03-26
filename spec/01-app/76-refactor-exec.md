# Refactor: cmd/exec.go

## Problem
`exec.go` is 218 lines with two responsibilities: command orchestration with scope resolution (flag parsing, batch execution, progress tracking) and terminal display with single-repo execution (repo command runner, result printing, banner, summary formatting).

## Target Layout

### exec.go (~155 lines) — Command & Execution
Stays:
- `runExec()`
- `execAllReposTracked()`
- `execOneRepoTracked()`
- `execAllRepos()`
- `execOneRepo()`
- `parseExecFlags()`
- `loadExecByScope()`
- `loadExecRecordsJSON()`
- `loadExecRecords()`
- `truncate()`

### execprint.go (~82 lines) — Display & Formatting
Moves:
- `execInRepo()`
- `printExecResult()`
- `printExecOutput()`
- `printExecBanner()`
- `printExecSummary()`
- `buildExecSummaryParts()`

Imports: `fmt`, `os/exec`, `strings`, `constants`, `model`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
