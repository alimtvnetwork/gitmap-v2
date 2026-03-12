# gitmap export

Export the local database to a portable file.

## Alias

ex

## Usage

    gitmap export [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Export as JSON format |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Export database

    gitmap export

**Output:**

    Exporting 42 repos...
    ✓ Exported to gitmap-export.json

### Example 2: Export as JSON

    gitmap ex --json

**Output:**

    Exporting 42 repos...
    ✓ Exported to gitmap-export.json

## See Also

- [import](import.md) — Import repos from an export file
- [scan](scan.md) — Scan directories to populate the database
- [profile](profile.md) — Manage database profiles
