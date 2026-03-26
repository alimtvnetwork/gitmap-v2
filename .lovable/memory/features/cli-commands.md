# Memory: features/cli-commands
Updated: now

The CLI supports 55 subcommands with aliases: 'scan' (s), 'clone' (c), 'pull' (p), 'rescan' (rs), 'setup', 'status' (st), 'exec' (x), 'desktop-sync' (ds), 'release' (r), 'release-branch' (rb), 'release-pending' (rp), 'latest-branch' (lb), 'list' (ls), 'group' (g), 'multi-group' (mg), 'db-reset', 'version' (v), 'changelog' (cl), 'list-versions' (lv), 'list-releases' (lr), 'revert', 'doctor', 'update', 'seo-write' (sw), 'amend' (am), 'amend-list' (al), 'history' (hi), 'history-reset' (hr), 'stats' (ss), 'bookmark' (bk), 'export' (ex), 'import' (im), 'profile' (pf), 'cd' (go), 'watch' (w), 'diff-profiles' (dp), 'gomod' (gm), 'go-repos' (gr), 'node-repos' (nr), 'react-repos' (rr), 'cpp-repos' (cr), 'csharp-repos' (csr), 'alias' (a), 'zip-group' (z), 'completion' (cmp), 'interactive' (i), 'clear-release-json' (crj), 'update-cleanup', 'has-any-updates' (hau/hac), 'docs' (d), 'changelog-generate' (cg), 'ssh', 'prune' (pr), and 'temp-release' (tr). Current version: v2.36.3.

The release workflow now re-runs legacy directory migration after returning to the original branch, ensuring old `.release/` files are merged into `.gitmap/release/` and removed before auto-commit.


## temp-release (tr) — NEW
Lightweight temporary branch creation from recent commits. Creates branches from SHAs without checkout or tags. Supports batch creation with version pattern (`$$` placeholder), auto-increment sequencing, listing, and removal (single/range/all) with confirmation prompts. Tracked in `TempReleases` SQLite table. Spec: `spec/01-app/55-temp-release.md`.
