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

    my-api      github.com/user/my-api      1.22
    shared-lib  github.com/user/shared-lib  1.21
    3 Go projects detected

### Example 2: JSON output

    gitmap gr --json

**Output:**

    [{"repo":"my-api","module":"github.com/user/my-api","go_version":"1.22"}]

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [node-repos](node-repos.md) — List Node.js projects
- [react-repos](react-repos.md) — List React projects
- [gomod](gomod.md) — Rename Go module paths