# Install Scripts

## Overview

One-liner installer scripts that download, verify, and install the `gitmap`
binary from GitHub Releases. Each script supports version pinning, checksum
verification, and automatic PATH registration.

---

## Repository

| Field       | Value                                              |
|-------------|----------------------------------------------------|
| GitHub Repo | `alimtvnetwork/git-repo-navigator`                 |
| Binary Name | `gitmap` (`gitmap.exe` on Windows)                 |
| Asset Format| `gitmap-{os}-{arch}.zip` (Windows), `gitmap-{os}-{arch}.tar.gz` (Unix) |
| Checksums   | `checksums.txt` (SHA-256, one line per asset)      |

---

## Windows — `install.ps1`

### One-Liner (Full Bootstrap)

The recommended one-liner follows the Chocolatey install pattern: it bypasses
the execution policy for the current process, enforces TLS 1.2+, and
downloads-then-executes the installer script. This ensures the command works
on locked-down machines, older Windows versions, and fresh installs where
`irm` may not be available.

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.ps1'))
```

### Short-Form (PowerShell 5+ / Modern Systems)

If the machine already has TLS 1.2 defaults and unrestricted execution
policy (e.g., developer workstations), the short form also works:

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.ps1 | iex
```

### Why the Full Bootstrap?

| Concern                  | `irm \| iex`      | Full bootstrap           |
|--------------------------|--------------------|--------------------------|
| Execution policy blocked | Fails              | Bypasses (process scope) |
| TLS 1.2 not default      | May fail on old OS | Forces TLS 1.2+         |
| PowerShell 3.x compat   | No (`irm` = PS3+) | Yes (`WebClient` = PS2+) |
| Corporate firewalls      | May fail silently  | Explicit protocol set    |

### Parameters

| Parameter    | Type   | Default                        | Description                        |
|--------------|--------|--------------------------------|------------------------------------|
| `Version`    | string | latest (via GitHub API)        | Pin a specific release tag         |
| `InstallDir` | string | `$env:LOCALAPPDATA\gitmap`     | Target directory for the binary    |
| `Arch`       | string | auto-detect                    | Force `amd64` or `arm64`           |
| `NoPath`     | switch | false                          | Skip adding install dir to PATH    |

### Flow

1. Resolve version — fetch latest tag from GitHub API or use pinned value.
2. Resolve architecture — read `PROCESSOR_ARCHITECTURE` or use override.
3. Download `gitmap-windows-{arch}.zip` and `checksums.txt`.
4. Verify SHA-256 checksum against `checksums.txt`.
5. Extract zip to install directory (rename-first if binary is running).
6. Add install directory to user PATH (unless `--NoPath`).
7. Print installed version via `gitmap version`.

### File

`gitmap/scripts/install.ps1`

---

## Linux / macOS — `install.sh`

### One-Liner

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.sh | sh
```

### Parameters (Environment Variables)

| Variable        | Default              | Description                        |
|-----------------|----------------------|------------------------------------|
| `VERSION`       | latest (GitHub API)  | Pin a specific release tag         |
| `INSTALL_DIR`   | `$HOME/.local/bin`   | Target directory for the binary    |
| `ARCH`          | auto-detect          | Force `amd64` or `arm64`           |
| `NO_PATH`       | unset                | Set to `1` to skip PATH hint       |

### Flow

1. Detect OS (`linux` or `darwin`) and architecture.
2. Resolve version — query GitHub API or use `$VERSION`.
3. Download `gitmap-{os}-{arch}.tar.gz` and `checksums.txt`.
4. Verify SHA-256 checksum (`sha256sum` or `shasum -a 256`).
5. Extract tarball to install directory.
6. Set executable permission (`chmod +x`).
7. Print PATH hint if install directory is not already in PATH.
8. Print installed version via `gitmap version`.

### File

`gitmap/scripts/install.sh` *(planned)*

---

## Checksum Verification

Both scripts download `checksums.txt` from the same release. Each line
follows the format:

```
<sha256-hash>  <filename>
```

The script matches the downloaded asset filename, compares hashes, and
aborts with a clear error on mismatch.

---

## Architecture Detection

| Platform | Source                          | Mapping                          |
|----------|---------------------------------|----------------------------------|
| Windows  | `$env:PROCESSOR_ARCHITECTURE`   | `AMD64`/`x86` → `amd64`, `ARM64` → `arm64` |
| Linux    | `uname -m`                      | `x86_64` → `amd64`, `aarch64` → `arm64`    |
| macOS    | `uname -m`                      | `x86_64` → `amd64`, `arm64` → `arm64`      |

---

## PATH Registration

| Platform | Method                                          |
|----------|-------------------------------------------------|
| Windows  | `[Environment]::SetEnvironmentVariable` (User)  |
| Linux    | Print shell-rc append instruction               |
| macOS    | Print shell-rc append instruction               |

Windows modifies the registry-backed user PATH immediately. Unix scripts
print an instruction the user can copy, avoiding surprise dotfile edits.

---

## Constraints

- No external dependencies beyond `curl`/`PowerShell` and `tar`/`Expand-Archive`.
- Scripts exit non-zero on any failure (download, checksum, extract).
- No interactive prompts — fully automatable.
- Temp files cleaned up in all exit paths.

---

## Related

- [CLI Interface](02-cli-interface.md)
- [Install Bootstrap](83-install-bootstrap.md)
- [Build & Deploy](09-build-deploy.md)
- [Future Features](82-future-features.md)
- [Release Workflow](../../.github/workflows/release.yml)
- [Release Workflow](../../.github/workflows/release.yml)
