# gitmap pull

Pull a specific tracked repository by slug, group, or all at once.

## Alias

p

## Usage

    gitmap pull <repo-name> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --group \<name\> | — | Pull all repos in a group |
| --all | false | Pull all tracked repos |
| --verbose | false | Enable verbose logging |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Pull a single repo by slug

    gitmap pull my-api

**Output:**

    Pulling my-api...
    Already up to date.

### Example 2: Pull all repos in a group

    gitmap p --group work

**Output:**

    Pulling 5 repos in group 'work'...
    [1/5] billing-svc... updated (3 commits)
    ✓ 5 repos pulled

### Example 3: Pull all tracked repos

    gitmap pull --all --verbose

**Output:**

    Pulling 42 tracked repos...
    [1/42] my-api... updated
    ✓ 42 repos pulled (12 updated, 30 up to date)

## See Also

- [scan](scan.md) — Scan directories to populate the database
- [clone](clone.md) — Clone repos from output files
- [status](status.md) — Check repo statuses before pulling
- [group](group.md) — Manage groups for targeted pulls