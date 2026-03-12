# gitmap watch

Live-refresh repository status dashboard with configurable interval.

## Alias

w

## Usage

    gitmap watch [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --interval \<seconds\> | 30 | Refresh interval (min: 5) |
| --group \<name\> | — | Monitor only repos in a group |
| --no-fetch | false | Skip git fetch before status |
| --json | false | Output single snapshot as JSON |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Watch all repos

    gitmap watch

**Output:**

    Watching 42 repos (30s) — Ctrl+C to stop
    my-api  main  clean | web-app  develop  dirty
    Refreshing in 28s...

### Example 2: Watch a group with fast refresh

    gitmap w --group work --interval 10

**Output:**

    Watching 5 repos (group: work, 10s)
    billing-svc  main  clean
    auth-gateway main  clean

### Example 3: Single JSON snapshot

    gitmap watch --json --no-fetch

**Output:**

    [{"name":"my-api","branch":"main","status":"clean"}]

## See Also

- [status](status.md) — One-time status snapshot
- [scan](scan.md) — Scan directories to populate the database
- [group](group.md) — Manage repo groups for filtered watching