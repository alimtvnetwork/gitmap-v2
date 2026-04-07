# Clone Next

## Status

Pending behavior correction after spec approval.

## Command

```text
gitmap clone-next <v++|v+1|vN> [flags]
```

## Alias

```text
cn
```

## Responsibility

From inside an existing Git repository, derive the source repository from
`remote.origin.url`, resolve the next or explicit versioned target repository,
clone it into the
parent directory, register it with GitHub Desktop, and optionally remove the
current local folder. If `--create-remote` is passed, the command will also
create the target GitHub repository before cloning when it does not exist.

## Source of Truth

The command must use the Git remote as the authoritative source for repo name
resolution.

1. Read the current repo URL using:
   - `git config --get remote.origin.url`
   - fallback: `.git/config` origin parsing if needed
2. Parse the host, owner/org, and repo name from that remote URL.
3. Derive the base repo name and current version from the **remote repo name**.
4. Use the current local folder name only for:
   - determining the parent directory for the clone
   - prompting/removing the current folder after success

This avoids incorrect behavior when the local folder name and remote repo name
are not perfectly aligned.

## Terminology

| Term | Meaning |
|------|---------|
| Base name | Repo name without version suffix, e.g. `macro-ahk` |
| Current version | Version implied by the current remote repo name |
| Target version | Version requested by the user |
| Target repo | New repo name in the form `<base>-vN` |

## Version Arguments

| Argument | Meaning | Example |
|----------|---------|---------|
| `v++` | Increment current version by 1 | `macro-ahk-v11` → `macro-ahk-v12` |
| `v+1` | Alias for increment-by-one | `coding-guidelines-v7` → `coding-guidelines-v8` |
| `vN` | Jump directly to an explicit version | `macro-ahk-v12` + `v15` → `macro-ahk-v15` |

## Version Rules

1. `v++` and `v+1` mean the same thing.
2. `vN` must accept only positive integers (`v1`, `v2`, `v15`, ...).
3. `v0`, negative values, and malformed inputs must fail with a clear error.
4. If the current repo has no suffix, the unsuffixed repo is treated as the
   original repo and the first increment target is `-v2`.

### No-Suffix Behavior

| Current repo | Argument | Target repo |
|--------------|----------|-------------|
| `macro-ahk` | `v++` | `macro-ahk-v2` |
| `macro-ahk` | `v+1` | `macro-ahk-v2` |
| `macro-ahk` | `v15` | `macro-ahk-v15` |

## Target Resolution

After parsing the current remote:

1. Compute the target version.
2. Build the target repo name: `<base-name>-v<target-version>`.
3. Build the target local path: `<parent-of-current-working-directory>/<target-repo-name>`.
4. Build the target remote URL by preserving the same host, owner/org, and URL
   scheme as the current remote.

### URL Examples

| Current remote | Target remote |
|----------------|---------------|
| `https://github.com/alimtvnetwork/macro-ahk-v11.git` | `https://github.com/alimtvnetwork/macro-ahk-v12.git` |
| `git@github.com:alimtvnetwork/macro-ahk-v11.git` | `git@github.com:alimtvnetwork/macro-ahk-v12.git` |
| `https://github.com/alimtvnetwork/coding-guidelines-v7.git` | `https://github.com/alimtvnetwork/coding-guidelines-v8.git` |

## Optional GitHub Creation (`--create-remote`)

By default, `clone-next` assumes the target remote already exists and proceeds
directly to `git clone`. When the `--create-remote` flag is set, the command
checks whether the target GitHub repository exists and creates it if missing
**before** attempting to clone. This requires `GITHUB_TOKEN` to be set.

### Behavior when `--create-remote` is set

1. Check whether the target remote repository exists.
2. If it does not exist and the host is GitHub, create it under the same
   owner/org as the source repo.
3. The created repo should use the target repo name exactly.
4. The command must not attempt `git clone` first when the target repo is known
   to be missing.
5. If repo creation fails, stop with a clear error and do not prompt for local
   deletion.

### Visibility (when creating)

The preferred behavior is to inherit the visibility of the source repository.
If that cannot be determined safely, the command should fail with a clear error
instead of guessing.

## Workflow

1. Confirm the current directory is a Git repo.
2. Resolve `remote.origin.url`.
3. Parse the current remote repo name and current version.
4. Resolve the target version from `v++`, `v+1`, or `vN`.
5. Compute the target repo name and target local path in the parent directory.
6. Check that the local target directory does not already exist.
7. If `--create-remote` is set, check whether the target remote exists and
   create it if missing.
8. Clone the target repo into the parent directory.
9. Register the cloned repo with GitHub Desktop unless `--no-desktop` is set.
10. Change to the parent directory to release file locks on the current folder.
11. If clone succeeds, either:
    - remove the current folder automatically with `--delete`
    - keep it automatically with `--keep`
    - otherwise prompt the user interactively
12. If removal fails, scan for processes locking the folder:
    - On Windows: use Sysinternals `handle.exe` or PowerShell WMI query.
    - On Unix/macOS: use `lsof +D <path>`.
    - Display the list of locking processes with name and PID.
    - Prompt the user to terminate them.
    - If confirmed, kill each process and retry `RemoveAll` after a brief delay.
13. If the current folder was removed, change into the newly cloned directory
    and print a confirmation (`→ Now in <target>`).

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--delete` | false | Remove the current folder automatically after successful clone |
| `--keep` | false | Keep the current folder and skip the removal prompt |
| `--no-desktop` | false | Skip GitHub Desktop registration |
| `--create-remote` | false | Create the target GitHub repo if it does not exist (requires `GITHUB_TOKEN`) |
| `--ssh-key <name>` / `-K <name>` | (none) | Use a named SSH key for Git operations |
| `--verbose` | false | Show detailed clone-next diagnostics |

If neither `--delete` nor `--keep` is provided, the command must prompt after a
successful clone.

## Examples

### Example 1: Simple clone with `v+1`

```text
D:\wp-work\riseup-asia\coding-guidelines-v7> gitmap cn v+1

Cloning coding-guidelines-v8 into D:\wp-work\riseup-asia...
✓ Cloned coding-guidelines-v8
✓ Registered coding-guidelines-v8 with GitHub Desktop
Remove current folder coding-guidelines-v7? [y/N] n
```

### Example 2: Simple clone with `v++`

```text
D:\wp-work\riseup-asia\macro-ahk-v11> gitmap cn v++

Cloning macro-ahk-v12 into D:\wp-work\riseup-asia...
✓ Cloned macro-ahk-v12
✓ Registered macro-ahk-v12 with GitHub Desktop
Remove current folder macro-ahk-v11? [y/N] n
```

### Example 3: Repo without an existing suffix

```text
D:\wp-work\riseup-asia\macro-ahk> gitmap cn v++

Cloning macro-ahk-v2 into D:\wp-work\riseup-asia...
✓ Cloned macro-ahk-v2
✓ Registered macro-ahk-v2 with GitHub Desktop
Remove current folder macro-ahk? [y/N] y
✓ Removed macro-ahk
→ Now in macro-ahk-v2
```

### Example 4: Jump to an exact version with auto-delete

```text
D:\wp-work\riseup-asia\macro-ahk-v12> gitmap cn v15 --delete

Cloning macro-ahk-v15 into D:\wp-work\riseup-asia...
✓ Cloned macro-ahk-v15
✓ Registered macro-ahk-v15 with GitHub Desktop
✓ Removed macro-ahk-v12
→ Now in macro-ahk-v15
```

### Example 5: Create remote repo before clone

```text
D:\wp-work\riseup-asia\macro-ahk-v12> gitmap cn v15 --create-remote --delete

Creating GitHub repo macro-ahk-v15...
✓ Created GitHub repo macro-ahk-v15
Cloning macro-ahk-v15 into D:\wp-work\riseup-asia...
✓ Cloned macro-ahk-v15
✓ Registered macro-ahk-v15 with GitHub Desktop
✓ Removed macro-ahk-v12
→ Now in macro-ahk-v15
```

## Error Handling

| Condition | Required behavior |
|-----------|-------------------|
| Not inside a Git repo | Print a clear error and exit 1 |
| `remote.origin.url` missing | Print a clear error and exit 1 |
| Remote URL cannot be parsed | Print a clear error and exit 1 |
| Invalid version argument | Print a clear error and exit 1 |
| Local target directory already exists | Print a clear error and suggest `cd` into it |
| Target GitHub repo creation fails (`--create-remote`) | Print a clear error and stop before clone |
| Clone fails | Print a clear error and do not delete current folder |
| GitHub Desktop registration fails | Warn, but keep clone success |
| Folder deletion fails | Warn, but keep clone success |

## Implementation Scope

| Component | File |
|-----------|------|
| Command handler | `cmd/clonenext.go` |
| Version parser | `clonenext/version.go` |
| Completion hints | `completion/*` and command completion sources |
| Help text | `helptext/clone-next.md` |
| Command usage output | `cmd/rootusage.go` and constants |
| Spec | `spec/01-app/59-clone-next.md` |

## Acceptance Criteria

1. `gitmap cn v++` increments the current version correctly.
2. `gitmap cn v+1` behaves exactly like `v++`.
3. `gitmap cn v15` clones the exact target version.
4. The source repo name is derived from the Git remote, not guessed from the
   local folder name.
5. The local clone target is always the parent directory of the current repo.
6. `--create-remote` creates missing target GitHub repos before clone.
7. GitHub Desktop registration happens by default after a successful clone.
8. Current-folder removal happens only after a successful clone.
9. `--delete` and `--keep` override the interactive removal prompt.
10. Help text, completion hints, and tests cover `v++`, `v+1`, and `vN`.

## Deferred Implementation Phases

1. ~~Version parsing and resolution fixes~~ — done
2. ~~Target GitHub repo existence check and creation~~ — done (opt-in via `--create-remote`)
3. ~~Clone workflow hardening~~ — done (auto-cd before removal)
4. Help, completion, and automated test updates

