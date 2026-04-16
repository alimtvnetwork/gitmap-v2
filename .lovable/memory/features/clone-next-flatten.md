---
name: Clone-next flatten mode
description: Clone-next flattens by default — clones into base-name folder (no version suffix) and tracks version history in DB
type: feature
---

## Feature: Default Flatten for `clone-next`

### Behavior (Default — No Flag Required)

When `gitmap cn v++` is used on a repo like `macro-ahk-v15`:

1. Target clone folder is `macro-ahk/` (base name only, version suffix stripped)
2. If `macro-ahk/` already exists, remove it entirely first (no prompt)
3. Clone the target repo (e.g., `macro-ahk-v16`) into `macro-ahk/`
4. The remote URL still points to `macro-ahk-v16` on GitHub — only the local folder name is flattened
5. Update the database with the new version info
6. Record version transition in `RepoVersionHistory`

### `gitmap clone <url>` Auto-Flatten

When cloning a versioned URL without a custom folder name:
- `gitmap clone https://github.com/user/wp-onboarding-v13` → clones into `wp-onboarding/`
- `gitmap clone https://github.com/user/wp-onboarding-v13 my-folder` → clones into `my-folder/` (no flatten)

### `--delete` Flag

- Removes the *current versioned folder* (e.g., `macro-ahk-v15/`) after clone — only applies when current folder differs from flattened path

### Database Schema

#### Repos table — version columns

| Column | Type | Description |
|--------|------|-------------|
| `CurrentVersionTag` | `TEXT DEFAULT ''` | Full version string, e.g., "v16" |
| `CurrentVersionNum` | `INTEGER DEFAULT 0` | Integer version number, e.g., 16 |

#### `RepoVersionHistory` table

Tracks every version transition (e.g., v15→v16, v16→v17) as separate rows.

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
