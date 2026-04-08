# Pending Task Workflow

## Overview

Every delete or remove operation in gitmap (e.g. `clone-next --delete`,
future file/folder removal paths) must be recorded as a task in SQLite
**before** the actual operation is attempted. This ensures no work is
silently lost when deletion fails due to locks, permissions, or missing
targets.

## Database Schema

All tables use PascalCase. Primary keys are auto-incrementing integers.

### TaskType

Normalized lookup table for task categories.

```sql
CREATE TABLE IF NOT EXISTS TaskType (
    Id   INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
);
```

Seed values: `Delete`, `Remove`.

### PendingTask

Holds every task that has not yet completed successfully.

```sql
CREATE TABLE IF NOT EXISTS PendingTask (
    Id            INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskTypeId    INTEGER NOT NULL REFERENCES TaskType(Id),
    TargetPath    TEXT    NOT NULL,
    SourceCommand TEXT    NOT NULL,
    FailureReason TEXT    DEFAULT '',
    CreatedAt     TEXT    DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt     TEXT    DEFAULT CURRENT_TIMESTAMP
);
```

### CompletedTask

Archive of successfully executed tasks.

```sql
CREATE TABLE IF NOT EXISTS CompletedTask (
    Id             INTEGER PRIMARY KEY AUTOINCREMENT,
    OriginalTaskId INTEGER NOT NULL,
    TaskTypeId     INTEGER NOT NULL REFERENCES TaskType(Id),
    TargetPath     TEXT    NOT NULL,
    SourceCommand  TEXT    NOT NULL,
    CompletedAt    TEXT    DEFAULT CURRENT_TIMESTAMP,
    CreatedAt      TEXT    NOT NULL
);
```

### Relationships

- `PendingTask.TaskTypeId` → `TaskType.Id`
- `CompletedTask.TaskTypeId` → `TaskType.Id`
- `CompletedTask.OriginalTaskId` stores the original `PendingTask.Id`
  value captured before the row is removed from `PendingTask`.

## Execution Lifecycle

### 1. Task Creation

Before any delete/remove operation:

1. Resolve `TaskTypeId` from `TaskType` (e.g. `Delete`).
2. Insert row into `PendingTask` with `TargetPath`, `SourceCommand`.
3. Only after successful insert, attempt the actual operation.

### 2. Success Path

1. Insert row into `CompletedTask` copying all fields + `CompletedAt`.
2. Delete row from `PendingTask`.
3. Both steps inside a single transaction.

### 3. Failure Path

1. Row stays in `PendingTask`.
2. Update `FailureReason` with human-readable context.
3. Update `UpdatedAt` to current timestamp.

### 4. Missing Target

If the target file/folder does not exist at retry time, the task remains
pending. The user must explicitly decide whether to dismiss or force-complete
it. This prevents silent data loss when paths are relocated.

## CLI Commands

### `gitmap pending`

Display all rows in `PendingTask`.

Output columns: Id, Type, TargetPath, SourceCommand, FailureReason, CreatedAt.

### `gitmap do-pending` (alias `dp`)

Retry all pending tasks. Each success moves to `CompletedTask`.
Each failure updates `FailureReason` and remains pending.

### `gitmap do-pending <id>`

Retry a single pending task by its integer Id.

## Duplicate Prevention

If a `PendingTask` already exists with the same `TaskTypeId` and
`TargetPath`, do not create a duplicate. Log a message indicating the
existing pending task Id.

## Integration Points

### clone-next (`cn`)

When `--delete` flag triggers folder removal:

1. Create `PendingTask` (Type=Delete, SourceCommand=clone-next).
2. Attempt removal (with lock-check retry logic).
3. On success → `CompletedTask`. On failure → keep pending.

### Future Commands

Any command that removes files/folders must follow the same pattern.
The `pending` package provides `CreateTask`, `CompleteTask`, and
`FailTask` helpers.

## Help Integration

### Standard Help (`gitmap help`)

Add to the "Data" or "Maintenance" group:

```
  pending              List pending tasks
  do-pending (dp)      Retry pending tasks
```

### Detailed Help (`gitmap pending --help`)

Show task lifecycle explanation, output format, and examples.

### UI Help

The documentation site must include a "Pending Tasks" section explaining:
- Why deletes may remain pending (locks, permissions, missing targets).
- How to inspect pending items.
- How to retry all or by specific Id.

## Constants

All SQL, error messages, and format strings must be in
`constants/constants_pending_task.go`. No magic strings in logic files.

## Store Package

```
store/
├── pendingtask.go      Insert, list, complete, fail, find by Id
└── tasktype.go         Seed, lookup by name
```

## Model Package

```
model/
├── pendingtask.go      PendingTask, CompletedTask structs
└── tasktype.go         TaskType struct
```

## Acceptance Criteria

1. Every delete/remove inserts into `PendingTask` before execution.
2. No delete path bypasses task creation.
3. Failed deletions remain visible in `PendingTask`.
4. Successful deletions appear in `CompletedTask` and are removed from `PendingTask`.
5. `gitmap pending` lists all pending tasks with Id, type, path, reason.
6. `gitmap do-pending` retries all; `gitmap dp` is an alias.
7. `gitmap do-pending <id>` retries a single task.
8. Duplicate pending tasks for the same type+path are prevented.
9. All commands appear in standard help, detailed help, and UI help.
10. Database uses PascalCase, integer PKs, and FK constraints.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
