# gitmap status

Show a dashboard of repository statuses (branch, clean/dirty, ahead/behind).

## Alias

st

## Usage

    gitmap status [--group <name>] [--all]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --group \<name\> | — | Show status for repos in a group |
| --all | false | Show status for all tracked repos |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Status of all tracked repos

    gitmap status --all

**Output:**

    my-api       main     clean   0/0
    web-app      develop  dirty   2/1
    ✓ 3 repos (1 dirty, 2 clean)

### Example 2: Status of a group

    gitmap st --group work

**Output:**

    billing-svc  main  clean  0/0
    auth-gateway main  dirty  1/0
    ✓ 2 repos (group: work)

## See Also

- [watch](watch.md) — Live-refresh status dashboard
- [scan](scan.md) — Scan directories to populate the database
- [group](group.md) — Manage repo groups
- [pull](pull.md) — Pull repos to sync changes