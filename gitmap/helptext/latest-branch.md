# gitmap latest-branch

Find the most recently updated remote branch in the current repository.

## Alias

lb

## Usage

    gitmap latest-branch [--top N] [--format json|csv|terminal]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --top \<N\> | 1 | Number of branches to show |
| --format json\|csv\|terminal | terminal | Output format |
| --no-fetch | false | Skip git fetch before query |
| --sort date\|name | date | Sort order |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Show latest branch

    gitmap lb

**Output:**

    Branch: feature/auth-redesign
    Last commit: 2 hours ago
    Author: developer@example.com

### Example 2: Top 5 branches as CSV

    gitmap lb 5 --format csv

**Output:**

    branch,last_commit,author
    feature/auth-redesign,2025-03-10T14:30:00Z,dev@example.com
    bugfix/login-fix,2025-03-10T12:15:00Z,dev@example.com
    main,2025-03-09T18:00:00Z,dev@example.com

### Example 3: JSON output

    gitmap latest-branch --format json

**Output:**

    {"branch":"feature/auth-redesign","last_commit":"2h ago","author":"dev@example.com"}

## See Also

- [status](status.md) — View repo branch and status info
- [release-branch](release-branch.md) — Create a release branch
- [watch](watch.md) — Live-refresh status dashboard
