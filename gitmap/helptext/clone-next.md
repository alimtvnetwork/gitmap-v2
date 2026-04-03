# gitmap clone-next

Clone the next or a specific versioned iteration of the current repository into the parent directory.

## Alias

cn

## Usage

    gitmap clone-next <v++|vN> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --delete | false | Auto-remove current folder after clone |
| --keep | false | Keep current folder without prompting |
| --no-desktop | false | Skip GitHub Desktop registration |
| --ssh-key \<name\> | (none) | Use a named SSH key for the clone |
| --verbose | false | Write detailed debug log |
| --create-remote | false | Create target GitHub repo if missing (requires GITHUB_TOKEN) |

## Prerequisites

- Must be run inside a Git repository with a remote origin configured

## Examples

### Example 1: Increment version by one

    gitmap cn v++

**Output:**

    Cloning macro-ahk-v12 into D:\wp-work\riseup-asia...
    ✓ Cloned macro-ahk-v12
    ✓ Registered macro-ahk-v12 with GitHub Desktop
    Remove current folder macro-ahk-v11? [y/N] n

### Example 2: Jump to a specific version with auto-delete

    gitmap cn v15 --delete

**Output:**

    Cloning macro-ahk-v15 into D:\wp-work\riseup-asia...
    ✓ Cloned macro-ahk-v15
    ✓ Registered macro-ahk-v15 with GitHub Desktop
    ✓ Removed macro-ahk-v12

### Example 3: Repo without a version suffix

    gitmap clone-next v++

**Output:**

    Cloning macro-ahk-v2 into D:\wp-work\riseup-asia...
    ✓ Cloned macro-ahk-v2
    ✓ Registered macro-ahk-v2 with GitHub Desktop
    Remove current folder macro-ahk? [y/N] y
    ✓ Removed macro-ahk

## See Also

- [clone](clone.md) — Clone repos from output files
- [desktop-sync](desktop-sync.md) — Sync repos to GitHub Desktop
- [ssh](ssh.md) — Manage named SSH keys
