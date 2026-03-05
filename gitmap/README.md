# GitMap

> Scan directories for Git repositories, generate clone instructions, and re-clone them anywhere.

## Quick Start

### Build

```powershell
# From the repo root:
.\run.ps1

# Skip git pull:
.\run.ps1 -NoPull

# Build only, no deploy:
.\run.ps1 -NoPull -NoDeploy

# Deploy to custom path:
.\run.ps1 -DeployPath "D:\tools"

# Build and run immediately:
.\run.ps1 -Run
.\run.ps1 -Run -RunPath "D:\projects"
.\run.ps1 -Run -RunArgs "--mode ssh"
```

The binary and `data/` config folder are output to `./bin/`. By default, the binary is also copied to the deploy path in `powershell.json` (default: `E:\bin-run`).

### Manual Build

```bash
cd gitmap
go build -o ../bin/gitmap.exe .
```

---

## Usage

### Scan a directory

Every scan **always produces all outputs** — terminal, CSV, JSON, and a folder structure Markdown file. They are written to a `gitmap-output/` folder at the root of the scanned directory.

```bash
# Scan current directory (outputs everything to ./gitmap-output/)
gitmap scan

# Scan a specific folder with SSH URLs
gitmap scan ./projects --mode ssh

# Scan and add repos to GitHub Desktop
gitmap scan ./projects --github-desktop

# Custom output directory
gitmap scan ./projects --output-path ./my-exports
```

### Output files

When you run `gitmap scan ./projects`, the following is created:

```
projects/
└── gitmap-output/
    ├── gitmap.csv              # All repos in CSV format
    ├── gitmap.json             # All repos in JSON format
    └── folder-structure.md     # Tree view of repo hierarchy
```

The **folder-structure.md** shows a visual tree of all discovered repos:

```
# Folder Structure

Git repositories discovered by gitmap.

├── 📦 **my-app** (`main`) — https://github.com/user/my-app.git
├── libs/
│   ├── 📦 **core-lib** (`develop`) — https://github.com/user/core-lib.git
│   └── 📦 **utils** (`main`) — https://github.com/user/utils.git
└── 📦 **docs** (`main`) — https://github.com/user/docs.git
```

### Output path behavior

| Flag | Behavior |
|------|----------|
| No flags | Creates `gitmap-output/` inside the scanned directory |
| `--output-path ./exports` | Writes to `./exports/` |
| `--out-file report.csv` | Overrides CSV file path only |

### Clone from a previous scan

```bash
# Clone from JSON (preserves original folder structure)
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored

# Clone from CSV
gitmap clone ./gitmap-output/gitmap.csv --target-dir ./restored
```

The clone command recreates the exact folder hierarchy from the `relativePath` field in each record.

---

## Configuration

### `data/config.json`

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": "./gitmap-output",
  "excludeDirs": [".cache", "node_modules", "vendor", ".venv"],
  "notes": ""
}
```

CLI flags override config values.

### `powershell.json`

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "gitmap.exe",
  "copyData": true
}
```

---

## CLI Reference

### `gitmap scan [dir]`

| Flag | Description | Default |
|------|-------------|---------|
| `--config <path>` | Config file path | `./data/config.json` |
| `--mode ssh\|https` | Clone URL style | `https` |
| `--output-path <dir>` | Output directory | `gitmap-output/` in scan dir |
| `--out-file <path>` | Exact CSV output file path | — |
| `--github-desktop` | Add discovered repos to GitHub Desktop | `false` |

### `gitmap clone <source>`

| Flag | Description | Default |
|------|-------------|---------|
| `--target-dir <path>` | Base clone directory | `.` |

---

## CSV Output Columns

`repoName, httpsUrl, sshUrl, branch, relativePath, absolutePath, cloneInstruction, notes`

---

## Project Structure

```
gitmap/
├── main.go              # Entry point
├── cmd/                  # CLI commands
│   ├── root.go           # Routing & flags
│   ├── scan.go           # Scan command
│   └── clone.go          # Clone command
├── config/               # Config loading
├── constants/            # All shared string literals
├── scanner/              # Directory walking
├── gitutil/              # Git command wrappers
├── mapper/               # Record building
├── formatter/            # Output (terminal, CSV, JSON, folder structure)
│   ├── terminal.go
│   ├── csv.go
│   ├── json.go
│   └── structure.go      # Folder tree Markdown
├── desktop/              # GitHub Desktop integration
├── cloner/               # Re-clone logic
├── model/                # Data structures
├── data/                 # Default config
│   └── config.json
├── powershell.json       # Build/deploy config
└── go.mod
```

## Specs

See [spec/01-app/](../spec/01-app/) for detailed specifications.
