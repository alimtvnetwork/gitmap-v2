# GitMap

[![CI](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml/badge.svg)](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml)
[![golangci-lint](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml)
[![GitHub Release](https://img.shields.io/github/v/release/alimtvnetwork/git-repo-navigator?style=flat-square&label=version)](https://github.com/alimtvnetwork/git-repo-navigator/releases)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey?style=flat-square)](https://github.com/alimtvnetwork/git-repo-navigator)
[![License](https://img.shields.io/badge/license-proprietary-red?style=flat-square)](./LICENSE)

A command-line tool that scans directory trees for Git repositories, extracts clone URLs and branch info, and outputs structured data. Every scan produces **all outputs** automatically:

- **Terminal** — formatted table to stdout
- **CSV** — `gitmap.csv`
- **JSON** — `gitmap.json`
- **Folder Structure** — `folder-structure.md` (tree view of discovered repos)

All files are written to a `gitmap-output/` folder at the root of the scanned directory.

**→ [GitMap Documentation](./gitmap/README.md)**
**→ [Specifications](./spec/01-app/)**

## Installation

### One-Liner Install (recommended)

**Windows (PowerShell — full bootstrap, works on any machine):**

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.ps1'))
```

**Windows (short form, PowerShell 5+):**

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.ps1 | iex
```

**Linux / macOS (Bash):**

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.sh | sh
```

Options: pin a version with `-Version v2.48.3`, choose install dir with `-InstallDir C:\tools\gitmap`, or skip PATH with `-NoPath`.

### Clone & Setup (Development)

```bash
git clone https://github.com/alimtvnetwork/git-repo-navigator.git
cd git-repo-navigator
./setup.sh
```

The setup script installs the pre-commit hook (golangci-lint), verifies your Go toolchain, and downloads dependencies.

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make all` | Lint → Test → Build (default) |
| `make setup` | Install hooks and dev tools |
| `make lint` | Run golangci-lint |
| `make vet` | Run go vet |
| `make test` | Run all tests |
| `make build` | Compile for current platform |
| `make vulncheck` | Scan dependencies for CVEs |
| `make release BUMP=patch` | Lint, test, then release (patch/minor/major) |
| `make release-dry` | Preview release without executing |
| `make clean` | Remove build artifacts |

### Build from Source (manual)

```bash
cd gitmap
go build -o ../gitmap .
```

### Build via run.ps1 (Windows)

```powershell
# Full pipeline: pull, build, deploy
.\run.ps1

# Build + scan parent folder
.\run.ps1 -R scan

# Build + scan with SSH mode
.\run.ps1 -R scan D:\repos --mode ssh
```

---

## Command Reference

Every command supports `--help` (or `-h`) for detailed usage with examples:

```bash
gitmap scan --help
gitmap cd -h
```

### Scanning & Cloning

| Command | Alias | Description |
|---------|-------|-------------|
| `scan` | `s` | Scan directory for Git repos |
| `clone` | `c` | Re-clone repos from structured file |
| `pull` | `p` | Pull a specific repo by name |
| `rescan` | `rs` | Re-scan previously scanned directories |
| `desktop-sync` | `ds` | Sync tracked repos with GitHub Desktop |

```bash
# Scan a directory and output as JSON with SSH URLs
gitmap scan ~/projects --output json --mode ssh

# Clone from JSON output into a target directory
gitmap clone json --target-dir ./restored

# Pull all repos in a group
gitmap pull --group work --all
```

→ Full details: [scan](gitmap/helptext/scan.md) · [clone](gitmap/helptext/clone.md) · [pull](gitmap/helptext/pull.md) · [rescan](gitmap/helptext/rescan.md) · [desktop-sync](gitmap/helptext/desktop-sync.md)

### Monitoring & Status

| Command | Alias | Description |
|---------|-------|-------------|
| `status` | `st` | Show repo status dashboard |
| `watch` | `w` | Live-refresh repo status dashboard |
| `exec` | `x` | Run git command across all repos |
| `latest-branch` | `lb` | Find most recently updated remote branch |

```bash
# Watch repos with 10s refresh
gitmap watch --interval 10 --group work

# Run git fetch --prune across all repos
gitmap exec fetch --prune

# Show top 5 recently updated branches
gitmap lb 5 --format csv
```

→ Full details: [status](gitmap/helptext/status.md) · [watch](gitmap/helptext/watch.md) · [exec](gitmap/helptext/exec.md) · [latest-branch](gitmap/helptext/latest-branch.md)

### Release & Versioning

| Command | Alias | Description |
|---------|-------|-------------|
| `release` | `r` | Create release branch, tag, and push |
| `release-branch` | `rb` | Create a release branch without tagging |
| `release-pending` | `rp` | Show unreleased commits since last tag |
| `changelog` | `cl` | Show release notes |
| `list-versions` | `lv` | List all available Git release tags |
| `list-releases` | `lr` | List release metadata from the database |
| `revert` | — | Revert to a specific release version |

```bash
# Auto-bump patch version and release
gitmap release --bump patch

# Preview release without executing
gitmap release --bump minor --dry-run

# Show unreleased commits
gitmap release-pending

# List all versions
gitmap list-versions --json --limit 5
```

→ Full details: [release](gitmap/helptext/release.md) · [release-branch](gitmap/helptext/release-branch.md) · [release-pending](gitmap/helptext/release-pending.md) · [changelog](gitmap/helptext/changelog.md) · [list-versions](gitmap/helptext/list-versions.md) · [list-releases](gitmap/helptext/list-releases.md) · [revert](gitmap/helptext/revert.md)

### Navigation & Organization

| Command | Alias | Description |
|---------|-------|-------------|
| `cd` | `go` | Navigate to a tracked repo directory |
| `list` | `ls` | Show all tracked repos with slugs (supports type filtering) |
| `group` | `g` | Manage repo groups / activate a group for batch ops |
| `multi-group` | `mg` | Select multiple groups for batch operations |
| `diff-profiles` | `dp` | Compare repos across two profiles |

```bash
# Navigate to a repo
gitmap cd my-api

# Interactive repo picker filtered by group
gitmap cd repos --group work

# Create a group and add repos
gitmap group create work --desc "Work repos"
gitmap group add work my-api web-app

# Activate a group for batch operations
gitmap g work                  # activate
gitmap g pull                  # pull all repos in active group
gitmap g status                # status for active group
gitmap g exec fetch --prune    # exec across active group
gitmap g clear                 # deactivate

# List repos filtered by project type
gitmap ls go                   # list Go projects
gitmap ls node                 # list Node.js projects
gitmap ls groups               # list all groups

# Multi-group batch operations
gitmap mg backend,frontend     # select multiple groups
gitmap mg pull                 # pull from all selected groups
gitmap mg clear                # clear selection

# Compare two profiles
gitmap diff-profiles home work
```

→ Full details: [cd](gitmap/helptext/cd.md) · [list](gitmap/helptext/list.md) · [group](gitmap/helptext/group.md) · [multi-group](gitmap/helptext/multi-group.md) · [diff-profiles](gitmap/helptext/diff-profiles.md)

### History & Stats

| Command | Alias | Description |
|---------|-------|-------------|
| `history` | `hi` | Show CLI command execution history |
| `history-reset` | `hr` | Clear command execution history |
| `stats` | `ss` | Show aggregated usage and performance metrics |
| `amend` | `am` | Rewrite commit author info |
| `amend-list` | `al` | List previous author amendments |

```bash
# Show last 10 commands
gitmap history --limit 10

# Show usage stats as JSON
gitmap stats --json

# Amend commit author
gitmap amend --name "John Doe" --email "john@example.com"
```

→ Full details: [history](gitmap/helptext/history.md) · [history-reset](gitmap/helptext/history-reset.md) · [stats](gitmap/helptext/stats.md) · [amend](gitmap/helptext/amend.md) · [amend-list](gitmap/helptext/amend-list.md)

### Project Detection

| Command | Alias | Description |
|---------|-------|-------------|
| `go-repos` | `gr` | List detected Go projects |
| `node-repos` | `nr` | List detected Node.js projects |
| `react-repos` | `rr` | List detected React projects |
| `cpp-repos` | `cr` | List detected C++ projects |
| `csharp-repos` | `csr` | List detected C# projects |

```bash
# List all Go projects
gitmap go-repos

# List C# projects as JSON
gitmap csharp-repos --json
```

→ Full details: [go-repos](gitmap/helptext/go-repos.md) · [node-repos](gitmap/helptext/node-repos.md) · [react-repos](gitmap/helptext/react-repos.md) · [cpp-repos](gitmap/helptext/cpp-repos.md) · [csharp-repos](gitmap/helptext/csharp-repos.md)

### Data & Profiles

| Command | Alias | Description |
|---------|-------|-------------|
| `export` | `ex` | Export database to file |
| `import` | `im` | Import repos from file |
| `profile` | `pf` | Manage database profiles |
| `bookmark` | `bk` | Save and run bookmarked commands |
| `db-reset` | — | Reset the local SQLite database |

```bash
# Export and import
gitmap export
gitmap import gitmap-export.json

# Manage profiles
gitmap profile create work
gitmap profile switch work

# Save a bookmark
gitmap bookmark save "daily" "scan ~/projects --quiet"
gitmap bookmark run daily
```

→ Full details: [export](gitmap/helptext/export.md) · [import](gitmap/helptext/import.md) · [profile](gitmap/helptext/profile.md) · [bookmark](gitmap/helptext/bookmark.md) · [db-reset](gitmap/helptext/db-reset.md)

### Utilities

| Command | Alias | Description |
|---------|-------|-------------|
| `setup` | — | Interactive first-time configuration wizard |
| `doctor` | — | Diagnose PATH, deploy, and version issues |
| `update` | — | Self-update from source repo |
| `version` | `v` | Show version number |
| `seo-write` | `sw` | Auto-commit SEO messages |
| `gomod` | `gm` | Rename Go module path across repo |
| `ssh` | — | Generate and manage SSH keys for Git authentication |
| `prune` | `pr` | Delete stale release branches that have been tagged |

```bash
# Run diagnostics
gitmap doctor --fix-path

# Self-update
gitmap update

# Rename Go module path (dry-run)
gitmap gomod "github.com/neworg/project" --dry-run

# SEO writes from CSV
gitmap seo-write --csv data.csv --max-commits 5
```

→ Full details: [setup](gitmap/helptext/setup.md) · [doctor](gitmap/helptext/doctor.md) · [update](gitmap/helptext/update.md) · [version](gitmap/helptext/version.md) · [seo-write](gitmap/helptext/seo-write.md) · [gomod](gitmap/helptext/gomod.md) · [ssh](gitmap/helptext/ssh.md) · [prune](gitmap/helptext/prune.md)

### Visualization

| Command | Alias | Description |
|---------|-------|-------------|
| `dashboard` | `db` | Generate an interactive HTML dashboard for a repo |

```bash
# Generate a full dashboard
gitmap dashboard

# Last 100 commits, open in browser
gitmap db --limit 100 --open

# Commits since a date, no merges
gitmap dashboard --since 2025-01-01 --no-merges
```

→ Full details: [dashboard](gitmap/helptext/dashboard.md)

---

## Build & Deploy

The project uses a single PowerShell script (`run.ps1`) at the repo root for the full lifecycle: pull, build, deploy, and optionally run.

| Flag | Description |
|------|-------------|
| `-NoPull` | Skip `git pull` |
| `-NoDeploy` | Skip deploy step |
| `-Update` | Update mode: full pipeline with post-update validation |
| `-R` | Run gitmap after build (trailing args forwarded) |

Configuration lives in `gitmap/powershell.json`.

## Web UI Dashboard

GitMap includes a React-based documentation and dashboard UI. To run it locally:

```bash
# Install dependencies
npm install

# Start the dev server (opens at http://localhost:5173)
npm run dev
```

You can also generate a per-repo HTML dashboard directly from the CLI:

```bash
# Generate and auto-open in browser
gitmap dashboard --open

# Last 50 commits, exclude merges
gitmap db --limit 50 --no-merges --open

# Commits since a specific date
gitmap dashboard --since 2025-06-01
```

### Tech Stack

- Vite
- TypeScript
- React
- shadcn-ui
- Tailwind CSS

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)

## License

This project is proprietary software. All rights reserved.
