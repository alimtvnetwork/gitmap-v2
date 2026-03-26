# Issue 11 — Automatic Legacy Directory Migration

## Summary

Legacy repo-local directories (`gitmap-output/`, `.release/`, `.deployed/`) are
automatically migrated to `.gitmap/` subdirectories when detected. The migration
runs on every command (except `version`) via `PersistentPreRun`.

## Migration Map

| Legacy Directory   | Target              |
|--------------------|---------------------|
| `gitmap-output/`   | `.gitmap/output/`   |
| `.release/`        | `.gitmap/release/`  |
| `.deployed/`       | `.gitmap/deployed/` |

## Rules

1. Detect legacy directory at working directory root.
2. Create `.gitmap/` if missing.
3. If target does **not** exist → rename (move) legacy directory to target.
4. If target **already exists** → merge files from legacy into target
   (skip files that already exist), then **remove the legacy directory**.
5. Print summary message per migration.
6. Database (`data/`) is **not affected**.

## Implementation

- File: `cmd/migrate.go` with `migrateLegacyDirs()` and `mergeAndRemoveLegacy()`.
- Called early in root command's `PersistentPreRun`.
- Skipped for `version` / `v` commands to keep stdout clean.
- Constants: `LegacyOutputDir`, `LegacyReleaseDir`, `LegacyDeployedDir` in `constants.go`.
- Messages: `MsgMigrated`, `MsgMergedAndRemoved`, `ErrMigrationFailed` in `constants_messages.go`.

## Status

Complete — v2.36.1+.
