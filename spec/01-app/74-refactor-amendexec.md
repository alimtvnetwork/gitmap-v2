# Refactor: cmd/amendexec.go

## Problem
`amendexec.go` is 223 lines with two responsibilities: git operations for amending (listing commits, parsing output, filter-branch, checkout) and output/display helpers (env-filter construction, author string building, force push, progress/dry-run printing).

## Target Layout

### amendexec.go (~130 lines) — Git Operations
Stays:
- `listCommitsForAmend()`
- `parseCommitLines()`
- `detectPreviousAuthor()`
- `getCurrentBranch()`
- `switchBranch()`
- `runFilterBranch()`
- `runAmendHead()`

### amendexecprint.go (~100 lines) — Output & Display
Moves:
- `buildEnvFilter()`
- `buildAuthorString()`
- `runForcePush()`
- `printAmendHeader()`
- `printAmendProgress()`
- `printAmendDryRun()`

Imports: `fmt`, `os`, `os/exec`, `strings`, `constants`, `model`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
