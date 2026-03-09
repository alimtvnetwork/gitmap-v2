# Development Plan

## Completed Work

### v1.1.0 → v1.1.3
- ✅ Self-update handoff, direct SSH clone output, deploy retry logic
- ✅ Desktop-sync command, enhanced terminal clone hints

### v2.0.0 → v2.1.0
- ✅ Removed GitHub Release integration (Git-only + local metadata)
- ✅ Nested deploy structure, update enhancements, update-cleanup command
- ✅ Generic spec files in `spec/02-general/`

### v2.2.0 → v2.9.0
- ✅ Release-pending, changelog, doctor, latest-branch commands
- ✅ Date formatting, sort/filter flags, CSV/JSON output formats
- ✅ Database with repos, groups, group management commands
- ✅ Self-update hardening (rename-first, stale-process fallback)

### v2.10.0 — Compliance Audit
- ✅ Full compliance audit (Wave 1 + Wave 2): all 75+ source files pass code style rules
- ✅ Trimmed oversized files, fixed negation/switch violations, extracted constants

### v2.11.0 — Constants Inventory
- ✅ Added constants inventory audit section documenting ~280 constants

### v2.12.0 — New Commands
- ✅ `list-versions` (`lv`): show all release tags sorted highest-first with changelog
- ✅ `revert <version>`: checkout tag + handoff rebuild (same mechanism as update)

### v2.13.0 — Changelog Enhancements
- ✅ Release metadata JSON includes changelog field from CHANGELOG.md
- ✅ `list-versions` shows changelog notes as sub-points (terminal + JSON)

### v2.14.0 — Database Releases & PascalCase
- ✅ Added `Releases` table to SQLite for persistent release metadata
- ✅ Converted all DB table/column names from snake_case to PascalCase
- ✅ Release workflow auto-persists to database after successful releases

### v2.18.0 — SEO-Write & Test Infrastructure
- ✅ `seo-write` (`sw`) command with CSV and template input modes
- ✅ `CommitTemplates` SQLite table with 25 titles and 20 descriptions seed
- ✅ Placeholder substitution, rotation mode, dry-run, and graceful shutdown
- ✅ Unit test suite (50+ tests) across cmd, store, and constants packages
- ✅ PowerShell test runner (`run.ps1 -t`) with report output to `data/unit-test-reports/`

## Pending Work

- ⬜ **Frontend documentation site**: Currently a placeholder React app
- ⬜ **Cross-platform support**: Currently Windows-only (PowerShell scripts)
