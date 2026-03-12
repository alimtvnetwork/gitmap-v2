# gitmap doctor

Diagnose PATH, deployment, and version issues with the gitmap installation.

## Alias

None

## Usage

    gitmap doctor [--fix-path]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --fix-path | false | Attempt to fix PATH issues automatically |

## Prerequisites

- None

## Examples

### Example 1: Run diagnostics

    gitmap doctor

**Output:**

    ✓ Git installed (v2.43.0)
    ✓ Binary at E:\bin-run\gitmap.exe
    ✗ Stale binary in C:\old\ — 1 issue

### Example 2: Fix PATH issues

    gitmap doctor --fix-path

**Output:**

    ✓ Git installed (v2.43.0)
    ✓ Fixed: removed stale binary
    All checks passed.