# gitmap cd

Navigate to a tracked repository directory using its slug or an interactive picker.

## Alias

go

## Usage

    gitmap cd <repo-name|repos> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --group \<name\> | — | Filter picker to a specific group |
| --pick | false | Force interactive picker |
| --default \<slug\> | — | Set or clear the default repo |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Navigate to a repo by slug

    gitmap cd my-api

**Output:**

    Changed directory to ~/projects/my-api
    Branch: main | Status: clean

### Example 2: Interactive repo picker

    gitmap cd repos

**Output:**

    1. my-api      ~/projects/my-api
    2. web-app     ~/projects/web-app
    Enter number: _

### Example 3: Pick from a specific group

    gitmap cd repos --group work

**Output:**

    1. billing-svc   ~/work/billing-svc
    2. auth-gateway  ~/work/auth-gateway
    Enter number: _