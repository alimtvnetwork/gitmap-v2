# Development Plan

## Completed Work

### v1.1.0 → v1.1.1
- ✅ **Self-update handoff**: Implemented copy-and-handoff mechanism to avoid Windows file locks during `gitmap update`
- ✅ **Direct SSH clone output**: Added `direct-clone-ssh.ps1` with raw SSH `git clone` commands
- ✅ **Version bump**: 1.1.0 → 1.1.1

### v1.1.1 → v1.1.2
- ✅ **Deploy retry logic**: Added 20-attempt retry with 500ms delay in `run.ps1` for locked binary
- ✅ **Update delay**: Added 1.2s delay before rebuild in update handoff
- ✅ **Version command docs**: Updated all spec docs for `version` command and build output
- ✅ **Spec updates**: Documented direct-clone-ssh.ps1, copy-and-handoff update, deploy retry, version display
- ✅ **Version bump**: 1.1.1 → 1.1.2

### v1.1.2 → v1.1.3
- ✅ **`desktop-sync` command**: New command that syncs repos from gitmap-output/gitmap.json to GitHub Desktop
- ✅ **Enhanced terminal clone hints**: Show both HTTPS and SSH clone commands, script shortcuts, and desktop-sync at end of output
- ✅ **Specs**: Created spec/01-app/11-desktop-sync.md, updated CLI interface and output formatting specs
- ✅ **Version bump**: 1.1.2 → 1.1.3

### v2.0.0 → v2.1.0
- ✅ **Removed GitHub Release integration**: Release command now handles Git branch/tag/push and local metadata only
- ✅ **Nested deploy structure**: Deploy to `bin-run/gitmap/` subfolder instead of flat directory
- ✅ **Update enhancements**: Skip-if-current, version comparison before/after, rollback safety with `.old` backups
- ✅ **`update-cleanup` command**: Explicit cleanup of temp copies and `.old` backups, auto-runs at end of update
- ✅ **Proactive temp cleanup removed**: Replaced auto-startup cleanup with explicit `update-cleanup` command
- ✅ **Generic spec files**: All `spec/02-general/` files made fully generic and shareable (no gitmap-specific references)

### General Guidelines (spec/02-general/)
- ✅ **CLI design patterns**: Subcommand routing, flag parsing, version command, constants, help output, error handling
- ✅ **PowerShell build/deploy**: Step-based scripts, logging, config, retry-on-lock, nested deploy, -R flag forwarding
- ✅ **Self-update mechanism**: Copy-and-handoff, skip-if-current, version comparison, rollback, cleanup command
- ✅ **Output & formatting**: Multi-format strategy, terminal reports, templates, CSV/JSON/Markdown
- ✅ **Config pattern**: Three-layer merge (defaults → JSON → CLI flags)
- ✅ **Code style rules**: Positive conditionals, function/file limits, no magic strings, naming

## Pending Work

- ⬜ **Verify update flow end-to-end**: Run `gitmap update` and confirm full cycle (skip-if-current, rollback, auto-cleanup)
- ⬜ **Verify direct-clone-ssh.ps1**: Run scan and confirm SSH output file is generated correctly
- ⬜ **Verify desktop-sync**: Run `gitmap desktop-sync` end-to-end
- ⬜ **Frontend documentation site**: Currently a placeholder React app — needs actual content
- ⬜ **Cross-platform support**: Currently Windows-only (PowerShell scripts, `.exe` binary)
