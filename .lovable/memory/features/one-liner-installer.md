# One-Liner Installer

The project provides cross-platform, URL-based one-liner installer and updater scripts. Both `install.ps1` (Windows/PowerShell) and `install.sh` (Linux/macOS/Bash) are CI-generated, version-pinned release assets — created during the GitHub Actions release pipeline via `sed` placeholder substitution and attached alongside the compiled binaries.

## Windows — `install.ps1`

- Suppresses progress bar UI via `$ProgressPreference = "SilentlyContinue"` to prevent terminal crashes during `irm | iex` execution.
- Uses regex-based matching (e.g., `^gitmap-v[\d.]+-windows-.*\.exe$`) to identify versioned binaries from CI.
- Verifies SHA256 checksums against `checksums.txt`.
- Includes `try/catch` error handling for graceful failure recovery with manual download fallback.
- CLI flags: `-Version`, `-InstallDir`, `-Arch`, `-NoPath`.

## Linux / macOS — `install.sh`

- Supports `bash`, `zsh`, and `fish` with shell-aware PATH registration (updating `.bashrc`, `.zshrc`, or `config.fish`).
- Handles both `.tar.gz` and `.zip` archive formats with 4-priority binary detection.
- Verifies SHA256 checksums via `sha256sum` or `shasum`.
- Conditional error handling for graceful failure recovery.
- CLI flags: `--version`, `--dir`, `--arch`, `--no-path`.

## CI Integration

Both scripts are generated in `.github/workflows/release.yml` step 5, with version placeholders replaced by the release tag. They are uploaded as release assets alongside the 12 cross-compiled binaries and `checksums.txt`. The release body includes both PowerShell and Bash one-liner install instructions.
