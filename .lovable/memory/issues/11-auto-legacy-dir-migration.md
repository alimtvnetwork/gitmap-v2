# Issue 11 — Automatic Legacy Directory Migration

## Summary

Legacy repo-local directories (`gitmap-output/`, `.release/`, `.deployed/`) must be
automatically migrated to `.gitmap/` subdirectories when detected, instead of only
warning via `doctor`.

## Migration Map

| Legacy Directory   | Target              |
|--------------------|---------------------|
| `gitmap-output/`   | `.gitmap/output/`   |
| `.release/`        | `.gitmap/release/`  |
| `.deployed/`       | `.gitmap/deployed/` |

## Rules

1. Detect legacy directory at working directory root.
2. Create `.gitmap/` if missing.
3. Rename (move) legacy directory to target.
4. Skip if target already exists — print warning, leave both.
5. Print `Migrated <old>/ → .gitmap/<new>/` per successful move.
6. Database (`data/`) is **not affected**.

## Implementation

- New file: `cmd/migrate.go` with `migrateLegacyDirs()`.
- Called early in root command's `PersistentPreRun`.
- Constants for legacy directory names in `constants/constants.go`.

## Status

Spec and plan updated. Implementation pending.
