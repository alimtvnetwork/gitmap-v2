# gitmap release-pending

Show unreleased commits since the last tag.

## Alias

rp

## Usage

    gitmap release-pending [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Must be inside a Git repository with at least one tag

## Examples

### Example 1: Show pending commits

    gitmap release-pending

**Output:**

    Unreleased since v2.3.10:
    abc1234 Add user auth endpoint
    3 commits pending release

### Example 2: JSON output

    gitmap rp --json

**Output:**

    [{"hash":"abc1234","message":"Add user auth endpoint","date":"2025-03-10"}]