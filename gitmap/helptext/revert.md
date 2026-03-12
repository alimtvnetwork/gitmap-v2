# gitmap revert

Revert the repository to a specific release version by checking out the tag.

## Alias

None

## Usage

    gitmap revert <version>

## Flags

None.

## Prerequisites

- Must be inside a Git repository with release tags
- Run `gitmap list-versions` to see available versions (see list-versions.md)

## Examples

### Example 1: Revert to a specific version

    gitmap revert v2.5.0

**Output:**

    Reverting to v2.5.0...
    Rebuilding... done
    ✓ Reverted to v2.5.0

### Example 2: Revert to an older version

    gitmap revert v2.3.7

**Output:**

    Reverting to v2.3.7...
    Rebuilding... done
    ✓ Reverted to v2.3.7

## See Also

- [list-versions](list-versions.md) — List available versions to revert to
- [release](release.md) — Create a new release
- [changelog](changelog.md) — View release notes before reverting