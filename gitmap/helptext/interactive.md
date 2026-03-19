# Interactive TUI

Launch a full-screen interactive terminal UI for browsing, searching,
and managing repositories.

## Alias

i

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

    ┌─ Repos ─────────────────────────────────────────┐
    │ Search: api_                                     │
    │                                                  │
    │ > my-api           main     clean    0/0         │
    │   api-gateway      main     dirty    1/0         │
    │   payments-api     develop  clean    0/2         │
    │                                                  │
    │ 3 matches (42 total) │ Space: select │ /: search │
    └──────────────────────────────────────────────────┘

- Type `/` to search, `j`/`k` or arrow keys to navigate
- `Space` to select repos, `a` to select all
- `Enter` to view detail

### Actions
Perform batch operations on selected repos:

    ┌─ Actions ───────────────────────────────────────┐
    │ 3 repos selected                                 │
    │                                                  │
    │   [p] Pull selected repos                        │
    │   [x] Execute git command across selected        │
    │   [s] Show status for selected                   │
    │   [g] Add selected to a group                    │
    │                                                  │
    │ Press a key to perform action                    │
    └──────────────────────────────────────────────────┘

### Groups
Manage repository groups:

    ┌─ Groups ────────────────────────────────────────┐
    │ GROUP           REPOS   DESCRIPTION              │
    │ > backend       5       All backend services     │
    │   frontend      3       React frontend apps      │
    │   infra         2       Infrastructure           │
    │                                                  │
    │ c: create │ d: delete │ Enter: show members      │
    └──────────────────────────────────────────────────┘

### Status
Live dashboard showing dirty/clean status, branch, ahead/behind:

    ┌─ Status ────────────────────────────────────────┐
    │ 42 repos │ Refreshing in 25s │ r: refresh now    │
    │                                                  │
    │ REPO             BRANCH     STATUS  AHEAD/BEHIND │
    │ my-api           main       clean   0/0          │
    │ web-app          develop    dirty   2/1          │
    │ billing-svc      main       clean   0/0          │
    │ auth-gateway     feature/x  dirty   5/0          │
    │                                                  │
    │ 2 dirty │ 40 clean                               │
    └──────────────────────────────────────────────────┘

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
