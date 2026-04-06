<#
.SYNOPSIS
    One-liner installer for gitmap CLI.

.DESCRIPTION
    Downloads the latest gitmap release from GitHub, verifies checksums,
    extracts to a local directory, and adds it to PATH.

.PARAMETER Version
    Install a specific version (e.g. v2.48.0). Default: latest.

.PARAMETER InstallDir
    Target directory. Default: $env:LOCALAPPDATA\gitmap

.PARAMETER NoPath
    Skip adding to PATH.

.PARAMETER Arch
    Force architecture (amd64, arm64). Default: auto-detect.

.EXAMPLE
    irm https://raw.githubusercontent.com/alimtvnetwork/git-repo-navigator/main/gitmap/scripts/install.ps1 | iex

.EXAMPLE
    & ./install.ps1 -Version v2.48.0

.NOTES
    Repository: https://github.com/alimtvnetwork/git-repo-navigator
#>

param(
    [string]$Version = "",
    [string]$InstallDir = "",
    [string]$Arch = "",
    [switch]$NoPath
)

$ErrorActionPreference = "Stop"

$Repo = "alimtvnetwork/git-repo-navigator"
$BinaryName = "gitmap.exe"

# --- Logging helpers ---

function Write-Step([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Cyan
}

function Write-OK([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Green
}

function Write-Err([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Red
}

# --- Resolve install directory ---

function Resolve-InstallDir([string]$dir) {
    if ($dir -ne "") { return $dir }
    return Join-Path $env:LOCALAPPDATA "gitmap"
}

# --- Detect architecture ---

function Resolve-Arch([string]$arch) {
    if ($arch -ne "") { return $arch }

    $cpu = $env:PROCESSOR_ARCHITECTURE
    switch ($cpu) {
        "AMD64"   { return "amd64" }
        "ARM64"   { return "arm64" }
        "x86"     { return "amd64" }
        default   { return "amd64" }
    }
}

# --- Resolve version (latest or pinned) ---

function Resolve-Version([string]$version) {
    if ($version -ne "") { return $version }

    Write-Step "Fetching latest release..."
    $url = "https://api.github.com/repos/$Repo/releases/latest"

    try {
        $release = Invoke-RestMethod -Uri $url -UseBasicParsing
        return $release.tag_name
    }
    catch {
        Write-Err "Failed to fetch latest release: $_"
        exit 1
    }
}

# --- Download asset ---

function Get-Asset([string]$version, [string]$arch) {
    $assetName = "gitmap-windows-${arch}.zip"
    $baseUrl = "https://github.com/$Repo/releases/download/$version"
    $assetUrl = "$baseUrl/$assetName"
    $checksumUrl = "$baseUrl/checksums.txt"

    $tmpDir = Join-Path $env:TEMP "gitmap-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

    $zipPath = Join-Path $tmpDir $assetName
    $checksumPath = Join-Path $tmpDir "checksums.txt"

    Write-Step "Downloading $assetName ($version)..."

    try {
        Invoke-WebRequest -Uri $assetUrl -OutFile $zipPath -UseBasicParsing
        Invoke-WebRequest -Uri $checksumUrl -OutFile $checksumPath -UseBasicParsing
    }
    catch {
        Write-Err "Download failed: $_"
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }

    # Verify checksum
    Write-Step "Verifying checksum..."
    $expectedLine = (Get-Content $checksumPath | Where-Object { $_ -match $assetName })
    if (-not $expectedLine) {
        Write-Err "Asset not found in checksums.txt"
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }

    $expectedHash = ($expectedLine -split '\s+')[0]
    $actualHash = (Get-FileHash -Path $zipPath -Algorithm SHA256).Hash.ToLower()

    if ($actualHash -ne $expectedHash) {
        Write-Err "Checksum mismatch!"
        Write-Err "  Expected: $expectedHash"
        Write-Err "  Got:      $actualHash"
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }

    Write-OK "Checksum verified."
    return @{ ZipPath = $zipPath; TmpDir = $tmpDir }
}

# --- Extract and install ---

function Install-Binary([string]$zipPath, [string]$installDir) {
    Write-Step "Installing to $installDir..."

    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }

    $targetPath = Join-Path $installDir $BinaryName

    # Rename-first strategy for running binary
    if (Test-Path $targetPath) {
        $oldPath = "$targetPath.old"
        if (Test-Path $oldPath) { Remove-Item $oldPath -Force }
        Rename-Item $targetPath $oldPath -Force
    }

    Expand-Archive -Path $zipPath -DestinationPath $installDir -Force

    # Handle nested structure: if extracted into a subfolder
    $extracted = Get-ChildItem -Path $installDir -Filter $BinaryName -Recurse | Select-Object -First 1
    if ($extracted -and $extracted.DirectoryName -ne $installDir) {
        Move-Item $extracted.FullName $targetPath -Force
    }

    # Cleanup .old
    $oldPath = "$targetPath.old"
    if (Test-Path $oldPath) {
        Remove-Item $oldPath -Force -ErrorAction SilentlyContinue
    }

    Write-OK "Installed $BinaryName to $installDir"
}

# --- Add to PATH ---

function Test-PathEntry([string]$pathValue, [string]$dir) {
    if ([string]::IsNullOrWhiteSpace($pathValue)) {
        return $false
    }

    $parts = $pathValue -split ";"

    foreach ($part in $parts) {
        if ($part.Trim() -ieq $dir) {
            return $true
        }
    }

    return $false
}

function Rebuild-SessionPath([string]$dir) {
    # Rebuild session PATH from registry (Machine + User) to pick up any changes
    $machinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $parts = @()
    if ($machinePath) { $parts += $machinePath.TrimEnd(";") }
    if ($userPath) { $parts += $userPath.TrimEnd(";") }
    $rebuilt = $parts -join ";"

    # Ensure install dir is present even if not yet persisted
    if (-not (Test-PathEntry $rebuilt $dir)) {
        $rebuilt = $rebuilt.TrimEnd(";") + ";" + $dir
    }

    return $rebuilt
}

function Broadcast-EnvironmentChange {
    Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;

public static class GitMapEnvNative {
    [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd,
        uint Msg,
        UIntPtr wParam,
        string lParam,
        uint fuFlags,
        uint uTimeout,
        out UIntPtr lpdwResult
    );
}
"@ -ErrorAction SilentlyContinue | Out-Null

    $HWND_BROADCAST = [IntPtr]0xffff
    $WM_SETTINGCHANGE = 0x001A
    $SMTO_ABORTIFHUNG = 0x0002
    [UIntPtr]$result = [UIntPtr]::Zero

    [void][GitMapEnvNative]::SendMessageTimeout(
        $HWND_BROADCAST,
        $WM_SETTINGCHANGE,
        [UIntPtr]::Zero,
        "Environment",
        $SMTO_ABORTIFHUNG,
        5000,
        [ref]$result
    )
}

function Add-ToPath([string]$dir) {
    $currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $userHasDir = Test-PathEntry $currentUserPath $dir

    if (-not $userHasDir) {
        if ([string]::IsNullOrWhiteSpace($currentUserPath)) {
            $newPath = $dir
        }
        else {
            $newPath = $currentUserPath.TrimEnd(";") + ";" + $dir
        }

        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Broadcast-EnvironmentChange
        Write-OK "Added to user PATH."
    }
    else {
        Write-Step "Already in user PATH."
    }
}

# --- Main ---

function Main {
    Write-Host ""
    Write-Host "  gitmap installer" -ForegroundColor White
    Write-Host "  github.com/$Repo" -ForegroundColor DarkGray
    Write-Host ""

    $resolvedArch = Resolve-Arch $Arch
    $resolvedVersion = Resolve-Version $Version
    $resolvedDir = Resolve-InstallDir $InstallDir

    $result = Get-Asset $resolvedVersion $resolvedArch

    try {
        Install-Binary $result.ZipPath $resolvedDir
    }
    finally {
        Remove-Item $result.TmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }

    if (-not $NoPath) {
        Add-ToPath $resolvedDir
    }

    # Also try Chocolatey refreshenv if available
    $refreshCmd = Get-Command refreshenv -ErrorAction SilentlyContinue
    if ($refreshCmd) {
        try { refreshenv | Out-Null } catch {}
    }

    # Verify
    $binPath = Join-Path $resolvedDir $BinaryName
    if (Test-Path $binPath) {
        Write-Host ""
        $versionOutput = & $binPath version 2>&1
        Write-OK "gitmap $versionOutput"
    }

    Write-Host ""
    Write-OK "Done! Run 'gitmap --help' to get started."
    Write-Host ""
}

Main
