# gitmap csharp-repos

List all detected C# projects across tracked repositories.

## Alias

csr

## Usage

    gitmap csharp-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List C# projects

    gitmap csharp-repos

**Output:**

    billing-svc  BillingSvc.sln  net8.0
    auth-api     AuthApi.sln     net7.0
    2 C# projects detected

### Example 2: JSON output

    gitmap csr --json

**Output:**

    [{"repo":"billing-svc","solution":"BillingSvc.sln","target":"net8.0"}]

## See Also

- [scan](scan.md) — Scan directories to detect projects
- [cpp-repos](cpp-repos.md) — List C++ projects
- [go-repos](go-repos.md) — List Go projects
- [node-repos](node-repos.md) — List Node.js projects