# Project Type Detection

## Overview

During `scan` and `rescan`, gitmap detects project types inside each
discovered Git repository. Detection results are written to dedicated
JSON files and persisted in the SQLite database. Users can query
detected projects by type using dedicated commands.

---

## Supported Project Types

| Type   | Key         | Description                   |
|--------|-------------|-------------------------------|
| Go     | `go`        | Go modules / packages         |
| Node   | `node`      | Node.js projects              |
| React  | `react`     | React applications            |
| C++    | `cpp`       | C/C++ projects                |
| C#     | `csharp`    | .NET / C# projects            |

---

## Detection Rules

### Go

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `go.mod` file exists               | High       |
| Secondary | `go.sum` file exists               | Medium     |
| Secondary | `*.go` files present               | Low        |

**Project name:** Parsed from `module` directive in `go.mod` (first
line starting with `module `). Fall back to directory name.

**False positive prevention:** Ignore `go.mod` inside `vendor/` or
`testdata/` directories.

### Node.js

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `package.json` file exists         | High       |
| Secondary | `package-lock.json` exists         | Medium     |
| Secondary | `yarn.lock` exists                 | Medium     |
| Secondary | `pnpm-lock.yaml` exists            | Medium     |
| Secondary | `bun.lockb` or `bun.lock` exists   | Medium     |

**Project name:** Parsed from `name` field in `package.json`.
Fall back to directory name.

**False positive prevention:** Ignore `package.json` inside
`node_modules/`, `vendor/`, or `testdata/` directories.

**Classification upgrade:** If `package.json` contains React or
ReactJS dependencies, the project is classified as `react` instead
of `node`. See React rules below.

### React

A Node.js project is reclassified as React when **any** of these
conditions are true:

| Condition                                           | Check Location        |
|-----------------------------------------------------|-----------------------|
| `react` in `dependencies`                           | `package.json`        |
| `react` in `devDependencies`                        | `package.json`        |
| `@types/react` in `devDependencies`                 | `package.json`        |
| `react-scripts` in `dependencies`                   | `package.json` (CRA)  |
| `next` in `dependencies`                            | `package.json` (Next) |
| `gatsby` in `dependencies`                          | `package.json`        |
| `remix` or `@remix-run/react` in `dependencies`     | `package.json`        |

**Project name:** Same as Node.js (from `package.json` `name` field).

**Note:** React projects do **not** also appear as Node.js. They are
classified exclusively as `react`.

### C++

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `CMakeLists.txt` file exists       | High       |
| Primary   | `*.vcxproj` file exists            | High       |
| Primary   | `meson.build` file exists          | High       |
| Secondary | `Makefile` with C++ content        | Medium     |
| Secondary | `*.cpp`, `*.cc`, `*.cxx` files     | Medium     |
| Secondary | `*.hpp`, `*.hh`, `*.hxx` files     | Medium     |
| Tertiary  | `conanfile.txt` or `conanfile.py`  | Medium     |
| Tertiary  | `vcpkg.json` exists                | Medium     |

**Project name:** For CMake projects, parsed from `project()`
directive in `CMakeLists.txt`. For others, fall back to directory name.

**Makefile disambiguation:** A `Makefile` alone does not trigger C++
detection. It must be accompanied by at least one C++ source file
(`*.cpp`, `*.cc`, `*.cxx`, `*.hpp`, `*.hh`, `*.hxx`) or a
`CMakeLists.txt` / `*.vcxproj` in the same directory tree.

**False positive prevention:** Ignore `build/`, `cmake-build-*/`,
`out/`, and `target/` directories.

### C#

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `*.csproj` file exists             | High       |
| Primary   | `*.sln` file exists                | High       |
| Secondary | `*.fsproj` file exists             | Medium     |
| Secondary | `global.json` file exists          | Medium     |
| Secondary | `*.cs` source files present        | Low        |

**Project name:** For `.csproj` projects, parsed from the filename
(e.g., `MyApp.csproj` → `MyApp`). For `.sln` solutions, parsed from
the filename. Fall back to directory name.

**False positive prevention:** Ignore `bin/`, `obj/`, and `packages/`
directories.

---

## Detection Scope

### Where Detection Runs

Detection scans **inside** each discovered Git repository. The scanner
walks the repo directory tree (excluding standard exclusion dirs) and
looks for indicator files.

### Monorepo Handling

A single Git repository may contain **multiple** detected projects.
For example, a monorepo with `backend/` (Go) and `frontend/` (React)
produces two separate `DetectedProject` records, both linked to the
same repo.

### Nested Projects

If a Node.js project at `./` contains a React project at `./web/`,
both are recorded. The more specific classification wins at each
path level — the root is `node`, the `web/` subdirectory is `react`.

### Exclusion Directories

The project detector skips the following directories:

| Directory            | Reason                          |
|----------------------|---------------------------------|
| `node_modules`       | Dependencies, not source        |
| `vendor`             | Vendored dependencies           |
| `.git`               | Git internals                   |
| `dist`               | Build output                    |
| `build`              | Build output                    |
| `target`             | Build output (Rust/Java/C++)    |
| `bin`                | Binary output                   |
| `obj`                | .NET build output               |
| `out`                | Generic build output            |
| `cmake-build-*`      | CMake build directories         |
| `testdata`           | Test fixtures                   |
| `packages`           | NuGet packages                  |
| `.venv`              | Python virtual environments     |
| `.cache`             | Cache directories               |

---

## Data Model

### DetectedProjects Table

| Column          | Type    | Constraints                           | Notes                                    |
|-----------------|---------|---------------------------------------|------------------------------------------|
| Id              | TEXT    | PRIMARY KEY                           | UUID                                     |
| RepoId          | TEXT    | NOT NULL, FK → Repos(Id) ON DELETE CASCADE | Link to parent repo               |
| ProjectType     | TEXT    | NOT NULL                              | `go`, `node`, `react`, `cpp`, `csharp`   |
| ProjectName     | TEXT    | NOT NULL                              | Parsed from manifest or dir name         |
| AbsolutePath    | TEXT    | NOT NULL                              | Full filesystem path to project root     |
| RepoPath        | TEXT    | NOT NULL                              | Absolute path of the Git repo root       |
| RelativePath    | TEXT    | NOT NULL                              | Path relative to repo root               |
| PrimaryIndicator| TEXT    | NOT NULL                              | File that triggered detection            |
| DetectedAt      | TEXT    | DEFAULT CURRENT_TIMESTAMP             |                                          |

**Unique constraint:** `(RepoId, ProjectType, RelativePath)` — one
entry per project type per path per repo.

**Upsert strategy:** On scan, match by the unique constraint. If a
row exists, update `ProjectName`, `AbsolutePath`, `PrimaryIndicator`,
and `DetectedAt`. Otherwise, insert a new row.

### Stale Entry Cleanup

On each scan, after upserting all detected projects for a repo, delete
any `DetectedProjects` rows for that `RepoId` that were **not**
upserted in the current scan. This handles removed projects.

---

## JSON Output

### Output Files

Each project type produces a dedicated JSON file in `gitmap-output/`:

| File                    | Contents                     |
|-------------------------|------------------------------|
| `go-projects.json`      | All detected Go projects     |
| `node-projects.json`    | All detected Node.js projects|
| `react-projects.json`   | All detected React projects  |
| `cpp-projects.json`     | All detected C++ projects    |
| `csharp-projects.json`  | All detected C# projects     |

### JSON Record Schema

```json
{
  "id": "uuid",
  "repoId": "uuid",
  "repoName": "my-api",
  "projectType": "go",
  "projectName": "github.com/user/my-api",
  "absolutePath": "/home/user/repos/my-api",
  "repoPath": "/home/user/repos/my-api",
  "relativePath": ".",
  "primaryIndicator": "go.mod",
  "detectedAt": "2026-03-11T09:54:00Z"
}
```

### Write Behavior

- Files are **overwritten** on each scan (not merged).
- Records are **sorted** by `repoName` then `relativePath`.
- Files are **only written** if at least one project of that type
  was detected. Empty files are not created.
- Write failures are logged to stderr but do not abort the scan.

---

## Commands

### Integrated into Scan

Project detection runs automatically as part of `scan` and `rescan`.
No additional flags are needed. After repo discovery and record
building, the scan pipeline adds a project detection phase.

### Query Commands

| Command             | Alias | Description                        |
|---------------------|-------|------------------------------------|
| `gitmap go-repos`   | `gr`  | List repos containing Go projects  |
| `gitmap node-repos` | `nr`  | List repos containing Node projects|
| `gitmap react-repos`| `rr`  | List repos containing React projects|
| `gitmap cpp-repos`  | `cr`  | List repos containing C++ projects |
| `gitmap csharp-repos`| `sr` | List repos containing C# projects  |

### Query Command Output

Terminal output for each detected project:

```
  go  github.com/user/my-api
      Path: /home/user/repos/my-api
      Indicator: go.mod

  go  github.com/user/my-cli/tools/linter
      Path: /home/user/repos/my-cli/tools/linter
      Indicator: go.mod
```

### Query Command Flags

| Flag       | Default    | Description                       |
|------------|------------|-----------------------------------|
| `--json`   | false      | Output as JSON instead of terminal|
| `--count`  | false      | Print count only                  |

### Query Command Data Source

Query commands read from the SQLite database. If the database does
not exist, print: `"No database found. Run 'gitmap scan' first."`

---

## Scan Pipeline Integration

### Current Scan Flow

```
1. Parse flags
2. Load config
3. ScanDir (discover repos)
4. BuildRecords (extract git metadata)
5. Write outputs (terminal, CSV, JSON, scripts)
6. Save scan cache
7. Upsert to DB
8. Import releases
9. Add to desktop
10. Open folder
```

### Extended Scan Flow

```
1. Parse flags
2. Load config
3. ScanDir (discover repos)
4. BuildRecords (extract git metadata)
5. DetectProjects (NEW — scan inside each repo)
6. Write outputs (terminal, CSV, JSON, scripts)
7. Write project JSON files (NEW)
8. Save scan cache
9. Upsert repos to DB
10. Upsert detected projects to DB (NEW)
11. Cleanup stale projects (NEW)
12. Import releases
13. Add to desktop
14. Open folder
```

### Detection Flow per Repo

```
repo root
  │
  ├─ Walk directory tree (skip exclusion dirs)
  │
  ├─ For each directory:
  │   ├─ Check for go.mod         → classify as "go"
  │   ├─ Check for package.json   → read contents
  │   │   ├─ Has react dep?       → classify as "react"
  │   │   └─ No react dep?        → classify as "node"
  │   ├─ Check for CMakeLists.txt → classify as "cpp"
  │   ├─ Check for *.vcxproj      → classify as "cpp"
  │   ├─ Check for meson.build    → classify as "cpp"
  │   ├─ Check for *.csproj       → classify as "csharp"
  │   ├─ Check for *.sln          → classify as "csharp"
  │   └─ No match                 → continue
  │
  └─ Collect all DetectedProject records
```

---

## Package Structure

### New Package

| Package    | Responsibility                              |
|------------|---------------------------------------------|
| `detector` | Walk repo trees and classify project types  |

### New Files

| File                          | Contents                               |
|-------------------------------|----------------------------------------|
| `detector/detector.go`        | Walk repo, collect detected projects   |
| `detector/rules.go`           | Detection rules per project type       |
| `detector/parser.go`          | Parse manifest files (go.mod, package.json) |
| `cmd/projectrepos.go`         | Query commands (go-repos, node-repos, etc.) |
| `cmd/projectreposoutput.go`   | Terminal and JSON output for queries   |
| `store/project.go`            | DetectedProject CRUD operations        |
| `model/project.go`            | DetectedProject struct                 |
| `constants/constants_project.go` | Project detection constants         |

### Modified Files

| File                          | Change                                 |
|-------------------------------|----------------------------------------|
| `cmd/scan.go`                 | Call detector after BuildRecords       |
| `cmd/scanoutput.go`           | Write project JSON files               |
| `cmd/root.go`                 | Register query commands in dispatch    |
| `store/store.go`              | Add DetectedProjects migration         |
| `constants/constants_store.go`| Add SQL for DetectedProjects table     |
| `constants/constants_cli.go`  | Add command names and aliases          |

---

## Model

### DetectedProject Struct

```go
type DetectedProject struct {
    ID               string `json:"id"`
    RepoID           string `json:"repoId"`
    RepoName         string `json:"repoName"`
    ProjectType      string `json:"projectType"`
    ProjectName      string `json:"projectName"`
    AbsolutePath     string `json:"absolutePath"`
    RepoPath         string `json:"repoPath"`
    RelativePath     string `json:"relativePath"`
    PrimaryIndicator string `json:"primaryIndicator"`
    DetectedAt       string `json:"detectedAt"`
}
```

---

## Constants

### Project Types

```go
const (
    ProjectTypeGo     = "go"
    ProjectTypeNode   = "node"
    ProjectTypeReact  = "react"
    ProjectTypeCpp    = "cpp"
    ProjectTypeCSharp = "csharp"
)
```

### Detection Indicator Files

```go
const (
    IndicatorGoMod        = "go.mod"
    IndicatorGoSum        = "go.sum"
    IndicatorPackageJSON  = "package.json"
    IndicatorCMakeLists   = "CMakeLists.txt"
    IndicatorMesonBuild   = "meson.build"
    IndicatorGlobalJSON   = "global.json"
)
```

### Output File Names

```go
const (
    FileGoProjects     = "go-projects.json"
    FileNodeProjects   = "node-projects.json"
    FileReactProjects  = "react-projects.json"
    FileCppProjects    = "cpp-projects.json"
    FileCSharpProjects = "csharp-projects.json"
)
```

### Command Names

```go
const (
    CmdGoRepos      = "go-repos"
    CmdGoReposAlias = "gr"
    CmdNodeRepos      = "node-repos"
    CmdNodeReposAlias = "nr"
    CmdReactRepos      = "react-repos"
    CmdReactReposAlias = "rr"
    CmdCppRepos      = "cpp-repos"
    CmdCppReposAlias = "cr"
    CmdCSharpRepos      = "csharp-repos"
    CmdCSharpReposAlias = "sr"
)
```

---

## SQL Statements

### Create Table

```sql
CREATE TABLE IF NOT EXISTS DetectedProjects (
    Id               TEXT PRIMARY KEY,
    RepoId           TEXT NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
    ProjectType      TEXT NOT NULL,
    ProjectName      TEXT NOT NULL,
    AbsolutePath     TEXT NOT NULL,
    RepoPath         TEXT NOT NULL,
    RelativePath     TEXT NOT NULL,
    PrimaryIndicator TEXT NOT NULL,
    DetectedAt       TEXT DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(RepoId, ProjectType, RelativePath)
)
```

### Upsert

```sql
INSERT INTO DetectedProjects (Id, RepoId, ProjectType, ProjectName,
    AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RepoId, ProjectType, RelativePath) DO UPDATE SET
    ProjectName=excluded.ProjectName,
    AbsolutePath=excluded.AbsolutePath,
    PrimaryIndicator=excluded.PrimaryIndicator,
    DetectedAt=CURRENT_TIMESTAMP
```

### Query by Type

```sql
SELECT dp.Id, dp.RepoId, dp.ProjectType, dp.ProjectName,
    dp.AbsolutePath, dp.RepoPath, dp.RelativePath,
    dp.PrimaryIndicator, dp.DetectedAt,
    r.RepoName
FROM DetectedProjects dp
JOIN Repos r ON dp.RepoId = r.Id
WHERE dp.ProjectType = ?
ORDER BY r.RepoName, dp.RelativePath
```

### Cleanup Stale

```sql
DELETE FROM DetectedProjects
WHERE RepoId = ? AND Id NOT IN (?, ?, ...)
```

### Drop

```sql
DROP TABLE IF EXISTS DetectedProjects
```

---

## Error Handling

| Scenario                        | Behavior                              |
|---------------------------------|---------------------------------------|
| Manifest file unreadable        | Skip project, log warning to stderr   |
| JSON parse failure              | Skip project, log warning to stderr   |
| DB upsert failure               | Log error, continue with next project |
| JSON file write failure         | Log error, continue scan              |
| No projects of a type found     | Skip JSON file creation for that type |
| No database for query command   | Print message, exit 1                 |

---

## Acceptance Criteria

### Detection

1. Go projects detected by `go.mod` presence.
2. Node.js projects detected by `package.json` presence.
3. React projects detected by `package.json` with `react` dependency.
4. C++ projects detected by `CMakeLists.txt`, `*.vcxproj`, or
   `meson.build` presence.
5. C# projects detected by `*.csproj` or `*.sln` presence.
6. Multiple projects in one repo are all detected.
7. Exclusion directories are never scanned.

### JSON Export

1. Each type produces a dedicated JSON file.
2. Records contain all specified fields.
3. No duplicates on repeated scans.
4. Empty types do not produce files.

### SQLite Storage

1. All detected projects saved in `DetectedProjects` table.
2. Upsert by `(RepoId, ProjectType, RelativePath)`.
3. Stale entries cleaned up after each scan.
4. Foreign key to `Repos` table enforced.

### Commands

1. `gitmap go-repos` returns Go projects from DB.
2. `gitmap node-repos` returns Node.js projects from DB.
3. `gitmap react-repos` returns React projects from DB.
4. `gitmap cpp-repos` returns C++ projects from DB.
5. `gitmap csharp-repos` returns C# projects from DB.
6. `--json` flag outputs JSON format.
7. `--count` flag outputs count only.

### Reliability

1. Scan completes even if one repo's detection fails.
2. Errors logged with repo path and indicator file.
3. Excluded directories skipped.
4. Extensible for future project types (add rule to `rules.go`).

---

## Constraints

- All code style rules from `spec/02-general/06-code-style-rules.md`.
- Functions 8–15 lines. Files under 200 lines.
- All string literals in `constants` package.
- Positive conditions only.
- Blank line before `return`.
- PascalCase for DB table/column names.

---

## Optional Enhancements (Future)

1. `gitmap projects` — unified command listing all types grouped.
2. `--type` flag for filtering: `gitmap projects --type go,react`.
3. Summary line after scan: `"Detected: 5 Go, 3 Node, 2 React"`.
4. Confidence score per detection.
5. Configurable detection rules via `config.json`.
6. Monorepo workspace detection (npm/yarn/pnpm workspaces).
7. Dry-run mode: `gitmap scan --detect-only`.
