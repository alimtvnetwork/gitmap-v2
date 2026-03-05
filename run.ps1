<#
.SYNOPSIS
    Build, deploy, and run gitmap CLI from the repo root.
.DESCRIPTION
    Pulls latest code, resolves Go dependencies, builds the binary
    into ./bin, copies data folder, deploys to a target directory,
    and optionally runs gitmap on a specified path.
.EXAMPLES
    .\run.ps1                        # pull, build, deploy
    .\run.ps1 -NoPull                # skip git pull
    .\run.ps1 -NoDeploy              # skip deploy step
    .\run.ps1 -Run                   # build + run on parent folder
    .\run.ps1 -Run -RunPath "D:\repos"  # build + run on specific path
    .\run.ps1 -Run -RunArgs "--mode ssh --output csv"
.NOTES
    Configuration is read from gitmap/powershell.json.
#>

param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [switch]$Run,
    [string]$RunPath,
    [string]$RunArgs,
    [string]$DeployPath
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$GitMapDir = Join-Path $RepoRoot "gitmap"

# ── Logging helpers ───────────────────────────────────────────
function Write-Step {
    param([string]$Step, [string]$Message)
    Write-Host ""
    Write-Host "  [$Step] " -ForegroundColor Magenta -NoNewline
    Write-Host $Message -ForegroundColor White
    Write-Host "  $('─' * 50)" -ForegroundColor DarkGray
}

function Write-Success {
    param([string]$Message)
    Write-Host "  ✓ " -ForegroundColor Green -NoNewline
    Write-Host $Message -ForegroundColor Green
}

function Write-Info {
    param([string]$Message)
    Write-Host "  → " -ForegroundColor Cyan -NoNewline
    Write-Host $Message -ForegroundColor Gray
}

function Write-Warn {
    param([string]$Message)
    Write-Host "  ⚠ " -ForegroundColor Yellow -NoNewline
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Fail {
    param([string]$Message)
    Write-Host "  ✗ " -ForegroundColor Red -NoNewline
    Write-Host $Message -ForegroundColor Red
}

# ── Banner ────────────────────────────────────────────────────
function Show-Banner {
    Write-Host ""
    Write-Host "  ╔══════════════════════════════════════╗" -ForegroundColor DarkCyan
    Write-Host "  ║         " -ForegroundColor DarkCyan -NoNewline
    Write-Host "gitmap builder" -ForegroundColor Cyan -NoNewline
    Write-Host "              ║" -ForegroundColor DarkCyan
    Write-Host "  ╚══════════════════════════════════════╝" -ForegroundColor DarkCyan
    Write-Host ""
}

# ── Load config ───────────────────────────────────────────────
function Load-Config {
    $configPath = Join-Path $GitMapDir "powershell.json"
    if (Test-Path $configPath) {
        Write-Info "Config loaded from powershell.json"

        return Get-Content $configPath | ConvertFrom-Json
    }
    Write-Warn "No powershell.json found, using defaults"

    return @{
        deployPath  = "E:\bin-run"
        buildOutput = "./bin"
        binaryName  = "gitmap.exe"
        copyData    = $true
    }
}

# ── Git pull ──────────────────────────────────────────────────
function Invoke-GitPull {
    Write-Step "1/4" "Pulling latest changes"
    Push-Location $RepoRoot
    try {
        $output = git pull 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Fail "Git pull failed"
            Write-Host "  $output" -ForegroundColor Red
            exit 1
        }
        foreach ($line in $output) {
            Write-Info $line
        }
        Write-Success "Pull complete"
    } finally {
        Pop-Location
    }
}

# ── Resolve dependencies ─────────────────────────────────────
function Resolve-Dependencies {
    Write-Step "2/4" "Resolving Go dependencies"
    Push-Location $GitMapDir
    try {
        $tidyOutput = go mod tidy 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Fail "go mod tidy failed"
            Write-Host "  $tidyOutput" -ForegroundColor Red
            exit 1
        }
        Write-Success "Dependencies resolved"
    } finally {
        Pop-Location
    }
}

# ── Build binary ──────────────────────────────────────────────
function Build-Binary {
    param($Config)

    Write-Step "3/4" "Building $($Config.binaryName)"

    $binDir  = Join-Path $RepoRoot $Config.buildOutput
    $outPath = Join-Path $binDir $Config.binaryName

    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
        Write-Info "Created bin directory"
    }

    Push-Location $GitMapDir
    try {
        $buildOutput = go build -o $outPath . 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Fail "Go build failed"
            foreach ($line in $buildOutput) {
                Write-Host "  $line" -ForegroundColor Red
            }
            exit 1
        }
    } finally {
        Pop-Location
    }

    if ($Config.copyData) {
        Copy-DataFolder -BinDir $binDir
    }

    $size = (Get-Item $outPath).Length / 1MB
    Write-Success ("Binary built ({0:N2} MB) → $outPath" -f $size)

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
        Write-Info "Copied data folder to bin"
    }
}

# ── Deploy to target directory ────────────────────────────────
function Deploy-Binary {
    param($Config, $BinaryPath, $OverridePath)

    Write-Step "4/4" "Deploying"

    $target = $Config.deployPath
    if ($OverridePath.Length -gt 0) {
        $target = $OverridePath
    }

    Write-Info "Target: $target"

    if (-not (Test-Path $target)) {
        New-Item -ItemType Directory -Path $target -Force | Out-Null
        Write-Info "Created deploy directory"
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
        Write-Info "Copied data folder to deploy target"
    }

    Write-Success "Deployed to $target"
    $cmdName = $Config.binaryName -replace '\.exe$', ''
    Write-Info "You can now run: $cmdName"
}

# ── Run gitmap ────────────────────────────────────────────────
function Invoke-Run {
    param($Config, $BinaryPath, $TargetPath, $ExtraArgs)

    Write-Host ""
    Write-Step "RUN" "Executing gitmap"

    $scanPath = $TargetPath
    if ($scanPath.Length -eq 0) {
        $scanPath = Split-Path $RepoRoot -Parent
        Write-Info "No path provided, scanning parent: $scanPath"
    } else {
        Write-Info "Scanning: $scanPath"
    }

    $fullCmd = "& `"$BinaryPath`" scan `"$scanPath`""
    if ($ExtraArgs.Length -gt 0) {
        $fullCmd = "$fullCmd $ExtraArgs"
        Write-Info "Extra args: $ExtraArgs"
    }

    Write-Host "  $('─' * 50)" -ForegroundColor DarkGray
    Write-Host ""

    Invoke-Expression $fullCmd

    Write-Host ""
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Run complete"
    } else {
        Write-Fail "gitmap exited with code $LASTEXITCODE"
    }
}

# ── Main ──────────────────────────────────────────────────────
Show-Banner
$config = Load-Config

if (-not $NoPull) {
    Invoke-GitPull
} else {
    Write-Info "Skipping git pull (--NoPull)"
}

Resolve-Dependencies
$binaryPath = Build-Binary -Config $config

if (-not $NoDeploy) {
    Deploy-Binary -Config $config -BinaryPath $binaryPath -OverridePath $DeployPath
} else {
    Write-Info "Skipping deploy (--NoDeploy)"
}

if ($Run) {
    Invoke-Run -Config $config -BinaryPath $binaryPath -TargetPath $RunPath -ExtraArgs $RunArgs
}

Write-Host ""
Write-Success "All done!"
Write-Host ""
