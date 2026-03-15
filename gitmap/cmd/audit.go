package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

// recordAuditStart inserts a new history record at command start.
func recordAuditStart(command string, args []string) (string, time.Time) {
	start := time.Now()
	id := generateAuditID()
	alias, flags, positional := classifyArgs(command, args)

	record := model.CommandHistoryRecord{
		ID:        id,
		Command:   command,
		Alias:     alias,
		Args:      positional,
		Flags:     flags,
		StartedAt: start.Format(time.RFC3339),
	}

	db, err := openAuditDB()
	if err != nil {
		return id, start
	}
	defer db.Close()

	_ = db.InsertHistory(record)

	return id, start
}

// recordAuditEnd updates a history record with completion details.
func recordAuditEnd(id string, start time.Time, exitCode int, summary string, repoCount int) {
	end := time.Now()
	duration := end.Sub(start).Milliseconds()

	record := model.CommandHistoryRecord{
		ID:         id,
		FinishedAt: end.Format(time.RFC3339),
		DurationMs: duration,
		ExitCode:   exitCode,
		Summary:    summary,
		RepoCount:  repoCount,
	}

	db, err := openAuditDB()
	if err != nil {
		return
	}
	defer db.Close()

	_ = db.UpdateHistory(record)
}

// openAuditDB opens the database silently (no error output).
func openAuditDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}

	_ = db.Migrate()

	return db, nil
}

// generateAuditID creates a timestamp-based unique ID.
func generateAuditID() string {
	return fmt.Sprintf("hist-%d", time.Now().UnixNano())
}

// classifyArgs separates flags from positional arguments.
func classifyArgs(command string, args []string) (string, string, string) {
	alias := resolveAlias(command)
	var flags, positional []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}

	return alias, strings.Join(flags, " "), strings.Join(positional, " ")
}

// resolveAlias returns the alias if the command was invoked by alias.
func resolveAlias(command string) string {
	return command
}
