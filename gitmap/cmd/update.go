package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/verbose"
)

// runUpdate handles the "update" subcommand.
// It copies the current binary to a temp file and re-launches update from
// that copy. The parent exits immediately so the original binary lock is
// released before deploy overwrites gitmap.exe.
func runUpdate() {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrNoRepoPath)
		os.Exit(1)
	}

	// Check for --verbose flag anywhere in remaining args.
	verboseMode := hasFlag("--verbose")

	// If we're already running from the update copy, do the actual update.
	if hasFlag("--from-copy") {
		if verboseMode {
			log, err := verbose.Init()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not create verbose log: %v\n", err)
			} else {
				defer log.Close()
				log.Log("update --from-copy starting, repo=%s", repoPath)
			}
		}
		fmt.Printf(constants.MsgUpdateStarting)
		fmt.Printf(constants.MsgUpdateRepoPath, repoPath)
		executeUpdate(repoPath)
		return
	}

	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable: %v\n", err)
		os.Exit(1)
	}

	copyPath := filepath.Join(os.TempDir(), fmt.Sprintf("gitmap-update-%d.exe", os.Getpid()))
	if err := copyFile(selfPath, copyPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating update copy: %v\n", err)
		os.Exit(1)
	}

	// Re-launch from the copy and immediately exit parent to release lock.
	copyArgs := []string{"update", "--from-copy"}
	if verboseMode {
		copyArgs = append(copyArgs, "--verbose")
	}
	cmd := exec.Command(copyPath, copyArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
		os.Exit(1)
	}
	os.Exit(0)
}

// runUpdateCleanup handles the "update-cleanup" subcommand.
// Removes leftover temp binaries from %TEMP% and .old backup files
// from the deploy directory.
func runUpdateCleanup() {
	fmt.Println("\n  Cleaning up update artifacts...")

	tempCleaned := cleanupTempCopies()
	oldCleaned := cleanupOldBackups()

	total := tempCleaned + oldCleaned
	if total > 0 {
		fmt.Printf("  ✓ Removed %d file(s)\n\n", total)
	} else {
		fmt.Println("  ✓ Nothing to clean up\n")
	}
}

// cleanupTempCopies removes leftover temp binaries from previous updates.
func cleanupTempCopies() int {
	pattern := filepath.Join(os.TempDir(), "gitmap-update-*.exe")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0
	}

	selfPath, _ := os.Executable()
	cleaned := 0
	for _, match := range matches {
		// Don't delete ourselves if we're running as the update copy.
		if selfPath != "" && filepath.Clean(match) == filepath.Clean(selfPath) {
			continue
		}
		if os.Remove(match) == nil {
			fmt.Printf("  → Removed temp copy: %s\n", filepath.Base(match))
			cleaned++
		}
	}

	return cleaned
}

// cleanupOldBackups removes .old backup binaries from the deploy directory.
func cleanupOldBackups() int {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		return 0
	}

	// Try to find deploy path from powershell.json
	configPath := filepath.Join(repoPath, "gitmap", "powershell.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return 0
	}

	// Simple extraction of deployPath from JSON
	deployPath := extractJSONString(data, "deployPath")
	if len(deployPath) == 0 {
		return 0
	}

	appDir := filepath.Join(deployPath, "gitmap")
	pattern := filepath.Join(appDir, "*.old")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0
	}

	cleaned := 0
	for _, match := range matches {
		if os.Remove(match) == nil {
			fmt.Printf("  → Removed backup: %s\n", filepath.Base(match))
			cleaned++
		}
	}

	return cleaned
}

// extractJSONString extracts a string value from JSON bytes by key.
// Simple parser to avoid importing encoding/json in cmd package.
func extractJSONString(data []byte, key string) string {
	s := string(data)
	needle := fmt.Sprintf(`"%s"`, key)
	idx := 0
	for {
		i := indexOf(s[idx:], needle)
		if i < 0 {
			return ""
		}
		idx += i + len(needle)
		// Skip whitespace and colon
		for idx < len(s) && (s[idx] == ' ' || s[idx] == ':' || s[idx] == '\t') {
			idx++
		}
		if idx < len(s) && s[idx] == '"' {
			end := indexOf(s[idx+1:], `"`)
			if end >= 0 {
				return s[idx+1 : idx+1+end]
			}
		}
	}
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}

// hasFlag checks if a flag is present in os.Args[2:].
func hasFlag(name string) bool {
	for _, arg := range os.Args[2:] {
		if arg == name {
			return true
		}
	}

	return false
}

// copyFile copies src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// executeUpdate writes a temp PS1 script and runs it.
func executeUpdate(repoPath string) {
	scriptPath := writeUpdateScript(repoPath)
	defer os.Remove(scriptPath)

	log := verbose.Get()
	if log != nil {
		log.Log("executing update script: %s", scriptPath)
	}

	runUpdateScript(scriptPath)
}

// writeUpdateScript creates a temporary PowerShell script for self-update.
// Writes with UTF-8 BOM so PowerShell correctly handles Unicode characters.
func writeUpdateScript(repoPath string) string {
	runPS1 := filepath.Join(repoPath, "run.ps1")
	script := buildUpdateScript(repoPath, runPS1)
	tmpFile := filepath.Join(os.TempDir(), "gitmap-update.ps1")

	bom := []byte{0xEF, 0xBB, 0xBF}
	content := append(bom, []byte(script)...)
	os.WriteFile(tmpFile, content, constants.DirPermission)

	return tmpFile
}

// buildUpdateScript generates the PowerShell script content.
// Ensures active PATH binary is synced and prints changelog bullets.
func buildUpdateScript(repoPath, runPS1 string) string {
	return fmt.Sprintf(`# gitmap self-update script (auto-generated)
Set-Location "%s"

# Detect deployed binary from powershell.json
$deployedBinary = $null
$oldVersion = "unknown"
$configPath = Join-Path "%s" "gitmap\powershell.json"
if (Test-Path $configPath) {
    $cfg = Get-Content $configPath | ConvertFrom-Json
    if ($cfg.deployPath) {
        $deployedBinary = Join-Path $cfg.deployPath "gitmap\gitmap.exe"
        if (Test-Path $deployedBinary) {
            $oldVersion = & $deployedBinary version 2>&1
        }
    }
}

# Detect source version from constants.go
$sourceVersion = "unknown"
$constantsPath = Join-Path "%s" "gitmap\constants\constants.go"
if (Test-Path $constantsPath) {
    $match = Select-String -Path $constantsPath -Pattern 'const Version = "([^"]+)"' | Select-Object -First 1
    if ($match -and $match.Matches.Count -gt 0) {
        $sourceVersion = "gitmap v" + $match.Matches[0].Groups[1].Value
    }
}

# Detect active gitmap on PATH (what the user invokes)
$activeBinary = $null
$activeVersion = "unknown"
$cmd = Get-Command gitmap -ErrorAction SilentlyContinue
if ($cmd) {
    $activeBinary = $cmd.Source
    if (Test-Path $activeBinary) {
        $activeVersion = & $activeBinary version 2>&1
    }
}

Write-Host ""
Write-Host "  Current deployed version: $oldVersion" -ForegroundColor DarkGray
Write-Host "  Current source version:   $sourceVersion" -ForegroundColor DarkGray
if ($activeBinary) {
    Write-Host "  Current PATH binary: $activeBinary" -ForegroundColor DarkGray
    Write-Host "  Current PATH version: $activeVersion" -ForegroundColor DarkGray
}

# Check if there are changes to pull
$prevPref = $ErrorActionPreference
$ErrorActionPreference = "Continue"
$pullOutput = git pull 2>&1
$ErrorActionPreference = $prevPref
$pullText = ($pullOutput | Out-String).Trim()

if ($pullText -match "Already up to date") {
    $needsRebuild = $false
    if ($sourceVersion -ne "unknown" -and $oldVersion -ne "unknown" -and $oldVersion -ne $sourceVersion) {
        $needsRebuild = $true
    }

    if ($needsRebuild -eq $false) {
        if ($activeBinary -and $deployedBinary -and (Test-Path $activeBinary) -and (Test-Path $deployedBinary)) {
            $activeResolved = (Resolve-Path $activeBinary).Path
            $deployedResolved = (Resolve-Path $deployedBinary).Path
            if ($activeResolved -ne $deployedResolved) {
                Write-Host "  [WARN] PATH points to a different gitmap binary." -ForegroundColor Yellow
                Write-Host "         Active:   $activeResolved" -ForegroundColor Yellow
                Write-Host "         Deployed: $deployedResolved" -ForegroundColor Yellow
                try {
                    Copy-Item $deployedBinary $activeBinary -Force -ErrorAction Stop
                    Write-Host "  [OK] Synced active PATH binary with deployed build." -ForegroundColor Green
                    $activeVersion = & $activeBinary version 2>&1
                    Write-Host "  Active PATH version is now: $activeVersion" -ForegroundColor Green
                } catch {
                    Write-Host "  [WARN] Could not sync active PATH binary: $_" -ForegroundColor Yellow
                }
            }
        }

        Write-Host "  Source is already up to date." -ForegroundColor Yellow
        Write-Host ""
        Write-Host "  No update needed - you are running the latest source." -ForegroundColor Green

        $changelogBinary = $activeBinary
        if ((-not $changelogBinary) -or (-not (Test-Path $changelogBinary))) {
            $changelogBinary = $deployedBinary
        }
        if ($changelogBinary -and (Test-Path $changelogBinary)) {
            Write-Host ""
            Write-Host "  Latest changelog:" -ForegroundColor Cyan
            & $changelogBinary changelog --latest
        }

        Write-Host ""
        exit 0
    }

    Write-Host "  Source is up to date, but binary version mismatch detected. Rebuilding..." -ForegroundColor Yellow
}

Write-Host "  Changes detected, rebuilding..." -ForegroundColor Cyan
Write-Host ""
Start-Sleep -Milliseconds 1200
& "%s" -NoPull
Write-Host ""

# Compare deployed versions before/after update
$newVersion = "unknown"
if ($deployedBinary -and (Test-Path $deployedBinary)) {
    $newVersion = & $deployedBinary version 2>&1
}

if ($oldVersion -eq $newVersion) {
    Write-Host "  [WARN] Version unchanged after update ($newVersion)" -ForegroundColor Yellow
    Write-Host "  The source changed but the version constant may not have been bumped." -ForegroundColor Yellow
} else {
    Write-Host "  [OK] Updated: $oldVersion -> $newVersion" -ForegroundColor Green
}

if ($deployedBinary) {
    Write-Host "  Deployed binary: $deployedBinary" -ForegroundColor DarkGray
}

# Sync active PATH binary when it differs
if ($activeBinary -and $deployedBinary -and (Test-Path $activeBinary) -and (Test-Path $deployedBinary)) {
    $activeResolved = (Resolve-Path $activeBinary).Path
    $deployedResolved = (Resolve-Path $deployedBinary).Path
    if ($activeResolved -ne $deployedResolved) {
        Write-Host "  [WARN] PATH points to a different gitmap binary." -ForegroundColor Yellow
        Write-Host "         Active:   $activeResolved" -ForegroundColor Yellow
        Write-Host "         Deployed: $deployedResolved" -ForegroundColor Yellow
        try {
            Copy-Item $deployedBinary $activeBinary -Force -ErrorAction Stop
            Write-Host "  [OK] Synced active PATH binary with deployed build." -ForegroundColor Green
            $activeVersion = & $activeBinary version 2>&1
            Write-Host "  Active PATH version is now: $activeVersion" -ForegroundColor Green
        } catch {
            Write-Host "  [WARN] Could not sync active PATH binary: $_" -ForegroundColor Yellow
        }
    }
}

# Show changelog bullets for updated version
$changelogBinary = $activeBinary
if ((-not $changelogBinary) -or (-not (Test-Path $changelogBinary))) {
    $changelogBinary = $deployedBinary
}
if ($changelogBinary -and (Test-Path $changelogBinary)) {
    Write-Host ""
    Write-Host "  Changelog:" -ForegroundColor Cyan
    if ($newVersion -and $newVersion -ne "unknown") {
        & $changelogBinary changelog $newVersion
    } else {
        & $changelogBinary changelog --latest
    }
}

# Run update-cleanup to remove temp copies and .old backups
$cleanupBinary = Join-Path "%s" "bin\gitmap.exe"
if ($deployedBinary -and (Test-Path $deployedBinary)) {
    $cleanupBinary = $deployedBinary
}

if (Test-Path $cleanupBinary) {
    Write-Host ""
    Write-Host "  Cleaning up update artifacts..." -ForegroundColor DarkGray
    & $cleanupBinary update-cleanup
}

Write-Host ""
exit 0
`, repoPath, repoPath, repoPath, runPS1, repoPath)
}

// runUpdateScript executes the PowerShell script with output piped to terminal.
func runUpdateScript(scriptPath string) {
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass",
		"-NoProfile", "-NoLogo", "-NonInteractive", "-File", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	log := verbose.Get()
	if log != nil {
		log.Log("update script exited: err=%v", err)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
		os.Exit(1)
	}
}
