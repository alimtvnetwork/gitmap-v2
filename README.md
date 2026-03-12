# GitMap

A command-line tool that scans directory trees for Git repositories, extracts clone URLs and branch info, and outputs structured data. Every scan produces **all outputs** automatically:

- **Terminal** — formatted table to stdout
- **CSV** — `gitmap.csv`
- **JSON** — `gitmap.json`
- **Folder Structure** — `folder-structure.md` (tree view of discovered repos)

All files are written to a `gitmap-output/` folder at the root of the scanned directory.

**→ [GitMap Documentation](./gitmap/README.md)**
**→ [Specifications](./spec/01-app/)**

## Quick Start

```powershell
# From the repo root:
.\run.ps1                # Pull, build, deploy to E:\bin-run
.\run.ps1 -R scan        # Build + scan parent folder
.\run.ps1 -R scan D:\repos --mode ssh   # Build + scan with SSH mode
```

Or manually:

```bash
cd gitmap
go build -o ../bin/gitmap.exe .
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
| `list` | `ls` | Show all tracked repos with slugs |
| `group` | `g` | Manage repo groups |
| `diff-profiles` | `dp` | Compare repos across two profiles |

```bash
# Navigate to a repo
gitmap cd my-api

# Interactive repo picker filtered by group
gitmap cd repos --group work

# Create a group and add repos
gitmap group create work --desc "Work repos"
gitmap group add work my-api web-app

# Compare two profiles
gitmap diff-profiles home work
```

→ Full details: [cd](gitmap/helptext/cd.md) · [list](gitmap/helptext/list.md) · [group](gitmap/helptext/group.md) · [diff-profiles](gitmap/helptext/diff-profiles.md)

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

→ Full details: [go-repos](gitmap/help/go-repos.md) · [node-repos](gitmap/help/node-repos.md) · [react-repos](gitmap/help/react-repos.md) · [cpp-repos](gitmap/help/cpp-repos.md) · [csharp-repos](gitmap/help/csharp-repos.md)

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

→ Full details: [export](gitmap/help/export.md) · [import](gitmap/help/import.md) · [profile](gitmap/help/profile.md) · [bookmark](gitmap/help/bookmark.md) · [db-reset](gitmap/help/db-reset.md)

### Utilities

| Command | Alias | Description |
|---------|-------|-------------|
| `setup` | — | Interactive first-time configuration wizard |
| `doctor` | — | Diagnose PATH, deploy, and version issues |
| `update` | — | Self-update from source repo |
| `version` | `v` | Show version number |
| `seo-write` | `sw` | Auto-commit SEO messages |
| `gomod` | `gm` | Rename Go module path across repo |

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

→ Full details: [setup](gitmap/help/setup.md) · [doctor](gitmap/help/doctor.md) · [update](gitmap/help/update.md) · [version](gitmap/help/version.md) · [seo-write](gitmap/help/seo-write.md) · [gomod](gitmap/help/gomod.md)

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

## Web Frontend (React + Vite)

A React + TypeScript + Tailwind CSS web application scaffold.

```sh
npm i
npm run dev
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
