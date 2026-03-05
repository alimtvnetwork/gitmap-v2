# GitMap — Overview

## Purpose

GitMap is a CLI tool that scans a directory tree for Git repositories,
extracts clone URLs and branch information, and outputs structured data
(terminal, CSV, JSON, folder-structure Markdown). It can also re-clone
repositories from that structured data, preserving the original folder
hierarchy, and optionally register repos with GitHub Desktop.

## Working Name

`gitmap`

## Code Style Constraints

| Constraint            | Rule                                                                 |
|-----------------------|----------------------------------------------------------------------|
| `if` conditions       | Always positive — no `!`, no `!=`                                    |
| Function length       | 8–15 lines                                                           |
| File length           | 100–200 lines max                                                    |
| Package granularity   | One responsibility per package                                       |
| Newline before return | Always add a blank line before `return`, unless the `return` is the only line inside an `if` block |
| No magic strings      | All string literals used for comparison, format templates, default values, and file extensions must be defined as constants in a dedicated `constants` package |

## High-Level Components

1. **Constants** — all shared string literals, formats, ANSI colors, and default values.
2. **Config loader** — reads JSON config, merges with CLI flags.
3. **Scanner** — walks directories, detects `.git` folders.
4. **Mapper** — converts raw Git data into output records.
5. **Formatter** — renders records to terminal (colored), CSV, JSON, and folder-structure Markdown.
6. **Cloner** — re-clones repos from a previously generated file.
7. **Desktop** — registers repos with GitHub Desktop.

## Assumptions

- Remote URL is extracted from `origin` remote only.
- Symlinked directories are not followed.
- "Text file" input means one clone command per line.
