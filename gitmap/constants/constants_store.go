package constants

// Database location.
const (
	DBDir  = "data"
	DBFile = "gitmap.db"
)

// Table names.
const (
	TableRepos     = "repos"
	TableGroups    = "groups"
	TableGroupRepo = "group_repos"
)

// SQL: create repos table.
const SQLCreateRepos = `CREATE TABLE IF NOT EXISTS repos (
	id               TEXT PRIMARY KEY,
	slug             TEXT NOT NULL,
	repo_name        TEXT NOT NULL,
	https_url        TEXT NOT NULL,
	ssh_url          TEXT NOT NULL,
	branch           TEXT NOT NULL,
	relative_path    TEXT NOT NULL,
	absolute_path    TEXT NOT NULL,
	clone_instruction TEXT NOT NULL,
	notes            TEXT DEFAULT '',
	created_at       TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at       TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create groups table.
const SQLCreateGroups = `CREATE TABLE IF NOT EXISTS groups (
	id          TEXT PRIMARY KEY,
	name        TEXT NOT NULL UNIQUE,
	description TEXT DEFAULT '',
	color       TEXT DEFAULT '',
	created_at  TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create group_repos join table.
const SQLCreateGroupRepos = `CREATE TABLE IF NOT EXISTS group_repos (
	group_id TEXT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
	repo_id  TEXT NOT NULL REFERENCES repos(id) ON DELETE CASCADE,
	PRIMARY KEY (group_id, repo_id)
)`

// SQL: enable foreign keys.
const SQLEnableFK = "PRAGMA foreign_keys = ON"

// SQL: repo operations.
const (
	SQLUpsertRepo = `INSERT INTO repos (id, slug, repo_name, https_url, ssh_url, branch, relative_path, absolute_path, clone_instruction, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			slug=excluded.slug, repo_name=excluded.repo_name, https_url=excluded.https_url,
			ssh_url=excluded.ssh_url, branch=excluded.branch, relative_path=excluded.relative_path,
			absolute_path=excluded.absolute_path, clone_instruction=excluded.clone_instruction,
			notes=excluded.notes, updated_at=CURRENT_TIMESTAMP`

	SQLSelectAllRepos = "SELECT id, slug, repo_name, https_url, ssh_url, branch, relative_path, absolute_path, clone_instruction, notes FROM repos ORDER BY slug"

	SQLSelectRepoBySlug = "SELECT id, slug, repo_name, https_url, ssh_url, branch, relative_path, absolute_path, clone_instruction, notes FROM repos WHERE slug = ?"

	SQLSelectRepoByPath = "SELECT id, slug, repo_name, https_url, ssh_url, branch, relative_path, absolute_path, clone_instruction, notes FROM repos WHERE absolute_path = ?"
)

// SQL: upsert by absolute_path (spec requirement).
const SQLUpsertRepoByPath = `INSERT INTO repos (id, slug, repo_name, https_url, ssh_url, branch, relative_path, absolute_path, clone_instruction, notes)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(absolute_path) DO UPDATE SET
		slug=excluded.slug, repo_name=excluded.repo_name, https_url=excluded.https_url,
		ssh_url=excluded.ssh_url, branch=excluded.branch, relative_path=excluded.relative_path,
		clone_instruction=excluded.clone_instruction, notes=excluded.notes, updated_at=CURRENT_TIMESTAMP`

// SQL: create unique index on absolute_path for upsert-by-path.
const SQLCreateAbsPathIndex = "CREATE UNIQUE INDEX IF NOT EXISTS idx_repos_absolute_path ON repos(absolute_path)"

// SQL: group operations.
const (
	SQLInsertGroup = "INSERT INTO groups (id, name, description, color) VALUES (?, ?, ?, ?)"

	SQLSelectAllGroups = "SELECT id, name, description, color, created_at FROM groups ORDER BY name"

	SQLSelectGroupByName = "SELECT id, name, description, color, created_at FROM groups WHERE name = ?"

	SQLDeleteGroup = "DELETE FROM groups WHERE name = ?"

	SQLInsertGroupRepo = "INSERT OR IGNORE INTO group_repos (group_id, repo_id) VALUES (?, ?)"

	SQLDeleteGroupRepo = "DELETE FROM group_repos WHERE group_id = ? AND repo_id = ?"

	SQLSelectGroupRepos = `SELECT r.id, r.slug, r.repo_name, r.https_url, r.ssh_url, r.branch,
		r.relative_path, r.absolute_path, r.clone_instruction, r.notes
		FROM repos r JOIN group_repos gr ON r.id = gr.repo_id WHERE gr.group_id = ? ORDER BY r.slug`

	SQLCountGroupRepos = "SELECT COUNT(*) FROM group_repos WHERE group_id = ?"
)

// SQL: reset operations.
const (
	SQLDropGroupRepos = "DROP TABLE IF EXISTS group_repos"
	SQLDropGroups     = "DROP TABLE IF EXISTS groups"
	SQLDropRepos      = "DROP TABLE IF EXISTS repos"
)

// Store error messages.
const (
	ErrDBOpen        = "failed to open database: %v"
	ErrDBMigrate     = "failed to initialize tables: %v"
	ErrDBUpsert      = "failed to upsert repo: %v"
	ErrDBQuery       = "failed to query repos: %v"
	ErrDBNoMatch     = "no repo matches slug: %s\n"
	ErrDBCreateDir   = "failed to create database directory: %v"
	ErrDBGroupCreate = "failed to create group: %v"
	ErrDBGroupQuery  = "failed to query groups: %v"
	ErrDBGroupAdd    = "failed to add repo to group: %v"
	ErrDBGroupRemove = "failed to remove repo from group: %v"
	ErrDBGroupDelete = "failed to delete group: %v"
	ErrDBGroupNone   = "no group found: %s"
	ErrDBGroupExists = "group already exists: %s"
)
