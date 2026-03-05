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

### General Guidelines (spec/02-general/)
- ✅ **CLI design patterns**: Subcommand routing, flag parsing, version command, constants, help output, error handling
- ✅ **PowerShell build/deploy**: Step-based scripts, logging, config, retry-on-lock, -R flag forwarding
- ✅ **Self-update mechanism**: Copy-and-handoff, delayed rebuild, file lock avoidance
- ✅ **Output & formatting**: Multi-format strategy, terminal reports, templates, CSV/JSON/Markdown
- ✅ **Config pattern**: Three-layer merge (defaults → JSON → CLI flags)
- ✅ **Code style rules**: Positive conditionals, function/file limits, no magic strings, naming

## Pending Work

- ⬜ **Verify update flow end-to-end**: Run `gitmap update` and confirm deploy succeeds without file-lock errors
- ⬜ **Verify direct-clone-ssh.ps1**: Run scan and confirm SSH output file is generated correctly
- ⬜ **Frontend documentation site**: Currently a placeholder React app — needs actual content
- ⬜ **Cross-platform support**: Currently Windows-only (PowerShell scripts, `.exe` binary)
