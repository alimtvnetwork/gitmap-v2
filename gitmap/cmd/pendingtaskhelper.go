package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// createPendingTask inserts a pending task into the database.
// Returns the task ID and DB handle (caller must close), or 0 on failure.
func createPendingTask(typeName, targetPath, sourceCmd string) (int64, *store.DB) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingDBOpen, err)

		return 0, nil
	}

	typeID, err := db.GetTaskTypeID(typeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingTypeLookup, err)
		db.Close()

		return 0, nil
	}

	existing := db.FindPendingTaskDuplicate(typeID, targetPath)
	if existing > 0 {
		fmt.Fprintf(os.Stderr, constants.ErrPendingTaskExists, typeName, targetPath, existing)

		return existing, db
	}

	taskID, err := db.InsertPendingTask(typeID, targetPath, sourceCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingInsertFailed, err)
		db.Close()

		return 0, nil
	}

	return taskID, db
}

// completePendingTask moves a pending task to the completed table.
func completePendingTask(db *store.DB, taskID int64) {
	if db == nil || taskID == 0 {
		return
	}

	err := db.CompleteTask(taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingCompleteFail, taskID, err)
	}
}

// failPendingTask updates the failure reason for a pending task.
func failPendingTask(db *store.DB, taskID int64, reason string) {
	if db == nil || taskID == 0 {
		return
	}

	err := db.FailTask(taskID, reason)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnPendingFailUpdate, taskID, err)
	}
}
