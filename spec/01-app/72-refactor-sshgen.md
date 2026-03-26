# Refactor: cmd/sshgen.go

## Problem
`sshgen.go` is 224 lines with two responsibilities: command orchestration (flag parsing, key generation, DB storage, user prompts) and utility functions (keygen validation, email resolution, fingerprint reading, path helpers).

## Target Layout

### sshgen.go (~160 lines) — Command Orchestration
Stays:
- `runSSHGenerate()`
- `parseSSHGenFlags()`
- `handleExistingKey()`
- `generateAndStore()`

### sshgenutil.go (~78 lines) — Utilities
Moves:
- `validateSSHKeygen()`
- `resolveGitEmail()`
- `readFingerprint()`
- `removeKeyFiles()`
- `defaultSSHKeyPath()`
- `expandHome()`
- `ensureSSHDir()`

Imports: `os`, `os/exec`, `path/filepath`, `strings`, `constants`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
