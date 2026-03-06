package cmd

import (
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

// runFixPath syncs the active PATH binary from the deployed binary.
func runFixPath() {
	fmt.Println()
	fmt.Printf(constants.DoctorFixBannerFmt, constants.Version)
	fmt.Println(constants.DoctorBannerRule)

	absActive, activeVersion := resolveActiveBinary()
	if len(absActive) == 0 {
		return
	}

	absDeployed, deployedVersion := resolveDeployedForSync()
	if len(absDeployed) == 0 {
		return
	}

	printFixPathInfo(absActive, activeVersion, absDeployed, deployedVersion)
	syncBinaries(absActive, activeVersion, absDeployed, deployedVersion)
}

// resolveActiveBinary finds and validates the active PATH binary.
func resolveActiveBinary() (string, string) {
	activePath, activeErr := exec.LookPath("gitmap")
	if activeErr != nil {
		printIssue(constants.DoctorNotOnPath, constants.DoctorNoSync)
		printFix(constants.DoctorAddPathFix)

		return "", ""
	}

	absActive, _ := filepath.Abs(activePath)

	return absActive, getBinaryVersion(absActive)
}

// resolveDeployedForSync finds the deployed binary for syncing.
func resolveDeployedForSync() (string, string) {
	deployedPath, deployedErr := resolveDeployedBinary()
	if deployedErr != nil {
		printIssue(constants.DoctorCannotResolve, deployedErr.Error())

		return "", ""
	}

	absDeployed, _ := filepath.Abs(deployedPath)

	return absDeployed, getBinaryVersion(absDeployed)
}

// printFixPathInfo displays the active and deployed binary paths.
func printFixPathInfo(absActive, activeVersion, absDeployed, deployedVersion string) {
	fmt.Printf(constants.DoctorActivePathFmt, absActive, activeVersion)
	fmt.Printf(constants.DoctorDeployedFmt, absDeployed, deployedVersion)
}

// syncBinaries orchestrates the 3-layer sync strategy.
func syncBinaries(absActive, activeVersion, absDeployed, deployedVersion string) {
	if absActive == absDeployed {
		printOK(constants.DoctorAlreadySynced)

		return
	}

	if activeVersion == deployedVersion {
		printOK(constants.DoctorVersionsMatch, activeVersion)

		return
	}

	fmt.Println()
	fmt.Printf(constants.DoctorSyncingFmt, absDeployed, absActive)
	attemptSync(absDeployed, absActive, deployedVersion)
}

// attemptSync tries copy, rename fallback, and kill strategies.
func attemptSync(absDeployed, absActive, deployedVersion string) {
	if tryCopyWithRetry(absDeployed, absActive, 20, 500*time.Millisecond) {
		verifySync(absActive, deployedVersion)

		return
	}

	if tryRenameFallback(absDeployed, absActive) {
		verifySync(absActive, deployedVersion)

		return
	}

	if tryKillAndCopy(absDeployed, absActive) {
		verifySync(absActive, deployedVersion)

		return
	}

	printSyncFailure(absDeployed, absActive)
}

// printSyncFailure reports that all sync strategies failed.
func printSyncFailure(absDeployed, absActive string) {
	fmt.Println()
	printIssue(constants.DoctorSyncFailTitle, constants.DoctorSyncFailDetail)
	printFix(constants.DoctorSyncFailFix1)
	printFix(fmt.Sprintf(constants.DoctorSyncFailFix2Fmt, absDeployed, absActive))
}

// resolveDeployedBinary finds the deployed binary path from powershell.json.
func resolveDeployedBinary() (string, error) {
	if len(constants.RepoPath) == 0 {
		return "", fmt.Errorf(constants.DoctorResolveNoRepo)
	}

	configPath := filepath.Join(constants.RepoPath, "gitmap", "powershell.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf(constants.DoctorResolveNoRead, err)
	}

	return resolveDeployedPath(data)
}

// resolveDeployedPath extracts and validates the deployed path from config data.
func resolveDeployedPath(data []byte) (string, error) {
	deployPath := extractJSONString(data, "deployPath")
	if len(deployPath) == 0 {
		return "", fmt.Errorf(constants.DoctorResolveNoDeploy)
	}

	binaryName := extractJSONString(data, "binaryName")
	if len(binaryName) == 0 {
		binaryName = constants.DoctorDefaultBinary
	}

	deployed := filepath.Join(deployPath, "gitmap", binaryName)
	if _, err := os.Stat(deployed); err != nil {
		return "", fmt.Errorf(constants.DoctorResolveNotFound, deployed)
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
			fmt.Printf(constants.DoctorRetryFmt, attempt, maxAttempts)
			time.Sleep(delay)
		}
	}

	return false
}

// tryRenameFallback renames the locked target to .old, then copies.
func tryRenameFallback(src, dst string) bool {
	backup := dst + ".old"
	os.Remove(backup)

	err := os.Rename(dst, backup)
	if err != nil {
		return false
	}

	fmt.Println(constants.DoctorRenamedMsg)
	err = copyFileOverwrite(src, dst)
	if err != nil {
		os.Rename(backup, dst)

		return false
	}

	return true
}

// tryKillAndCopy finds stale gitmap processes and terminates them.
func tryKillAndCopy(src, dst string) bool {
	if runtime.GOOS == constants.OSWindows {
		return tryKillWindows(src, dst)
	}

	return false
}

// tryKillWindows kills stale gitmap processes on Windows and retries copy.
func tryKillWindows(src, dst string) bool {
	fmt.Println(constants.DoctorKillingMsg)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
		fmt.Sprintf(`Get-CimInstance Win32_Process -Filter "Name='gitmap.exe'" | `+
			`Where-Object { $_.ExecutablePath -and (Resolve-Path $_.ExecutablePath -ErrorAction SilentlyContinue).Path -eq '%s' -and $_.ProcessId -ne %d } | `+
			`ForEach-Object { Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue; $_.ProcessId }`,
			dst, os.Getpid()))

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	reportKilledProcesses(string(out))

	return copyFileOverwrite(src, dst) == nil
}

// reportKilledProcesses logs which processes were stopped.
func reportKilledProcesses(output string) {
	killed := strings.TrimSpace(output)
	if len(killed) > 0 {
		fmt.Printf(constants.DoctorKilledFmt, killed)
		time.Sleep(500 * time.Millisecond)
	}
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
		printOK(constants.DoctorOKPathFmt, actualVersion)

		return
	}

	printWarn(fmt.Sprintf(constants.DoctorWarnSyncFmt, actualVersion, expectedVersion))
}
