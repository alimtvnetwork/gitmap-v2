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

### v2.14.0 — Go Release Assets, Compression & Checksums
- ✅ `--compress` flag: archive assets (.zip for Windows, .tar.gz for others)
- ✅ `--checksums` flag: generate SHA256 `checksums.txt` for all release assets
- ✅ Go cross-compilation pipeline: auto-detect `go.mod` + `cmd/` entries
- ✅ Builds 6 default targets (windows/linux/darwin × amd64/arm64) with `CGO_ENABLED=0`
- ✅ `--no-assets` flag to skip binary compilation
- ✅ `--targets` flag for custom OS/arch matrix
- ✅ Native GitHub API asset upload with retry (no external CLI dependencies)
- ✅ Dry-run support showing planned binaries, compression, and checksums
- ✅ Shell completion updated for all new flags (bash/zsh/powershell)
- ✅ Enhanced `list` (`ls`) output with labeled fields and inline `cd` hints
- ✅ Spec document: `spec/01-app/41-go-release-assets.md`
- ✅ Config-driven release targets: `release.targets` in `config.json` overrides default OS/arch matrix
- ✅ Config-driven `release.checksums` and `release.compress` booleans (CLI flags still win)

### v2.15.0 — Cross-Platform & CI/CD
- ✅ Full documentation site with real gitmap command docs, examples, and architecture pages
- ✅ `run.sh` build script with full parity to `run.ps1`
- ✅ Makefile with corrected flags matching `run.sh`
- ✅ GitHub Actions CI workflow: test on push, cross-compile 6 targets
- ✅ GitHub Actions Release workflow: auto-release on `v*` tags with compression and checksums
- ✅ Spec document: `spec/01-app/42-cross-platform.md`

### v2.16.0 — Interactive TUI
- ✅ Bubble Tea TUI with 4 views: repo browser, batch actions, group management, status dashboard
- ✅ Fuzzy search via `sahilm/fuzzy`, multi-select, keyboard navigation
- ✅ `tui/` package (7 files): tui.go, browser.go, actions.go, groups.go, dashboard.go, keys.go, styles.go
- ✅ `interactive` (i) command wired into CLI dispatch, help, and shell completion
- ✅ Spec document: `spec/01-app/43-interactive-tui.md`
- ✅ Constants in `constants/constants_tui.go`, help text in `helptext/interactive.md`

### v2.15.1 — Database Path Fix
- ✅ Fixed DB path resolution: database now at `<binary-dir>/data/` instead of CWD-relative
- ✅ Added `store/location.go` with `BinaryDataDir()`, `OpenDefault()`, `OpenDefaultProfile()`
- ✅ Updated all 13 database callers across the codebase
- ✅ Issue post-mortem: `.lovable/memory/issues/04-database-path-resolution.md`
- ✅ Spec update: `spec/01-app/16-database.md` updated with binary-relative path

## Pending Work

- ⬜ **Wire real git status** into TUI dashboard (dirty/clean, ahead/behind via gitutil)
