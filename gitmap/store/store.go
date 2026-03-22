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
	conn *sql.DB
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
	}

	for _, stmt := range statements {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf(constants.ErrDBMigrate, err)
		}
	}

	db.migrateSourceColumn()
	db.migrateNotesColumn()
	db.migrateZipGroupItemPaths()

	return db.SeedProjectTypes()
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

// Reset drops all tables and recreates them for a fresh start.
func (db *DB) Reset() error {
	drops := []string{
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
		constants.SQLDropZipGroups,
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

// buildDBPath constructs the full path to the default database file.
func buildDBPath(outputDir string) string {
	dbFile := ActiveProfileDBFile(outputDir)

	return filepath.Join(outputDir, constants.DBDir, dbFile)
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
