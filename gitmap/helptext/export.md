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
