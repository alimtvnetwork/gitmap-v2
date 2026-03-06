// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/user/gitmap/constants"
)

// runDoctor handles the 'doctor' command.
func runDoctor() {
	fixPath := parseDoctorFlags(os.Args[2:])

	if fixPath {
		runFixPath()
		return
	}

	fmt.Printf("\n  gitmap doctor (v%s)\n", constants.Version)
	fmt.Println("  ──────────────────────────────────────────")

	issues := 0

	issues += checkRepoPath()
	issues += checkActiveBinary()
	issues += checkDeployedBinary()
	issues += checkVersionMismatch()
	issues += checkGit()
	issues += checkGo()
	issues += checkChangelogFile()

	fmt.Println()
	if issues > 0 {
		fmt.Printf("  Found %d issue(s). See recommendations above.\n", issues)
		fmt.Printf("  Tip: run 'gitmap doctor --fix-path' to auto-sync the PATH binary.\n\n")
	} else {
		fmt.Println("  All checks passed.\n")
	}
}

// parseDoctorFlags parses flags for the doctor command.
func parseDoctorFlags(args []string) (fixPath bool) {
	fs := flag.NewFlagSet(constants.CmdDoctor, flag.ExitOnError)
	fixPathFlag := fs.Bool("fix-path", false, "Sync the active PATH binary from the deployed binary")
	fs.Parse(args)

	return *fixPathFlag
}

// runFixPath syncs the active PATH binary from the deployed binary
// using retry, rename fallback, and stale-process termination.
func runFixPath() {
	fmt.Println()
	fmt.Printf("  gitmap doctor --fix-path (v%s)\n", constants.Version)
	fmt.Println("  ──────────────────────────────────────────")

	activePath, activeErr := exec.LookPath("gitmap")
	if activeErr != nil {
		printIssue("gitmap not found on PATH", "Cannot sync — no active binary to replace.")
		printFix("Add your deploy directory to PATH first.")
		return
	}

	absActive, _ := filepath.Abs(activePath)
	activeVersion := getBinaryVersion(absActive)

	deployedPath, deployedErr := resolveDeployedBinary()
	if deployedErr != nil {
		printIssue("Cannot resolve deployed binary", deployedErr.Error())
		return
	}

	absDeployed, _ := filepath.Abs(deployedPath)
	deployedVersion := getBinaryVersion(absDeployed)

	fmt.Printf("  Active PATH:  %s (%s)\n", absActive, activeVersion)
	fmt.Printf("  Deployed:     %s (%s)\n", absDeployed, deployedVersion)

	if absActive == absDeployed {
		printOK("PATH already points to the deployed binary. Nothing to sync.")
		return
	}

	if activeVersion == deployedVersion {
		printOK("Versions already match (%s). No sync needed.", activeVersion)
		return
	}

	fmt.Println()
	fmt.Printf("  Syncing %s -> %s...\n", absDeployed, absActive)

	// Layer 1: Direct copy with retries.
	if tryCopyWithRetry(absDeployed, absActive, 20, 500*time.Millisecond) {
		verifySync(absActive, deployedVersion)
		return
	}

	// Layer 2: Rename fallback (Windows allows rename of locked .exe).
	if tryRenameFallback(absDeployed, absActive) {
		verifySync(absActive, deployedVersion)
		return
	}

	// Layer 3: Kill stale gitmap processes and retry.
	if tryKillAndCopy(absDeployed, absActive) {
		verifySync(absActive, deployedVersion)
		return
	}

	fmt.Println()
	printIssue("Could not sync PATH binary after all fallback attempts",
		"The file is still locked by another process.")
	printFix("Close all terminals and apps using gitmap, then run:")
	printFix(fmt.Sprintf("  Copy-Item \"%s\" \"%s\" -Force", absDeployed, absActive))
}

// resolveDeployedBinary finds the deployed binary path from powershell.json.
func resolveDeployedBinary() (string, error) {
	if len(constants.RepoPath) == 0 {
		return "", fmt.Errorf("RepoPath not embedded — rebuild with run.ps1")
	}

	configPath := filepath.Join(constants.RepoPath, "gitmap", "powershell.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("cannot read powershell.json: %v", err)
	}

	deployPath := extractJSONString(data, "deployPath")
	if len(deployPath) == 0 {
		return "", fmt.Errorf("no deployPath in powershell.json")
	}

	binaryName := extractJSONString(data, "binaryName")
	if len(binaryName) == 0 {
		binaryName = "gitmap.exe"
	}

	deployed := filepath.Join(deployPath, "gitmap", binaryName)
	if _, err := os.Stat(deployed); err != nil {
		return "", fmt.Errorf("deployed binary not found: %s", deployed)
	}

	return deployed, nil
}

// tryCopyWithRetry attempts to copy src to dst with retries.
func tryCopyWithRetry(src, dst string, maxAttempts int, delay time.Duration) bool {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := copyFileOverwrite(src, dst)
		if err == nil {
			return true
		}
		if attempt < maxAttempts {
			fmt.Printf("  [%d/%d] File in use, retrying...\n", attempt, maxAttempts)
			time.Sleep(delay)
		}
	}

	return false
}

// tryRenameFallback renames the locked target to .old, then copies.
func tryRenameFallback(src, dst string) bool {
	backup := dst + ".old"

	// Clean up previous backup.
	os.Remove(backup)

	err := os.Rename(dst, backup)
	if err != nil {
		return false
	}

	fmt.Println("  Renamed locked binary to .old, copying fresh...")
	err = copyFileOverwrite(src, dst)
	if err != nil {
		// Rollback: restore from backup.
		os.Rename(backup, dst)
		return false
	}

	return true
}

// tryKillAndCopy finds stale gitmap processes on the target path
// and terminates them before retrying the copy.
func tryKillAndCopy(src, dst string) bool {
	if runtime.GOOS != "windows" {
		return false
	}

	fmt.Println("  Attempting to stop stale gitmap processes...")

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
		fmt.Sprintf(`Get-CimInstance Win32_Process -Filter "Name='gitmap.exe'" | `+
			`Where-Object { $_.ExecutablePath -and (Resolve-Path $_.ExecutablePath -ErrorAction SilentlyContinue).Path -eq '%s' -and $_.ProcessId -ne %d } | `+
			`ForEach-Object { Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue; $_.ProcessId }`,
			dst, os.Getpid()))

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	killed := strings.TrimSpace(string(out))
	if len(killed) > 0 {
		fmt.Printf("  Stopped process(es): %s\n", killed)
		time.Sleep(500 * time.Millisecond)
	}

	return copyFileOverwrite(src, dst) == nil
}

// copyFileOverwrite copies src to dst, overwriting dst if it exists.
func copyFileOverwrite(src, dst string) error {
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

// verifySync checks that the synced binary reports the expected version.
func verifySync(path, expectedVersion string) {
	fmt.Println()
	actualVersion := getBinaryVersion(path)
	if actualVersion == expectedVersion {
		printOK("PATH binary synced successfully: %s", actualVersion)
	} else {
		printWarn(fmt.Sprintf("Synced but version mismatch: got %s, expected %s", actualVersion, expectedVersion))
	}
}

// checkRepoPath reports whether RepoPath is embedded.
func checkRepoPath() int {
	if len(constants.RepoPath) == 0 {
		printIssue("RepoPath not embedded", "Binary was not built with run.ps1. Self-update will not work.")
		printFix("Rebuild with: .\\run.ps1")
		return 1
	}

	printOK("RepoPath: %s", constants.RepoPath)
	return 0
}

// checkActiveBinary reports the gitmap binary on PATH.
func checkActiveBinary() int {
	path, err := exec.LookPath("gitmap")
	if err != nil {
		printIssue("gitmap not found on PATH", "The gitmap binary is not accessible from your terminal.")
		printFix("Add your deploy directory to PATH (e.g., E:\\bin-run\\gitmap)")
		return 1
	}

	absPath, _ := filepath.Abs(path)
	version := getBinaryVersion(absPath)
	printOK("PATH binary: %s (%s)", absPath, version)
	return 0
}

// checkDeployedBinary reports the deployed binary from powershell.json.
func checkDeployedBinary() int {
	if len(constants.RepoPath) == 0 {
		return 0
	}

	configPath := filepath.Join(constants.RepoPath, "gitmap", "powershell.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		printIssue("Cannot read powershell.json", "Deploy path detection unavailable.")
		return 1
	}

	deployPath := extractJSONString(data, "deployPath")
	if len(deployPath) == 0 {
		printIssue("No deployPath in powershell.json", "Deploy target not configured.")
		return 1
	}

	binaryName := extractJSONString(data, "binaryName")
	if len(binaryName) == 0 {
		binaryName = "gitmap.exe"
	}

	deployedBinary := filepath.Join(deployPath, "gitmap", binaryName)
	if _, err := os.Stat(deployedBinary); err != nil {
		printIssue("Deployed binary not found", deployedBinary)
		printFix("Run: .\\run.ps1")
		return 1
	}

	version := getBinaryVersion(deployedBinary)
	printOK("Deployed binary: %s (%s)", deployedBinary, version)
	return 0
}

// checkVersionMismatch compares PATH vs deployed vs source versions.
func checkVersionMismatch() int {
	sourceVersion := fmt.Sprintf("gitmap v%s", constants.Version)
	issues := 0

	activePath, activeErr := exec.LookPath("gitmap")
	activeVersion := ""
	if activeErr == nil {
		absPath, _ := filepath.Abs(activePath)
		activeVersion = getBinaryVersion(absPath)
	}

	deployedVersion := ""
	deployedPath := ""
	if len(constants.RepoPath) > 0 {
		configPath := filepath.Join(constants.RepoPath, "gitmap", "powershell.json")
		data, err := os.ReadFile(configPath)
		if err == nil {
			dp := extractJSONString(data, "deployPath")
			bn := extractJSONString(data, "binaryName")
			if len(bn) == 0 {
				bn = "gitmap.exe"
			}
			if len(dp) > 0 {
				deployedPath = filepath.Join(dp, "gitmap", bn)
				deployedVersion = getBinaryVersion(deployedPath)
			}
		}
	}

	if len(activeVersion) > 0 && activeVersion != sourceVersion {
		printIssue("PATH binary version mismatch",
			fmt.Sprintf("PATH: %s, Source: %s", activeVersion, sourceVersion))
		printFix("Run: gitmap update  or  gitmap doctor --fix-path")
		issues++
	}

	if len(deployedVersion) > 0 && deployedVersion != sourceVersion {
		printIssue("Deployed binary version mismatch",
			fmt.Sprintf("Deployed: %s, Source: %s", deployedVersion, sourceVersion))
		printFix("Run: .\\run.ps1 -NoPull")
		issues++
	}

	if len(activeVersion) > 0 && len(deployedVersion) > 0 && activeVersion != deployedVersion {
		absActive, _ := filepath.Abs(activePath)
		absDeployed, _ := filepath.Abs(deployedPath)
		if absActive != absDeployed {
			printIssue("PATH and deployed binaries differ",
				fmt.Sprintf("PATH: %s (%s), Deployed: %s (%s)", absActive, activeVersion, absDeployed, deployedVersion))
			printFix("Run: gitmap doctor --fix-path")
			issues++
		}
	}

	if issues == 0 {
		printOK("Source version: %s (all binaries match)", sourceVersion)
	}

	return issues
}

// checkGit verifies git is available.
func checkGit() int {
	path, err := exec.LookPath(constants.GitBin)
	if err != nil {
		printIssue("git not found on PATH", "Git is required for most gitmap commands.")
		return 1
	}

	cmd := exec.Command(constants.GitBin, "--version")
	out, err := cmd.Output()
	if err != nil {
		printOK("Git: %s (version unknown)", path)
		return 0
	}

	printOK("Git: %s (%s)", path, strings.TrimSpace(string(out)))
	return 0
}

// checkGo verifies Go is available for building.
func checkGo() int {
	path, err := exec.LookPath("go")
	if err != nil {
		printWarn("Go not found on PATH (needed only for building from source)")
		return 0
	}

	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		printOK("Go: %s (version unknown)", path)
		return 0
	}

	printOK("Go: %s", strings.TrimSpace(string(out)))
	return 0
}

// checkChangelogFile verifies CHANGELOG.md exists.
func checkChangelogFile() int {
	if _, err := os.Stat(constants.ChangelogFile); err != nil {
		printWarn("CHANGELOG.md not found (changelog command will not work)")
		return 0
	}

	printOK("CHANGELOG.md present")
	return 0
}

// getBinaryVersion runs a binary with "version" and returns the output.
func getBinaryVersion(path string) string {
	if _, err := os.Stat(path); err != nil {
		return "not found"
	}

	cmd := exec.Command(path, "version")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}

// printOK prints a green check.
func printOK(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s[OK]%s %s\n", colorGreen(), colorReset(), msg)
}

// printIssue prints a red issue.
func printIssue(title, detail string) {
	fmt.Printf("  %s[!!]%s %s\n", colorRed(), colorReset(), title)
	fmt.Printf("       %s\n", detail)
}

// printFix prints a fix recommendation.
func printFix(fix string) {
	fmt.Printf("       %sFix:%s %s\n", colorCyan(), colorReset(), fix)
}

// printWarn prints a yellow warning.
func printWarn(msg string) {
	fmt.Printf("  %s[--]%s %s\n", colorYellow(), colorReset(), msg)
}

func colorGreen() string  { return constants.ColorGreen }
func colorRed() string    { return "\033[31m" }
func colorCyan() string   { return constants.ColorCyan }
func colorYellow() string { return constants.ColorYellow }
func colorReset() string  { return constants.ColorReset }
