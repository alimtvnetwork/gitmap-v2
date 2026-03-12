# gitmap history

Show CLI command execution history with timestamps.

## Alias

hi

## Usage

    gitmap history [--limit N] [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --limit \<N\> | 20 | Number of entries to show |
| --json | false | Output as structured JSON |

## Prerequisites

- None (history is recorded automatically)

## Examples

### Example 1: Show recent history

    gitmap history

**Output:**

    1  scan ~/projects       2025-03-10 14:30
    2  clone json            2025-03-10 14:32
    3  status --all          2025-03-10 15:00

### Example 2: Last 5 entries as JSON

    gitmap hi --limit 5 --json

**Output:**

    [{"id":1,"command":"scan ~/projects","timestamp":"2025-03-10T14:30:00Z"}]

## See Also

- [history-reset](history-reset.md) — Clear command history
- [stats](stats.md) — View aggregated usage metrics
- [bookmark](bookmark.md) — Save commands for re-execution