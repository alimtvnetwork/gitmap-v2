# gitmap release-branch

Create a release branch without tagging or publishing.

## Alias

rb

## Usage

    gitmap release-branch [version] [--bump major|minor|patch]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --bump major\|minor\|patch | — | Auto-increment version |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Create release branch with minor bump

    gitmap release-branch --bump minor

**Output:**

    Current version: v2.21.0
    v2.21.0 → v2.22.0
    Creating branch release/v2.22.0... done
    Pushing release/v2.22.0 to origin... done
    Switched to branch 'release/v2.22.0'
    → Ready to finalize: gitmap release v2.22.0

### Example 2: Create release branch with explicit version

    gitmap rb v3.0.0

**Output:**

    Creating branch release/v3.0.0 from main... done
    Pushing release/v3.0.0 to origin... done
    Switched to branch 'release/v3.0.0'
    → Ready to finalize: gitmap release v3.0.0

### Example 3: Release branch with major bump

    gitmap release-branch --bump major

**Output:**

    Current version: v2.22.0
    v2.22.0 → v3.0.0
    Creating branch release/v3.0.0... done
    Pushing release/v3.0.0 to origin... done
    Switched to branch 'release/v3.0.0'

## See Also

- [release](release.md) — Create a full release with tag and push
- [release-pending](release-pending.md) — Show unreleased commits
- [latest-branch](latest-branch.md) — Find most recently updated branch
