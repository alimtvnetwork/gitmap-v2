<#
.SYNOPSIS
    Build, deploy, and run gitmap CLI from the repo root.
.DESCRIPTION
    Pulls latest code, resolves Go dependencies, builds the binary
    into ./bin, copies data folder, deploys to a target directory,
    and optionally runs gitmap with any arguments.
.EXAMPLES
    .\run.ps1                                    # pull, build, deploy
    .\run.ps1 -NoPull                            # skip git pull
    .\run.ps1 -NoDeploy                          # skip deploy step
    .\run.ps1 -R scan                            # build + scan parent folder
    .\run.ps1 -R scan D:\repos                   # build + scan specific path
    .\run.ps1 -R scan D:\repos --mode ssh        # build + scan with flags
    .\run.ps1 -R clone .\gitmap-output\gitmap.json --target-dir .\restored
    .\run.ps1 -R help                            # build + show help
    .\run.ps1 -NoPull -NoDeploy -R scan          # just build and scan
.NOTES
    Configuration is read from gitmap/powershell.json.
    -R accepts ALL gitmap CLI arguments after it (scan, clone, help, flags, paths).
    If -R is used with no arguments, it defaults to: scan <parent folder>
#>

param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [string]$DeployPath = "",
    [switch]$R,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs
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
    param($Config, $BinaryPath, [string[]]$CliArgs)

    Write-Host ""
    Write-Step "RUN" "Executing gitmap"

    $resolvedArgs = Resolve-RunArgs -CliArgs $CliArgs
    $argString = $resolvedArgs -join ' '
    $currentDir = (Get-Location).Path
    Write-Info "Runner CWD: $currentDir"
    Write-Info "Repo root: $RepoRoot"
    Write-Info "Command: gitmap $argString"
    if ($resolvedArgs.Count -ge 2 -and $resolvedArgs[0] -eq "scan") {
        Write-Info "Scan target: $($resolvedArgs[1])"
    }
    Write-Host "  $('─' * 50)" -ForegroundColor DarkGray
    Write-Host ""

    $proc = Start-Process -FilePath $BinaryPath -ArgumentList $resolvedArgs -WorkingDirectory $currentDir -NoNewWindow -Wait -PassThru

    Write-Host ""
    if ($proc.ExitCode -eq 0) {
        Write-Success "Run complete"
    } else {
        Write-Fail "gitmap exited with code $($proc.ExitCode)"
    }
}

# ── Resolve run arguments ─────────────────────────────────────
function Resolve-RunArgs {
    param([string[]]$CliArgs)

    if ($CliArgs.Count -eq 0) {
        $parentDir = Split-Path $RepoRoot -Parent
        Write-Info "No args provided, defaulting to: scan $parentDir"

        return @("scan", $parentDir)
    }

    # Resolve relative paths to absolute so Start-Process always receives correct targets
    $baseDir = (Get-Location).Path
    $resolved = @()
    foreach ($arg in $CliArgs) {
        if ($arg -match '^(\.\.[\\/]|\.[\\/]|\.\.?$)' -and -not $arg.StartsWith('-')) {
            $path = Resolve-Path -LiteralPath $arg -ErrorAction SilentlyContinue
            if ($path) {
                $resolved += $path.Path
            } else {
                $resolved += [System.IO.Path]::GetFullPath((Join-Path $baseDir $arg))
            }
        } else {
            $resolved += $arg
        }
    }

    return $resolved
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

if ($R) {
    Invoke-Run -Config $config -BinaryPath $binaryPath -CliArgs $RunArgs
}

Write-Host ""
Write-Success "All done!"
Write-Host ""
