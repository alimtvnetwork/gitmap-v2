# Project Overview

## What is gitmap?

gitmap is a portable Go CLI tool that scans directory trees for Git repositories, extracts clone URLs and branch information, and outputs structured data in multiple formats. It can also re-clone repositories from that data, preserving the original folder hierarchy.

## Current Version

**v2.1.0** (defined in `gitmap/constants/constants.go`)

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
| `spec/02-general/` | Reusable design patterns & guidelines (generic, shareable) |
| `src/` | React frontend (documentation site) |
| `.lovable/memory/` | AI memory and tracking |

## CLI Commands

| Command | Alias | Description | Status |
|---------|-------|-------------|--------|
| `scan [dir]` | `s` | Scan directory for Git repos, output all formats | ✅ Done |
| `clone <source>` | `c` | Re-clone from CSV/JSON/text preserving hierarchy | ✅ Done |
| `pull <name>` | `p` | Pull a specific repo by name | ✅ Done |
| `rescan` | `rs` | Re-run last scan with cached flags | ✅ Done |
| `desktop-sync` | `ds` | Sync repos to GitHub Desktop from scan output | ✅ Done |
| `setup` | — | Configure Git global settings from JSON | ✅ Done |
| `status` | `st` | Show dirty/clean, ahead/behind for all repos | ✅ Done |
| `exec <args>` | `x` | Run any git command across all repos | ✅ Done |
| `release [ver]` | `r` | Create release branch, tag, push | ✅ Done |
| `release-branch` | `rb` | Complete release from existing branch | ✅ Done |
| `update` | — | Self-update via copy-and-handoff + auto-cleanup | ✅ Done |
| `update-cleanup` | — | Remove update temp files and .old backups | ✅ Done |
| `version` | `v` | Print version string and exit | ✅ Done |
| `help` | — | Show usage information | ✅ Done |

## Output Files (per scan)

All written to `gitmap-output/` inside the scanned directory:

| File | Description |
|------|-------------|
| Terminal output | Colored banner, repo list, folder tree, clone instructions |
| `gitmap.csv` | CSV with repo data |
| `gitmap.json` | JSON with repo data |
| `folder-structure.md` | Markdown folder tree |
| `clone.ps1` | PowerShell clone script with comments |
| `direct-clone.ps1` | Plain `git clone` commands (HTTPS) |
| `direct-clone-ssh.ps1` | Plain `git clone` commands (SSH) |
| `register-desktop.ps1` | GitHub Desktop registration script |

## Deploy Structure

```
bin-run/
└── gitmap/
    ├── gitmap.exe
    └── data/
        └── config.json
```

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
