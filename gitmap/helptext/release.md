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
| --compress | false | Wrap assets in .zip (Windows) or .tar.gz archives |
| --checksums | false | Generate SHA256 checksums.txt for assets |

## Prerequisites

- Must be inside a Git repository with at least one commit
- GitHub CLI (`gh`) recommended for publishing

## Examples

### Example 1: Release with auto-bump

    gitmap release --bump patch

**Output:**

    v2.3.9 → v2.3.10
    Creating tag v2.3.10... done
    ✓ Released v2.3.10

### Example 2: Dry-run preview

    gitmap r --bump minor --dry-run

**Output:**

    [DRY RUN] v2.3.10 → v2.4.0
    [DRY RUN] Would create tag and push
    No changes made.

### Example 3: Release with assets

    gitmap release v3.0.0 --assets ./dist/gitmap.exe

**Output:**

    Creating tag v3.0.0... done
    Attaching assets... done
    ✓ Released v3.0.0

## See Also

- [release-branch](release-branch.md) — Create a release branch without tagging
- [release-pending](release-pending.md) — Show unreleased commits
- [changelog](changelog.md) — View release notes
- [list-versions](list-versions.md) — List release tags
- [list-releases](list-releases.md) — List stored release metadata
- [revert](revert.md) — Revert to a previous release