package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
)

// verifyInstallation confirms a tool is accessible after install.
func verifyInstallation(tool string) {
	fmt.Printf(constants.MsgInstallVerifying, tool)
	binary := toolBinaryName(tool)

	version := getInstalledVersion(binary)
	if version == "" {
		fmt.Fprintf(os.Stderr, constants.ErrInstallVerifyFailed, tool)

		return
	}

	fmt.Printf(constants.MsgInstallSuccess, tool)
	runPostInstall(tool)
}

// detectInstalledVersion checks if a tool is already installed.
func detectInstalledVersion(tool string) string {
	binary := toolBinaryName(tool)

	return getInstalledVersion(binary)
}

// getInstalledVersion runs --version and returns the output.
func getInstalledVersion(binary string) string {
	path, err := exec.LookPath(binary)
	if err != nil {
		return ""
	}

	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// toolBinaryName maps tool names to their binary/executable names.
func toolBinaryName(tool string) string {
	binaryMap := map[string]string{
		constants.ToolVSCode:        "code",
		constants.ToolNodeJS:        "node",
		constants.ToolYarn:          "yarn",
		constants.ToolBun:           "bun",
		constants.ToolPnpm:          "pnpm",
		constants.ToolPython:        "python3",
		constants.ToolGo:            "go",
		constants.ToolGit:           "git",
		constants.ToolGitLFS:        "git-lfs",
		constants.ToolGHCLI:         "gh",
		constants.ToolGitHubDesktop: "github-desktop",
		constants.ToolCPP:           "g++",
		constants.ToolPHP:           "php",
		constants.ToolPowerShell:    "pwsh",
	}

	binary, exists := binaryMap[tool]
	if exists {
		return binary
	}

	return tool
}

// runPostInstall executes post-install actions for specific tools.
func runPostInstall(tool string) {
	if tool == constants.ToolGitLFS {
		runPostInstallGitLFS()

		return
	}
	if tool == constants.ToolGit {
		runPostInstallGit()

		return
	}
}

// runPostInstallGitLFS runs git lfs install.
func runPostInstallGitLFS() {
	cmd := exec.Command("git", "lfs", "install")
	_ = cmd.Run()
}

// runPostInstallGit configures git longpaths.
func runPostInstallGit() {
	cmd := exec.Command("git", "config", "--global", "core.longpaths", "true")
	_ = cmd.Run()
}