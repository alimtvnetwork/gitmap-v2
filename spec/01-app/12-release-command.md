# Release Command

## Overview

The `gitmap release` command automates Git release workflows: creating
release branches, tags, GitHub releases with changelogs, and tracking
release history. It supports full semver, partial versions (auto-padded),
pre-release suffixes, draft mode, dry-run preview, and auto-increment.

---

## Commands

### `gitmap release [version]` (alias: `r`)

Create a release branch, Git tag, and GitHub release.

- Version can be full (`v1.2.3`), partial (`v1.2`, `v1`), or omitted.
- Partial versions are zero-padded: `v1` → `v1.0.0`, `v1.2` → `v1.2.0`.
- If no version is provided, reads from `version.json` in the project root.
- If neither exists, exits with an error.

### `gitmap release-branch <branch>` (alias: `rb`)

Complete a release from an existing `release/vX.Y.Z` branch.
Creates the tag and GitHub release if not already done.

---

## Flags

### Release Flags

| Flag                         | Description                                      | Default     |
|------------------------------|--------------------------------------------------|-------------|
| `--assets <path>`            | Directory or file to attach to the release       | (none)      |
| `--commit <sha>`             | Create release from a specific commit            | (none)      |
| `--branch <name>`            | Create release from latest commit of a branch    | (none)      |
| `--bump major\|minor\|patch` | Auto-increment from the latest released version  | (none)      |
| `--draft`                    | Create an unpublished draft release              | `false`     |
| `--dry-run`                  | Preview release steps without executing          | `false`     |
| `--verbose`                  | Write detailed debug log                         | `false`     |

### Release-Branch Flags

| Flag               | Description                         | Default |
|--------------------|-------------------------------------|---------|
| `--assets <path>`  | Directory or file to attach         | (none)  |
| `--draft`          | Create an unpublished draft release | `false` |
| `--dry-run`        | Preview steps without executing     | `false` |
| `--verbose`        | Write detailed debug log            | `false` |

---

## Mutual Exclusivity Rules

The following flag combinations are invalid and cause an immediate error:

| Conflict                         | Error Message                                              |
|----------------------------------|------------------------------------------------------------|
| `--bump` + version argument      | `--bump cannot be used with an explicit version argument.` |
| `--commit` + `--branch`          | `--commit and --branch are mutually exclusive.`            |

---

## Version Resolution

Version is resolved in priority order:

1. **CLI argument** — `gitmap release v1.2.3`
2. **`--bump` flag** — reads latest from `.release/latest.json`, increments
3. **`version.json`** — `{ "version": "1.2.3" }` in project root
4. **Error** — no version source found

### Partial Version Padding

| Input   | Resolved  | Branch              | Tag       |
|---------|-----------|---------------------|-----------|
| `v1`    | `v1.0.0`  | `release/v1.0.0`    | `v1.0.0`  |
| `v1.2`  | `v1.2.0`  | `release/v1.2.0`    | `v1.2.0`  |
| `v1.2.3`| `v1.2.3`  | `release/v1.2.3`    | `v1.2.3`  |

### Pre-Release Versions

Pre-release suffixes are preserved and not padded:

| Input           | Resolved        | Tag             |
|-----------------|-----------------|-----------------|
| `v1.0.0-rc.1`  | `v1.0.0-rc.1`  | `v1.0.0-rc.1`  |
| `v1.0.0-beta`  | `v1.0.0-beta`  | `v1.0.0-beta`  |
| `v2.0.0-alpha.3`| `v2.0.0-alpha.3`| `v2.0.0-alpha.3`|

Pre-release versions are **never** marked as `latest`.

---

## Source Resolution

The commit used as the release base is resolved in order:

1. **`--commit <sha>`** — exact commit
2. **`--branch <name>`** — latest commit on that branch
3. **Current HEAD** — default

---

## Dry-Run Mode

When `--dry-run` is passed, each step is printed with a `[dry-run]`
prefix. No branches, tags, pushes, or GitHub releases are created.
Metadata files are not written.

```
  [dry-run] Create branch release/v1.2.3 from main
  [dry-run] Create tag v1.2.3
  [dry-run] Push branch and tag to origin
  [dry-run] Use CHANGELOG.md as release body
  [dry-run] Attach README.md
  [dry-run] Create GitHub release
  [dry-run] Write metadata to .release/v1.2.3.json
  [dry-run] Mark v1.2.3 as latest
```

---

## GitHub Release Creation

The tool uses a two-tier strategy:

1. **Preferred**: GitHub CLI (`gh release create`) — requires `gh` to be
   installed and authenticated.
2. **Fallback**: GitHub REST API via HTTP — requires a `GITHUB_TOKEN`
   environment variable.

### Token Resolution

For HTTP fallback, `GITHUB_TOKEN` is read from the environment.
If neither `gh` nor `GITHUB_TOKEN` is available, the release aborts
with a clear error message.

### Attachments

The following files are auto-detected and attached if present:

| File             | Behavior                           |
|------------------|------------------------------------|
| `CHANGELOG.md`   | Attached as release notes body     |
| `README.md`      | Attached as a release asset        |

Additional assets from `--assets <path>`:

- If path is a directory, all files inside are attached individually.
- If path is a file, that single file is attached.

---

## Auto-Increment (`--bump`)

Reads the latest version from `.release/latest.json` and increments:

| Current Latest | `--bump patch` | `--bump minor` | `--bump major` |
|----------------|----------------|----------------|----------------|
| `1.2.3`        | `1.2.4`        | `1.3.0`        | `2.0.0`        |
| `0.9.1`        | `0.9.2`        | `0.10.0`       | `1.0.0`        |

If no `latest.json` exists and `--bump` is used, exits with an error
instructing the user to create an initial release first.

`--bump` is mutually exclusive with a version argument.

---

## Duplicate Detection

Before creating a release, the tool checks:

1. **`.release/vX.Y.Z.json`** — if the metadata file exists, abort.
2. **Git tags** — if the tag already exists locally or remotely, abort.

Error message:
```
Version v1.2.3 is already released. See .release/v1.2.3.json for details.
```

---

## Error Scenarios

| Scenario                        | Behavior                                                    |
|---------------------------------|-------------------------------------------------------------|
| Invalid version string          | `'abc' is not a valid version.`                             |
| `--commit` SHA not found        | `commit abc123 not found.`                                  |
| `--branch` does not exist       | `branch develop does not exist.`                            |
| Push to remote fails            | `failed to push to remote: <detail>`                        |
| GitHub release creation fails   | `failed to create GitHub release: <detail>`                 |
| `gh` not installed, no token    | `no GITHUB_TOKEN found. Set GITHUB_TOKEN or install 'gh'.`  |
| Metadata write fails            | `could not write release metadata: <detail>`                |
| `version.json` unreadable       | `could not read version.json: <detail>`                     |

### Rollback Strategy

If a step fails after partial execution:

- **Branch/tag created but push fails**: error is reported; user must
  manually delete the local branch and tag.
- **Push succeeds but GitHub release fails**: branch and tag remain on
  remote; metadata is not written; user can retry with
  `gitmap release-branch release/vX.Y.Z`.
- **Metadata write fails**: GitHub release exists but local tracking
  is incomplete; user should manually create the `.release/` file.

No automatic rollback is performed. The error message includes the
failed step so the user knows exactly what to clean up.

---

## version.json Behavior

- `version.json` is **read-only** from the tool's perspective.
- The tool never auto-updates `version.json` after a release.
- Users manage `version.json` manually or via their own CI scripts.

---

## Workflow: Release from HEAD / Branch / Commit

```
 1. Resolve version (CLI → --bump → version.json → error)
 2. Pad partial version to full semver
 3. Check .release/ and git tags for duplicates → abort if exists
 4. Resolve source commit (--commit / --branch / HEAD)
 5. Create branch release/vX.Y.Z from source
 6. Create git tag vX.Y.Z
 7. Push branch + tag to origin
 8. Detect CHANGELOG.md → use as release body
 9. Detect README.md → attach as asset
10. Collect --assets contents → attach as assets
11. Create GitHub release (gh CLI → HTTP fallback)
12. Write .release/vX.Y.Z.json
13. Update .release/latest.json if highest stable version
```

## Workflow: Release from Existing Branch

```
1. Validate release/vX.Y branch exists
2. Extract version from branch name, pad to semver
3. Check if tag/release already exists → abort if so
4. Checkout the release branch
5. Create tag, push, create GitHub release (steps 6–13 above)
```

---

## Package Layout

```
release/
├── semver.go       # Version parsing, padding, comparison, bumping
├── metadata.go     # Read/write .release/*.json, latest.json, version.json
├── gitops.go       # Branch, tag, push, checkout Git operations
├── github.go       # GitHub release creation (gh CLI + HTTP fallback)
└── workflow.go     # Orchestration: Execute(), ExecuteFromBranch()
```

Each file stays under 200 lines. `workflow.go` is the entry point;
all other files are pure helpers with no cross-dependencies.

---

## CLI Examples

```bash
# Full semver release from HEAD
gitmap release v1.2.3

# Partial version (padded to v1.0.0)
gitmap release v1

# With assets
gitmap release v2.0.0 --assets ./dist

# Alias
gitmap r v1.5.0

# From specific commit
gitmap release v1.2.3 --commit abc123def

# From specific branch
gitmap release v1.0.0 --branch develop

# Auto-increment
gitmap release --bump patch
gitmap release --bump minor --assets ./bin

# Draft release
gitmap release v3.0.0-rc.1 --draft

# Dry-run preview
gitmap release v1.0.0 --dry-run

# No version — reads version.json
gitmap release

# Complete release from existing release branch
gitmap release-branch release/v1.2.0
gitmap rb release/v1.2.0

# Dry-run from branch
gitmap release-branch release/v1.2.0 --dry-run
```

---

## Acceptance Criteria

- **Given** `gitmap release v1.0.0`, **then** branch `release/v1.0.0`,
  tag `v1.0.0`, and GitHub release are created.
- **Given** `--assets ./dist`, **then** dist folder contents are attached.
- **Given** `CHANGELOG.md` exists, **then** its content is used as the
  release body.
- **Given** `README.md` exists, **then** it is attached as an asset.
- **Given** version already released, **then** abort with clear message.
- **Given** no version arg and `version.json` exists, **then** version is
  read from it.
- **Given** `--commit <sha>`, **then** release branch starts from that commit.
- **Given** `--branch main`, **then** latest commit of `main` is used.
- **Given** `gitmap release-branch release/v1.2.0`, **then** tag + release
  is created from that branch.
- **Given** multiple releases, **then** `latest.json` points to the highest
  stable semver.
- **Given** `--bump patch` with latest `v1.2.3`, **then** releases `v1.2.4`.
- **Given** `--draft`, **then** GitHub release is created as a draft.
- **Given** pre-release version, **then** it is not marked as latest.
- **Given** `v1`, **then** padded to `v1.0.0`.
- **Given** `gh` is not installed, **then** falls back to GitHub HTTP API.
- **Given** `--dry-run`, **then** all steps are printed but nothing is
  executed; no branches, tags, or releases are created.
- **Given** `--bump` with a version argument, **then** abort with conflict
  error.
- **Given** `--commit` with `--branch`, **then** abort with mutual
  exclusivity error.
- **Given** push fails after branch/tag creation, **then** error message
  includes the failed step for manual cleanup.
- **Given** GitHub release fails after push, **then** user can retry via
  `gitmap release-branch`.
