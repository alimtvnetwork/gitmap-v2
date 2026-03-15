# Interactive TUI

Launch a full-screen interactive terminal UI for browsing, searching,
and managing repositories.

## Usage

    gitmap interactive [--refresh <seconds>]
    gitmap i [--refresh <seconds>]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--refresh` | config or 30 | Dashboard auto-refresh interval in seconds |

## Views

The TUI has four views, accessible via **Tab**:

### Repos (default)
Browse all tracked repositories with fuzzy search:
- Type `/` to search, `j`/`k` or arrow keys to navigate
- `Space` to select repos, `a` to select all
- `Enter` to view detail

### Actions
Perform batch operations on selected repos:
- `p` — Pull selected repos
- `x` — Execute a git command across selected
- `s` — Show status for selected
- `g` — Add selected to a group

### Groups
Manage repository groups:
- Browse existing groups with member counts
- `c` — Create a new group
- `d` — Delete a group

### Status
Live dashboard showing dirty/clean status, branch, ahead/behind
counts for all repos. Press `r` to refresh.

## Key Bindings

    q / Esc      Quit
    Tab          Switch view
    j / ↓        Move down
    k / ↑        Move up
    Space        Select/deselect repo
    a            Select all
    /            Focus search
    Enter        Show detail / confirm

## Requirements

Requires a terminal with alternate screen support. Falls back
to an error message if not available.
