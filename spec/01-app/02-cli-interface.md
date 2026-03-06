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

1. Copies itself to a temporary file (`gitmap-update-<pid>.exe`) in the same directory (fallback to `%TEMP%`).
2. Launches the copy with the hidden `update-runner` command using **foreground/blocking** execution.
3. The parent waits for the worker to complete, keeping the terminal session stable.
4. The `update-runner` spawns a temporary PowerShell script that:
   - Captures the currently deployed version.
   - Runs `run.ps1 -Update` (full pipeline: pull → build → deploy with `.old` rollback backup).
   - PATH sync uses rename-first (rename active to `.old`, copy new).
   - Compares old vs new version (warns if unchanged).
   - Runs `gitmap changelog --latest` from the updated binary.
   - Runs `gitmap update-cleanup` to remove temp copies and `.old` backups.

This two-step handoff ensures the deploy step can overwrite `gitmap.exe`
without encountering a "file in use" lock (rename-first handles the locked binary).

**Critical rules:**
- Parent MUST use `cmd.Run()` (foreground/blocking), NEVER `cmd.Start()` + `os.Exit(0)` (async breaks terminal).
- PATH sync MUST use rename-first in update mode.
- Generated scripts MUST NOT contain `Read-Host` or interactive prompts.

### `gitmap update-cleanup`

Remove leftover artifacts from the update process:

- **Temp update copies** — `%TEMP%\gitmap-update-*.exe` files from
  previous copy-and-handoff operations.
- **Old backup binaries** — `*.old` files in the deploy directory
  created as rollback backups during deploy.

This command runs automatically at the end of a successful `gitmap update`,
but can also be invoked manually for ad-hoc cleanup.

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

### `gitmap release [version]` (alias: `r`)

Create a release branch, Git tag, and push to remote. Version can be
full (`v1.2.3`), partial (`v1`, `v1.2` — zero-padded), or omitted
(reads from `version.json`). Supports pre-release suffixes (`-rc.1`,
`-beta`) and draft mode.

- Checks `.release/` and Git tags to prevent duplicate releases.
- Records assets from `--assets` in release metadata.
- Writes release metadata to `.release/vX.Y.Z.json`.
- Updates `.release/latest.json` for the highest stable version.

See [12-release-command.md](./12-release-command.md) for full details.

### `gitmap release-branch <branch>` (alias: `rb`)

Complete a release from an existing `release/vX.Y.Z` branch. Creates
the tag and pushes if not already done. Useful when the release
branch was created manually or by a previous incomplete release.

### `gitmap release-pending` (alias: `rp`)

Release all `release/v*` branches that are missing tags. Scans local
branches for `release/vX.Y.Z` patterns, checks whether the
corresponding `vX.Y.Z` tag already exists, and creates+pushes tags
for any that are untagged.

- Supports `--assets`, `--draft`, `--dry-run`, and `--verbose`.
- Useful for catching up on releases after manual branch creation.

### `gitmap changelog [version]` (alias: `cl`)

Display concise, CLI-friendly release notes from `CHANGELOG.md`.

- **No args** — prints the last 5 versions (configurable via `--limit`).
- **`--latest`** — prints only the most recent version's notes.
- **`<version>`** — prints notes for a specific version (e.g., `gitmap changelog v2.3.0`).
- **`--open`** — opens `CHANGELOG.md` in the default system application.
- **`changelog.md`** (as command) — shorthand for `changelog --open`.

The `gitmap update` command automatically runs `gitmap changelog --latest`
after a successful update to show the user what changed.

### `gitmap doctor [--fix-path]` (no alias)

Diagnose environment and deployment health. Runs a series of checks
and prints `[OK]`, `[!!]`, or `[--]` for each:

1. **RepoPath embedded** — confirms binary was built with `run.ps1`.
2. **PATH binary** — finds `gitmap` on PATH and reports its location/version.
3. **Deployed binary** — reads `powershell.json` to find the deploy target.
4. **Version mismatch** — compares source, PATH, and deployed versions;
   prints exact `Copy-Item` fix commands when they differ.
5. **Git available** — checks `git --version`.
6. **Go available** — checks `go version` (warning only, needed for building).
7. **CHANGELOG.md present** — confirms changelog command will work.

If issues are found, each is accompanied by a recommended fix command.

**`--fix-path` flag:**

When passed, skips the diagnostic checks and instead directly syncs
the active PATH binary from the deployed binary. Uses a three-layer
fallback strategy:

1. **Direct copy with retries** — 20 attempts × 500ms delay.
2. **Rename fallback** — renames the locked `.exe` to `.old`, copies
   the deployed binary in its place (with rollback on failure).
3. **Stale-process termination** — finds and kills `gitmap.exe`
   processes bound to the old PATH location, then retries.

Prints clear confirmation with version verification after sync.

If issues are found, each is accompanied by a recommended fix command.

### `gitmap latest-branch` (alias: `lb`)

Find the most recently updated remote branch by commit date. Fetches
all remotes, reads tip commits, sorts by date, and resolves the branch
name via `--points-at`.

See [14-latest-branch.md](./14-latest-branch.md) for full details.

### `gitmap version` (alias: `v`)

Prints the current version number (e.g., `gitmap v2.4.0`) and exits.

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
| `status`         | `st`  |
| `exec`           | `x`   |
| `release`        | `r`   |
| `release-branch` | `rb`  |
| `release-pending`| `rp`  |
| `changelog`      | `cl`  |
| `latest-branch`  | `lb`  |
| `version`        | `v`   |
| `update`         | —     |
| `update-cleanup` | —     |
| `doctor`         | —     |

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

## Setup Flags

| Flag                   | Description                          | Default                    |
|------------------------|--------------------------------------|----------------------------|
| `--config <path>`      | Path to git-setup.json config file   | `./data/git-setup.json`    |
| `--dry-run`            | Preview changes without applying     | `false`                    |

## Release Flags

| Flag                          | Description                                      | Default |
|-------------------------------|--------------------------------------------------|---------|
| `--assets <path>`             | Directory or file to attach to the release       | (none)  |
| `--commit <sha>`              | Create release from a specific commit            | (none)  |
| `--branch <name>`             | Create release from latest commit of a branch    | (none)  |
| `--bump major\|minor\|patch`  | Auto-increment from latest released version      | (none)  |
| `--draft`                     | Create an unpublished draft release              | `false` |
| `--verbose`                   | Write detailed debug log                         | `false` |

## Release-Branch Flags

| Flag              | Description                         | Default |
|-------------------|-------------------------------------|---------|
| `--assets <path>` | Directory or file to attach         | (none)  |
| `--draft`         | Create an unpublished draft release | `false` |
| `--verbose`       | Write detailed debug log            | `false` |

## Release-Pending Flags

| Flag              | Description                              | Default |
|-------------------|------------------------------------------|---------|
| `--assets <path>` | Directory or file to attach              | (none)  |
| `--draft`         | Mark release metadata as draft           | `false` |
| `--dry-run`       | Preview steps without executing          | `false` |
| `--verbose`       | Write detailed debug log                 | `false` |

## Changelog Flags

| Flag              | Description                              | Default |
|-------------------|------------------------------------------|---------|
| `--latest`        | Show only the most recent version        | `false` |
| `--limit <n>`     | Max number of versions to display        | `5`     |
| `--open`          | Open CHANGELOG.md in default application | `false` |

## Latest-Branch Flags

| Flag                    | Description                                          | Default  |
|-------------------------|------------------------------------------------------|----------|
| `--remote <name>`       | Remote to filter branches against                    | `origin` |
| `--all-remotes`         | Include branches from all remotes                    | `false`  |
| `--contains-fallback`   | Fall back to `--contains` if `--points-at` is empty  | `false`  |
| `--top <n>`             | Show top N most recently updated branches            | `0`      |
| `--json`                | Output structured JSON instead of plain text         | `false`  |

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

# Configure Git global settings (preview first)
gitmap setup --dry-run
gitmap setup

# Show repo status dashboard
gitmap status
gitmap st            # alias

# Run git fetch across all repos
gitmap exec fetch --prune
gitmap x status -s   # alias

# Self-update from source repo
gitmap update

# Clean up leftover update artifacts manually
gitmap update-cleanup

# Create a release from HEAD
gitmap release v1.2.3
gitmap r v1.0.0      # alias

# Partial version (padded to v1.0.0)
gitmap release v1

# Release with assets
gitmap release v2.0.0 --assets ./dist

# Release from specific commit or branch
gitmap release v1.2.3 --commit abc123
gitmap release v1.0.0 --branch develop

# Auto-increment version
gitmap release --bump patch
gitmap release --bump minor --assets ./bin

# Draft / pre-release
gitmap release v3.0.0-rc.1 --draft

# Read version from version.json
gitmap release

# Complete release from existing release branch
gitmap release-branch release/v1.2.0
gitmap rb release/v1.2.0

# Release all untagged release branches
gitmap release-pending
gitmap rp            # alias
gitmap release-pending --dry-run

# View changelog
gitmap changelog             # last 5 versions
gitmap cl --latest           # most recent only
gitmap changelog v2.3.0      # specific version
gitmap changelog --open      # open CHANGELOG.md
gitmap changelog.md          # shorthand for --open

# Diagnose environment issues
gitmap doctor

# Print version number
gitmap version
gitmap v             # alias
```
