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

    Watching 42 repos (interval: 30s)  Press Ctrl+C to stop
    ─────────────────────────────────────────────
    Repository        Branch   Status   Behind
    my-api            main     clean    0
    web-app           develop  dirty    2
    Refreshing in 28s...

### Example 2: Watch a group with fast refresh

    gitmap w --group work --interval 10

**Output:**

    Watching 5 repos (group: work, interval: 10s)
    ─────────────────────────────────────────────
    billing-svc       main     clean    0
    auth-gateway      main     clean    0

### Example 3: Single JSON snapshot

    gitmap watch --json --no-fetch

**Output:**

    [{"name":"my-api","branch":"main","status":"clean","behind":0}]
