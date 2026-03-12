# gitmap list-versions

List all available Git release tags in the repository.

## Alias

lv

## Usage

    gitmap list-versions [--json] [--limit N]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |
| --limit \<N\> | 0 | Show only the top N versions (0 = all) |

## Prerequisites

- Must be inside a Git repository with tags

## Examples

### Example 1: List all versions

    gitmap list-versions

**Output:**

    v2.8.0  2025-03-10
    v2.5.0  2025-03-01
    4 versions found

### Example 2: Top 3 as JSON

    gitmap lv --json --limit 3

**Output:**

    [{"version":"v2.8.0","date":"2025-03-10"},...]