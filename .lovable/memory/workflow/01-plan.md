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

## Pending Work

- ⬜ **Frontend documentation site**: Currently a placeholder React app — needs real gitmap docs
- ⬜ **Cross-platform support**: Currently Windows-only (PowerShell scripts); add shell scripts, cross-compile binary
- ⬜ **Version bump to post-v2.14.0**: Next feature set will determine minor/patch
