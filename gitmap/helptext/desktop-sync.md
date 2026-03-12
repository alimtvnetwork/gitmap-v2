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
    [2/42] web-app... already registered
    ✓ 42 repos synced (15 new, 27 existing)

### Example 2: Using alias

    gitmap ds

**Output:**

    Syncing 42 repos to GitHub Desktop...
    ✓ 42 repos synced
