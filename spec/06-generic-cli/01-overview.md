# Generic CLI Creation Guidelines — Overview

## Purpose

This specification is a **complete, self-contained blueprint** for
building production-quality CLI tools. Hand it to any AI assistant
or developer and they can implement a well-structured CLI from scratch.

These guidelines are language-agnostic in principle but use Go for
concrete examples. Adapt syntax to your target language.

---

## Design Philosophy

| Principle | Detail |
|-----------|--------|
| Consistency over cleverness | Predictable patterns across all commands |
| Convention over configuration | Sensible defaults; config is optional |
| Fail fast, fail clearly | Bad input → immediate error with actionable message |
| One responsibility per unit | Each file, function, and package does one thing |
| No magic strings | Every literal in a constants package |
| Self-documenting | Help text, version, and examples built into the binary |

---

## Document Index

| File | Topic |
|------|-------|
| `01-overview.md` | This document — philosophy, scope, index |
| `02-project-structure.md` | Package layout, file organization, naming |
| `03-subcommand-architecture.md` | Routing, dispatch, handler pattern |
| `04-flag-parsing.md` | Per-command flags, defaults, validation |
| `05-configuration.md` | Three-layer config (defaults → file → flags) |
| `06-output-formatting.md` | Terminal, CSV, JSON, Markdown, scripts |
| `07-error-handling.md` | Exit codes, error messages, batch errors |
| `08-code-style.md` | Function length, file length, naming, conditionals |
| `09-help-system.md` | Embedded help files, `--help` interception |
| `10-database.md` | Local persistence, schema, upsert patterns |
| `11-build-deploy.md` | Build scripts, deploy, self-update |
| `12-testing.md` | Test structure, conventions, coverage |
| `13-checklist.md` | Step-by-step implementation checklist for AI |
| `14-date-formatting.md` | Centralized date display format |
| `15-constants-reference.md` | Every constant category with naming patterns |
| `16-verbose-logging.md` | Verbose/debug logging pattern with `--verbose` flag |
| `17-progress-tracking.md` | Progress reporting pattern for batch operations |
| `18-batch-execution.md` | Exec command pattern for running commands across repos |
| `19-shell-completion.md` | Shell tab-completion for PowerShell, Bash, and Zsh |

---

## How to Use This Spec

1. **Start with `13-checklist.md`** — it gives a sequenced plan.
2. **Reference individual docs** as you implement each layer.
3. **Every code example is a pattern** — adapt names, not structure.
4. **All constraints are mandatory** unless explicitly marked optional.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
