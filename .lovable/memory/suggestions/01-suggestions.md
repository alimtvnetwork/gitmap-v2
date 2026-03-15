# Suggestions Tracker

## Completed Suggestions

- ✅ Add `direct-clone-ssh.ps1` output
- ✅ Implement copy-and-handoff for `gitmap update`
- ✅ Add deploy retry logic in `run.ps1`
- ✅ Document `version` command in specs
- ✅ Bump version on every code change
- ✅ Update all spec docs for new features
- ✅ Create `spec/02-general/` with reusable design guidelines
- ✅ Add `desktop-sync` command
- ✅ Enhanced terminal output with HTTPS and SSH clone instructions
- ✅ Remove GitHub Release integration
- ✅ Nested deploy structure
- ✅ Update enhancements: skip-if-current, version comparison, rollback safety
- ✅ `update-cleanup` command with auto-run
- ✅ Made all `spec/02-general/` files fully generic
- ✅ Full compliance audit (Wave 1 + Wave 2)
- ✅ Constants inventory documentation
- ✅ `list-versions` and `revert` commands
- ✅ Changelog in release metadata JSON
- ✅ Releases table in SQLite database
- ✅ PascalCase for all DB table/column names
- ✅ `seo-write` command with templates, CSV, rotation, and dry-run
- ✅ Unit test infrastructure with PowerShell runner (`run.ps1 -t`)
- ✅ `--compress` flag for release assets (.zip/.tar.gz)
- ✅ `--checksums` flag for SHA256 checksums.txt generation
- ✅ Go cross-compilation pipeline (6 targets, auto-detect, GitHub upload)
- ✅ `--no-assets` and `--targets` flags for release customization
- ✅ Shell completion for all release flags (bash/zsh/powershell)
- ✅ Enhanced `list` output with labeled fields and inline cd hints
- ✅ Spec 41: Go Release Assets specification document
- ✅ Config-driven release targets: `release.targets` in `config.json` overrides default matrix
- ✅ Config-driven `release.checksums` and `release.compress` booleans
- ✅ `--list-targets` flag: prints resolved target matrix with source label

## Completed (Recent)

- ✅ **Build documentation site**: Replace placeholder React frontend with actual gitmap docs (commands, examples, architecture)
- ✅ **Add Linux/macOS support**: Shell scripts alongside PowerShell, cross-compile binary, GitHub Actions CI/CD
- ✅ **Add progress bar for clone**: Real-time `[current/total]` counter with repo name, duration, and success/failure summary

## Pending Suggestions

- ⬜ **Add real git status checks in TUI dashboard**: Wire up `gitutil.WatchStatus` to populate dirty/clean, ahead/behind in the dashboard view
