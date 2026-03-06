# Database & Repo Grouping

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

### groups Table

| Column      | Type    | Constraints     | Notes                        |
|-------------|---------|-----------------|------------------------------|
| id          | TEXT    | PRIMARY KEY     | UUID                         |
| name        | TEXT    | NOT NULL UNIQUE | Group display name           |
| description | TEXT    | DEFAULT ''      | Optional description         |
| color       | TEXT    | DEFAULT ''      | Terminal color (e.g. "green") |
| created_at  | TEXT    | DEFAULT CURRENT_TIMESTAMP |                    |

### group_repos Table (Join)

| Column   | Type | Constraints                              |
|----------|------|------------------------------------------|
| group_id | TEXT | NOT NULL, FK → groups(id) ON DELETE CASCADE |
| repo_id  | TEXT | NOT NULL, FK → repos(id) ON DELETE CASCADE  |
| PRIMARY KEY | | (group_id, repo_id)                      |

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

## CLI Commands

### `gitmap list` (alias: `ls`)

Show all tracked repos with slugs.

```
SLUG                 REPO NAME
──────────────────────────────────────────
my-api               My API
my-api               My API (personal)
dashboard            Dashboard
auth-service         Auth Service
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--group` | `-g` | Filter by group name |
| `--verbose` | `-V` | Show full paths and URLs |

### `gitmap group create <name>` (alias: `g create`)

Create a new group.

```
gitmap group create backend
gitmap group create backend --description "All backend services"
gitmap group create backend --color cyan
```

### `gitmap group add <group> <slug...>` (alias: `g add`)

Add one or more repos to a group by slug.

```
gitmap group add backend my-api auth-service
```

If a slug is duplicated, disambiguation triggers (interactive or
`slug@path`).

### `gitmap group remove <group> <slug...>` (alias: `g rm`)

Remove repos from a group.

```
gitmap group remove backend my-api
```

### `gitmap group list` (alias: `g ls`)

List all groups with repo counts.

```
GROUP           REPOS   DESCRIPTION
──────────────────────────────────────────
backend         3       All backend services
frontend        2       UI applications
```

### `gitmap group show <name>` (alias: `g show`)

Show repos in a specific group.

```
Group: backend (3 repos)
  my-api           /home/user/work/my-api
  auth-service     /home/user/work/auth-service
  gateway          /home/user/work/gateway
```

### `gitmap group delete <name>` (alias: `g del`)

Delete a group (does not delete repos, only the grouping).

---

## Batch Operations on Groups

All existing repo-level commands support a `--group` (`-g`) flag and
an `--all` flag to target repos in bulk:

| Command | `--group` Example | `--all` Example |
|---------|-------------------|-----------------|
| `pull` | `gitmap pull --group backend` | `gitmap pull --all` |
| `exec` | `gitmap exec --group backend "git fetch --all"` | `gitmap exec --all "git fetch --all"` |
| `status` | `gitmap status --group backend` | `gitmap status --all` |
| `release` | `gitmap release --group backend` | — |
| `clone` | `gitmap clone json --group backend` | — |

### Selective Multi-Repo

Select specific repos by slug (comma-separated or repeated flag):

```
gitmap pull my-api,auth-service
gitmap pull my-api auth-service
gitmap status my-api auth-service
```

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

### Group (model/group.go)

```go
type Group struct {
    ID          string
    Name        string
    Description string
    Color       string
    CreatedAt   string
}
```

---

## Package Structure

### New Packages

| Package | Responsibility |
|---------|----------------|
| `store` | SQLite database init, connection, CRUD operations |

### New Files

| File | Contents |
|------|----------|
| `store/store.go` | DB init, open, close, migration |
| `store/repo.go` | Repo CRUD (upsert, list, find by slug) |
| `store/group.go` | Group CRUD (create, list, add/remove repos) |
| `model/group.go` | Group and GroupRepo structs |
| `cmd/list.go` | `list` command handler |
| `cmd/group.go` | `group` command routing |
| `cmd/groupcreate.go` | `group create` handler |
| `cmd/groupadd.go` | `group add` handler |
| `cmd/groupremove.go` | `group remove` handler |
| `cmd/grouplist.go` | `group list` handler |
| `cmd/groupshow.go` | `group show` handler |
| `cmd/groupdelete.go` | `group delete` handler |
| `constants/constants_store.go` | DB path, table names, SQL statements |

### Updated Files

| File | Change |
|------|--------|
| `cmd/root.go` | Register `list`, `group` commands |
| `cmd/pull.go` | Add `--group`, `--all` flags; DB-first lookup |
| `cmd/exec.go` | Add `--group`, `--all` flags; DB-first lookup |
| `cmd/status.go` | Add `--group`, `--all` flags; DB-first lookup |
| `cmd/scan.go` | Trigger DB upsert after scan |
| `model/record.go` | Add `Slug` field to `ScanRecord` |
| `mapper/mapper.go` | Populate `Slug` in `BuildRecords` |
| `spec/01-app/07-data-model.md` | Document new tables and `Slug` field |
| `constants/constants_cli.go` | New command names, aliases, help text |

---

## Error Handling

| Scenario | Behavior |
|----------|----------|
| DB file cannot be created | Print error, exit 1 |
| Slug not found | `"No repo matches slug: %s"` |
| Group not found | `"No group found: %s"` |
| No database and DB-required flag | `"No database found. Run 'gitmap scan' first."` |
| Duplicate slug without qualifier | Interactive prompt (or error in non-TTY) |
| Repo already in group | Silent no-op |
| Group name already exists | `"Group already exists: %s"` |

---

## Constraints

- SQLite driver must be **CGo-free** (`modernc.org/sqlite`).
- All string literals in `constants` package.
- All files under 200 lines.
- All functions 8–15 lines.
- Positive conditions only (no negation).
- Blank line before `return`.
