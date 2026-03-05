# Project Overview

## What is gitmap?

gitmap is a portable Go CLI tool that scans directory trees for Git repositories, extracts clone URLs and branch information, and outputs structured data in multiple formats. It can also re-clone repositories from that data, preserving the original folder hierarchy.

## Current Version

**v1.1.2** (defined in `gitmap/constants/constants.go`)

## Tech Stack

- **CLI**: Go (compiled to `gitmap.exe`)
- **Build/Deploy**: PowerShell (`run.ps1`)
- **Frontend**: React + Vite + Tailwind (documentation site, currently placeholder)
- **Config**: JSON (`powershell.json`, `data/config.json`)

## Key Directories

| Directory | Purpose |
|-----------|---------|
| `gitmap/` | Go source code for the CLI |
| `spec/01-app/` | App-specific specification documents |
| `spec/02-general/` | Reusable design patterns & guidelines (AI-trainable) |
| `src/` | React frontend (documentation site) |
| `.lovable/memory/` | AI memory and tracking |

## CLI Commands

| Command | Description | Status |
|---------|-------------|--------|
| `scan [dir]` | Scan directory for Git repos, output all formats | ✅ Done |
| `clone <source>` | Re-clone from CSV/JSON/text preserving hierarchy | ✅ Done |
| `update` | Self-update via copy-and-handoff mechanism | ✅ Done |
| `version` | Print version string and exit | ✅ Done |
| `help` | Show usage information | ✅ Done |

## Output Files (per scan)

All written to `gitmap-output/` inside the scanned directory:

| File | Description |
|------|-------------|
| Terminal output | Colored banner, repo list, folder tree |
| `gitmap.csv` | CSV with repo data |
| `gitmap.json` | JSON with repo data |
| `folder-structure.md` | Markdown folder tree |
| `clone.ps1` | PowerShell clone script with comments |
| `direct-clone.ps1` | Plain `git clone` commands (HTTPS) |
| `direct-clone-ssh.ps1` | Plain `git clone` commands (SSH) |
| `register-desktop.ps1` | GitHub Desktop registration script |

## Code Style Rules

- No negation in `if` conditions (no `!`, no `!=`)
- Functions: 8–15 lines
- Files: 100–200 lines max
- One responsibility per package
- Blank line before `return` (unless sole line in `if` block)
- All string literals in `constants` package (no magic strings)

## Version Policy

- **Bump on every code change** that alters behavior or output
- Follows SemVer (`MAJOR.MINOR.PATCH`)
- Displayed in terminal banner, `help`, and `version` command
- `run.ps1` prints version after each build
