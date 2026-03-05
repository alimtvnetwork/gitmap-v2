# Suggestions Tracker

## Completed Suggestions

- ✅ Add `direct-clone-ssh.ps1` output (plain SSH clone commands, one per line)
- ✅ Implement copy-and-handoff for `gitmap update` to avoid file-lock errors
- ✅ Add deploy retry logic in `run.ps1` (20 attempts, 500ms delay)
- ✅ Document `version` command in specs
- ✅ Bump version on every code change
- ✅ Update all spec docs for new features
- ✅ Create `spec/02-general/` with reusable AI-trainable design guidelines (6 files)
- ✅ Add `desktop-sync` command to sync repos to GitHub Desktop from scan output
- ✅ Enhanced terminal output with both HTTPS and SSH clone instructions
- ✅ Remove GitHub Release integration (release command now Git-only + local metadata)
- ✅ Nested deploy structure (`bin-run/gitmap/` subfolder)
- ✅ Update enhancements: skip-if-current, version comparison, rollback safety (.old backups)
- ✅ `update-cleanup` command with auto-run at end of update cycle
- ✅ Made all `spec/02-general/` files fully generic (no gitmap-specific references)

## Pending Suggestions

- ⬜ **Fix dispatch pattern inaccuracy**: `01-cli-design-patterns.md` shows `switch` but code uses chained `if/return`
- ⬜ **Add missing details**: UTF-8 BOM in self-update, tree-building algorithm, chained if/return pattern
- ⬜ **Add missing pattern files**: Batch operations, external tool integration, directory walking, testing conventions

- ⬜ **Verify update flow**: Run `gitmap update` end-to-end, confirm skip-if-current + rollback + auto-cleanup
- ⬜ **Verify SSH output**: Run scan, check `direct-clone-ssh.ps1` contains correct SSH URLs
- ⬜ **Verify desktop-sync**: Run `gitmap desktop-sync` end-to-end
- ⬜ **Build documentation site**: Replace placeholder React frontend with actual gitmap docs
- ⬜ **Add Linux/macOS support**: Shell scripts alongside PowerShell, cross-compile binary
- ⬜ **Add `--dry-run` flag**: Preview scan/clone output without writing files
- ⬜ **Add progress bar for clone**: Show progress during multi-repo clone operations
