# Acceptance Criteria

## Scan Feature

- **Given** a directory with 3 nested Git repos,
  **when** `gitmap scan ./dir` is run,
  **then** all 3 repos appear in colored terminal output with clone instructions.

- **Given** `--mode ssh`,
  **then** clone instructions use `git@github.com:…` format.

- **Given** any scan,
  **then** a `gitmap-output/` folder is created inside the scanned directory
  containing `gitmap.csv`, `gitmap.json`, `folder-structure.md`, `clone.ps1`,
  and `register-desktop.ps1`.

- **Given** a folder with no `.git`,
  **then** it is skipped silently.

- **Given** a repo with no remote,
  **then** URL fields are empty, notes say "no remote configured."

- **Given** `--github-desktop`,
  **then** all discovered repos are registered with GitHub Desktop.

## Clone Feature

- **Given** a valid CSV or JSON from scan,
  **when** `gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored`,
  **then** all repos are cloned into correct relative paths preserving
  the original folder hierarchy.

- **Given** a repo that fails to clone,
  **then** it is logged and remaining repos continue.
  Summary shows N succeeded, M failed.

- **Given** `--github-desktop` on clone,
  **then** successfully cloned repos are added to GitHub Desktop.

## Config Feature

- **Given** no `--config` flag,
  **then** `./data/config.json` is loaded if it exists.

- **Given** CLI flags that conflict with config,
  **then** CLI flags take precedence.

## Terminal Output

- Terminal output shows a colored banner, repo list (name + path + clone instruction),
  folder tree, and clone help instructions for another machine.

## Build & Deploy

- **Given** `.\run.ps1` is run from the repo root,
  **then** it pulls, resolves deps, builds, and deploys the binary.

- **Given** `.\run.ps1 -R scan ../projects`,
  **then** it builds, resolves the relative path to absolute, and runs gitmap.

- **Given** `.\run.ps1 -R scan ../..`,
  **then** `../..` is resolved to an absolute path before being passed to gitmap.

## Code Quality

- No `if` condition uses negation.
- Every function is 8–15 lines.
- Every file is 100–200 lines.
- Each package has a single clear responsibility.
- All string literals are centralized in the `constants` package.
- A blank line precedes every `return` (unless sole line in `if` block).
