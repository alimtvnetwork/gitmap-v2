# gitmap amend-list

List previous author amendment records.

## Alias

al

## Usage

    gitmap amend-list [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap amend` at least once (see amend.md)

## Examples

### Example 1: List amendments

    gitmap amend-list

**Output:**

    1  abc1234  old@email.com → john@example.com  2025-03-10
    2  def5678  other@email.com → jane@co.com     2025-03-09
    2 amendments

### Example 2: JSON output

    gitmap al --json

**Output:**

    [{"commit":"abc1234","old_email":"old@email.com","new_email":"john@example.com"}]

## See Also

- [amend](amend.md) — Rewrite commit author information
- [history](history.md) — View command history