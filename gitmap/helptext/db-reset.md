# gitmap db-reset

Reset the local SQLite database, removing all tracked repos and metadata.

## Alias

None

## Usage

    gitmap db-reset [--confirm]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --confirm | false | Skip confirmation prompt |

## Prerequisites

- None

## Examples

### Example 1: Reset with confirmation

    gitmap db-reset

**Output:**

    This will delete all data in the current profile. Continue? [y/N]: y
    ✓ Database reset (42 repos, 3 groups removed)

### Example 2: Reset without prompt

    gitmap db-reset --confirm

**Output:**

    ✓ Database reset (42 repos, 3 groups removed)

## See Also

- [scan](scan.md) — Re-scan to rebuild the database
- [history-reset](history-reset.md) — Clear command history only
- [profile](profile.md) — Manage profiles (reset affects current profile)
