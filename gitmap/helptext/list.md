# gitmap list

Show all tracked repositories with their slugs and paths.

## Alias

ls

## Usage

    gitmap list [--group <name>] [--verbose]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --group \<name\> | — | Filter to a specific group |
| --verbose | false | Show full paths and metadata |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: List all tracked repos

    gitmap list

**Output:**

    my-api       ~/projects/my-api
    web-app      ~/projects/web-app
    3 repos tracked

### Example 2: List repos in a group

    gitmap ls --group work --verbose

**Output:**

    billing-svc   ~/work/billing-svc   main  origin
    auth-gateway  ~/work/auth-gateway  main  origin
    2 repos in group 'work'

## See Also

- [cd](cd.md) — Navigate to a tracked repo
- [group](group.md) — Manage repo groups
- [scan](scan.md) — Scan directories to populate the database
- [status](status.md) — View repo statuses