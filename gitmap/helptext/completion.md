# gitmap completion

Generate or install shell tab-completion scripts for gitmap commands and repo names.

## Alias

cmp

## Usage

    gitmap completion <powershell|bash|zsh>

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --list-repos | — | Print repo slugs, one per line (for script use) |
| --list-groups | — | Print group names, one per line (for script use) |
| --list-commands | — | Print all command names, one per line (for script use) |

## Prerequisites

- Run `gitmap scan` first for repo name completion
- Shell profile must be writable for auto-install via `gitmap setup`

## Examples

### Example 1: Print PowerShell completion script

    gitmap completion powershell

**Output:**

    Register-ArgumentCompleter -CommandName gitmap ...

### Example 2: Print Bash completion script

    gitmap completion bash

**Output:**

    _gitmap_completions() { ... }
    complete -F _gitmap_completions gitmap

### Example 3: List repo slugs for scripting

    gitmap completion --list-repos

**Output:**

    my-api
    web-app
    shared-lib

## See Also

- [setup](setup.md) — Auto-installs completions during setup
- [cd](cd.md) — Navigate to repos using tab-completed slugs
- [group](group.md) — Group names are also tab-completed
