package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// retryPendingTask executes a single pending task based on its type.
func retryPendingTask(db *store.DB, taskID int64, typeName, targetPath, workDir, cmdArgs string) {
	if isDeleteTaskType(typeName) {
		retryDeleteTask(db, taskID, targetPath)

		return
	}

	if isReplayableTaskType(typeName) {
		retryReplayTask(db, taskID, workDir, cmdArgs)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrTaskTypeNotFound, typeName)
}

// isDeleteTaskType returns true for Delete or Remove task types.
func isDeleteTaskType(typeName string) bool {
	return typeName == constants.TaskTypeDelete || typeName == constants.TaskTypeRemove
}

// isReplayableTaskType returns true for task types that can be replayed via CLI.
func isReplayableTaskType(typeName string) bool {
	return typeName == constants.TaskTypeScan ||
		typeName == constants.TaskTypeClone ||
		typeName == constants.TaskTypePull ||
		typeName == constants.TaskTypeExec
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

// retryReplayTask re-executes a stored CLI command.
func retryReplayTask(db *store.DB, taskID int64, workDir, cmdArgs string) {
	args := strings.Fields(cmdArgs)
	if len(args) == 0 {
		reason := fmt.Sprintf(constants.ReasonReplayFailed, "empty command args")
		fmt.Printf(constants.MsgPendingTaskFailed, taskID, reason)
		_ = db.FailTask(taskID, reason)

		return
	}

	fmt.Printf(constants.MsgPendingReplaying, cmdArgs)

	binaryPath, err := os.Executable()
	if err != nil {
		binaryPath = "gitmap"
	}

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if workDir != "" {
		cmd.Dir = workDir
	}

	err = cmd.Run()
	if err != nil {
		reason := fmt.Sprintf(constants.ReasonReplayFailed, err)
		fmt.Printf(constants.MsgPendingTaskFailed, taskID, reason)
		_ = db.FailTask(taskID, reason)

		return
	}

	fmt.Printf(constants.MsgPendingTaskCompleted, taskID, cmdArgs)
	_ = db.CompleteTask(taskID)
}
