# gitmap bookmark

Save and run bookmarked gitmap commands for quick re-execution.

## Alias

bk

## Usage

    gitmap bookmark <save|list|run|delete> [args]

## Flags

None.

## Prerequisites

- None

## Examples

### Example 1: Save a bookmark

    gitmap bookmark save "daily-scan" "scan ~/projects --quiet"

**Output:**

    ✓ Bookmark 'daily-scan' saved

### Example 2: List bookmarks

    gitmap bk list

**Output:**

    daily-scan  scan ~/projects --quiet
    work-pull   pull --group work
    2 bookmarks

### Example 3: Run a bookmark

    gitmap bookmark run daily-scan

**Output:**

    Running 'daily-scan'...
    Found 42 repositories
    ✓ Scan complete