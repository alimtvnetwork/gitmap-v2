# CLI Interface

## Commands

### `gitmap scan [dir]`

Scan `dir` recursively for Git repositories.
Default: current working directory.

Every scan **always produces all outputs** — terminal, CSV, JSON,
folder-structure Markdown, clone script (`clone.ps1`), and desktop
registration script (`register-desktop.ps1`) — written to a
`gitmap-output/` folder at the root of the scanned directory.

### `gitmap clone <source-file>`

Re-clone repositories from a CSV, JSON, or text file.

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
```
