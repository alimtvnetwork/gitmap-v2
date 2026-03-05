<#
.SYNOPSIS
    Build and deploy gitmap CLI from the repo root.
.DESCRIPTION
    Pulls latest code, builds the Go binary into ./bin,
    copies data folder alongside it, and optionally deploys
    to a configured directory (e.g. E:\bin-run).
.NOTES
    Configuration is read from gitmap/powershell.json.
#>

param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [string]$DeployPath
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$GitMapDir = Join-Path $RepoRoot "gitmap"

# ── Load config ───────────────────────────────────────────────
function Load-Config {
    $configPath = Join-Path $GitMapDir "powershell.json"
    if (Test-Path $configPath) {
        return Get-Content $configPath | ConvertFrom-Json
    }
    return @{
        deployPath  = "E:\bin-run"
        buildOutput = "./bin"
        binaryName  = "gitmap.exe"
        copyData    = $true
    }
}

# ── Git pull ──────────────────────────────────────────────────
function Invoke-GitPull {
    Write-Host "Pulling latest changes..." -ForegroundColor Cyan
    Push-Location $RepoRoot
    try {
        git pull
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Git pull failed."
            exit 1
        }
    } finally {
        Pop-Location
    }
    Write-Host "Pull complete." -ForegroundColor Green
}

# ── Build binary ──────────────────────────────────────────────
function Build-Binary {
    param($Config)

    $binDir  = Join-Path $RepoRoot $Config.buildOutput
    $outPath = Join-Path $binDir $Config.binaryName

    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    }

    Write-Host "Building $($Config.binaryName)..." -ForegroundColor Cyan
    Push-Location $GitMapDir
    try {
        go build -o $outPath .
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Go build failed."
            exit 1
        }
    } finally {
        Pop-Location
    }

    if ($Config.copyData) {
        Copy-DataFolder -BinDir $binDir
    }

    Write-Host "Build complete: $outPath" -ForegroundColor Green
    return $outPath
}

# ── Copy data folder ─────────────────────────────────────────
function Copy-DataFolder {
    param($BinDir)

    $dataSource = Join-Path $GitMapDir "data"
    $dataDest   = Join-Path $BinDir "data"

    if (Test-Path $dataSource) {
        if (Test-Path $dataDest) {
            Remove-Item $dataDest -Recurse -Force
        }
        Copy-Item $dataSource $dataDest -Recurse
        Write-Host "Copied data folder to bin." -ForegroundColor Gray
    }
}

# ── Deploy to target directory ────────────────────────────────
function Deploy-Binary {
    param($Config, $BinaryPath, $OverridePath)

    $target = $Config.deployPath
    if ($OverridePath.Length -gt 0) {
        $target = $OverridePath
    }

    if (-not (Test-Path $target)) {
        New-Item -ItemType Directory -Path $target -Force | Out-Null
    }

    $destFile = Join-Path $target $Config.binaryName
    Copy-Item $BinaryPath $destFile -Force

    $binDir   = Split-Path $BinaryPath -Parent
    $dataDir  = Join-Path $binDir "data"
    $dataDest = Join-Path $target "data"
    if (Test-Path $dataDir) {
        if (Test-Path $dataDest) {
            Remove-Item $dataDest -Recurse -Force
        }
        Copy-Item $dataDir $dataDest -Recurse
    }

    Write-Host "Deployed to $target" -ForegroundColor Green
    Write-Host "You can now run: $($Config.binaryName -replace '\.exe$','')" -ForegroundColor Yellow
}

# ── Main ──────────────────────────────────────────────────────
$config = Load-Config

if (-not $NoPull) {
    Invoke-GitPull
}

$binaryPath = Build-Binary -Config $config

if (-not $NoDeploy) {
    Deploy-Binary -Config $config -BinaryPath $binaryPath -OverridePath $DeployPath
}

Write-Host "`nDone!" -ForegroundColor Green
