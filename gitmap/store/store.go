// Package store manages the SQLite database for gitmap.
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection.
type DB struct {
	conn   *sql.DB
	dbDir  string
}

// Open creates or opens the SQLite database for the active profile.
func Open(outputDir string) (*DB, error) {
	dbFile := ActiveProfileDBFile(outputDir)
	dbPath := filepath.Join(outputDir, constants.DBDir, dbFile)

	return openDBAt(dbPath)
}

// OpenProfile opens the database for a specific named profile.
func OpenProfile(outputDir, profileName string) (*DB, error) {
	dbFile := ProfileDBFile(profileName)
	dbPath := filepath.Join(outputDir, constants.DBDir, dbFile)

	return openDBAt(dbPath)
}

// openDBAt opens a database at an exact path.
func openDBAt(dbPath string) (*DB, error) {
	dbDir := filepath.Dir(dbPath)
	if err := ensureDir(dbDir); err != nil {
		return nil, fmt.Errorf(constants.ErrDBCreateDir, dbDir, err)
	}

	if err := acquireLock(dbDir); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		releaseLock(dbDir)

		return nil, fmt.Errorf(constants.ErrDBOpen, dbPath, err)
	}

	// SQLite does not support concurrent writes; pin to one connection
	// so PRAGMAs (foreign_keys, etc.) persist across all operations.
	conn.SetMaxOpenConns(1)

	err = enableFK(conn)
	if err != nil {
		conn.Close()
		releaseLock(dbDir)

		return nil, err
	}

	return &DB{conn: conn, dbDir: dbDir}, nil
}

// Migrate creates all required tables if they don't exist.
func (db *DB) Migrate() error {
	db.migrateLegacyIDs()

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
		constants.SQLCreateProjectTypes,
		constants.SQLCreateDetectedProjects,
		constants.SQLCreateGoProjectMetadata,
		constants.SQLCreateGoRunnableFiles,
		constants.SQLCreateCSharpProjectMeta,
		constants.SQLCreateCSharpProjectFiles,
		constants.SQLCreateCSharpKeyFiles,
		constants.SQLCreateSettings,
		constants.SQLCreateAliases,
		constants.SQLCreateZipGroups,
		constants.SQLCreateZipGroupItems,
		constants.SQLCreateSSHKeys,
		constants.SQLCreateTempReleases,
		constants.SQLCreateInstalledTools,
		constants.SQLCreateTaskType,
		constants.SQLCreatePendingTask,
		constants.SQLCreateCompletedTask,
	}

	for _, stmt := range statements {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	db.migrateSourceColumn()
	db.migrateNotesColumn()
	db.migrateZipGroupItemPaths()
	db.migrateTRCommitSha()

	if err := db.SeedProjectTypes(); err != nil {
		return err
	}

	return db.SeedTaskTypes()
}

// migrateSourceColumn adds the Source column to existing Releases tables.
func (db *DB) migrateSourceColumn() {
	_, _ = db.conn.Exec(constants.SQLAddSourceColumn)
}

// migrateNotesColumn adds the Notes column to existing Releases tables.
func (db *DB) migrateNotesColumn() {
	_, _ = db.conn.Exec(constants.SQLAddNotesColumn)
}

// migrateZipGroupItemPaths adds RepoPath, RelativePath, FullPath columns
// to existing ZipGroupItems tables and copies Path into FullPath.
func (db *DB) migrateZipGroupItemPaths() {
	_, _ = db.conn.Exec(constants.SQLMigrateZGIRepoPath)
	_, _ = db.conn.Exec(constants.SQLMigrateZGIRelativePath)
	_, _ = db.conn.Exec(constants.SQLMigrateZGIFullPath)
	_, _ = db.conn.Exec(constants.SQLMigrateZGICopyPath)
}

// migrateTRCommitSha renames the Commit column to CommitSha in TempReleases.
func (db *DB) migrateTRCommitSha() {
	_, _ = db.conn.Exec(constants.SQLMigrateTRCommitSha)
}

// Reset drops all tables and recreates them for a fresh start.
func (db *DB) Reset() error {
	drops := []string{
		constants.SQLDropCompletedTask,
		constants.SQLDropPendingTask,
		constants.SQLDropTaskType,
		constants.SQLDropSettings,
		constants.SQLDropGoRunnableFiles,
		constants.SQLDropGoProjectMetadata,
		constants.SQLDropCSharpKeyFiles,
		constants.SQLDropCSharpProjectFiles,
		constants.SQLDropCSharpProjectMeta,
		constants.SQLDropDetectedProjects,
		constants.SQLDropProjectTypes,
		constants.SQLDropGroupRepos,
		constants.SQLDropGroups,
		constants.SQLDropReleases,
		constants.SQLDropAmendments,
		constants.SQLDropCommitTemplates,
		constants.SQLDropCommandHistory,
		constants.SQLDropBookmarks,
		constants.SQLDropAliases,
		constants.SQLDropZipGroupItems,
		constants.SQLDropTempReleases,
		constants.SQLDropZipGroups,
		constants.SQLDropInstalledTools,
		constants.SQLDropRepos,
	}

	for _, stmt := range drops {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	return db.Migrate()
}

// Close closes the database connection and releases the lock.
func (db *DB) Close() error {
	releaseLock(db.dbDir)

	return db.conn.Close()
}

// Conn returns the underlying sql.DB for advanced queries.
func (db *DB) Conn() *sql.DB {
	return db.conn
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
