# gitmap stats

Show aggregated usage and performance metrics for gitmap commands.

## Alias

ss

## Usage

    gitmap stats [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- None (stats are recorded automatically)

## Examples

### Example 1: Show stats summary

    gitmap stats

**Output:**

    scan   15 runs  2.3s avg
    clone   8 runs  12.1s avg
    Total: 65 executions

### Example 2: Stats as JSON

    gitmap ss --json

**Output:**

    [{"command":"scan","runs":15,"avg_ms":2300}]