# Installer PATH Not Active in Current Shell After Install

## Ticket

After installing gitmap via `curl | bash` on macOS, running `gitmap` immediately
returns `zsh: command not found: gitmap`. The user must manually source their
profile or open a new terminal.

## Symptoms

1. User runs `curl -fsSL .../install.sh | bash` on macOS (zsh default).
2. Installer adds PATH entry to `~/.zprofile` (or `~/.zshrc`).
3. Installer prints `export PATH=...` and reload instructions.
4. User types `gitmap` → `zsh: command not found: gitmap`.
5. User types `gitmap --help` → same error.

## Root Cause

Two separate issues compound:

### 1. Subshell isolation (unfixable)

When invoked via `curl | bash`, the installer runs in a **child process**.
The `export PATH=...` on line 375 only affects that subshell — the parent
interactive shell never receives the updated PATH. This is fundamental POSIX
behavior and cannot be fixed from inside the script.

### 2. Single-profile write (fixable)

The installer writes the PATH entry to only **one** profile file, chosen by
`$SHELL` detection. Users who open a different shell (bash, sh, fish) or whose
terminal reads a different profile (`.zprofile` vs `.zshrc`) won't find gitmap.

On macOS, the default `$SHELL` is zsh but `.zshrc` may not exist on a fresh
system, causing the installer to write to `.zprofile` instead. While
Terminal.app reads `.zprofile` for login shells, many terminal emulators
(iTerm2, VS Code terminal) open interactive non-login shells that only
read `.zshrc`.

### 3. No immediate activation instruction

The post-install message says "Open a new terminal or run: . ~/.zprofile"
but users expect `gitmap` to just work immediately. The instruction to source
the profile is buried after several lines of output and easy to miss.

## Fix

### Phase 1: Multi-profile PATH registration

Write the `export PATH=...` line to **all** detected profile files that exist
or are standard for the platform:

| Shell | Profiles written |
|-------|-----------------|
| zsh   | `~/.zshrc` AND `~/.zprofile` (create `.zshrc` if neither exists) |
| bash  | `~/.bashrc` AND `~/.bash_profile` (or `~/.profile`) |
| fish  | `~/.config/fish/config.fish` |

Additionally, always write to `~/.profile` as a catch-all for POSIX `sh`.

Each profile write is idempotent — skipped if the directory entry already
exists in the file.

### Phase 2: Immediate activation guidance

After installation, print a **prominent**, shell-specific activation command
that the user can copy-paste immediately:

```
  ✓ Installed! To start using gitmap right now, run:

      source ~/.zshrc

  Or open a new terminal window.
```

The activation command uses the **primary** profile file (the one most likely
to be read by the current terminal), not `.zprofile`.

### Phase 3: Session PATH export

The script already does `export PATH="${PATH}:${dir}"` which works when the
script is sourced directly (`source install.sh`) but not via pipe. Document
this limitation clearly.

## Prevention

1. All installer scripts must write PATH to **multiple** profile files to
   cover login shells, interactive shells, and POSIX sh.
2. Post-install output must show a single, copy-pasteable activation command
   **prominently** (not buried in a summary block).
3. The installer must explicitly state that `curl | bash` cannot modify the
   parent shell's environment.

## Related

- `spec/02-app-issues/20-path-not-available-in-other-shells.md` — cross-shell visibility
- `spec/01-app/82-install-script.md` — installer specification
- `gitmap/scripts/install.sh` — implementation
