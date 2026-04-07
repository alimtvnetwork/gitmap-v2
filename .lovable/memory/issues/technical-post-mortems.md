# Technical Post-Mortems

Recent critical issues and their resolutions:

1. **Windows File Locks**: Resolved via a 'rename-first' strategy.
2. **Database Fragmentation**: Fixed by anchoring `data/` to the binary location.
3. **Path Double-Nesting**: Corrected `store` package path resolution.
4. **CLI Changelog Sync**: Synchronized release notes with compiled constants.
5. **Zip Group Silent Failure**: Added explicit error reporting.
6. **Auto-commit Push Rejection**: Implemented `git pull --rebase` recovery.
7. **List-Releases Resolution**: Prioritized local metadata.
8. **Legacy UUID Detection**: Provided recovery instructions.
9. **Auto Legacy Dir Migration**: Consolidated folders at startup.
10. **Legacy ID Migration**: Rebuilt Repos table.
11. **One-liner PATH Propagation**: Broadcasted system changes.
12. **Binary Extraction Failure**: Implemented flexible filename mapping in installer.
13. **Install Script 404**: Fixed relative paths.
14. **Latest Tag Disconnect**: Updated CI to explicitly use `make_latest` for stable releases.
15. **Octal Literal Style**: Switched to `0o644` in Go for `gocritic` compliance.
16. **Redundant Newline**: Switched to `fmt.Fprint` for constants with trailing newlines.
17. **Constant Redeclaration**: Resolved by centralizing all command IDs in `constants_cli.go`.
18. **Unchecked Errors**: Added lint-compliant error checking for `db.RemoveInstalledTool`, `dev.Process.Kill()`, and `cmd.Start()`.
19. **Release Pipeline Directory Error**: Resolved `cd: dist` failure by setting explicit `working-directory: gitmap/dist` for compression and checksums. See `spec/02-app-issues/13-release-pipeline-dist-directory.md`.
20. **G305 Zip Path Traversal**: Fixed `installnpp.go` to validate extracted file paths stay within target directory. See `spec/02-app-issues/14-security-hardening-gosec-fixes.md`.
21. **G110 Decompression Bomb**: Replaced `io.Copy` with `io.LimitReader` capped at 10 MB. See `spec/02-app-issues/14-security-hardening-gosec-fixes.md`.
22. **Format Verb Mismatch**: Fixed `fmt.Fprintf` argument count at `tasksync.go:138`; audited ~140 call sites. See `spec/02-app-issues/14-security-hardening-gosec-fixes.md`.
23. **Code Red Error Audit**: Standardized 35+ error constants with mandatory path, operation, and reason context. See `spec/02-app-issues/error-management-file-path-and-missing-file-code-red-rule.md`.
24. **CI Passthrough Gate Pattern**: Job-level `if` skipping caused cached SHA runs to show grey "Skipped" in GitHub UI. Replaced with step-level conditionals so every job reports ✅ Success. See `spec/02-app-issues/16-ci-passthrough-gate-pattern.md`.
25. **Go Flag Ordering — Silent Flag Drop**: Go's `flag` package stops parsing at the first positional argument, silently dropping flags like `-y`. Fixed with `reorderFlagsBeforeArgs()`. See `spec/02-app-issues/17-go-flag-ordering-silent-drop.md`.
26. **CI Release Branch Cancellation Protection**: Unconditional `cancel-in-progress: true` cancelled release branch runs on rapid pushes. Fixed with conditional expression protecting `release/**` branches. See `spec/02-app-issues/18-ci-release-branch-cancellation-protection.md`.
