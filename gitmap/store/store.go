// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// DB wraps the SQLite database connection.
type DB struct {
	conn *sql.DB
}

// Open creates or opens the SQLite database at the standard location.
func Open(outputDir string) (*DB, error) {
	dbPath := buildDBPath(outputDir)

	if err := ensureDir(filepath.Dir(dbPath)); err != nil {
		return nil, fmt.Errorf(constants.ErrDBCreateDir, err)
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBOpen, err)
	}

	return &DB{conn: conn}, enableFK(conn)
}

// Migrate creates all required tables if they don't exist.
func (db *DB) Migrate() error {
	statements := []string{
		constants.SQLCreateRepos,
		constants.SQLCreateAbsPathIndex,
		constants.SQLCreateGroups,
		constants.SQLCreateGroupRepos,
		constants.SQLCreateReleases,
		constants.SQLCreateCommitTemplates,
		constants.SQLCreateAmendments,
		constants.SQLCreateCommandHistory,
		constants.SQLCreateBookmarks,
	}

	for _, stmt := range statements {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	db.migrateSourceColumn()

	return nil
}

// migrateSourceColumn adds the Source column to existing Releases tables.
func (db *DB) migrateSourceColumn() {
	_, _ = db.conn.Exec(constants.SQLAddSourceColumn)
}

// Reset drops all tables and recreates them for a fresh start.
func (db *DB) Reset() error {
	drops := []string{
		constants.SQLDropGroupRepos,
		constants.SQLDropGroups,
		constants.SQLDropReleases,
		constants.SQLDropAmendments,
		constants.SQLDropCommitTemplates,
		constants.SQLDropCommandHistory,
		constants.SQLDropBookmarks,
		constants.SQLDropRepos,
	}

	for _, stmt := range drops {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	return db.Migrate()
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying sql.DB for advanced queries.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// buildDBPath constructs the full path to the database file.
func buildDBPath(outputDir string) string {
	return filepath.Join(outputDir, constants.DBDir, constants.DBFile)
}

// ensureDir creates the directory tree if it doesn't exist.
func ensureDir(dir string) error {
	return os.MkdirAll(dir, constants.DirPermission)
}

// enableFK turns on SQLite foreign key enforcement.
func enableFK(conn *sql.DB) error {
	_, err := conn.Exec(constants.SQLEnableFK)

	return err
}
