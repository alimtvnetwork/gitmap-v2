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
// It creates a handoff copy of the active binary and re-launches a hidden
// worker command from that copy. The parent exits immediately so file locks
// are released before deploy overwrites gitmap.exe.
func runUpdate() {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrNoRepoPath)
		os.Exit(1)
	}

	verboseMode := hasFlag("--verbose")

	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable: %v\n", err)
		os.Exit(1)
	}

	copyPath := filepath.Join(filepath.Dir(selfPath), fmt.Sprintf("gitmap-update-%d.exe", os.Getpid()))
	if err := copyFile(selfPath, copyPath); err != nil {
		fallbackPath := filepath.Join(os.TempDir(), fmt.Sprintf("gitmap-update-%d.exe", os.Getpid()))
		if errFallback := copyFile(selfPath, fallbackPath); errFallback != nil {
			fmt.Fprintf(os.Stderr, "Error creating update copy: %v\n", err)
			os.Exit(1)
		}
		copyPath = fallbackPath
	}

	fmt.Printf("  → Active: %s\n  → Handoff: %s\n", selfPath, copyPath)

	copyArgs := []string{constants.CmdUpdateRunner}
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

	// Parent must exit immediately to release file lock on active binary.
	os.Exit(0)
}

// runUpdateRunner is a hidden command that performs the real update work.
func runUpdateRunner() {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrNoRepoPath)
		os.Exit(1)
	}

	verboseMode := hasFlag("--verbose")
	if verboseMode {
		log, err := verbose.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create verbose log: %v\n", err)
		} else {
			defer log.Close()
			log.Log("update-runner starting, repo=%s", repoPath)
		}
	}

	fmt.Printf(constants.MsgUpdateStarting)
	fmt.Printf(constants.MsgUpdateRepoPath, repoPath)
	executeUpdate(repoPath)
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

// cleanupTempCopies removes leftover handoff binaries from previous updates.
func cleanupTempCopies() int {
	selfPath, _ := os.Executable()
	patterns := []string{
		filepath.Join(os.TempDir(), "gitmap-update-*.exe"),
	}
	if selfPath != "" {
		patterns = append(patterns, filepath.Join(filepath.Dir(selfPath), "gitmap-update-*.exe"))
	}

	seen := map[string]bool{}
	cleaned := 0
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			cleanPath := filepath.Clean(match)
			if seen[cleanPath] {
				continue
			}
			seen[cleanPath] = true

			if selfPath != "" && cleanPath == filepath.Clean(selfPath) {
				continue
			}
			if os.Remove(match) == nil {
				fmt.Printf("  → Removed temp copy: %s\n", filepath.Base(match))
				cleaned++
			}
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
	scriptPath, err := writeUpdateScript(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
		os.Exit(1)
	}
	defer os.Remove(scriptPath)

	log := verbose.Get()
	if log != nil {
		log.Log("executing update script: %s", scriptPath)
	}

	runUpdateScript(scriptPath)
}

// writeUpdateScript creates a temporary PowerShell script for self-update.
// Writes with UTF-8 BOM so PowerShell correctly handles Unicode characters.
func writeUpdateScript(repoPath string) (string, error) {
	runPS1 := filepath.Join(repoPath, "run.ps1")
	script := buildUpdateScript(repoPath, runPS1)

	tmpFile, err := os.CreateTemp(os.TempDir(), "gitmap-update-*.ps1")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := tmpFile.Write(bom); err != nil {
		return "", err
	}
	if _, err := tmpFile.WriteString(script); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// buildUpdateScript generates the PowerShell script content.
// It delegates update execution to run.ps1 -Update and verifies active PATH sync.
func buildUpdateScript(repoPath, runPS1 string) string {
	return fmt.Sprintf(`# gitmap self-update script (auto-generated)
Set-Location "%s"

$configPath = Join-Path "%s" "gitmap\powershell.json"
$deployedBinary = $null
if (Test-Path $configPath) {
    $cfg = Get-Content $configPath | ConvertFrom-Json
    if ($cfg.deployPath) {
        $deployedBinary = Join-Path $cfg.deployPath "gitmap\gitmap.exe"
    }
}

$activeBinary = $null
$activeBefore = "unknown"
$cmdBefore = Get-Command gitmap -ErrorAction SilentlyContinue
if ($cmdBefore -and (Test-Path $cmdBefore.Source)) {
    $activeBinary = $cmdBefore.Source
    $activeBefore = & $activeBinary version 2>&1
}

Write-Host ""
Write-Host "  Starting update via run.ps1 -Update" -ForegroundColor Cyan
& "%s" -Update
$runExit = $LASTEXITCODE
if (($runExit -ne 0) -and ($runExit -ne $null)) {
    exit $runExit
}

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

Write-Host ""
Write-Host "  Version before:   $activeBefore" -ForegroundColor DarkGray
Write-Host "  Version active:   $activeAfter" -ForegroundColor DarkGray
Write-Host "  Version deployed: $deployedAfter" -ForegroundColor DarkGray

if (($activeAfter -eq "unknown") -or ($deployedAfter -eq "unknown") -or ($activeAfter -ne $deployedAfter)) {
    Write-Host "  [FAIL] Active PATH version does not match deployed version." -ForegroundColor Red
    exit 1
}

Write-Host "  [OK] Active PATH binary matches deployed version." -ForegroundColor Green

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
`, repoPath, repoPath, runPS1)
}

// runUpdateScript executes the PowerShell script with output piped to terminal.
func runUpdateScript(scriptPath string) {
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass",
		"-NoProfile", "-NoLogo", "-File", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
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
