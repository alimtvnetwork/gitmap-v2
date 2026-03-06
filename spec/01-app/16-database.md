# Database & Repo Storage

## Overview

After scanning, gitmap persists all discovered repositories in a local
SQLite database. The database enables slug-based lookup, repo grouping,
and batch operations across selected repos.

## SQLite Setup

### Package

Use a **CGo-free** SQLite driver. The recommended package is
`modernc.org/sqlite` (pure Go, no C compiler required).

### Database Location

| Item | Value |
|------|-------|
| Directory | `gitmap-output/data/` (created automatically) |
| File name | `gitmap.db` |
| Full path | `gitmap-output/data/gitmap.db` |

### Auto-Creation

On every `scan` completion, gitmap:

1. Checks if `gitmap-output/data/gitmap.db` exists.
2. If missing, creates the database and initializes all tables.
3. Upserts all scanned repos into the `repos` table.

---

## Data Model

### repos Table

| Column           | Type    | Constraints          | Notes                            |
|------------------|---------|----------------------|----------------------------------|
| id               | TEXT    | PRIMARY KEY          | UUID from ScanRecord             |
| slug             | TEXT    | NOT NULL             | Derived from GitHub repo name    |
| repo_name        | TEXT    | NOT NULL             | Display name                     |
| https_url        | TEXT    | NOT NULL             |                                  |
| ssh_url          | TEXT    | NOT NULL             |                                  |
| branch           | TEXT    | NOT NULL             |                                  |
| relative_path    | TEXT    | NOT NULL             |                                  |
| absolute_path    | TEXT    | NOT NULL             |                                  |
| clone_instruction| TEXT    | NOT NULL             |                                  |
| notes            | TEXT    | DEFAULT ''           |                                  |
| created_at       | TEXT    | DEFAULT CURRENT_TIMESTAMP |                             |
| updated_at       | TEXT    | DEFAULT CURRENT_TIMESTAMP |                             |

**Slug derivation:** Extract the repository name from the HTTPS URL
(e.g. `https://github.com/user/my-api.git` → `my-api`). The slug is
**not unique** — duplicate slugs are expected when the same repo name
exists under different orgs or paths.

**Upsert strategy:** On scan, match by `absolute_path`. If a row with
that path exists, update all fields. Otherwise, insert a new row.

---

## Slug Generation

`mapper.BuildRecords` populates the `Slug` field on every `ScanRecord`
during scan, so it is available in both JSON output and DB upsert.

Extract from HTTPS URL:

```
https://github.com/user/my-api.git  →  my-api
https://github.com/org/my-api.git   →  my-api  (duplicate allowed)
```

**Algorithm:**

1. Parse the HTTPS URL.
2. Take the last path segment.
3. Strip `.git` suffix if present.
4. Lowercase the result.

If the HTTPS URL is empty, fall back to `repoName`.

---

## Slug Disambiguation

Since slugs can be duplicated, gitmap must disambiguate when a command
targets a non-unique slug.

### Interactive Mode (Default)

When a duplicate slug is selected, display a numbered prompt:

```
Multiple repos match "my-api":

  1. my-api  →  /home/user/work/org-a/my-api
  2. my-api  →  /home/user/personal/my-api

Select [1-2]:
```

### Path Qualifier (Scripting)

For non-interactive or scripted use, support a `slug@path` syntax:

```
gitmap pull my-api@/home/user/work/org-a/my-api
```

If `@path` is provided and matches exactly one repo, skip the prompt.
If it matches none, print an error and exit.

---

## DB-First Lookup with JSON Fallback

Commands that resolve repos by slug (`pull`, `exec`, `status`) use a
two-tier lookup strategy:

1. **Try the database first.** Open `gitmap-output/data/gitmap.db` and
   query the `repos` table.
2. **Fall back to JSON.** If the database does not exist (no prior scan
   with DB support), load `gitmap-output/gitmap.json` and match by
   repo name as before.

This ensures backward compatibility with older scan outputs while
preferring the richer DB data when available.

### No-Database Error

When a command requires the database (e.g. `--group`, `--all`) and no
`gitmap.db` exists, print:

```
No database found. Run 'gitmap scan' first.
```

Exit with code 1.

---

## Model Additions

### ScanRecord Update

Add `Slug` field (populated by `mapper.BuildRecords`):

```go
type ScanRecord struct {
    // ... existing fields ...
    Slug string `json:"slug" csv:"slug"`
}
```

---

## Package Structure (Database)

### New Packages

| Package | Responsibility |
|---------|----------------|
| `store` | SQLite database init, connection, CRUD operations |

### New Files

| File | Contents |
|------|----------|
| `store/store.go` | DB init, open, close, migration |
| `store/repo.go` | Repo CRUD (upsert, list, find by slug) |
| `constants/constants_store.go` | DB path, table names, SQL statements |

### Updated Files

| File | Change |
|------|--------|
| `cmd/scan.go` | Trigger DB upsert after scan |
| `cmd/pull.go` | DB-first lookup with JSON fallback |
| `cmd/exec.go` | DB-first lookup with JSON fallback |
| `cmd/status.go` | DB-first lookup with JSON fallback |
| `model/record.go` | Add `Slug` field to `ScanRecord` |
| `mapper/mapper.go` | Populate `Slug` in `BuildRecords` |
| `spec/01-app/07-data-model.md` | Document `Slug` field |

---

## Error Handling (Database)

| Scenario | Behavior |
|----------|----------|
| DB file cannot be created | Print error, exit 1 |
| Slug not found | `\"No repo matches slug: %s\"` |
| No database and DB-required flag | `\"No database found. Run 'gitmap scan' first.\"` |
| Duplicate slug without qualifier | Interactive prompt (or error in non-TTY) |

---

## Constraints

- SQLite driver must be **CGo-free** (`modernc.org/sqlite`).
- All string literals in `constants` package.
- All files under 200 lines.
- All functions 8–15 lines.
- Positive conditions only (no negation).
- Blank line before `return`.
