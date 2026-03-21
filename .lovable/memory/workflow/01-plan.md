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

### v2.17.0 → v2.23.0
- ✅ Enhanced group management (scoped commands, multi-group operations)
- ✅ `gomod` command for Go module rename/replace operations
- ✅ `diff-profiles` command for comparing scan profiles
- ✅ `watch` command for live repository monitoring
- ✅ `zip-group` command with Level 9 Deflate compression
- ✅ `alias` command with global `-A` flag resolution across pull, exec, status, cd
- ✅ Enriched CLI help with terminal simulations, standardized headers, 3-example limits
- ✅ Documentation site expansion with dedicated spec pages and reciprocal See Also navigation
- ✅ Unit/integration tests: location resolution, remote origin parsing, alias CRUD, SEO write
- ✅ Shell completion and cross-platform build parity

### v2.24.0 — Release Workflow Restructure
- ✅ Metadata committed on original branch: release branch only contains code and tags
- ✅ Verified 10-step release sequence: resolve → pad → check → source → branch → tag → push/assets → return → metadata → auto-commit
- ✅ `release-branch` and `release-pending` skip metadata steps 9–10
- ✅ `--notes` / `-N` flag for custom release titles, persisted in DB and JSON
- ✅ `--no-commit` flag to skip post-release auto-commit
- ✅ `--skip-meta` flag to bypass metadata writing
- ✅ Dry-run integration tests validating workflow step order
- ✅ TUI expanded to 6 views: Repos, Actions, Groups, Status, Zip Groups, Aliases

## Pending Work

### TUI Enhancements
- ✅ **Wire real git status** into TUI dashboard — unreachable detection, file counts, dashformat.go refactor
- ✅ **TUI release view**: trigger releases, view release history, and inspect metadata from within the TUI
- ✅ **TUI log viewer**: browse recent command history with detail view, duration formatting, and refresh
- ⬜ **TUI log filtering**: filter logs by command name or exit code within the TUI

### CLI Robustness
- ⬜ **Partial failure rollback**: auto-cleanup branch/tag on push failure in the release workflow
- ⬜ **Retry logic for GitHub uploads**: exponential backoff on transient API errors during asset upload
- ⬜ **Graceful offline mode**: detect no-network and skip remote operations (push, upload) with clear warnings
- ⬜ **Lock file for concurrent runs**: prevent two `gitmap release` processes from conflicting

### New Commands & Features
- ⬜ **`docs` command**: open the documentation website in the default browser from the CLI
- ⬜ **`prune` command**: detect and clean up stale release branches with no matching tags
- ✅ **`changelog generate`**: auto-generate changelog entries from commit messages between tags
- ⬜ **Config validation**: `doctor` check that warns on invalid or deprecated `config.json` fields

### Testing Coverage
- ⬜ **SkipMeta integration test**: verify metadata writing is suppressed when the flag is true
- ⬜ **Alias suggest tests**: cover auto-suggestion and conflict detection during scan
- ⬜ **TUI interaction tests**: automated key-press simulation using a TUI testing framework
- ⬜ **Release rollback test**: simulate push failure and verify branch/tag cleanup
- ⬜ **End-to-end release test**: full cycle from bump through metadata commit on a temp repo
