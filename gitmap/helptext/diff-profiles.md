# gitmap diff-profiles

Compare repositories across two profiles to find differences.

## Alias

dp

## Usage

    gitmap diff-profiles <profileA> <profileB> [--all] [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --all | false | Show all repos, not just differences |
| --json | false | Output as structured JSON |

## Prerequisites

- At least two profiles must exist (see profile.md)

## Examples

### Example 1: Compare two profiles

    gitmap diff-profiles home work

**Output:**

    Only in 'home':
      personal-blog     ~/home/personal-blog
    Only in 'work':
      billing-svc       ~/work/billing-svc
    Common: 12 repos

### Example 2: Full comparison as JSON

    gitmap dp home work --all --json

**Output:**

    {"only_a":["personal-blog"],"only_b":["billing-svc"],"common":["my-api","web-app"]}
