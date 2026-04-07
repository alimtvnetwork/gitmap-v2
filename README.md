# GitMap

[![CI](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml/badge.svg)](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml)
[![golangci-lint](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/ci.yml)
[![Vulncheck](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/vulncheck.yml/badge.svg)](https://github.com/alimtvnetwork/git-repo-navigator/actions/workflows/vulncheck.yml)
[![GitHub Release](https://img.shields.io/github/v/release/alimtvnetwork/git-repo-navigator?style=flat-square&label=version)](https://github.com/alimtvnetwork/git-repo-navigator/releases)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey?style=flat-square)](https://github.com/alimtvnetwork/git-repo-navigator)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](./LICENSE)

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

### Installer Options

| Flag | Description | Example |
|------|-------------|---------|
| `-Version` | Pin a specific release | `-Version v2.49.1` |
| `-InstallDir` | Custom install directory | `-InstallDir C:\tools\gitmap` |
| `-Arch` | Force architecture (`amd64`, `arm64`) | `-Arch arm64` |
| `-NoPath` | Skip adding to user PATH | `-NoPath` |

**Custom directory install (one-liner):**

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.ps1 | iex; Install-Gitmap -InstallDir "D:\DevTools\gitmap"
```

**Custom directory install (downloaded script):**

```powershell
.\install.ps1 -InstallDir "D:\DevTools\gitmap"
```

**Pinned version + custom directory:**

```powershell
.\install.ps1 -Version v2.49.1 -InstallDir "C:\tools\gitmap"
```

> **Tip for other installers:** Use `-InstallDir` and `-NoPath` together to integrate gitmap into your own package layout without modifying the user's PATH.

### Clone & Setup (Development)

```bash
git clone https://github.com/alimtvnetwork/git-repo-navigator.git
cd git-repo-navigator
./setup.sh
```

The setup script installs the pre-commit hook (golangci-lint), verifies your Go toolchain, and downloads dependencies. See [CONTRIBUTING.md](CONTRIBUTING.md) for the full development workflow, coding standards, and PR requirements.

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

#### Release Workflow

| Command | Alias | Description |
|---------|-------|-------------|
| `release` | `r` | Create release branch, tag, and push |
| `release-self` | `rs` | Release gitmap itself from any directory |
| `release-branch` | `rb` | Create a release branch without tagging |

```bash
# Auto-bump patch version and release
gitmap release --bump patch

# Release with Go binary cross-compilation
gitmap release --bump minor --bin

# Preview release without executing
gitmap release --bump minor --dry-run

# Release with notes, compressed assets, and checksums
gitmap release v3.0.0 --bin --compress --checksums -N "Major redesign"

# Release gitmap itself from any directory
gitmap release-self --bump patch

# Create a release branch without tagging
gitmap release-branch v3.0.0-rc1 --branch develop
```

#### Release History & Info

| Command | Alias | Description |
|---------|-------|-------------|
| `changelog` | `cl` | Show release notes |
| `list-versions` | `lv` | List all available Git release tags |
| `list-releases` | `lr` | List release metadata from the database |
| `release-pending` | `rp` | Show unreleased commits since last tag |
| `revert` | — | Revert to a specific release version |
| `clear-release-json` | `crj` | Remove orphaned release metadata files |

```bash
# View changelog for a version
gitmap changelog v2.49.0

# Show unreleased commits
gitmap release-pending

# List all versions as JSON (last 5)
gitmap list-versions --json --limit 5

# List stored release records
gitmap list-releases --limit 10

# Revert to a previous release
gitmap revert v2.48.0

# Clean up orphaned release metadata
gitmap clear-release-json v2.30.0
```

> **CI Pipeline:** Pushing a `release/*` branch or `v*` tag triggers GitHub Actions to cross-compile all 6 targets, generate checksums, attach a version-pinned install script, and create a GitHub release with changelog, metadata, and install instructions. See [spec/01-app/12-release-command.md](spec/01-app/12-release-command.md#ci-release-pipeline-github-actions) for details.

→ Full details: [release](gitmap/helptext/release.md) · [release-self](gitmap/helptext/release-self.md) · [release-branch](gitmap/helptext/release-branch.md) · [release-pending](gitmap/helptext/release-pending.md) · [changelog](gitmap/helptext/changelog.md) · [list-versions](gitmap/helptext/list-versions.md) · [list-releases](gitmap/helptext/list-releases.md) · [revert](gitmap/helptext/revert.md) · [clear-release-json](gitmap/helptext/clear-release-json.md)

### Navigation & Organization

| Command | Alias | Description |
|---------|-------|-------------|
| `cd` | `go` | Navigate to a tracked repo directory |
| `list` | `ls` | Show all tracked repos with slugs (supports type filtering) |
| `group` | `g` | Manage repo groups / activate a group for batch ops |
| `multi-group` | `mg` | Select multiple groups for batch operations |
| `alias` | `a` | Assign short names to repos for quick access |
| `diff-profiles` | `dp` | Compare repos across two profiles |

```bash
# Navigate to a repo
gitmap cd my-api

# Interactive repo picker filtered by group
gitmap cd repos --group work

# Set a default path for a repo name (skip picker)
gitmap cd set-default my-api D:\repos\api-gateway
gitmap cd clear-default my-api

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

# Create and use repo aliases
gitmap alias set api github/user/api-gateway
gitmap pull -A api             # pull using alias
gitmap cd -A api               # navigate using alias
gitmap alias suggest           # auto-suggest aliases for unaliased repos
gitmap alias suggest --apply   # auto-accept all suggestions
gitmap alias list              # list all aliases with paths

# Compare two profiles
gitmap diff-profiles home work
```

→ Full details: [cd](gitmap/helptext/cd.md) · [list](gitmap/helptext/list.md) · [group](gitmap/helptext/group.md) · [multi-group](gitmap/helptext/multi-group.md) · [alias](gitmap/helptext/alias.md) · [diff-profiles](gitmap/helptext/diff-profiles.md)

### History, Stats & Amend

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

# Amend commit author (all commits)
gitmap amend --name "John Doe" --email "john@example.com"

# Amend a specific commit
gitmap amend abc123 --name "John Doe" --email "john@example.com"

# Dry-run (preview without changing)
gitmap amend --name "John" --email "john@co.com" --dry-run

# Amend and force push
gitmap amend --name "John" --email "john@co.com" --force-push

# View amendment history
gitmap amend-list --json --limit 5
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

### Data, Profiles & Bookmarks

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

# Save a bookmark and run it later
gitmap bookmark save daily scan ~/projects
gitmap bookmark run daily

# List and delete bookmarks
gitmap bookmark list --json
gitmap bookmark delete daily
```

→ Full details: [export](gitmap/helptext/export.md) · [import](gitmap/helptext/import.md) · [profile](gitmap/helptext/profile.md) · [bookmark](gitmap/helptext/bookmark.md) · [db-reset](gitmap/helptext/db-reset.md)

### Cloning & Versioned Repos

| Command | Alias | Description |
|---------|-------|-------------|
| `clone` | `c` | Re-clone repos from structured file |
| `clone-next` | `cn` | Clone next versioned iteration of current repo |

```bash
# Clone from JSON output into a target directory
gitmap clone json --target-dir ./restored

# Clone using a specific SSH key
gitmap clone repos.json --ssh-key work

# Increment repo version by one (my-app-v3 → my-app-v4)
gitmap cn v++

# Jump to a specific version and auto-delete current folder
gitmap cn v15 --delete

# Keep the current folder without prompting
gitmap cn v++ --keep

# Clone with remote creation if missing
gitmap cn v++ --create-remote
```

→ Full details: [clone](gitmap/helptext/clone.md) · [clone-next](gitmap/helptext/clone-next.md)

### SSH Key Management

| Command | Alias | Description |
|---------|-------|-------------|
| `ssh` | — | Generate and manage SSH keys for Git authentication |

```bash
# Generate a new SSH key
gitmap ssh --name work --path ~/.ssh/id_rsa_work

# Display the public key
gitmap ssh cat --name work

# List all stored keys
gitmap ssh list

# Clone repos using a named SSH key
gitmap clone repos.json --ssh-key work

# Regenerate ~/.ssh/config managed entries
gitmap ssh config
```

→ Full details: [ssh](gitmap/helptext/ssh.md)

### Zip Groups (Release Archives)

| Command | Alias | Description |
|---------|-------|-------------|
| `zip-group` | `z` | Manage named collections of files for release archives |

```bash
# Create a zip group with paths
gitmap z create "chrome extension" chrome-extension/dist

# Create and add items separately
gitmap z create docs-bundle
gitmap z add docs-bundle ./README.md ./CHANGELOG.md ./docs/

# Set custom archive name
gitmap z create extras --archive extra-files.zip
gitmap z add extras ./config/ ./scripts/deploy.sh

# List all zip groups
gitmap z list

# Show items with dynamic folder expansion
gitmap z show docs-bundle

# Include in a release
gitmap release v3.0.0 --zip-group docs-bundle
```

→ Full details: [zip-group](gitmap/helptext/zip-group.md)

### Environment & Tool Installation

| Command | Alias | Description |
|---------|-------|-------------|
| `env` | `ev` | Manage persistent environment variables and PATH |
| `install` | `in` | Install developer tools via platform package manager |

```bash
# Set and retrieve environment variables
gitmap env set GOPATH "/home/user/go"
gitmap env get GOPATH
gitmap env list

# Manage PATH entries
gitmap env path add /usr/local/go/bin
gitmap env path remove /usr/local/go/bin
gitmap env path list

# Preview changes
gitmap env set NODE_ENV production --dry-run

# Install developer tools
gitmap install node
gitmap install go --check          # check if installed
gitmap install python --dry-run    # preview install command
gitmap install vscode --manager winget
gitmap install --list              # list all supported tools
```

→ Full details: [env](gitmap/helptext/env.md) · [install](gitmap/helptext/install.md)

### Temp Releases & File-Sync Tasks

| Command | Alias | Description |
|---------|-------|-------------|
| `temp-release` | `tr` | Create lightweight temp release branches |
| `task` | `tk` | Manage file-sync watch tasks |

```bash
# Create 10 temp release branches from last 10 commits
gitmap tr 10 v1.$$ -s 5

# Preview without creating
gitmap tr 5 v2.$$$ --dry-run

# List and clean up temp release branches
gitmap tr list
gitmap tr remove v1.05 to v1.10
gitmap tr remove all

# Create a file-sync task
gitmap task create my-sync --src ./src --dest ./backup

# Run sync task with verbose output
gitmap tk run my-sync --interval 10 --verbose

# List and manage tasks
gitmap task list
gitmap task delete my-sync
```

→ Full details: [temp-release](gitmap/helptext/temp-release.md) · [task](gitmap/helptext/task.md)

### Utilities

| Command | Alias | Description |
|---------|-------|-------------|
| `setup` | — | Interactive first-time configuration wizard |
| `doctor` | — | Diagnose PATH, deploy, and version issues |
| `update` | — | Self-update from source repo or via gitmap-updater |
| `version` | `v` | Show version number |
| `completion` | `cmp` | Generate shell tab-completion scripts |
| `interactive` | `i` | Launch full-screen interactive TUI |
| `has-any-updates` | `hau` | Check if remote has new commits |
| `docs` | `d` | Open documentation website in browser |
| `seo-write` | `sw` | Auto-commit SEO messages |
| `gomod` | `gm` | Rename Go module path across repo |
| `changelog-generate` | `cg` | Auto-generate changelog from commits |
| `prune` | `pr` | Delete stale release branches that have been tagged |

```bash
# Run diagnostics
gitmap doctor --fix-path

# Self-update (from source or via gitmap-updater fallback)
gitmap update
gitmap update --repo-path C:\gitmap-src

# Generate shell completions
gitmap completion powershell    # output script
gitmap completion bash
gitmap completion --list-repos  # for scripting

# Launch interactive TUI
gitmap interactive
gitmap i --refresh 10

# Check for remote updates
gitmap hau

# Auto-generate changelog
gitmap changelog-generate                       # between latest two tags
gitmap cg --from v2.22.0                        # from tag to HEAD
gitmap cg --from v2.23.0 --to v2.24.0 --write  # write to CHANGELOG.md

# Rename Go module path (dry-run)
gitmap gomod "github.com/neworg/project" --dry-run

# SEO writes from CSV
gitmap seo-write --csv data.csv --max-commits 5
```

→ Full details: [setup](gitmap/helptext/setup.md) · [doctor](gitmap/helptext/doctor.md) · [update](gitmap/helptext/update.md) · [version](gitmap/helptext/version.md) · [completion](gitmap/helptext/completion.md) · [interactive](gitmap/helptext/interactive.md) · [has-any-updates](gitmap/helptext/has-any-updates.md) · [docs](gitmap/helptext/docs.md) · [seo-write](gitmap/helptext/seo-write.md) · [gomod](gitmap/helptext/gomod.md) · [changelog-generate](gitmap/helptext/changelog-generate.md) · [prune](gitmap/helptext/prune.md)

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

## Author

### Md. Alim Ul Karim

**Creator & Lead Architect** | Chief Software Engineer, Riseup Asia LLC

A system architect with **20+ years** of professional software engineering experience across enterprise, fintech, and distributed systems. His technology stack spans **.NET/C# (18+ years)**, **JavaScript (10+ years)**, **TypeScript (6+ years)**, and **Golang (4+ years)**.

Recognized as a **top 1% talent at Crossover** and one of the top software architects globally. He is also the **CEO of Riseup Asia LLC** and maintains an active presence on **Stack Overflow** (2,452+ reputation, member since 2010) and **LinkedIn** (12,500+ followers).

His architectural philosophy — _consistency over cleverness, convention over configuration_ — is the driving force behind every design decision in this project. His published writings on clean function design and meaningful naming directly inform the coding principles encoded in this specification system.

|  |  |
|---|---|
| **Website** | [alimkarim.com](https://alimkarim.com/) · [my.alimkarim.com](https://my.alimkarim.com/) |
| **LinkedIn** | [linkedin.com/in/alimkarim](https://linkedin.com/in/alimkarim) |
| **Google** | [Alim Ul Karim](https://www.google.com/search?q=Alim+Ul+Karim) |
| **Role** | Chief Software Engineer, Riseup Asia LLC |

### Riseup Asia LLC

Top Leading Software Company in WY (2026)

| | |
|---|---|
| **Website** | [riseup-asia.com](https://riseup-asia.com) |
| **Facebook** | [riseupasia.talent](https://www.facebook.com/riseupasia.talent/) |
| **LinkedIn** | [Riseup Asia](https://www.linkedin.com/company/105304484/) |
| **YouTube** | [@riseup-asia](https://www.youtube.com/@riseup-asia) |

## License

This project is licensed under the [MIT License](./LICENSE).
