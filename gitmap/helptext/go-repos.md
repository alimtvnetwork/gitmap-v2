# gitmap go-repos

List all detected Go projects across tracked repositories.

## Alias

gr

## Usage

    gitmap go-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List Go projects

    gitmap go-repos

**Output:**

    Repository        Module Path                    Go Version
    my-api            github.com/user/my-api         1.22
    shared-lib        github.com/user/shared-lib     1.21
    gitmap            github.com/user/gitmap         1.22
    3 Go projects detected

### Example 2: JSON output

    gitmap gr --json

**Output:**

    [{"repo":"my-api","module":"github.com/user/my-api","go_version":"1.22"}]
