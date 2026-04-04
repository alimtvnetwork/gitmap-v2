package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
)

// installTool dispatches to the platform-specific installer.
func installTool(opts installOptions) {
	manager := resolvePackageManager(opts.Manager)
	installCmd := buildInstallCommand(manager, opts.Tool, opts.Version)

	if opts.DryRun {
		fmt.Printf(constants.MsgInstallDryCmd, strings.Join(installCmd, " "))

		return
	}

	fmt.Printf(constants.MsgInstallInstalling, opts.Tool)
	runInstallCommand(installCmd, opts.Verbose)
	verifyInstallation(opts.Tool)
}

// buildInstallCommand builds the install command for a given manager and tool.
func buildInstallCommand(manager, tool, version string) []string {
	pkgName := resolvePackageName(manager, tool)

	if manager == constants.PkgMgrChocolatey {
		return buildChocoCommand(pkgName, version)
	}
	if manager == constants.PkgMgrWinget {
		return buildWingetCommand(pkgName, version)
	}
	if manager == constants.PkgMgrApt {
		return buildAptCommand(pkgName, version)
	}
	if manager == constants.PkgMgrBrew {
		return buildBrewCommand(tool, pkgName)
	}

	return buildChocoCommand(pkgName, version)
}

// buildChocoCommand builds a Chocolatey install command.
func buildChocoCommand(pkg, version string) []string {
	args := []string{"choco", "install", pkg, "-y"}

	if version != "" {
		args = append(args, "--version", version)
	}

	return args
}

// buildWingetCommand builds a Winget install command.
func buildWingetCommand(pkg, version string) []string {
	args := []string{"winget", "install", pkg, "--accept-package-agreements", "--accept-source-agreements"}

	if version != "" {
		args = append(args, "--version", version)
	}

	return args
}

// buildAptCommand builds an apt install command.
func buildAptCommand(pkg, version string) []string {
	target := pkg

	if version != "" {
		target = pkg + "=" + version
	}

	return []string{"sudo", "apt", "install", "-y", target}
}

// buildBrewCommand builds a Homebrew install command.
func buildBrewCommand(tool, pkg string) []string {
	if isBrewCaskTool(tool) {
		return []string{"brew", "install", "--cask", pkg}
	}

	return []string{"brew", "install", pkg}
}

// isBrewCaskTool returns true for GUI apps needing --cask.
func isBrewCaskTool(tool string) bool {
	if tool == constants.ToolVSCode {
		return true
	}
	if tool == constants.ToolGitHubDesktop {
		return true
	}
	if tool == constants.ToolPowerShell {
		return true
	}

	return false
}

// runInstallCommand executes the install command.
func runInstallCommand(args []string, verbose bool) {
	cmd := exec.Command(args[0], args[1:]...)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrInstallFailed, args[1], err)
		os.Exit(1)
	}
}

// resolvePackageName maps tool name to package ID for a manager.
func resolvePackageName(manager, tool string) string {
	if manager == constants.PkgMgrWinget {
		return resolveWingetPackage(tool)
	}

	return resolveChocoPackage(tool)
}

// resolveChocoPackage maps tool names to Chocolatey package IDs.
func resolveChocoPackage(tool string) string {
	chocoMap := map[string]string{
		constants.ToolVSCode:        constants.ChocoPkgVSCode,
		constants.ToolNodeJS:        constants.ChocoPkgNodeJS,
		constants.ToolYarn:          constants.ChocoPkgYarn,
		constants.ToolBun:           constants.ChocoPkgBun,
		constants.ToolPnpm:          constants.ChocoPkgPnpm,
		constants.ToolPython:        constants.ChocoPkgPython,
		constants.ToolGo:            constants.ChocoPkgGo,
		constants.ToolGit:           constants.ChocoPkgGit,
		constants.ToolGitLFS:        constants.ChocoPkgGitLFS,
		constants.ToolGHCLI:         constants.ChocoPkgGHCLI,
		constants.ToolGitHubDesktop: constants.ChocoPkgGitHubDesktop,
		constants.ToolCPP:           constants.ChocoPkgCPP,
		constants.ToolPHP:           constants.ChocoPkgPHP,
	}

	pkg, exists := chocoMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveWingetPackage maps tool names to Winget package IDs.
func resolveWingetPackage(tool string) string {
	wingetMap := map[string]string{
		constants.ToolVSCode:    constants.WingetPkgVSCode,
		constants.ToolPowerShell: constants.WingetPkgPowerShell,
	}

	pkg, exists := wingetMap[tool]
	if exists {
		return pkg
	}

	return tool
}