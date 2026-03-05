# Build, Deploy & Run

## Overview

The project uses a single PowerShell script (`run.ps1`) at the repo root
to pull, build, deploy, and optionally run the gitmap CLI.
Build configuration lives in `gitmap/powershell.json`.

## Build Script — `run.ps1`

| Step | Description |
|------|-------------|
| 1. Git Pull | Pulls latest changes from remote |
| 2. Resolve Deps | Runs `go mod tidy` in `gitmap/` |
| 3. Build | Compiles binary to `./bin/gitmap.exe` |
| 4. Deploy | Copies binary + `data/` to deploy target |

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-NoPull` | Skip `git pull` | pull enabled |
| `-NoDeploy` | Skip deploy step | deploy enabled |
| `-DeployPath <dir>` | Override deploy directory | from `powershell.json` |
| `-R <args...>` | Run gitmap after build with given CLI arguments | — |

### Examples

```powershell
# Full pipeline: pull, build, deploy
.\run.ps1

# Build only, no pull or deploy
.\run.ps1 -NoPull -NoDeploy

# Build and scan parent folder
.\run.ps1 -R scan

# Build and scan specific folder with SSH mode
.\run.ps1 -R scan D:\repos --mode ssh

# Build and clone from JSON
.\run.ps1 -R clone .\gitmap-output\gitmap.json --target-dir .\restored

# Build and clone with GitHub Desktop registration
.\run.ps1 -R clone .\gitmap.json --github-desktop

# Deploy to custom path
.\run.ps1 -DeployPath "D:\tools"
```

## Configuration — `gitmap/powershell.json`

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "gitmap.exe",
  "copyData": true
}
```

| Field | Description | Default |
|-------|-------------|---------|
| `deployPath` | Directory where binary is deployed | `E:\bin-run` |
| `buildOutput` | Local build output directory | `./bin` |
| `binaryName` | Name of the compiled binary | `gitmap.exe` |
| `copyData` | Whether to copy `data/` alongside binary | `true` |

## Build Output

After a successful build, the `./bin/` directory contains:

```
bin/
├── gitmap.exe
└── data/
    └── config.json
```

The same structure is replicated at the deploy target.

## `-R` Flag Behavior

The `-R` flag passes **all remaining arguments** directly to the gitmap
binary. It accepts any valid gitmap command and flags.

- If `-R` is used with no arguments, it defaults to `scan <parent-folder>`.
- `-R` runs after build and deploy steps complete.

### Path Resolution

Relative path arguments (e.g., `..`, `../..`, `./projects`) are
automatically resolved to **absolute paths** before being passed to the
gitmap binary. This ensures the binary receives correct paths regardless
of the working directory used by `Start-Process`.

```powershell
# User runs:
.\run.ps1 -R scan "../.."

# Script resolves "../.." to absolute, e.g.:
# gitmap scan D:\wp-work
```

## Deploy Target

The default deploy path (`E:\bin-run`) is assumed to be on the system
`PATH`, so after deployment the user can run `gitmap` from any terminal
without specifying the full path.

## Logging

The script uses colored, step-numbered output:

- **Magenta** — step headers (`[1/4]`, `[2/4]`, etc.)
- **Green** — success messages (✓)
- **Cyan** — informational messages (→)
- **Yellow** — warnings (⚠)
- **Red** — errors (✗)
