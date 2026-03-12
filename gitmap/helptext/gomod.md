# gitmap gomod

Rename the Go module path across the entire repository with branch safety.

## Alias

gm

## Usage

    gitmap gomod <new-module-path> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --ext \<exts\> | all files | Comma-separated extensions to filter |
| --dry-run | false | Preview changes without modifying |
| --no-merge | false | Stay on feature branch after commit |
| --no-tidy | false | Skip go mod tidy after replacement |
| --verbose | false | Print each file path as modified |

## Prerequisites

- Must be inside a Go project with go.mod

## Examples

### Example 1: Rename module path

    gitmap gomod "github.com/neworg/myproject"

**Output:**

    Replacing in 24 files... done
    Running go mod tidy... done
    ✓ Module renamed

### Example 2: Dry-run with specific extensions

    gitmap gm "github.com/new/name" --ext "*.go,*.md" --dry-run

**Output:**

    [DRY RUN] 18 .go files, 3 .md files
    No changes made.

### Example 3: Rename without merge

    gitmap gomod "github.com/new/name" --no-merge --verbose

**Output:**

    Replacing in cmd/root.go... done
    Committed on gomod-rename (not merged)
    ✓ 24 files updated

## See Also

- [go-repos](go-repos.md) — List detected Go projects
- [scan](scan.md) — Scan directories to detect Go projects