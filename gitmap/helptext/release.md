# gitmap release

Create a release: tag, push, and optionally publish a GitHub release.

## Alias

r

## Usage

    gitmap release [version] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --assets \<path\> | — | Attach files to release |
| --commit \<sha\> | HEAD | Release from specific commit |
| --branch \<name\> | current | Release from branch |
| --bump major\|minor\|patch | — | Auto-increment version |
| --draft | false | Create unpublished draft |
| --dry-run | false | Preview without executing |

## Prerequisites

- Must be inside a Git repository with at least one commit
- GitHub CLI (`gh`) recommended for publishing

## Examples

### Example 1: Release with auto-bump

    gitmap release --bump patch

**Output:**

    Current version: v2.3.9
    New version: v2.3.10
    Creating tag v2.3.10... done
    Pushing to origin... done
    ✓ Released v2.3.10

### Example 2: Dry-run preview

    gitmap r --bump minor --dry-run

**Output:**

    [DRY RUN] Current version: v2.3.10
    [DRY RUN] New version: v2.4.0
    [DRY RUN] Would create tag v2.4.0
    [DRY RUN] Would push to origin
    No changes made.

### Example 3: Release with assets

    gitmap release v3.0.0 --assets ./dist/gitmap.exe

**Output:**

    Creating tag v3.0.0... done
    Attaching assets... done
    Pushing to origin... done
    ✓ Released v3.0.0
