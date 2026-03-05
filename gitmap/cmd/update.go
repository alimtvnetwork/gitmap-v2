package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/user/gitmap/constants"
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

	// If we're already running from the update copy, do the actual update.
	if len(os.Args) > 2 && os.Args[2] == "--from-copy" {
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
	cmd := exec.Command(copyPath, "update", "--from-copy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
		os.Exit(1)
	}
	os.Exit(0)
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
func buildUpdateScript(repoPath, runPS1 string) string {
	return fmt.Sprintf(`# gitmap self-update script (auto-generated)
Set-Location "%s"
Write-Host ""
Write-Host "  Pulling and rebuilding gitmap..." -ForegroundColor Cyan
Write-Host ""
& "%s"
Write-Host ""
$newBinary = Join-Path "%s" "bin\gitmap.exe"
if (Test-Path $newBinary) {
    $version = & $newBinary help 2>&1 | Select-String -Pattern "v\d+\.\d+\.\d+" | ForEach-Object { $_.Matches[0].Value }
    Write-Host "  [OK] Updated to gitmap $version" -ForegroundColor Green
} else {
    Write-Host "  [OK] Update complete" -ForegroundColor Green
}
Write-Host ""
`, repoPath, runPS1, repoPath)
}

// runUpdateScript executes the PowerShell script with output piped to terminal.
func runUpdateScript(scriptPath string) {
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass",
		"-NoProfile", "-File", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
		os.Exit(1)
	}
}
