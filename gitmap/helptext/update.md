# gitmap update

Self-update gitmap from the source repository. Pulls latest, rebuilds, and deploys.

## Alias

None

## Usage

    gitmap update [--repo-path <path>] [--verbose]

## Flags

| Flag | Description |
|------|-------------|
| `--repo-path <path>` | Override the source repository path for this run |
| `--verbose` | Enable verbose logging to file |

## Prerequisites

- Git must be installed
- Source repository must be accessible

## Examples

### Example 1: Update to a newer version

    gitmap update

**Output:**

    ■ Checking for updates...
    Current version: v2.19.0
    Latest version:  v2.22.0
    v2.19.0 → v2.22.0
    ■ Pulling latest source...
    ■ Building gitmap.exe...
    ■ Deploying to E:\bin-run\gitmap.exe...
    ✓ Updated to v2.22.0
    → Run 'gitmap changelog --latest' to see what's new

### Example 2: Already up to date

    gitmap update

**Output:**

    ■ Checking for updates...
    Current version: v2.22.0
    Latest version:  v2.22.0
    ✓ Already up to date (v2.22.0)

### Example 3: Update with custom repo path

    gitmap update --repo-path C:\Projects\git-repo-navigator

**Output:**

    → Repo path: C:\Projects\git-repo-navigator
    ■ Pulling latest source...
    ■ Building gitmap.exe...
    ✓ Updated to v2.49.1

### Example 4: Update with network error

    gitmap update

**Output:**

    ■ Checking for updates...
    ✗ Failed to pull latest: network timeout
    → Check your internet connection and try again

## See Also

- [version](version.md) — Check current version
- [doctor](doctor.md) — Diagnose installation issues
- [changelog](changelog.md) — View release notes for new version
