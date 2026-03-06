# Memory: features/latest-branch

The 'latest-branch' (alias 'lb') command identifies the most recently updated remote-tracking branches. It fetches remotes, sorts branches by commit date, and resolves names using 'git for-each-ref --points-at' with a '--contains-fallback' option.

## Flags
- `--remote <name>` — filter by remote (default: origin)
- `--all-remotes` — include all remotes
- `--contains-fallback` — use `git branch -r --contains` if exact SHA resolution fails
- `--top <n>` — show top N most recently updated branches
- `--json` — structured JSON output

## Positional integer shorthand
A bare integer positional argument acts as shorthand for `--top`:
`gitmap lb 3` is equivalent to `gitmap lb --top 3`.

## Files
- `cmd/latestbranch.go` — CLI handler, flag parsing, output
- `gitutil/latestbranch.go` — Git operations (fetch, list, log, resolve)
- `constants/constants.go` — command name, alias, messages, flags
- `spec/01-app/14-latest-branch.md` — full specification
