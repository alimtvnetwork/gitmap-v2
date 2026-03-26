package store

import (
	"fmt"

	"github.com/user/gitmap/constants"
)

// migrateLegacyIDs detects if the Repos table uses TEXT IDs (legacy UUIDs)
// and rebuilds it with INTEGER PRIMARY KEY AUTOINCREMENT. All tables with
// foreign keys to Repos.Id are dropped and recreated by the normal migration.
func (db *DB) migrateLegacyIDs() {
	if !db.hasLegacyTextID(constants.TableRepos) {
		return
	}

	fmt.Println(constants.MsgLegacyIDMigrationStart)

	db.dropProjectTables()
	db.dropGroupRepos()
	db.rebuildReposTable()

	fmt.Println(constants.MsgLegacyIDMigrationDone)
}

// hasLegacyTextID checks if the Id column of a table is TEXT via PRAGMA.
func (db *DB) hasLegacyTextID(table string) bool {
	query := fmt.Sprintf("PRAGMA table_info(%s)", table)
	rows, err := db.conn.Query(query)
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dflt interface{}
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dflt, &pk); err != nil {
			continue
		}
		if name == "Id" && colType == "TEXT" {
			return true
		}
	}

	return false
}

// dropProjectTables removes all project detection tables (FK to Repos).
func (db *DB) dropProjectTables() {
	drops := []string{
		constants.SQLDropGoRunnableFiles,
		constants.SQLDropGoProjectMetadata,
		constants.SQLDropCSharpKeyFiles,
		constants.SQLDropCSharpProjectFiles,
		constants.SQLDropCSharpProjectMeta,
		constants.SQLDropDetectedProjects,
		constants.SQLDropProjectTypes,
	}

	for _, stmt := range drops {
		_, _ = db.conn.Exec(stmt)
	}
}

// dropGroupRepos removes the GroupRepos join table (FK to Repos).
func (db *DB) dropGroupRepos() {
	_, _ = db.conn.Exec(constants.SQLDropGroupRepos)
}

// rebuildReposTable recreates Repos with INTEGER PRIMARY KEY AUTOINCREMENT,
// preserving all data except the old UUID IDs.
func (db *DB) rebuildReposTable() {
	_, _ = db.conn.Exec("ALTER TABLE Repos RENAME TO Repos_legacy")

	_, _ = db.conn.Exec(constants.SQLCreateRepos)
	_, _ = db.conn.Exec(constants.SQLCreateAbsPathIndex)

	_, _ = db.conn.Exec(`INSERT INTO Repos (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes)
		SELECT Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes
		FROM Repos_legacy`)

	_, _ = db.conn.Exec("DROP TABLE IF EXISTS Repos_legacy")
}
