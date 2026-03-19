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

### Example 1: Run full diagnostics

    gitmap doctor

**Output:**

    ■ Checking system...
    ✓ Git installed (v2.43.0)
    ✓ Go installed (go1.22.1)
    ✓ GitHub CLI (gh) installed (v2.44.1)
    ■ Checking binary...
    ✓ Binary at E:\bin-run\gitmap.exe
    ✗ Stale binary found at C:\old\gitmap.exe
    ■ Checking database...
    ✓ Database OK (42 repos, 3 groups)
    ■ Result: 1 issue found
      → Run 'gitmap doctor --fix-path' to resolve

### Example 2: Fix PATH issues automatically

    gitmap doctor --fix-path

**Output:**

    ■ Checking system...
    ✓ Git installed (v2.43.0)
    ✓ Go installed (go1.22.1)
    ■ Checking binary...
    ✓ Binary at E:\bin-run\gitmap.exe
    ✓ Fixed: removed stale binary at C:\old\gitmap.exe
    ■ Checking database...
    ✓ Database OK (42 repos, 3 groups)
    ✓ All checks passed.

### Example 3: Clean installation

    gitmap doctor

**Output:**

    ■ Checking system...
    ✓ Git installed (v2.43.0)
    ✓ Go installed (go1.22.1)
    ✓ GitHub CLI (gh) installed (v2.44.1)
    ■ Checking binary...
    ✓ Binary at E:\bin-run\gitmap.exe
    ■ Checking database...
    ✓ Database OK (42 repos, 3 groups)
    ✓ All checks passed. No issues found.

## See Also

- [setup](setup.md) — Run first-time configuration
- [update](update.md) — Update gitmap to the latest version
- [version](version.md) — Check current version
