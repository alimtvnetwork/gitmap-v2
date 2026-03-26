# gitmap clone

Re-clone repositories from a structured output file (JSON, CSV, or text).

## Alias

c

## Usage

    gitmap clone <source|json|csv|text> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --target-dir \<dir\> | current directory | Base directory for clones |
| --safe-pull | false | Pull existing repos with retry + diagnostics |
| --verbose | false | Write detailed debug log |

## Prerequisites

- Run `gitmap scan` first to generate output files (see scan.md)

## Examples

### Example 1: Clone from JSON output

    gitmap clone json --target-dir D:\projects

**Output:**

    Cloning from .gitmap/output/gitmap.json...
    [1/12] Cloning my-api... done
    [2/12] Cloning web-app... done
    [3/12] Cloning billing-svc... done
    [4/12] Cloning auth-gateway... done
    ...
    ✓ 12 repositories cloned to D:\projects

### Example 2: Clone with safe-pull for existing repos

    gitmap c csv --safe-pull

**Output:**

    [1/8] my-api exists → pulling... Already up to date.
    [2/8] web-app exists → pulling... Updated (3 new commits)
    [3/8] Cloning billing-svc... done
    [4/8] Cloning auth-gateway... done
    ...
    ✓ 8 repositories processed (2 pulled, 6 cloned)

### Example 3: Clone from text file with verbose logging

    gitmap clone text --verbose

**Output:**

    [verbose] Log file: gitmap-debug-2025-03-10T14-30.log
    Cloning from .gitmap/output/gitmap.txt...
    [1/5] Cloning https://github.com/user/my-api.git... done
    [2/5] Cloning https://github.com/user/web-app.git... done
    ...
    ✓ 5 repositories cloned
    [verbose] Debug log written to gitmap-debug-2025-03-10T14-30.log

## See Also

- [scan](scan.md) — Scan directories to generate output files
- [pull](pull.md) — Pull individual or grouped repos
- [desktop-sync](desktop-sync.md) — Sync repos to GitHub Desktop
