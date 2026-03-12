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

    Repository        Branch   Status   Ahead  Behind
    my-api            main     clean    0      0
    web-app           develop  dirty    2      1
    shared-lib        main     clean    0      3
    ✓ 3 repos (1 dirty, 2 clean)

### Example 2: Status of a group

    gitmap st --group work

**Output:**

    Repository        Branch   Status   Ahead  Behind
    billing-svc       main     clean    0      0
    auth-gateway      main     dirty    1      0
    ✓ 2 repos (group: work)
