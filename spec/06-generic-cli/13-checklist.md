# Implementation Checklist

## Instructions for AI

This is a sequenced implementation plan. Execute each phase in order.
Reference the numbered spec files for detailed patterns.
All constraints from `08-code-style.md` apply to every file you write.

---

## Phase 1: Scaffold (Do First)

- [ ] Create `main.go` with minimal entry point calling `cmd.Run()`
- [ ] Create `constants/constants.go` with `Version`, tool name
- [ ] Create `constants/constants_cli.go` with all command names + aliases
- [ ] Create `constants/constants_messages.go` with error messages
- [ ] Create `cmd/root.go` with `Run()` and `dispatch()`
- [ ] Create `cmd/rootusage.go` with `printUsage()` help text
- [ ] Implement `version` command (print version, exit 0)
- [ ] Implement `help` command (print usage, exit 0)

**Verify:** `go build && ./toolname version && ./toolname help`

---

## Phase 2: Configuration

- [ ] Create `model/` package with core data structs
- [ ] Create `config/config.go` with three-layer merge logic
- [ ] Create `data/config.json` with default settings
- [ ] Wire config loading into first real command

**Verify:** Tool reads config, flag overrides work

---

## Phase 3: Core Command

- [ ] Create `scanner/scanner.go` (or domain-specific logic package)
- [ ] Create `mapper/mapper.go` for data transformation
- [ ] Implement first real command (e.g., `scan`) with flag parsing
- [ ] Create `cmd/helpcheck.go` with `checkHelp()` function
- [ ] Add `checkHelp` call to command handler

**Verify:** `./toolname scan <input> && ./toolname scan --help`

---

## Phase 4: Output Formatting

- [ ] Create `formatter/terminal.go` with colored output
- [ ] Create `formatter/csv.go` with CSV writer
- [ ] Create `formatter/json.go` with JSON writer
- [ ] Create `formatter/structure.go` with Markdown tree
- [ ] Add `--output` flag to main command

**Verify:** Each format produces correct output

---

## Phase 5: Database

- [ ] Create `store/store.go` with DB init and migration
- [ ] Create `store/repo.go` (or domain CRUD file)
- [ ] Create `constants/constants_store.go` with SQL + paths
- [ ] Wire DB upsert into main command's output flow
- [ ] Implement `db-reset` command

**Verify:** Data persists across runs, reset clears it

---

## Phase 6: Additional Commands

- [ ] Implement each remaining command in its own file
- [ ] Add flag parsing per command
- [ ] Add `checkHelp` to every handler
- [ ] Create help files in `helptext/*.md`
- [ ] Wire commands into dispatch

**Verify:** Each command works with `--help` and normal execution

---

## Phase 7: Help System

- [ ] Create `helptext/print.go` with `go:embed` and `Print()`
- [ ] Create one `.md` file per command (see `09-help-system.md`)
- [ ] Ensure every handler calls `checkHelp` as first line
- [ ] Verify `--help` and `-h` both work on every command

---

## Phase 8: Build & Deploy

- [ ] Create build script (`run.ps1` or `Makefile`)
- [ ] Add `-ldflags` for embedded variables
- [ ] Add deploy step with retry logic
- [ ] Add version verification after build
- [ ] Implement `update` command with self-update

---

## Phase 9: Testing

- [ ] Add unit tests for `mapper`, `config`, `formatter`
- [ ] Add integration tests under `tests/`
- [ ] Verify all tests pass: `go test ./...`

---

## Phase 10: Polish

- [ ] Update `README.md` with grouped command reference
- [ ] Verify all files ≤ 200 lines
- [ ] Verify all functions ≤ 15 lines
- [ ] Verify no magic strings (all in `constants`)
- [ ] Verify positive conditionals only
- [ ] Verify blank line before every `return`
- [ ] Final version bump

---

## Quick Reference — File Counts

| Phase | Files Created |
|-------|--------------|
| Scaffold | ~6 |
| Configuration | ~3 |
| Core Command | ~4 |
| Formatting | ~5 |
| Database | ~4 |
| Commands | 1 per command |
| Help | 1 per command + print.go |
| Build | 1-2 scripts |
| Tests | 1 per testable package |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
