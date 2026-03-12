# gitmap exec

Run a git command across all tracked repositories.

## Alias

x

## Usage

    gitmap exec <git-args...>

## Flags

None (all arguments are passed directly to git).

## Prerequisites

- Run `gitmap scan` first to populate the database (see scan.md)

## Examples

### Example 1: Fetch and prune across all repos

    gitmap exec fetch --prune

**Output:**

    [my-api] git fetch --prune... done
    [web-app] git fetch --prune... done
    ✓ 3 repos processed

### Example 2: Check remote URLs

    gitmap x remote -v

**Output:**

    [my-api] origin https://github.com/user/my-api.git
    [web-app] origin https://github.com/user/web-app.git
    ✓ 2 repos processed