# gitmap multi-group

Select multiple groups for batch operations (pull, status, exec).

## Alias

mg

## Usage

    gitmap multi-group <group1,group2,...|clear|pull|status|exec>

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)
- Create groups with `gitmap group create` (see group.md)

## Examples

### Example 1: Select multiple groups

    gitmap mg backend,frontend

**Output:**

    Multi-group set: backend,frontend

### Example 2: Pull repos from all selected groups

    gitmap mg pull

**Output:**

    Pulling my-api (main)...
    ✓ my-api is up to date.

### Example 3: Clear multi-group selection

    gitmap mg clear

**Output:**

    Multi-group selection cleared.

## See Also

- [group](group.md) — Manage and activate single groups
- [pull](pull.md) — Pull a specific repo
- [status](status.md) — View repo statuses
- [exec](exec.md) — Run git across repos
