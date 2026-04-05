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
    irm https://raw.githubusercontent.com/alimtvnetwork/gitmap/main/scripts/install.ps1 | iex

.EXAMPLE
    & ./install.ps1 -Version v2.48.0

.NOTES
    Repository: https://github.com/alimtvnetwork/gitmap
#>

param(
    [string]$Version = "",
    [string]$InstallDir = "",
    [string]$Arch = "",
    [switch]$NoPath
)

$ErrorActionPreference = "Stop"

$Repo = "alimtvnetwork/gitmap"
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

function Add-ToPath([string]$dir) {
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $parts = $currentPath -split ";"

    foreach ($part in $parts) {
        if ($part.Trim() -ieq $dir) {
            Write-Step "Already in PATH."
            return
        }
    }

    $newPath = $currentPath.TrimEnd(";") + ";" + $dir
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    $env:PATH = $env:PATH + ";" + $dir
    Write-OK "Added to PATH (restart terminal for full effect)."
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
