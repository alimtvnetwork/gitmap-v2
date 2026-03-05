# Suggestions Tracker

## Completed Suggestions

- ✅ Add `direct-clone-ssh.ps1` output (plain SSH clone commands, one per line)
- ✅ Implement copy-and-handoff for `gitmap update` to avoid file-lock errors
- ✅ Add deploy retry logic in `run.ps1` (20 attempts, 500ms delay)
- ✅ Document `version` command in specs
- ✅ Bump version on every code change
- ✅ Update all spec docs for new features
- ✅ Create `spec/02-general/` with reusable AI-trainable design guidelines (6 files)

## Pending Suggestions

- ⬜ **Fix dispatch pattern inaccuracy**: `01-cli-design-patterns.md` shows `switch` but code uses chained `if/return`
- ⬜ **Add missing details**: UTF-8 BOM in self-update, tree-building algorithm, chained if/return pattern
- ⬜ **Add missing pattern files**: Batch operations, external tool integration, directory walking, testing conventions

- ⬜ **Verify update flow**: Run `gitmap update` end-to-end, confirm no file-lock errors
- ⬜ **Verify SSH output**: Run scan, check `direct-clone-ssh.ps1` contains correct SSH URLs
- ⬜ **Build documentation site**: Replace placeholder React frontend with actual gitmap docs
- ⬜ **Add Linux/macOS support**: Shell scripts alongside PowerShell, cross-compile binary
- ⬜ **Add `--dry-run` flag**: Preview scan/clone output without writing files
- ⬜ **Add progress bar for clone**: Show progress during multi-repo clone operations
