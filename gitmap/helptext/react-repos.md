# gitmap react-repos

List all detected React projects across tracked repositories.

## Alias

rr

## Usage

    gitmap react-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List React projects

    gitmap react-repos

**Output:**

    web-app    @user/web-app  18.2.0
    docs-site  docs-site      18.2.0
    2 React projects detected

### Example 2: JSON output

    gitmap rr --json

**Output:**

    [{"repo":"web-app","package":"@user/web-app","react_version":"18.2.0"}]

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [node-repos](node-repos.md) — List Node.js projects
- [go-repos](go-repos.md) — List Go projects
- [csharp-repos](csharp-repos.md) — List C# projects