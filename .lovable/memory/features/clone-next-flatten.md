---
name: Clone-next flatten mode
description: New --flatten flag for clone-next that clones into base-name folder (no version suffix) and tracks version history in DB
type: feature
---

## Feature: `--flatten` flag for `clone-next`

### Behavior

When `gitmap cn v++ --flatten` is used on a repo like `macro-ahk-v15`:

1. Target clone folder is `macro-ahk/` (base name only, version suffix stripped)
2. If `macro-ahk/` already exists, remove it entirely first
3. Clone the target repo (e.g., `macro-ahk-v16`) into `macro-ahk/`
4. The remote URL still points to `macro-ahk-v16` on GitHub — only the local folder name is flattened
5. Update the database with the new version info

### Flag

- `--flatten` — Clone into base-name folder (without version suffix), replacing any existing folder with that name
- Separate from `--delete` (which removes the *current versioned* folder after clone)
- Can be combined: `--flatten --delete` would flatten AND remove the old versioned folder if different

### Database Schema Changes

#### Repos table — new columns

| Column | Type | Description |
|--------|------|-------------|
| `CurrentVersionTag` | `TEXT DEFAULT ''` | Full version string, e.g., "v16" |
| `CurrentVersionNum` | `INTEGER DEFAULT 0` | Integer version number, e.g., 16 |

#### New table: `RepoVersionHistory`

```sql
CREATE TABLE IF NOT EXISTS RepoVersionHistory (
    Id              INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoId          INTEGER NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
    FromVersionTag  TEXT NOT NULL,
    FromVersionNum  INTEGER NOT NULL,
    ToVersionTag    TEXT NOT NULL,
    ToVersionNum    INTEGER NOT NULL,
    FlattenedPath   TEXT DEFAULT '',
    CreatedAt       TEXT DEFAULT CURRENT_TIMESTAMP
);
```

Tracks every version transition (e.g., v15→v16, v16→v17) as separate rows.

### Workflow

1. Parse `--flatten` flag
2. Resolve target version (same as current logic)
3. Compute target folder = base name only (e.g., `macro-ahk/`)
4. If folder exists, remove it
5. `git clone <target-url> <base-name-folder>`
6. Update `Repos` row: set `CurrentVersionTag` and `CurrentVersionNum`
7. Insert row into `RepoVersionHistory` with from/to versions
8. Register with GitHub Desktop (using flattened path)

### Examples

```
# Before: in macro-ahk-v15/
gitmap cn v++ --flatten

# Clones macro-ahk-v16 into macro-ahk/
# DB: CurrentVersionTag="v16", CurrentVersionNum=16
# History: v15 → v16
```
