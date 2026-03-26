# Spec 56 — Unified `.gitmap/` Directory

## Overview

Consolidate all repo-local output directories under a single `.gitmap/` folder
at the repository root. This replaces the current split between `.release/` and
`gitmap-output/`.

## Current State

| Purpose            | Current Path          | Constant                  |
|--------------------|-----------------------|---------------------------|
| Release metadata   | `.release/`           | `DefaultReleaseDir`       |
| Scan output        | `gitmap-output/`      | `DefaultOutputDir`        |
| Output folder name | `gitmap-output`       | `DefaultOutputFolder`     |
| Verbose logs       | `gitmap-output/`      | `DefaultVerboseLogDir`    |

The SQLite database stays binary-relative (`data/`) — **not affected** by this change.

## Target State

```
.gitmap/
├── release/          ← was .release/
│   ├── latest.json
│   └── v2.34.0.json
└── output/           ← was gitmap-output/
    ├── gitmap.csv
    ├── gitmap.json
    ├── gitmap.txt
    ├── folder-structure.md
    └── clone.ps1
```

## New Constants

```go
// Root directory — all repo-local gitmap data lives here.
const GitMapDir = ".gitmap"

// Subdirectories.
const (
    ReleaseDirName = "release"
    OutputDirName  = "output"
)

// Resolved default paths.
var DefaultReleaseDir = filepath.Join(GitMapDir, ReleaseDirName)   // .gitmap/release
const DefaultOutputDir  = ".gitmap/output"
const DefaultOutputFolder = "output"
const DefaultVerboseLogDir = ".gitmap/output"
```

## Migration Checklist

### Phase 1 — Constants (1 file)

Update `constants/constants.go`:

| Old Value            | New Value              |
|----------------------|------------------------|
| `.release`           | `.gitmap/release`      |
| `./gitmap-output`    | `.gitmap/output`       |
| `gitmap-output`      | `.gitmap/output`       |

### Phase 2 — Code References (~27 files)

All Go source files already use `constants.DefaultReleaseDir`,
`constants.DefaultOutputDir`, etc. — no hardcoded paths to fix in
business logic. Files to verify:

- `release/metadata.go` — uses `DefaultReleaseDir` ✓
- `release/workflowfinalize.go` — uses `DefaultReleaseDir` ✓
- `release/autocommit.go` — prefix check on `DefaultReleaseDir` ✓
- `release/workflow.go` — uses `DefaultReleaseDir` ✓
- `cmd/scanimport.go` — uses `DefaultReleaseDir` ✓
- `cmd/clearreleasejson.go` — uses `DefaultReleaseDir` ✓
- `cmd/listreleases.go` — uses `DefaultReleaseDir` via release pkg ✓
- `config/config_test.go` — hardcoded `"./gitmap-output"` ⚠️ update
- `constants/constants_terminal.go` — hardcoded `gitmap-output` in display strings ⚠️ update
- `constants/constants_messages.go` — hardcoded `gitmap-output/` in error msg ⚠️ update
- `constants/constants_cli.go` — hardcoded in help text ⚠️ update

### Phase 3 — Display Strings & Help Text (~4 files)

Update user-facing messages that reference old paths:

- `constants_terminal.go` — `StatusRepoCountFmt`, `ExecRepoCountFmt`,
  `TermCloneCmd1`, `TermCloneCmd3b`
- `constants_messages.go` — `MsgNoOutputDir`
- `constants_cli.go` — `HelpOutputPath`
- `helptext/*.md` — any references to `gitmap-output/` or `.release/`

### Phase 4 — Tests

- `config/config_test.go` — update expected `OutputDir` value
- `release/metadata_test.go` — uses `DefaultReleaseDir` var (auto-inherits) ✓
- `release/workflow_test.go` — uses `DefaultReleaseDir` var ✓

### Phase 5 — Specs, Docs & Memory

- Update all spec files referencing `.release/` or `gitmap-output/`
- Update `helptext/` markdown files
- Update docs site (`src/data/`, `src/pages/`)
- Update `.lovable/memory/` files
- Add changelog entry

## .gitignore Consideration

Projects using gitmap should add `.gitmap/output/` to their `.gitignore`.
The `.gitmap/release/` directory should remain tracked (release metadata is
committed). The `.gitmap/` root itself should NOT be gitignored.

## Backward Compatibility

No automatic migration of existing `.release/` or `gitmap-output/` folders.
The `doctor` command should gain a check that warns if old directories exist
and suggests moving them.

## Not In Scope

- Database location (stays binary-relative)
- Config file location (stays in `data/`)
- Profile storage (stays in `data/`)
