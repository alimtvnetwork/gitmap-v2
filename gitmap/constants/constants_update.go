package constants

// Update handoff file patterns.
const (
	UpdateCopyFmt    = "gitmap-update-%d.exe"
	UpdateCopyGlob   = "gitmap-update-*.exe"
	UpdateScriptGlob = "gitmap-update-*.ps1"
)

// Update flags.
const FlagVerbose = "--verbose"

// Update UI messages.
const (
	MsgUpdateActive      = "  → Active: %s\n  → Handoff: %s\n"
	MsgUpdateCleanStart  = "\n  Cleaning up update artifacts..."
	MsgUpdateCleanDone   = "  ✓ Removed %d file(s)\n\n"
	MsgUpdateCleanNone   = "  ✓ Nothing to clean up\n"
	MsgUpdateTempRemoved = "  → Removed temp copy: %s\n"
	MsgUpdateOldRemoved  = "  → Removed backup: %s\n"
	UpdateRunnerLogStart = "update-runner starting, repo=%s"
	UpdateScriptLogExec  = "executing update script: %s"
	UpdateScriptLogExit  = "update script exited: err=%v"
)

// Update error messages.
const (
	ErrUpdateExecFind = "Error finding executable: %v\n"
	ErrUpdateCopyFail = "Error creating update copy: %v\n"
)

// Update PowerShell script template sections.
const (
	UpdatePSHeader = `# gitmap self-update script (auto-generated)
Set-Location "%s"
`
	UpdatePSDeployDetect = `
$configPath = Join-Path "%s" "gitmap\powershell.json"
$deployedBinary = $null
if (Test-Path $configPath) {
    $cfg = Get-Content $configPath | ConvertFrom-Json
    if ($cfg.deployPath) {
        $deployedBinary = Join-Path $cfg.deployPath "gitmap\gitmap.exe"
    }
}
`
	UpdatePSVersionBefore = `
$activeBinary = $null
$activeBefore = "unknown"
$cmdBefore = Get-Command gitmap -ErrorAction SilentlyContinue
if ($cmdBefore -and (Test-Path $cmdBefore.Source)) {
    $activeBinary = $cmdBefore.Source
    $activeBefore = & $activeBinary version 2>&1
}
`
	UpdatePSRunUpdate = `
Write-Host ""
Write-Host "  Starting update via run.ps1 -Update" -ForegroundColor Cyan
& "%s" -Update
$runExit = $LASTEXITCODE
if (($runExit -ne 0) -and ($runExit -ne $null)) {
    exit $runExit
}
`
	UpdatePSVersionAfter = `
$activeAfter = "unknown"
$deployedAfter = "unknown"
$cmdAfter = Get-Command gitmap -ErrorAction SilentlyContinue
if ($cmdAfter -and (Test-Path $cmdAfter.Source)) {
    $activeBinary = $cmdAfter.Source
    $activeAfter = & $activeBinary version 2>&1
}
if ($deployedBinary -and (Test-Path $deployedBinary)) {
    $deployedAfter = & $deployedBinary version 2>&1
}
`
	UpdatePSVerify = `
Write-Host ""
Write-Host "  Version before:   $activeBefore" -ForegroundColor DarkGray
Write-Host "  Version active:   $activeAfter" -ForegroundColor DarkGray
Write-Host "  Version deployed: $deployedAfter" -ForegroundColor DarkGray

$lastReleaseScript = Join-Path "%s" "gitmap" "scripts" "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    & $lastReleaseScript -BinaryPath $activeBinary -RepoRoot "%s"
}

if (($activeAfter -eq "unknown") -or ($deployedAfter -eq "unknown") -or ($activeAfter -ne $deployedAfter)) {
    Write-Host "  [FAIL] Active PATH version does not match deployed version." -ForegroundColor Red
    exit 1
}

Write-Host "  [OK] Active PATH binary matches deployed version." -ForegroundColor Green
`
	UpdatePSPostActions = `
if ($activeBinary -and (Test-Path $activeBinary)) {
    Write-Host ""
    Write-Host "  Latest changelog:" -ForegroundColor Cyan
    & $activeBinary changelog --latest

    Write-Host ""
    Write-Host "  Cleaning update artifacts..." -ForegroundColor DarkGray
    & $activeBinary update-cleanup
}

Write-Host ""
exit 0
`
)

// Revert PowerShell script template sections.
const (
	RevertPSHeader = `# gitmap revert script (auto-generated)
Set-Location "%s"
`
	RevertPSBuild = `
Write-Host ""
Write-Host "  Building from checked-out version..." -ForegroundColor Cyan
& "%s"
$runExit = $LASTEXITCODE
if (($runExit -ne 0) -and ($runExit -ne $null)) {
    exit $runExit
}
`
	RevertPSPostActions = `
$cmdAfter = Get-Command gitmap -ErrorAction SilentlyContinue
if ($cmdAfter -and (Test-Path $cmdAfter.Source)) {
    $activeAfter = & $cmdAfter.Source version 2>&1
    Write-Host "  Active version: $activeAfter" -ForegroundColor DarkGray

    Write-Host ""
    Write-Host "  Cleaning artifacts..." -ForegroundColor DarkGray
    & $cmdAfter.Source update-cleanup
}

Write-Host ""
exit 0
`
)

// Backup file extension glob.
const OldBackupGlob = "*.old"

// PowerShell execution arguments.
const (
	PSBin            = "powershell"
	PSExecPolicy     = "-ExecutionPolicy"
	PSBypass         = "Bypass"
	PSNoProfile      = "-NoProfile"
	PSNoLogo         = "-NoLogo"
	PSFile           = "-File"
	PSNonInteractive = "-NonInteractive"
	PSCommand        = "-Command"
)
