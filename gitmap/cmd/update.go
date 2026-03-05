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

// cleanupUpdateCopies removes leftover temp binaries from previous updates.
func cleanupUpdateCopies() {
	pattern := filepath.Join(os.TempDir(), "gitmap-update-*.exe")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	selfPath, _ := os.Executable()
	for _, match := range matches {
		// Don't delete ourselves if we're running as the update copy.
		if selfPath != "" && filepath.Clean(match) == filepath.Clean(selfPath) {
			continue
		}
		os.Remove(match)
	}
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
// Captures current version before update, checks if git pull brings changes,
// and compares versions after update to confirm it applied.
func buildUpdateScript(repoPath, runPS1 string) string {
	return fmt.Sprintf(`# gitmap self-update script (auto-generated)
Set-Location "%s"

# Capture current version before update
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

Write-Host ""
Write-Host "  Current version: $oldVersion" -ForegroundColor DarkGray

# Check if there are changes to pull
$prevPref = $ErrorActionPreference
$ErrorActionPreference = "Continue"
$pullOutput = git pull 2>&1
$ErrorActionPreference = $prevPref
$pullText = ($pullOutput | Out-String).Trim()

if ($pullText -match "Already up to date") {
    Write-Host "  Source is already up to date." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  No update needed — you are running the latest version." -ForegroundColor Green
    Write-Host ""
    exit 0
}

Write-Host "  Changes detected, rebuilding..." -ForegroundColor Cyan
Write-Host ""
Start-Sleep -Milliseconds 1200
& "%s" -NoPull
Write-Host ""

# Compare versions
$newBinary = Join-Path "%s" "bin\gitmap.exe"
if (Test-Path $newBinary) {
    $newVersion = & $newBinary version 2>&1
    if ($oldVersion -eq $newVersion) {
        Write-Host "  [WARN] Version unchanged after update ($newVersion)" -ForegroundColor Yellow
        Write-Host "  The source changed but the version constant was not bumped." -ForegroundColor Yellow
    } else {
        Write-Host "  [OK] Updated: $oldVersion -> $newVersion" -ForegroundColor Green
    }
} else {
    Write-Host "  [OK] Update complete" -ForegroundColor Green
}
Write-Host ""
exit 0
`, repoPath, repoPath, runPS1, repoPath)
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
