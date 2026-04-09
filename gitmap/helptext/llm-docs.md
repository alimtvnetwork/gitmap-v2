# gitmap llm-docs

Generate a consolidated LLM.md reference file that describes all gitmap
commands, flags, architecture, and usage patterns. Designed to be shared
with AI assistants so they can understand gitmap and generate CLI commands
from natural language.

## Alias

ld

## Usage

    gitmap llm-docs

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --stdout | false | Print to stdout instead of writing LLM.md file |

## Prerequisites

- None — works from any directory

## Examples

### Example 1: Generate LLM.md in the current directory

    gitmap llm-docs

**Output:**

    ↻ Generating LLM.md from command registry...
    ✓ LLM.md written to D:\repos\my-project\LLM.md

### Example 2: Using the alias

    gitmap ld

**Output:**

    ↻ Generating LLM.md from command registry...
    ✓ LLM.md written to /home/user/projects/LLM.md

### Example 3: Print to stdout and pipe to clipboard

    gitmap llm-docs --stdout | pbcopy    # macOS
    gitmap llm-docs --stdout | clip      # Windows
    gitmap llm-docs --stdout | xclip -sel clip  # Linux

### Example 4: Pipe directly to an AI tool

    gitmap ld --stdout | ai-chat --context -

## What's Included

The generated LLM.md contains:

- **Command Reference** — all 60+ commands with aliases, descriptions, examples
- **Global Flags** — flags that work across commands
- **Architecture** — module layout, database schema, project structure
- **Coding Conventions** — rules for modifying gitmap source code
- **Common Patterns** — ready-made command sequences for typical tasks
- **Version** — tagged with the current gitmap version

## See Also

- [docs](docs.md) — Open documentation website in browser
- [help-dashboard](help-dashboard.md) — Serve docs site locally
- [version](version.md) — Show version number
