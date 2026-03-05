# CLI Interface

## Commands

### `gitmap scan [dir]`

Scan `dir` recursively for Git repositories.
Default: current working directory.

Every scan **always produces all outputs** — terminal, CSV, JSON,
folder-structure Markdown, clone script (`clone.ps1`), and desktop
registration script (`register-desktop.ps1`) — written to a
`gitmap-output/` folder at the root of the scanned directory.

### `gitmap clone <source|json|csv>`

Re-clone repositories from a CSV, JSON, or text file.

**Shorthands:**
- `gitmap clone json` → resolves to `./gitmap-output/gitmap.json`
- `gitmap clone csv` → resolves to `./gitmap-output/gitmap.csv`

If the resolved file doesn't exist, an error instructs the user to run `gitmap scan` first.

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

### `gitmap version`

Prints the current version number (e.g., `gitmap v1.1.2`) and exits.

### `gitmap desktop-sync`

Sync previously scanned repos to GitHub Desktop without re-scanning.
Reads from `./gitmap-output/gitmap.json` in the current directory.

- Validates output directory and JSON file exist.
- Checks GitHub Desktop CLI is installed.
- Skips repos whose paths no longer exist on disk.
- Logs per-repo success/skip/failure and prints a summary.

### `gitmap help`

Display usage information for all commands and flags.

## Scan Flags

| Flag                   | Description                          | Default              |
|------------------------|--------------------------------------|----------------------|
| `--config <path>`      | Path to JSON config file             | `./data/config.json` |
| `--mode ssh \| https`  | Clone URL style                      | `https`              |
| `--output-path <dir>`  | Output directory                     | `gitmap-output/` in scan dir |
| `--out-file <path>`    | Exact CSV output file path           | auto                 |
| `--github-desktop`     | Add discovered repos to GitHub Desktop | `false`            |

## Clone Flags

| Flag                   | Description                          | Default |
|------------------------|--------------------------------------|---------|
| `--target-dir <path>`  | Base dir to recreate folder structure | `.`    |
| `--github-desktop`     | Add cloned repos to GitHub Desktop   | `false` |

## Examples

```bash
# Scan current directory — outputs terminal + CSV + JSON + folder-structure.md
gitmap scan

# Scan with SSH URLs
gitmap scan ./projects --mode ssh

# Scan and add repos to GitHub Desktop
gitmap scan ./projects --github-desktop

# Scan parent directory
gitmap scan ..

# Clone from JSON, preserving folder structure
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored

# Clone and register with GitHub Desktop
gitmap clone ./gitmap-output/gitmap.csv --target-dir ./restored --github-desktop

# Sync existing scan output to GitHub Desktop
gitmap desktop-sync

# Self-update from source repo
gitmap update

# Print version number
gitmap version
```
