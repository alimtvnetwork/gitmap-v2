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

### Example 1: Create release branch with bump

    gitmap release-branch --bump minor

**Output:**

    Current version: v2.3.10
    Creating branch release/v2.4.0... done
    Switched to branch 'release/v2.4.0'

### Example 2: Create release branch with explicit version

    gitmap rb v3.0.0

**Output:**

    Creating branch release/v3.0.0... done
    Switched to branch 'release/v3.0.0'
