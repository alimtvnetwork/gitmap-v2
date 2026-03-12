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

    GitMap Doctor
    ─────────────
    ✓ Git installed (v2.43.0)
    ✓ gitmap binary found at E:\bin-run\gitmap.exe
    ✓ PATH includes E:\bin-run
    ✗ Stale binary detected in C:\old\gitmap.exe
    1 issue found

### Example 2: Fix PATH issues

    gitmap doctor --fix-path

**Output:**

    GitMap Doctor
    ─────────────
    ✓ Git installed (v2.43.0)
    ✓ Fixed: removed stale binary from C:\old\
    All checks passed.
