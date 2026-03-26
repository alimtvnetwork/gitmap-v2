# gitmap list-releases

List release metadata from the current git repo or stored database.

When run inside a git repo with `.gitmap/release/v*.json` files, releases are read
directly from those files. Otherwise, releases are loaded from the gitmap
database.

## Alias

lr

## Usage

    gitmap list-releases [--json] [--limit N] [--source repo|release|import]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |
| --limit \<N\> | 0 | Show only the top N releases (0 = all) |
| --source \<type\> | — | Filter by release source (repo, release, or import) |

## Prerequisites

- Inside a git repo with `.gitmap/release/v*.json` files, **or**
- Run `gitmap scan` or `gitmap release` to populate the database

## Examples

### Example 1: List releases from current repo

    gitmap list-releases

**Output:**

    Releases (5 found)
    ────────────────────────────────────────────────────────────────────────
      VERSION    TAG          BRANCH              DRAFT  LATEST  SOURCE   DATE
      2.33.0     v2.33.0      release/v2.33.0     no     yes     repo     2026-03-26
      2.31.0     v2.31.0      release/v2.31.0     no     no      repo     2026-03-20
      2.30.0     v2.30.0      release/v2.30.0     no     no      repo     2026-03-15
      2.27.0     v2.27.0      release/v2.27.0     no     no      repo     2026-03-10
      2.26.0     v2.26.0      release/v2.26.0     no     no      repo     2026-03-05
      5 releases found

### Example 2: Show top 3 releases

    gitmap lr --limit 3

**Output:**

    Releases (3 found)
    ────────────────────────────────────────────────────────────────────────
      VERSION    TAG          BRANCH              DRAFT  LATEST  SOURCE   DATE
      2.33.0     v2.33.0      release/v2.33.0     no     yes     repo     2026-03-26
      2.31.0     v2.31.0      release/v2.31.0     no     no      repo     2026-03-20
      2.30.0     v2.30.0      release/v2.30.0     no     no      repo     2026-03-15

### Example 3: JSON output

    gitmap lr --json

**Output:**

    [
      {"version":"2.33.0","tag":"v2.33.0","branch":"release/v2.33.0","source":"repo","draft":false,"isLatest":true},
      {"version":"2.31.0","tag":"v2.31.0","branch":"release/v2.31.0","source":"repo","draft":false,"isLatest":false}
    ]

### Example 4: Filter by source (database releases only)

    gitmap lr --source release

## See Also

- [list-versions](list-versions.md) — List Git release tags
- [changelog](changelog.md) — View release notes
- [release](release.md) — Create a release
- [scan](scan.md) — Scan to import release data
