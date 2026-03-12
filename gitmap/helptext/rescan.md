# gitmap rescan

Re-scan previously scanned directories using cached scan parameters.

## Alias

rs

## Usage

    gitmap rescan [--output csv|json|terminal]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --output csv\|json\|terminal | terminal | Output format |

## Prerequisites

- Run `gitmap scan` at least once to create scan cache (see scan.md)

## Examples

### Example 1: Quick rescan

    gitmap rescan

**Output:**

    Re-scanning ~/projects (cached)...
    Found 44 repositories (+2 new)
    Output written to ./gitmap-output/

### Example 2: Rescan with JSON output

    gitmap rs --output json

**Output:**

    Re-scanning ~/projects (cached)...
    Found 44 repositories
    ✓ gitmap-repos.json

## See Also

- [scan](scan.md) — Initial directory scan
- [status](status.md) — View repo statuses
- [clone](clone.md) — Clone from scan output
