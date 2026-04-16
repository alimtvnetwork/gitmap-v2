# 87 — Clone-Next Flatten Mode

## Overview

The `--flatten` flag for `gitmap clone-next` (`cn`) clones a versioned
repository into a **consistent base-name folder** instead of the default
version-suffixed folder. This keeps a single, predictable local path
across version iterations (e.g., always `macro-ahk/` instead of
`macro-ahk-v15/`, `macro-ahk-v16/`, etc.).

The flag also enables **version tracking** in the gitmap database:
both the current active version and a full transition history.

---

## Command Syntax

```
gitmap cn <version-spec> --flatten [--delete] [--yes] [--dry-run] [--verbose]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--flatten` | `-f` | Clone into base-name folder (version suffix stripped) |
| `--delete` | `-d` | Remove the current versioned folder after clone |
| `--yes` | `-y` | Skip confirmation prompts |
| `--dry-run` | — | Preview actions without executing |
| `--verbose` | — | Print detailed progress |

### Flag Interactions

| Flags Used | Behavior |
|------------|----------|
| `--flatten` only | Clone into `macro-ahk/`, replacing it if it exists |
| `--delete` only | Clone into `macro-ahk-v16/`, then delete `macro-ahk-v15/` |
| `--flatten --delete` | Clone into `macro-ahk/`, then delete `macro-ahk-v15/` (if different path) |

`--flatten` and `--delete` are independent and composable. `--flatten`
controls the **target** folder name; `--delete` controls whether the
**source** (current) folder is removed afterward.

---

## Folder Name Resolution

### Base Name Extraction

Strip the version suffix from the current repository folder name:

```
macro-ahk-v15     → macro-ahk
my-tool-v2        → my-tool
project-v100      → project
some-repo         → some-repo  (no version suffix — unchanged)
```

**Rules:**

1. Match the pattern `-v<digits>` at the end of the folder name.
2. Strip the matched suffix to produce the base name.
3. If no version suffix is found, use the folder name as-is.
4. The base name must be non-empty after stripping.

### Regex Pattern

```
^(.+)-v(\d+)$
```

- Group 1: base name (e.g., `macro-ahk`)
- Group 2: version number as string (e.g., `15`)

### Target Folder

The target clone folder is the base name only, at the same parent
directory level as the current folder:

```
# Current: /projects/macro-ahk-v15/
# Target:  /projects/macro-ahk/
```

---

## Version Number Parsing

Two representations are stored for every version:

| Field | Type | Example | Source |
|-------|------|---------|--------|
| Version tag | `TEXT` | `v16` | Extracted from target repo name |
| Version number | `INTEGER` | `16` | Parsed integer from the tag |

### Parsing Rules

1. Extract digits from the version suffix: `v16` → `16`.
2. Parse as integer. If parsing fails, store `0` and log a warning.
3. The tag always retains the `v` prefix for display consistency.

---

## Database Schema Changes

### Repos Table — New Columns

Add two columns to the existing `Repos` table via idempotent migration:

```sql
ALTER TABLE Repos ADD COLUMN CurrentVersionTag TEXT DEFAULT '';
ALTER TABLE Repos ADD COLUMN CurrentVersionNum INTEGER DEFAULT 0;
```

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `CurrentVersionTag` | `TEXT` | `''` | Full version string, e.g., `v16` |
| `CurrentVersionNum` | `INTEGER` | `0` | Integer version number, e.g., `16` |

These columns are updated **only** when `--flatten` is used. For
non-flattened clones, they remain at their default values.

### New Table: `RepoVersionHistory`

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

| Column | Type | Description |
|--------|------|-------------|
| `Id` | `INTEGER PK` | Auto-incrementing primary key |
| `RepoId` | `INTEGER FK` | References `Repos(Id)`, cascade delete |
| `FromVersionTag` | `TEXT` | Previous version tag (e.g., `v15`) |
| `FromVersionNum` | `INTEGER` | Previous version number (e.g., `15`) |
| `ToVersionTag` | `TEXT` | New version tag (e.g., `v16`) |
| `ToVersionNum` | `INTEGER` | New version number (e.g., `16`) |
| `FlattenedPath` | `TEXT` | Absolute path of the flattened folder |
| `CreatedAt` | `TEXT` | ISO 8601 timestamp of the transition |

### Migration

The migration must use the idempotent `ALTER TABLE ADD COLUMN` pattern
with `isDuplicateColumnError` to silently skip if columns already exist.
The `CREATE TABLE IF NOT EXISTS` handles the new table idempotently.

See [13-release-data-model.md](./13-release-data-model.md) for the
existing migration pattern reference.

---

## Constants

All constants must follow the project's domain-specific file convention.

### `constants_cli.go`

```go
const CloneNextFlattenFlagName = "flatten"
const CloneNextFlattenFlagShort = "f"
```

### `constants_clone_next.go` (or existing clone-next constants file)

```go
const CloneNextFlattenFlagUsage = "Clone into base-name folder (version suffix stripped)"
const CloneNextFlattenConfirmMsg = "Flatten clone into '%s'? This will replace the existing folder."
const CloneNextFlattenRemovedMsg = "Removed existing folder: %s"
const CloneNextFlattenClonedMsg = "Cloned into flattened folder: %s"
const CloneNextFlattenDBUpdatedMsg = "Version tracking updated: %s -> %s"
const CloneNextVersionParseErrMsg = "Warning: could not parse version number from '%s'"
```

### Error Messages

```go
const ErrCloneNextFlattenEmptyBaseName = "base name is empty after stripping version suffix"
const ErrCloneNextFlattenRemoveFailed = "failed to remove existing folder '%s': %w"
const ErrCloneNextFlattenCloneFailed = "failed to clone into flattened folder '%s': %w"
const ErrCloneNextFlattenDBUpdateFailed = "failed to update version tracking in database: %w"
```

---

## Workflow

### Step-by-Step Execution

```
1. Parse flags (--flatten, --delete, --yes, --dry-run, --verbose)
2. Resolve target version (existing clone-next logic)
3. IF --flatten:
   a. Extract base name from current folder (strip -vN suffix)
   b. Validate base name is non-empty
   c. Compute target path = parent_dir / base_name
   d. Parse current version tag/num from current folder name
   e. Parse target version tag/num from target repo name
   f. IF target path exists:
      i.  Prompt for confirmation (unless --yes)
      ii. Remove target folder entirely
   g. git clone <target-url> <base-name-folder>
   h. Update Repos row: CurrentVersionTag, CurrentVersionNum
   i. INSERT into RepoVersionHistory (from → to)
   j. Register with GitHub Desktop (using flattened path)
4. IF --delete AND source folder != target folder:
   a. Remove the old versioned folder
5. Shell handoff: set GITMAP_SHELL_HANDOFF to the flattened path
```

### Dry-Run Output

When `--dry-run` is active, print all planned actions without executing:

```
[dry-run] Would clone macro-ahk-v16 into /projects/macro-ahk/
[dry-run] Would remove existing folder: /projects/macro-ahk/
[dry-run] Would update DB: CurrentVersionTag=v16, CurrentVersionNum=16
[dry-run] Would insert history: v15 → v16
[dry-run] Would delete old folder: /projects/macro-ahk-v15/
```

---

## Shell Handoff

The flattened path is used for the `GITMAP_SHELL_HANDOFF` environment
variable, ensuring the parent shell navigates to the correct folder
after clone-next completes.

See the [navigation helper](./09-go-to.md) spec for the shell wrapper
mechanism.

---

## Error Handling

All errors follow the project's zero-swallow policy. Every failure must
be logged to `os.Stderr` using the standardized format.

| Scenario | Behavior |
|----------|----------|
| Base name empty after stripping | Exit with error, do not clone |
| Target folder removal fails | Exit with error, do not proceed |
| Git clone fails | Exit with error, do not update DB |
| DB update fails | Log error to stderr, do not exit (clone succeeded) |
| Version number parse fails | Log warning, store `0` as version number |
| `--flatten` with no version suffix | Use folder name as-is, version tag = `""`, num = `0` |

---

## Examples

### Basic Flatten

```bash
# In /projects/macro-ahk-v15/
gitmap cn v+1 --flatten

# Result:
#   Cloned macro-ahk-v16 → /projects/macro-ahk/
#   DB: CurrentVersionTag="v16", CurrentVersionNum=16
#   History: v15 → v16
#   Shell navigates to /projects/macro-ahk/
```

### Flatten with Delete

```bash
# In /projects/macro-ahk-v15/
gitmap cn v+1 --flatten --delete

# Result:
#   Cloned macro-ahk-v16 → /projects/macro-ahk/
#   Deleted /projects/macro-ahk-v15/
#   DB and history updated as above
```

### Flatten with Dry Run

```bash
# In /projects/macro-ahk-v15/
gitmap cn v+1 --flatten --dry-run

# Prints planned actions, executes nothing
```

### Repeated Flatten (v16 → v17)

```bash
# In /projects/macro-ahk/ (flattened from v16)
gitmap cn v+1 --flatten

# Result:
#   Removes /projects/macro-ahk/
#   Clones macro-ahk-v17 → /projects/macro-ahk/
#   DB: CurrentVersionTag="v17", CurrentVersionNum=17
#   History: v16 → v17
```

### No Version Suffix

```bash
# In /projects/some-tool/
gitmap cn v+1 --flatten

# Result:
#   Base name = "some-tool" (unchanged)
#   Clones some-tool-v2 → /projects/some-tool/
#   DB: CurrentVersionTag="v2", CurrentVersionNum=2
#   History: "" → v2
```

---

## Acceptance Criteria

1. **Folder creation**: `--flatten` clones into base-name folder without
   version suffix at the same parent directory level.
2. **Folder replacement**: If the base-name folder already exists, it is
   fully removed before cloning (with user confirmation unless `--yes`).
3. **Version suffix parsing**: Correctly strips `-v<digits>` suffix from
   folder names; handles missing suffix gracefully.
4. **DB current version**: `Repos.CurrentVersionTag` and
   `Repos.CurrentVersionNum` are updated on every `--flatten` clone.
5. **DB history**: A new `RepoVersionHistory` row is inserted for each
   transition, with correct from/to version data and flattened path.
6. **Flag independence**: `--flatten` and `--delete` work independently
   and in combination without conflict.
7. **Dry-run**: All planned actions are printed; nothing is executed.
8. **Shell handoff**: `GITMAP_SHELL_HANDOFF` is set to the flattened path.
9. **GitHub Desktop**: The flattened path is registered correctly.
10. **Error handling**: All failures logged to stderr; clone failure
    prevents DB update; DB failure does not roll back a successful clone.
11. **Idempotent migration**: Schema changes use `isDuplicateColumnError`
    and `CREATE TABLE IF NOT EXISTS`.
12. **Constants**: All flag names, messages, and error strings are defined
    in the constants package — no magic strings in command logic.

---

## Component Mapping

| Component | File / Package | Responsibility |
|-----------|---------------|----------------|
| Flag registration | `cmd/clone_next.go` | Register `--flatten` and `-f` flags |
| Base name parser | `internal/pathutil/` or `cmd/clone_next.go` | Strip `-v<digits>` suffix, extract version |
| Folder operations | `cmd/clone_next.go` | Remove existing folder, clone into base name |
| DB migration | `store/migrations.go` | Add columns to `Repos`, create `RepoVersionHistory` |
| DB write (current) | `store/repos.go` | Update `CurrentVersionTag` and `CurrentVersionNum` |
| DB write (history) | `store/repo_version_history.go` | Insert transition row |
| Constants (flags) | `constants/constants_cli.go` | `CloneNextFlattenFlagName`, `CloneNextFlattenFlagShort` |
| Constants (messages) | `constants/constants_clone_next.go` | Usage text, confirmation, error messages |
| Shell handoff | `cmd/clone_next.go` | Set `GITMAP_SHELL_HANDOFF` to flattened path |
| GitHub Desktop | `cmd/clone_next.go` | Register flattened path |

---

## Cross-References

| Document | Relevance |
|----------|-----------|
| [09-go-to.md](./09-go-to.md) | Shell handoff via `GITMAP_SHELL_HANDOFF` |
| [13-release-data-model.md](./13-release-data-model.md) | Migration pattern, version parsing reference |
| [12-release-command.md](./12-release-command.md) | Version resolution logic, `--dry-run` pattern |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
