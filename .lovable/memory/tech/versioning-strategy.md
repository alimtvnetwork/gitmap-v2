# Versioning Strategy

The project follows Semantic Versioning (v2.20.0 current). The `release` system resolves versions using a three-tier priority: 1) Explicit CLI version argument, 2) --bump flag applied to a resolved baseline, 3) Current project version. Baseline resolution for bumps checks `.release/latest.json` first, falling back to scanning local Git tags (`v*`) if metadata is missing.

## IMPORTANT: .release/ Directory Policy

The `.release/` directory should **NOT** be committed to the repository. Release metadata JSON files (`.release/vX.Y.Z.json`, `.release/latest.json`) are local build artifacts, not source code. They must be added to `.gitignore`.

Use `gitmap clear-release-json <version>` (alias `crj`) to remove individual release files when needed.
