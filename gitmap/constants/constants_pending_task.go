package constants

// Pending task table names.
const (
	TableTaskType      = "TaskType"
	TablePendingTask   = "PendingTask"
	TableCompletedTask = "CompletedTask"
)

// Pending task type seed values.
const (
	TaskTypeDelete = "Delete"
	TaskTypeRemove = "Remove"
)

// SQL: create TaskType table.
const SQLCreateTaskType = `CREATE TABLE IF NOT EXISTS TaskType (
	Id   INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT NOT NULL UNIQUE
)`

// SQL: create PendingTask table.
const SQLCreatePendingTask = `CREATE TABLE IF NOT EXISTS PendingTask (
	Id            INTEGER PRIMARY KEY AUTOINCREMENT,
	TaskTypeId    INTEGER NOT NULL REFERENCES TaskType(Id),
	TargetPath    TEXT    NOT NULL,
	SourceCommand TEXT    NOT NULL,
	FailureReason TEXT    DEFAULT '',
	CreatedAt     TEXT    DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt     TEXT    DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create CompletedTask table.
const SQLCreateCompletedTask = `CREATE TABLE IF NOT EXISTS CompletedTask (
	Id             INTEGER PRIMARY KEY AUTOINCREMENT,
	OriginalTaskId INTEGER NOT NULL,
	TaskTypeId     INTEGER NOT NULL REFERENCES TaskType(Id),
	TargetPath     TEXT    NOT NULL,
	SourceCommand  TEXT    NOT NULL,
	CompletedAt    TEXT    DEFAULT CURRENT_TIMESTAMP,
	CreatedAt      TEXT    NOT NULL
)`

// SQL: seed TaskType values.
const SQLSeedTaskTypes = `INSERT OR IGNORE INTO TaskType (Name) VALUES ('Delete'), ('Remove')`

// SQL: pending task operations.
const (
	SQLInsertPendingTask = `INSERT INTO PendingTask (TaskTypeId, TargetPath, SourceCommand)
		VALUES (?, ?, ?)`

	SQLSelectAllPendingTasks = `SELECT p.Id, p.TaskTypeId, t.Name, p.TargetPath, p.SourceCommand,
		p.FailureReason, p.CreatedAt, p.UpdatedAt
		FROM PendingTask p JOIN TaskType t ON p.TaskTypeId = t.Id
		ORDER BY p.Id`

	SQLSelectPendingTaskByID = `SELECT p.Id, p.TaskTypeId, t.Name, p.TargetPath, p.SourceCommand,
		p.FailureReason, p.CreatedAt, p.UpdatedAt
		FROM PendingTask p JOIN TaskType t ON p.TaskTypeId = t.Id
		WHERE p.Id = ?`

	SQLSelectPendingTaskByTypePath = `SELECT p.Id FROM PendingTask p
		WHERE p.TaskTypeId = ? AND p.TargetPath = ?`

	SQLUpdatePendingTaskFailure = `UPDATE PendingTask
		SET FailureReason = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE Id = ?`

	SQLDeletePendingTask = `DELETE FROM PendingTask WHERE Id = ?`
)

// SQL: completed task operations.
const (
	SQLInsertCompletedTask = `INSERT INTO CompletedTask
		(OriginalTaskId, TaskTypeId, TargetPath, SourceCommand, CreatedAt)
		VALUES (?, ?, ?, ?, ?)`

	SQLSelectAllCompletedTasks = `SELECT c.Id, c.OriginalTaskId, c.TaskTypeId, t.Name,
		c.TargetPath, c.SourceCommand, c.CompletedAt, c.CreatedAt
		FROM CompletedTask c JOIN TaskType t ON c.TaskTypeId = t.Id
		ORDER BY c.CompletedAt DESC`
)

// SQL: task type lookup.
const SQLSelectTaskTypeByName = `SELECT Id FROM TaskType WHERE Name = ?`

// SQL: drop pending task tables.
const (
	SQLDropCompletedTask = "DROP TABLE IF EXISTS CompletedTask"
	SQLDropPendingTask   = "DROP TABLE IF EXISTS PendingTask"
	SQLDropTaskType      = "DROP TABLE IF EXISTS TaskType"
)

// Pending task error messages.
const (
	ErrPendingTaskInsert   = "failed to insert pending task: %v (operation: insert)"
	ErrPendingTaskQuery    = "failed to query pending tasks: %v (operation: query)"
	ErrPendingTaskComplete = "failed to complete task: %v (operation: complete)"
	ErrPendingTaskFail     = "failed to update task failure: %v (operation: update)"
	ErrPendingTaskNotFound = "pending task not found: %d\n"
	ErrTaskTypeNotFound    = "task type not found: %s"
	ErrPendingTaskExists   = "pending task already exists for %s at %s (Id %d)\n"
)

// Pending task warning messages.
const (
	WarnPendingDBOpen       = "Warning: could not open DB for task tracking: %v\n"
	WarnPendingTypeLookup   = "Warning: task type lookup failed: %v\n"
	WarnPendingInsertFailed = "Warning: could not record pending task: %v\n"
	WarnPendingCompleteFail = "Warning: could not mark task #%d complete: %v\n"
	WarnPendingFailUpdate   = "Warning: could not update task #%d failure: %v\n"
)

// Pending task failure reasons for FailureReason field.
const (
	ReasonLockScanFailed   = "lock scan failed: %v"
	ReasonNoLockingProcs   = "removal failed, no locking processes found: %v"
	ReasonUserDeclined     = "user declined to terminate locking processes"
	ReasonRetryFailed      = "retry removal failed: %v"
)

// Pending task help text.
const (
	HelpPending   = "  pending              List all pending tasks"
	HelpDoPending = "  do-pending (dp)      Retry pending tasks (all or by ID)"
)

// Pending task terminal messages.
const (
	MsgPendingTaskCreated   = "Task #%d created: %s %s\n"
	MsgPendingTaskCompleted = "Task #%d completed: %s %s\n"
	MsgPendingTaskFailed    = "Task #%d failed: %s\n"
	MsgPendingListHeader    = "Pending Tasks:\n"
	MsgPendingListRow       = "  #%-6d %-8s %-40s %s\n"
	MsgPendingListEmpty     = "No pending tasks.\n"
	MsgPendingRetryAll      = "Retrying %d pending task(s)...\n"
	MsgPendingRetryOne      = "Retrying task #%d...\n"
)
