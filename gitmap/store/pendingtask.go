// Package store — pendingtask.go manages PendingTask and CompletedTask CRUD.
package store

import (
	"database/sql"
	"fmt"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// InsertPendingTask creates a new pending task and returns its ID.
func (db *DB) InsertPendingTask(taskTypeID int64, targetPath, sourceCmd string) (int64, error) {
	result, err := db.conn.Exec(constants.SQLInsertPendingTask,
		taskTypeID, targetPath, sourceCmd)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrPendingTaskInsert, err)
	}

	return result.LastInsertId()
}

// FindPendingTaskDuplicate checks if a pending task already exists for the given type and path.
// Returns the existing task ID or 0 if none found.
func (db *DB) FindPendingTaskDuplicate(taskTypeID int64, targetPath string) int64 {
	row := db.conn.QueryRow(constants.SQLSelectPendingTaskByTypePath,
		taskTypeID, targetPath)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		return 0
	}

	return id
}

// ListPendingTasks returns all pending tasks ordered by ID.
func (db *DB) ListPendingTasks() ([]model.PendingTaskRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllPendingTasks)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
	}
	defer rows.Close()

	return scanPendingTaskRows(rows)
}

// FindPendingTaskByID returns a single pending task by ID.
func (db *DB) FindPendingTaskByID(id int64) (model.PendingTaskRecord, error) {
	row := db.conn.QueryRow(constants.SQLSelectPendingTaskByID, id)

	var r model.PendingTaskRecord

	err := row.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
		&r.SourceCommand, &r.FailureReason, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return r, fmt.Errorf(constants.ErrPendingTaskQuery, err)
	}

	return r, nil
}

// CompleteTask moves a pending task to the completed table in a transaction.
func (db *DB) CompleteTask(taskID int64) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	task, err := findPendingTaskInTx(tx, taskID)
	if err != nil {
		_ = tx.Rollback()

		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	_, err = tx.Exec(constants.SQLInsertCompletedTask,
		task.ID, task.TaskTypeId, task.TargetPath, task.SourceCommand, task.CreatedAt)
	if err != nil {
		_ = tx.Rollback()

		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	_, err = tx.Exec(constants.SQLDeletePendingTask, taskID)
	if err != nil {
		_ = tx.Rollback()

		return fmt.Errorf(constants.ErrPendingTaskComplete, err)
	}

	return tx.Commit()
}

// FailTask updates the failure reason for a pending task.
func (db *DB) FailTask(taskID int64, reason string) error {
	_, err := db.conn.Exec(constants.SQLUpdatePendingTaskFailure, reason, taskID)
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskFail, err)
	}

	return nil
}

// ListCompletedTasks returns all completed tasks ordered by completion time.
func (db *DB) ListCompletedTasks() ([]model.CompletedTaskRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllCompletedTasks)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
	}
	defer rows.Close()

	return scanCompletedTaskRows(rows)
}

// findPendingTaskInTx reads a pending task within an existing transaction.
func findPendingTaskInTx(tx *sql.Tx, id int64) (model.PendingTaskRecord, error) {
	row := tx.QueryRow(constants.SQLSelectPendingTaskByID, id)

	var r model.PendingTaskRecord

	err := row.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
		&r.SourceCommand, &r.FailureReason, &r.CreatedAt, &r.UpdatedAt)

	return r, err
}

// scanPendingTaskRows reads all rows into PendingTaskRecord slices.
func scanPendingTaskRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.PendingTaskRecord, error) {
	var results []model.PendingTaskRecord

	for rows.Next() {
		var r model.PendingTaskRecord

		err := rows.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
			&r.SourceCommand, &r.FailureReason, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
		}

		results = append(results, r)
	}

	return results, nil
}

// scanCompletedTaskRows reads all rows into CompletedTaskRecord slices.
func scanCompletedTaskRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.CompletedTaskRecord, error) {
	var results []model.CompletedTaskRecord

	for rows.Next() {
		var r model.CompletedTaskRecord

		err := rows.Scan(&r.ID, &r.OriginalTaskId, &r.TaskTypeId, &r.TaskTypeName,
			&r.TargetPath, &r.SourceCommand, &r.CompletedAt, &r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
		}

		results = append(results, r)
	}

	return results, nil
}
