# Memory: features/cli-commands
Updated: now

The CLI supports 56 subcommands with aliases: 'scan' (s), 'clone' (c), 'clone-next' (cn), 'pull' (p), 'rescan' (rs), 'setup', 'status' (st), 'exec' (x), 'desktop-sync' (ds), 'release' (r), 'release-branch' (rb), 'release-pending' (rp), 'latest-branch' (lb), 'list' (ls), 'group' (g), 'multi-group' (mg), 'db-reset', 'version' (v), 'changelog' (cl), 'list-versions' (lv), 'list-releases' (lr), 'revert', 'doctor', 'update', 'seo-write' (sw), 'amend' (am), 'amend-list' (al), 'history' (hi), 'history-reset' (hr), 'stats' (ss), 'bookmark' (bk), 'export' (ex), 'import' (im), 'profile' (pf), 'cd' (go), 'watch' (w), 'diff-profiles' (dp), 'gomod' (gm), 'go-repos' (gr), 'node-repos' (nr), 'react-repos' (rr), 'cpp-repos' (cr), 'csharp-repos' (csr), 'alias' (a), 'zip-group' (z), 'completion' (cmp), 'interactive' (i), 'clear-release-json' (crj), 'update-cleanup', 'has-any-updates' (hau/hac), 'docs' (d), 'changelog-generate' (cg), 'ssh', 'prune' (pr), 'temp-release' (tr), and 'dashboard' (db). Current version: v2.36.7.

The release workflow re-runs legacy directory migration after returning to the original branch, ensuring old `.release/` files are merged into `.gitmap/release/` and removed before auto-commit.

## Batch Operations

The `pull` and `exec` commands support `--stop-on-fail` to halt batch operations after the first failure. Failed items are tracked with `FailWithError` and reported via `PrintFailureReport`. Partial failures exit with code 3 (`ExitPartialFailure`).

## temp-release (tr)

Lightweight temporary branch creation from recent commits. Creates branches from SHAs without checkout or tags. Supports batch creation with version pattern (`$$` placeholder), auto-increment sequencing, listing, and removal (single/range/all) with confirmation prompts. Tracked in `TempReleases` SQLite table. Spec: `spec/01-app/55-temp-release.md`.

## interactive (i)

Full-screen TUI with 9 views: Repos, Actions, Groups, Status, Releases, Temp Releases, Zip Groups, Aliases, Logs. See `features/interactive-tui.md` for details.
