# Release Pipeline

The release system (v2.24.0) automates a 10-step lifecycle where metadata is written and committed on the original branch after the release process completes:

1. Resolve version (CLI arg, `--bump`, or `version.json`)
2. Pad version to three segments
3. Check for duplicate tags/branches
4. Resolve source ref (`--commit`, `--branch`, or HEAD)
5. Create release branch (`release/vX.Y.Z`)
6. Create annotated tag (`vX.Y.Z`)
7. Push branch and tag; build/upload assets (Go binaries, zip groups, GitHub release)
8. Return to the original branch
9. Write `.release/vX.Y.Z.json` and update `latest.json` on the original branch
10. Auto-commit and push the metadata files

The `release-branch` (rb) and `release-pending` (rp) commands skip steps 9–10 to avoid duplicate metadata artifacts. Native Go cross-compilation targets Windows, Linux, and Darwin (amd64/arm64) with static linking (CGO_ENABLED=0) and direct GitHub API uploads. Asset handling includes compression (.zip/.tar.gz), SHA256 checksums, and target matrix inspection via `--list-targets`.

## CRITICAL: .release/ Directory Policy

The `.release/` directory must **NEVER** be modified by the AI/editor. Release metadata JSON files are local build artifacts managed exclusively by the CLI tool. The AI must not create, edit, or delete any files in `.release/`. Add `.release/` to `.gitignore`. Use `gitmap clear-release-json <version>` (alias `crj`) to remove individual release files when needed.

## Cleanup Command

- `gitmap clear-release-json <version>` (alias: `crj`) — removes a single `.release/vX.Y.Z.json` file
- Only deletes the on-disk JSON file; does not affect Git branches, tags, or the database
- Exits with error if the file does not exist
