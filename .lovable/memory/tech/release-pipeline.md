# Release Pipeline

The release system (v2.20.0) automates a 10-step lifecycle via the `release` (r), `release-branch` (rb), and `release-pending` (rp) commands. It features native Go cross-compilation for Windows, Linux, and Darwin (amd64/arm64) with static linking (CGO_ENABLED=0) and direct GitHub API asset uploads. Automated asset handling includes compression (.zip for Windows, .tar.gz for Unix), SHA256 `checksums.txt` generation, and target matrix inspection via the `--list-targets` flag. Target resolution validates `GOOS/GOARCH` pairs against a hardcoded matrix.

## IMPORTANT: .release/ Directory Policy

The `.release/` directory should **NOT** be committed to the repository. Release metadata JSON files are local build artifacts, not source code. Add `.release/` to `.gitignore`. Use `gitmap clear-release-json <version>` (alias `crj`) to remove individual release files when needed.

## Cleanup Command

- `gitmap clear-release-json <version>` (alias: `crj`) — removes a single `.release/vX.Y.Z.json` file
- Only deletes the on-disk JSON file; does not affect Git branches, tags, or the database
- Exits with error if the file does not exist
