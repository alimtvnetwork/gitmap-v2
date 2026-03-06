# PowerShell Build & Deploy Patterns

## Overview

This document describes reusable patterns for PowerShell build scripts
that manage the lifecycle of compiled CLI tools: pull, build, deploy,
and run.

## Script Architecture

### Single Entry Point

One script (`run.ps1`) at the repo root handles the full lifecycle.
It reads configuration from a JSON file and exposes behavior through
switch/string parameters.

```powershell
[CmdletBinding(PositionalBinding=$false)]
param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [string]$DeployPath = "",
    [switch]$Update,
    [switch]$R,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs
)
```

### Step-Based Execution

Break the pipeline into numbered steps:

| Step | Action | Skippable |
|------|--------|-----------|
| 1/4 | Git pull | `-NoPull` |
| 2/4 | Resolve dependencies | No |
| 3/4 | Build binary | No |
| 4/4 | Deploy to target | `-NoDeploy` |

Each step is a dedicated function with clear responsibility.

## Configuration Pattern

### External JSON Config

Store build/deploy settings in a JSON file alongside the source:

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "toolname.exe",
  "copyData": true
}
```

### Config Loading

```powershell
function Load-Config {
    $configPath = Join-Path $ProjectDir "powershell.json"
    if (Test-Path $configPath) {
        return Get-Content $configPath | ConvertFrom-Json
    }
    # Return sensible defaults if file is missing
    return @{
        deployPath  = "E:\bin-run"
        buildOutput = "./bin"
        binaryName  = "toolname.exe"
        copyData    = $true
    }
}
```

### Rules

- CLI flags always override config file values.
- Missing config file is a warning, not an error.
- All paths in config are relative to the project root unless absolute.

## Logging Pattern

### Semantic Logging Functions

Use color-coded helper functions for consistent output:

| Function | Color | Prefix | Use Case |
|----------|-------|--------|----------|
| `Write-Step` | Magenta | `[N/M]` | Step headers |
| `Write-Success` | Green | `OK` | Successful operations |
| `Write-Info` | Cyan/Gray | `->` | Informational messages |
| `Write-Warn` | Yellow | `!!` | Non-fatal warnings |
| `Write-Fail` | Red | `XX` | Errors before exit |

### Banner

Display an ASCII banner at script start for visual identity:

```powershell
function Show-Banner {
    Write-Host "  +======================================+"
    Write-Host "  |         toolname builder             |"
    Write-Host "  +======================================+"
}
```

## Build Patterns

### Build with Embedded Variables

Use Go's `-ldflags` to embed values at compile time:

```powershell
$absRepoRoot = (Resolve-Path $RepoRoot).Path
$ldflags = "-X 'pkg/constants.RepoPath=$absRepoRoot'"
go build -ldflags $ldflags -o $outPath .
```

### Version Verification

After building, immediately run the binary with `version` to confirm:

```powershell
$versionOutput = & $binaryPath version 2>&1
Write-Info "Version: $versionOutput"
```

This catches build issues early — if the version doesn't match
expectations, the build is suspect.

### Data Folder Copy

If the binary needs companion data files, copy them alongside:

```powershell
if ($Config.copyData) {
    Copy-Item $dataSource $dataDest -Recurse
}
```

## Deploy Patterns

### Retry-on-Lock

When deploying to a target that may be in use (especially on Windows),
wrap `Copy-Item` in a retry loop:

```powershell
$maxAttempts = 20
$attempt = 1
while ($true) {
    try {
        Copy-Item $BinaryPath $destFile -Force -ErrorAction Stop
        break
    } catch {
        if ($attempt -ge $maxAttempts) { throw }
        Write-Warn "Target is in use; retrying ($attempt/$maxAttempts)..."
        Start-Sleep -Milliseconds 500
        $attempt++
    }
}
```

### Nested Deploy Structure

Deploy the binary into a named subfolder within the target directory.
This keeps the deploy target organized when multiple tools share the
same parent directory:

```
deploy-target/
└── toolname/
    ├── toolname.exe
    └── data/
        └── config.json
```

The subfolder (not the parent) should be added to the system `PATH`.

### Deploy Target on PATH

The deploy directory should be on the system `PATH` so the tool can
be run from any terminal without specifying the full path.

## Run Pattern (`-R` Flag)

### Forwarding Arguments

Use `ValueFromRemainingArguments` to capture all trailing arguments
after the `-R` switch and forward them to the built binary:

```powershell
[switch]$R,
[Parameter(ValueFromRemainingArguments=$true)]
[string[]]$RunArgs
```

### Path Resolution

Resolve relative paths to absolute before passing to the binary,
since `Start-Process` may run from a different working directory:

```powershell
foreach ($arg in $CliArgs) {
    if ($arg -match '^(\.\.[\\/]|\.[\\/]|\.\.?$)') {
        $path = Resolve-Path -LiteralPath $arg -ErrorAction SilentlyContinue
        if ($path) { $resolved += $path.Path }
        else { $resolved += [System.IO.Path]::GetFullPath((Join-Path $baseDir $arg)) }
    } else {
        $resolved += $arg
    }
}
```

### Default Behavior

If `-R` is used with no arguments, default to a sensible action
(e.g., process the parent folder of the repo).

### Context Logging

Before executing, print diagnostic info:

```
[RUN] Executing toolname
→ Runner CWD: D:\projects\my-tool
→ Command: toolname scan D:\projects
→ Scan target: D:\projects
```

## Self-Update Orchestration (Windows-Safe)

When a CLI updates itself from a PATH-managed executable, use a two-phase handoff so the active binary lock is released before deploy.

### Phase 1: Handoff from active binary
1. `tool update` creates a handoff copy in the same active binary directory (for example `toolname.exe.old` or `toolname-update-<pid>.exe`).
2. It launches the handoff copy with `update --from-copy`.
3. The parent exits immediately.

### Phase 2: Execute update from handoff copy
1. Resolve repo root from embedded/configured repo path.
2. Run `run.ps1 -Update` (pull, build, deploy).
3. Sync active PATH binary from deployed binary using retry loop first, then rename fallback.
4. Read and print versions from the binaries (before update and after update) using `tool version`.
5. Show latest notes using the updated binary (`tool changelog --latest`).
6. Run `tool update-cleanup` to remove handoff and `.old` artifacts.

### Required Validation
- Fail the update if active version still does not match deployed version after sync.
- Version/changelog output must come from the updated executable, not static constants.
- Cleanup must run after successful update so rollback artifacts exist during deploy.

## Error Handling

| Pattern | Implementation |
|---------|----------------|
| `$ErrorActionPreference = "Stop"` | Fail fast on uncaught errors |
| Check `$LASTEXITCODE` after external commands | Detect non-PowerShell failures |
| Print error details before `exit 1` | User sees what went wrong |
| Use `try/finally` with `Push-Location/Pop-Location` | Always restore working directory |
