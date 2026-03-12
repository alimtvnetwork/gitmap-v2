# gitmap node-repos

List all detected Node.js projects across tracked repositories.

## Alias

nr

## Usage

    gitmap node-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List Node.js projects

    gitmap node-repos

**Output:**

    web-app    @user/web-app  18.x
    docs-site  docs-site      20.x
    2 Node.js projects detected

### Example 2: JSON output

    gitmap nr --json

**Output:**

    [{"repo":"web-app","package":"@user/web-app","node_version":"18.x"}]