# gitmap list-releases

List release metadata stored in the local database.

## Alias

lr

## Usage

    gitmap list-releases [--json] [--source manual|scan]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |
| --source \<type\> | — | Filter by release source (manual or scan) |

## Prerequisites

- Run `gitmap scan` or `gitmap release` to populate release data (see scan.md, release.md)

## Examples

### Example 1: List all releases

    gitmap list-releases

**Output:**

    Version   Date         Source   Commits
    v2.8.0    2025-03-10   manual   12
    v2.5.0    2025-03-01   scan     8
    v2.4.0    2025-02-20   scan     15
    3 releases found

### Example 2: Filter by source

    gitmap lr --source scan --json

**Output:**

    [{"version":"v2.5.0","date":"2025-03-01","source":"scan","commits":8}]
