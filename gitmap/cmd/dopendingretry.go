package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// retryPendingTask executes a single pending task based on its type.
func retryPendingTask(db *store.DB, taskID int64, typeName, targetPath string) {
	if typeName == constants.TaskTypeDelete || typeName == constants.TaskTypeRemove {
		retryDeleteTask(db, taskID, targetPath)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrTaskTypeNotFound, typeName)
}

// retryDeleteTask attempts to delete the target path for a pending task.
func retryDeleteTask(db *store.DB, taskID int64, targetPath string) {
	err := os.RemoveAll(targetPath)
	if err != nil {
		reason := fmt.Sprintf(constants.ReasonRetryFailed, err)
		fmt.Printf(constants.MsgPendingTaskFailed, taskID, reason)
		_ = db.FailTask(taskID, reason)

		return
	}

	fmt.Printf(constants.MsgPendingTaskCompleted, taskID, targetPath)
	_ = db.CompleteTask(taskID)
}
