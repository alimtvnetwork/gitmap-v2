# gitmap history-reset

Clear the CLI command execution history.

## Alias

hr

## Usage

    gitmap history-reset [--confirm]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --confirm | false | Skip confirmation prompt |

## Prerequisites

- None

## Examples

### Example 1: Reset with confirmation

    gitmap history-reset

**Output:**

    This will delete all command history. Continue? [y/N]: y
    ✓ History cleared (42 entries removed)

### Example 2: Reset without prompt

    gitmap hr --confirm

**Output:**

    ✓ History cleared (42 entries removed)
