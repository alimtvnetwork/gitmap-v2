# gitmap clear-release-json

Remove a specific release metadata JSON file from the `.release/` directory.

## Alias

crj

## Usage

    gitmap clear-release-json <version>

## Flags

None

## Prerequisites

- A `.release/vX.Y.Z.json` file must exist for the given version.

## Examples

### Example 1: Remove a release JSON file

    gitmap clear-release-json v2.20.0

**Output:**

    ✓ Removed .release/v2.20.0.json

### Example 2: Version not found

    gitmap clear-release-json v9.9.9

**Output:**

    Error: no release file found for v9.9.9

## See Also

- [release](release.md) — Create a release
- [list-releases](list-releases.md) — Show stored releases
- [db-reset](db-reset.md) — Clear the entire database
