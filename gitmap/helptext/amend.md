# gitmap amend

Rewrite commit author information for one or more commits.

## Alias

am

## Usage

    gitmap amend [commit-hash] --name <name> --email <email> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --name \<name\> | — | New author name |
| --email \<email\> | — | New author email |
| --branch \<name\> | current | Target branch |
| --dry-run | false | Preview without executing |
| --force | false | Skip confirmation prompt |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Amend last commit author

    gitmap amend --name "John Doe" --email "john@example.com"

**Output:**

    Amending commit abc1234...
    Author: John Doe <john@example.com>
    ✓ 1 commit amended

### Example 2: Dry-run preview

    gitmap am abc1234 --name "Jane" --email "jane@co.com" --dry-run

**Output:**

    [DRY RUN] Would amend commit abc1234
    [DRY RUN] Old: Old Name <old@email.com>
    [DRY RUN] New: Jane <jane@co.com>
    No changes made.
