# CLI Interface

## Commands

### `gitmap scan [dir]` (alias: `s`)

Scan `dir` recursively for Git repositories.
Default: current working directory.

Every scan **always produces all outputs** — terminal, CSV, JSON,
folder-structure Markdown, clone script (`clone.ps1`), and desktop
registration script (`register-desktop.ps1`) — written to a
`gitmap-output/` folder at the root of the scanned directory.

After each scan, a **`last-scan.json`** cache file is written to
`gitmap-output/` so the scan can be replayed with `gitmap rescan`.

### `gitmap clone <source|json|csv>` (alias: `c`)

Re-clone repositories from a CSV, JSON, or text file.

**Shorthands:**
- `gitmap clone json` → resolves to `./gitmap-output/gitmap.json`
- `gitmap clone csv` → resolves to `./gitmap-output/gitmap.csv`
- `gitmap clone text` → resolves to `./gitmap-output/gitmap.txt`

If the resolved file doesn't exist, an error instructs the user to run `gitmap scan` first.

### `gitmap pull <repo-name>` (alias: `p`)

Pull a specific repo by its name (slug). The name is matched
against `repoName` values in `./gitmap-output/gitmap.json`.

- **Exact match** takes priority; falls back to partial/substring match (case-insensitive).
- Lists all available repo names if no match is found.
- Supports `--verbose` for debug logging.

### `gitmap rescan` (alias: `rs`)

Re-run the last scan using cached flags from `gitmap-output/last-scan.json`.
No flags are needed — all options from the previous scan are replayed exactly.

If no previous scan cache exists, an error instructs the user to run `gitmap scan` first.

### `gitmap update`

Self-update gitmap by pulling latest source and rebuilding. The binary
embeds the repo path at build time (via `-ldflags`). When invoked:

1. Copies itself to a temporary file (`gitmap-update-<pid>.exe`).
2. Launches the copy with `update --from-copy`.
3. The original process **exits immediately**, releasing the file lock.
4. The copy spawns a temporary PowerShell script that:
   - Changes to the embedded source repo directory.
   - Waits briefly for the parent to fully exit.
   - Runs `run.ps1` (pull → build → deploy).
   - Prints the new version on completion.

This two-step handoff ensures the deploy step can overwrite `gitmap.exe`
without encountering a "file in use" lock.

### `gitmap desktop-sync` (alias: `ds`)

Sync previously scanned repos to GitHub Desktop without re-scanning.
Reads from `./gitmap-output/gitmap.json` in the current directory.

- Validates output directory and JSON file exist.
- Checks GitHub Desktop CLI is installed.
- Skips repos whose paths no longer exist on disk.
- Logs per-repo success/skip/failure and prints a summary.

### `gitmap setup` (no alias)

Configure Git global settings — diff/merge tools, aliases, credential
helper, and core options — from a JSON config file.

- Reads `./data/git-setup.json` by default (override with `--config`).
- Compares each setting against the current `git config --global` value.
- Only applies settings that differ; unchanged values are skipped.
- Supports `--dry-run` to preview changes without writing anything.
- Color-coded output: ✓ applied, ⊘ unchanged, ✗ failed.

**`git-setup.json` format:**

```json
{
  "diffTool": {
    "name": "vscode",
    "cmd": "code --wait --diff $LOCAL $REMOTE"
  },
  "mergeTool": {
    "name": "vscode",
    "cmd": "code --wait $MERGED"
  },
  "aliases": {
    "co": "checkout",
    "st": "status",
    "br": "branch",
    "lg": "log --oneline --graph --all"
  },
  "credentialHelper": "manager",
  "core": {
    "autocrlf": "true",
    "longpaths": "true",
    "editor": "code --wait"
  }
}
```

Each top-level key maps to a section header in the output. All fields
are optional — omit a section to leave those settings untouched.

### `gitmap status` (alias: `st`)

Show a live dashboard of all scanned repos with current branch,
dirty/clean state, ahead/behind counts, stash entries, and file
change breakdown (staged/modified/untracked). Reads from
`./gitmap-output/gitmap.json`.

### `gitmap exec <git-args...>` (alias: `x`)

Run any git command across all repos from `./gitmap-output/gitmap.json`.
Arguments after `exec` are passed directly to `git` inside each repo directory.

- Skips repos whose paths no longer exist on disk.
- Shows per-repo success/failure with captured output.
- Prints a summary of succeeded/failed/missing counts.

### `gitmap version` (alias: `v`)

Prints the current version number (e.g., `gitmap v1.9.0`) and exits.

### `gitmap help`

Display usage information for all commands and flags.

---

## Command Aliases

All aliases are single-letter or short abbreviations for faster usage:

| Command          | Alias |
|------------------|-------|
| `scan`           | `s`   |
| `clone`          | `c`   |
| `pull`           | `p`   |
| `rescan`         | `rs`  |
| `desktop-sync`   | `ds`  |
| `version`        | `v`   |

---

## Auto Safe-Pull

When running `gitmap clone`, the tool automatically detects whether any
target directories already contain Git repositories. If existing repos
are found **and `--safe-pull` was not explicitly passed**, safe-pull is
enabled automatically and a message is printed:

```
Existing repos detected — safe-pull enabled automatically.
```

**Safe-pull behavior:**

1. Runs `git pull --ff-only` inside the existing repo directory.
2. On failure, retries up to **4 times** with a 600 ms delay between attempts.
3. On Windows, attempts to clear read-only file attributes on files
   reported in `unable to unlink` errors before retrying.
4. After all retries, produces a **diagnosis** covering:
   - File lock / read-only attribute issues
   - Windows path length risks (paths ≥ 240 characters)
   - OneDrive sync folder detection
5. When `--verbose` is enabled, every attempt, its stdout/stderr output,
   and the diagnosis are logged to a timestamped file in `gitmap-output/`.

This means users never need to remember to pass `--safe-pull` — it
activates whenever existing repos are detected during a clone operation.

---

## Scan Flags

| Flag                   | Description                          | Default              |
|------------------------|--------------------------------------|----------------------|
| `--config <path>`      | Path to JSON config file             | `./data/config.json` |
| `--mode ssh \| https`  | Clone URL style                      | `https`              |
| `--output-path <dir>`  | Output directory                     | `gitmap-output/` in scan dir |
| `--out-file <path>`    | Exact CSV output file path           | auto                 |
| `--github-desktop`     | Add discovered repos to GitHub Desktop | `false`            |
| `--open`               | Open output folder after scan completes | `false`           |

## Clone Flags

| Flag                   | Description                          | Default |
|------------------------|--------------------------------------|---------|
| `--target-dir <path>`  | Base dir to recreate folder structure | `.`    |
| `--safe-pull`          | Pull existing repos with retry + unlock diagnostics (auto-enabled) | `false` |
| `--github-desktop`     | Add cloned repos to GitHub Desktop   | `false` |
| `--verbose`            | Write detailed debug log to a timestamped file | `false` |

## Pull Flags

| Flag                   | Description                          | Default |
|------------------------|--------------------------------------|---------|
| `--verbose`            | Write detailed debug log to a timestamped file | `false` |

## Examples

```bash
# Scan current directory — outputs terminal + CSV + JSON + folder-structure.md
gitmap scan
gitmap s             # alias

# Scan with SSH URLs
gitmap scan ./projects --mode ssh

# Scan and add repos to GitHub Desktop
gitmap scan ./projects --github-desktop

# Scan parent directory
gitmap scan ..

# Re-run the last scan with the same flags
gitmap rescan
gitmap rs            # alias

# Clone using shorthand (auto-resolves to ./gitmap-output/gitmap.json)
gitmap clone json
gitmap c json        # alias

# Clone using CSV shorthand
gitmap clone csv

# Clone from JSON, preserving folder structure
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored

# Clone with verbose logging
gitmap clone json --verbose

# Clone and register with GitHub Desktop
gitmap clone ./gitmap-output/gitmap.csv --target-dir ./restored --github-desktop

# Pull a single repo by name
gitmap pull my-api-service
gitmap p my-api      # partial match works

# Sync existing scan output to GitHub Desktop
gitmap desktop-sync
gitmap ds            # alias

# Self-update from source repo
gitmap update

# Print version number
gitmap version
gitmap v             # alias
```
