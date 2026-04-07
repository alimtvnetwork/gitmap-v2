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
    → Now in macro-ahk-v15

### Example 3: Lock detection when folder is in use

    gitmap cn v++ --delete

**Output:**

    Cloning macro-ahk-v12 into D:\wp-work\riseup-asia...
    ✓ Cloned macro-ahk-v12
    ✓ Registered macro-ahk-v12 with GitHub Desktop
    Warning: could not remove macro-ahk-v11: unlinkat: access denied
    Checking for processes locking macro-ahk-v11...
    The following processes are using this folder:
      • Code.exe (PID 14320)
      • explorer.exe (PID 5928)
    Terminate these processes to allow deletion? [y/N] y
    Terminating Code.exe (PID 14320)...
    ✓ Terminated Code.exe
    Terminating explorer.exe (PID 5928)...
    ✓ Terminated explorer.exe
    Retrying folder removal...
    ✓ Removed macro-ahk-v11
    → Now in macro-ahk-v12

## See Also

- [clone](clone.md) — Clone repos from output files
- [desktop-sync](desktop-sync.md) — Sync repos to GitHub Desktop
- [ssh](ssh.md) — Manage named SSH keys
