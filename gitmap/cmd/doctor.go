// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/user/gitmap/constants"
)

// runDoctor handles the 'doctor' command.
func runDoctor() {
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
		fmt.Printf("  Found %d issue(s). See recommendations above.\n\n", issues)
	} else {
		fmt.Println("  All checks passed.\n")
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
		printFix("Run: gitmap update")
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
			printFix("Copy deployed binary to PATH location:")
			printFix(fmt.Sprintf("  Copy-Item \"%s\" \"%s\" -Force", absDeployed, absActive))
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

func colorGreen() string {
	if runtime.GOOS == "windows" {
		return constants.ColorGreen
	}
	return constants.ColorGreen
}

func colorRed() string    { return "\033[31m" }
func colorCyan() string   { return constants.ColorCyan }
func colorYellow() string { return constants.ColorYellow }
func colorReset() string  { return constants.ColorReset }
