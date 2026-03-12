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

    Command          Runs   Avg Time   Last Run
    scan             15     2.3s       2025-03-10
    clone            8      12.1s      2025-03-09
    status           42     0.8s       2025-03-10
    Total: 65 executions

### Example 2: Stats as JSON

    gitmap ss --json

**Output:**

    [{"command":"scan","runs":15,"avg_ms":2300,"last":"2025-03-10"}]
