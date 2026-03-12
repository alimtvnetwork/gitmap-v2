# gitmap import

Import repositories from an export file into the local database.

## Alias

im

## Usage

    gitmap import <file>

## Flags

None.

## Prerequisites

- An export file from `gitmap export` (see export.md)

## Examples

### Example 1: Import from file

    gitmap import gitmap-export.json

**Output:**

    Importing from gitmap-export.json...
    [1/42] my-api... added
    [2/42] web-app... added
    ✓ 42 repos imported

### Example 2: Import using alias

    gitmap im backup.json

**Output:**

    Importing from backup.json...
    ✓ 15 repos imported (3 skipped, already exist)
