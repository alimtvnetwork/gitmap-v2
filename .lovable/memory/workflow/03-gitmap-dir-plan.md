# Plan: Unified .gitmap/ Directory Migration

## Goal

Consolidate `.release/` and `gitmap-output/` under a single `.gitmap/` folder.

See: `spec/01-app/56-unified-gitmap-dir.md`

## Steps

### Step 1 — Update constants (constants.go)
Change `DefaultReleaseDir` from `.release` to `.gitmap/release`.
Change `DefaultOutputDir` from `./gitmap-output` to `.gitmap/output`.
Change `DefaultOutputFolder` from `gitmap-output` to `.gitmap/output`.
Change `DefaultVerboseLogDir` from `gitmap-output` to `.gitmap/output`.
Add `GitMapDir = ".gitmap"` root constant.

### Step 2 — Update hardcoded display strings
Fix `constants_terminal.go` — repo count formats, clone step hints.
Fix `constants_messages.go` — `MsgNoOutputDir` error text.
Fix `constants_cli.go` — `HelpOutputPath` flag description.

### Step 3 — Update tests
Fix `config/config_test.go` expected `OutputDir`.

### Step 4 — Update help text files
Search all `helptext/*.md` for `.release/` and `gitmap-output/` references.

### Step 5 — Update spec documents
Search all `spec/` files for `.release/` and `gitmap-output/` references.

### Step 6 — Update docs site
Update `src/data/` and `src/pages/` references.

### Step 7 — Update memory files
Update `.lovable/memory/` references.

### Step 8 — Add doctor check for old directories
Add warning in `doctor` if `.release/` or `gitmap-output/` exist at repo root.

### Step 9 — Bump version and changelog

## Status

Step 1 complete — constants, display strings, comments, and test updated.
Steps 2–9 remaining.
