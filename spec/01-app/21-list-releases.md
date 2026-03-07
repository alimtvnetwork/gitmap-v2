# 21 — list-releases Command

## Purpose

`gitmap list-releases` (`lr`) queries the `Releases` table in the SQLite database and displays stored release records, sorted from newest to oldest.

## Command Signature

```
gitmap list-releases [flags]
gitmap lr [flags]
```

## Flags

| Flag       | Short | Default | Description                              |
|------------|-------|---------|------------------------------------------|
| `--json`   |       | false   | Output as JSON array                     |
| `--limit`  |       | 0       | Show only the top N releases (0 = all)   |
| `--source` |       | (all)   | Filter by source: `release` or `import`  |

## Data Source

All data comes from the `Releases` table (populated by `gitmap release`).
No Git commands are executed. If the database is missing or empty, print an
informative message and exit 1.

## Behavior

1. Open the database via `store.Open()` and run `db.Migrate()`.
2. Call `db.ListReleases()` to fetch all records (ordered by `CreatedAt DESC`).
3. If `--limit N` is provided and N > 0, trim the result slice to N entries.
4. Render output in terminal or JSON format.
5. If no releases are found, print `"No releases found."` and exit 0.

## Terminal Output

Table format with columns: Version, Tag, Branch, Draft, Latest, Date.

```
Releases (3 found)
──────────────────────────────────────────────────────────
  VERSION    TAG        BRANCH       DRAFT  LATEST  DATE
  2.15.0     v2.15.0    release/v2   no     yes     2026-03-07
  2.14.0     v2.14.0    release/v2   no     no      2026-03-01
  2.13.0     v2.13.0    release/v2   yes    no      2026-02-20
```

## JSON Output Example

```json
[
  {
    "version": "2.15.0",
    "tag": "v2.15.0",
    "branch": "release/v2.15.0",
    "sourceBranch": "main",
    "commitSha": "abc123",
    "changelog": "Added --limit flag",
    "draft": false,
    "preRelease": false,
    "isLatest": true,
    "createdAt": "2026-03-07T10:00:00Z"
  }
]
```

## Error Handling

| Condition             | Message                                          | Exit |
|-----------------------|--------------------------------------------------|------|
| DB missing / no scan  | `"No database found. Run gitmap scan first.\n"`  | 1    |
| DB open/migrate error | `"failed to load releases: %v\n"`                | 1    |
| No releases           | `"No releases found.\n"`                         | 0    |

## Implementation Files

| File                            | Responsibility                               |
|---------------------------------|----------------------------------------------|
| `cmd/listreleases.go`           | Command handler, flag parsing, output        |
| `constants/constants_cli.go`    | `CmdListReleases`, `CmdListReleasesAlias`    |
| `constants/constants_messages.go` | Terminal output format strings             |
| `store/release.go`              | `ListReleases()` (already exists)            |

## Integration Points

- `cmd/root.go`: register `list-releases` / `lr` in `dispatchMisc`.
- Reuse `openDB()` helper from `cmd/list.go`.
- Reuse `store.ListReleases()` — no new DB queries needed.

## Code Style

All functions ≤ 15 lines. Positive logic. Blank line before every return.
No magic strings. No switch statements. PascalCase for SQL column names.
