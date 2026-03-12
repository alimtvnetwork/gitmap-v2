# gitmap desktop-sync

Sync all tracked repositories with GitHub Desktop.

## Alias

ds

## Usage

    gitmap desktop-sync

## Flags

None.

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)
- GitHub Desktop must be installed

## Examples

### Example 1: Sync repos to GitHub Desktop

    gitmap desktop-sync

**Output:**

    Syncing 42 repos to GitHub Desktop...
    [1/42] my-api... added
    ✓ 42 repos synced (15 new, 27 existing)

### Example 2: Using alias

    gitmap ds

**Output:**

    Syncing 42 repos...
    ✓ 42 repos synced

## See Also

- [scan](scan.md) — Scan directories to populate the database
- [clone](clone.md) — Clone repos from output files
- [list](list.md) — View tracked repos