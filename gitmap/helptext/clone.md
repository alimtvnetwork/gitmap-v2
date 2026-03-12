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

    gitmap clone json --target-dir ./projects

**Output:**

    Cloning from gitmap-output/gitmap.json...
    [1/12] Cloning my-api... done
    ✓ 12 repositories cloned

### Example 2: Clone with safe-pull for existing repos

    gitmap c csv --safe-pull

**Output:**

    [1/8] my-api exists, pulling... done
    [2/8] Cloning web-app... done
    ✓ 8 repositories processed

## See Also

- [scan](scan.md) — Scan directories to generate output files
- [pull](pull.md) — Pull individual or grouped repos
- [desktop-sync](desktop-sync.md) — Sync repos to GitHub Desktop