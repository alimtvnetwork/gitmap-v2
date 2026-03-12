# gitmap group

Manage repository groups (create, add, remove, list, show, delete).

## Alias

g

## Usage

    gitmap group <create|add|remove|list|show|delete> [args]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --desc \<text\> | — | Group description (for create) |
| --color \<hex\> | — | Group color (for create) |

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Create a group and add repos

    gitmap group create work --desc "Work repositories"
    gitmap group add work my-api web-app

**Output:**

    ✓ Group 'work' created
    ✓ Added my-api to 'work'
    ✓ Added web-app to 'work'

### Example 2: List all groups

    gitmap g list

**Output:**

    work      5 repos   Work repositories
    personal  3 repos   Personal projects
    2 groups

### Example 3: Show group details

    gitmap group show work

**Output:**

    Group: work (5 repos)
    billing-svc   ~/work/billing-svc
    auth-gateway  ~/work/auth-gateway