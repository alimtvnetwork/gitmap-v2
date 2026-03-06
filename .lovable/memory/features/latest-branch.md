# Memory: features/latest-branch

The 'latest-branch' (alias 'lb') command identifies the most recently updated remote-tracking branches. It fetches remotes, sorts branches by commit date, and resolves names using 'git for-each-ref --points-at' with a '--contains-fallback' option.

## Flags
- `--remote <name>` — filter by remote (default: origin)
- `--all-remotes` — include all remotes
- `--contains-fallback` — use `git branch -r --contains` if exact SHA resolution fails
- `--top <n>` — show top N most recently updated branches
- `--format <fmt>` — output format: `terminal` (default), `json`, `csv`
- `--json` — shorthand for `--format json`
- `--no-fetch` — skip `git fetch --all --prune` (use existing remote refs)

## Positional integer shorthand
A bare integer positional argument acts as shorthand for `--top`:
`gitmap lb 3` is equivalent to `gitmap lb --top 3`.

## Date display
All dates are formatted via `gitutil.FormatDisplayDate`: UTC → local timezone → `DD-Mon-YYYY hh:mm AM/PM` (e.g. `06-Mar-2026 09:17 AM`).

## Output formats
- **terminal** — human-readable key-value display + optional top-N table
- **json** — structured JSON to stdout
- **csv** — header row + data rows to stdout

## Files
- `cmd/latestbranch.go` — CLI handler, flag parsing, output
- `gitutil/latestbranch.go` — Git operations (fetch, list, log, resolve)
- `gitutil/dateformat.go` — centralized date formatting
- `constants/constants.go` — command name, alias, messages, flags, date layout
- `spec/01-app/14-latest-branch.md` — full specification
- `spec/01-app/15-date-display-format.md` — date format spec
- `spec/02-general/07-date-display-format.md` — general date pattern
