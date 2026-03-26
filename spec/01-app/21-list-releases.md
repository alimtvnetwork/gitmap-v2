# 21 — list-releases Command

## Purpose

`gitmap list-releases` (`lr`) displays release records, sorted from newest to
oldest. It reads from the **current git repo** first (`.release/v*.json`
files), falling back to the SQLite database only when no `.release/` directory
exists.

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
| `--source` |       | (all)   | Filter by source: `release`, `import`, or `repo` |

## Data Source (Resolution Order)

1. **Repo-local `.release/` files** (preferred): read all `.release/v*.json`
   files via `release.ListReleaseMetaFiles()`, convert each `ReleaseMeta` to
   a `ReleaseRecord` with `Source = "repo"`, sort by `CreatedAt DESC`, and
   mark the latest using `.release/latest.json`.
2. **Database fallback**: if no `.release/` files are found, open the SQLite
   database via `store.Open()` and call `db.ListReleases()`.

This ensures `gitmap lr` works inside any git repo that has `.release/`
metadata, without requiring a prior `gitmap scan`.

## Behavior

1. Call `release.ListReleaseMetaFiles()` to read `.release/v*.json`.
2. If results are found, convert to `[]model.ReleaseRecord`, sort by
   `CreatedAt DESC`, and mark `IsLatest` from `latest.json`.
3. If no `.release/` files exist, open the database and call
   `db.ListReleases()`.
4. Apply `--source` filter if provided.
5. Apply `--limit N` if provided and N > 0.
6. Render output in terminal or JSON format.
7. If no releases are found, print `"No releases found."` and exit 0.

## Terminal Output

Table format with columns: Version, Tag, Branch, Draft, Latest, Source, Date.

```
Releases (3 found)
──────────────────────────────────────────────────────────
  VERSION    TAG        BRANCH       DRAFT  LATEST  SOURCE   DATE
  2.33.0     v2.33.0    release/v2   no     yes     repo     2026-03-26
  2.31.0     v2.31.0    release/v2   no     no      repo     2026-03-20
  2.30.0     v2.30.0    release/v2   yes    no      repo     2026-03-15
```

## JSON Output Example

```json
[
  {
    "version": "2.33.0",
    "tag": "v2.33.0",
    "branch": "release/v2.33.0",
    "sourceBranch": "main",
    "commitSha": "abc123",
    "changelog": "Added --limit flag",
    "draft": false,
    "preRelease": false,
    "isLatest": true,
    "source": "repo",
    "createdAt": "2026-03-26T10:00:00Z"
  }
]
```

## Error Handling

| Condition              | Message                                          | Exit |
|------------------------|--------------------------------------------------|------|
| No .release/ + no DB   | `"No database found. Run gitmap scan first.\n"`  | 1    |
| DB open/migrate error  | `"failed to load releases: %v\n"`                | 1    |
| No releases            | `"No releases found.\n"`                         | 0    |

## Implementation Files

| File                            | Responsibility                               |
|---------------------------------|----------------------------------------------|
| `cmd/listreleases.go`           | Command handler, repo/DB loading, output     |
| `release/metadata.go`           | `ListReleaseMetaFiles()`, `ReadLatest()`     |
| `model/release.go`              | `ReleaseRecord`, `SourceRepo` constant       |
| `constants/constants_cli.go`    | `CmdListReleases`, `CmdListReleasesAlias`    |
| `constants/constants_messages.go` | Terminal output format strings             |
| `store/release.go`              | `ListReleases()` (DB fallback)               |

## Integration Points

- `cmd/root.go`: register `list-releases` / `lr` in `dispatchMisc`.
- Reuse `release.ListReleaseMetaFiles()` — no new filesystem code needed.
- DB path is only resolved when `.release/` files are absent.

## Code Style

All functions ≤ 15 lines. Positive logic. Blank line before every return.
No magic strings. No switch statements. PascalCase for SQL column names.
