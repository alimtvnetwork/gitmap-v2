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
| --no-assets | false | Skip Go binary cross-compilation |
| --targets \<list\> | all 6 | Cross-compile targets: windows/amd64,linux/arm64 |
| --list-targets | false | Print resolved target matrix and exit |
| --zip-group \<name\> | — | Include a persistent zip group as a release asset |
| -Z \<path\> | — | Add ad-hoc file or folder to zip as a release asset |
| --bundle \<name.zip\> | — | Bundle all -Z items into a single named archive |

## Prerequisites

- Must be inside a Git repository with at least one commit
- GitHub CLI (`gh`) recommended for publishing

## Orphaned Metadata Recovery

If a `.release/vX.Y.Z.json` file exists but neither the Git tag nor
the release branch is found, the command warns and prompts:

    ⚠ Release metadata exists for v2.3.10 but no tag or branch was found.
    → Do you want to remove the release JSON and proceed? (y/N):

Answering `y` removes the stale JSON file and proceeds with the release.
Answering `n` (or pressing Enter) aborts the release.

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

### Example 4: Release with a persistent zip group

    gitmap release v3.0.0 --zip-group docs-bundle

**Output:**

    Creating tag v3.0.0... done
    ✓ Compressed docs-bundle → docs-bundle_v3.0.0.zip
    ✓ Released v3.0.0

### Example 5: Ad-hoc zip with bundle

    gitmap release v3.0.0 -Z ./dist/report.pdf -Z ./dist/manual.pdf --bundle docs.zip

**Output:**

    Creating tag v3.0.0... done
    ✓ Compressed 2 items → docs.zip
    ✓ Released v3.0.0

## See Also

- [release-branch](release-branch.md) — Create a release branch without tagging
- [release-pending](release-pending.md) — Show unreleased commits
- [changelog](changelog.md) — View release notes
- [list-versions](list-versions.md) — List release tags
- [list-releases](list-releases.md) — List stored release metadata
- [revert](revert.md) — Revert to a previous release
- [zip-group](zip-group.md) — Manage zip group definitions