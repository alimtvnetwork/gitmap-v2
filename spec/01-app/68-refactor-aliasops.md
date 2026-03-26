# Refactor: cmd/aliasops.go

## Problem
`aliasops.go` is 235 lines with two responsibilities: core CRUD operations (set, remove, list, show) and the suggest workflow (auto-propose aliases for unaliased repos with interactive prompts).

## Target Layout

### aliasops.go (~155 lines) — CRUD Operations
Stays:
- `runAliasSet()`
- `executeAliasSet()`
- `runAliasRemove()`
- `runAliasList()`
- `printAliasList()`
- `runAliasShow()`
- `isLegacyDataError()`

### aliassuggest.go (~100 lines) — Suggest Workflow
Moves:
- `runAliasSuggest()`
- `parseAliasSuggestFlags()`
- `suggestAliases()`
- `promptAliasSuggestion()`
- `createSuggestedAlias()`

Imports: `bufio`, `flag`, `fmt`, `os`, `strings`, `constants`, `store`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
