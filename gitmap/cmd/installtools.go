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
	recordInstallation(opts.Tool, manager)
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
	if manager == constants.PkgMgrSnap {
		return buildSnapCommand(pkgName)
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

// buildSnapCommand builds a Snap install command.
func buildSnapCommand(pkg string) []string {
	return []string{"sudo", "snap", "install", pkg}
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
	if manager == constants.PkgMgrApt {
		return resolveAptPackage(tool)
	}
	if manager == constants.PkgMgrBrew {
		return resolveBrewPackage(tool)
	}
	if manager == constants.PkgMgrSnap {
		return resolveSnapPackage(tool)
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
		constants.ToolMySQL:         constants.ChocoPkgMySQL,
		constants.ToolMariaDB:       constants.ChocoPkgMariaDB,
		constants.ToolPostgreSQL:    constants.ChocoPkgPostgreSQL,
		constants.ToolSQLite:        constants.ChocoPkgSQLite,
		constants.ToolMongoDB:       constants.ChocoPkgMongoDB,
		constants.ToolCouchDB:       constants.ChocoPkgCouchDB,
		constants.ToolRedis:         constants.ChocoPkgRedis,
		constants.ToolNeo4j:         constants.ChocoPkgNeo4j,
		constants.ToolElasticsearch: constants.ChocoPkgElasticsearch,
		constants.ToolDuckDB:        constants.ChocoPkgDuckDB,
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
		constants.ToolVSCode:     constants.WingetPkgVSCode,
		constants.ToolPowerShell: constants.WingetPkgPowerShell,
	}

	pkg, exists := wingetMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveAptPackage maps tool names to apt package IDs.
func resolveAptPackage(tool string) string {
	aptMap := map[string]string{
		constants.ToolNodeJS:        constants.AptPkgNodeJS,
		constants.ToolPython:        constants.AptPkgPython,
		constants.ToolGo:            constants.AptPkgGo,
		constants.ToolGit:           constants.AptPkgGit,
		constants.ToolGitLFS:        constants.AptPkgGitLFS,
		constants.ToolCPP:           constants.AptPkgCPP,
		constants.ToolPHP:           constants.AptPkgPHP,
		constants.ToolMySQL:         constants.AptPkgMySQL,
		constants.ToolMariaDB:       constants.AptPkgMariaDB,
		constants.ToolPostgreSQL:    constants.AptPkgPostgreSQL,
		constants.ToolSQLite:        constants.AptPkgSQLite,
		constants.ToolMongoDB:       constants.AptPkgMongoDB,
		constants.ToolCouchDB:       constants.AptPkgCouchDB,
		constants.ToolRedis:         constants.AptPkgRedis,
		constants.ToolCassandra:     constants.AptPkgCassandra,
		constants.ToolElasticsearch: constants.AptPkgElasticsearch,
	}

	pkg, exists := aptMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveBrewPackage maps tool names to Homebrew package IDs.
func resolveBrewPackage(tool string) string {
	brewMap := map[string]string{
		constants.ToolNodeJS:        constants.BrewPkgNodeJS,
		constants.ToolPython:        constants.BrewPkgPython,
		constants.ToolGo:            constants.BrewPkgGo,
		constants.ToolGit:           constants.BrewPkgGit,
		constants.ToolGitLFS:        constants.BrewPkgGitLFS,
		constants.ToolGHCLI:         constants.BrewPkgGHCLI,
		constants.ToolCPP:           constants.BrewPkgCPP,
		constants.ToolPHP:           constants.BrewPkgPHP,
		constants.ToolMySQL:         constants.BrewPkgMySQL,
		constants.ToolMariaDB:       constants.BrewPkgMariaDB,
		constants.ToolPostgreSQL:    constants.BrewPkgPostgreSQL,
		constants.ToolSQLite:        constants.BrewPkgSQLite,
		constants.ToolMongoDB:       constants.BrewPkgMongoDB,
		constants.ToolCouchDB:       constants.BrewPkgCouchDB,
		constants.ToolRedis:         constants.BrewPkgRedis,
		constants.ToolNeo4j:         constants.BrewPkgNeo4j,
		constants.ToolElasticsearch: constants.BrewPkgElasticsearch,
		constants.ToolDuckDB:        constants.BrewPkgDuckDB,
	}

	pkg, exists := brewMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// resolveSnapPackage maps tool names to Snap package IDs.
func resolveSnapPackage(tool string) string {
	snapMap := map[string]string{
		constants.ToolCouchDB: constants.SnapPkgCouchDB,
		constants.ToolRedis:   constants.SnapPkgRedis,
	}

	pkg, exists := snapMap[tool]
	if exists {
		return pkg
	}

	return tool
}

// recordInstallation saves the install record to the database.
func recordInstallation(tool, manager string) {
	version := detectInstalledVersion(tool)

	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.SaveInstalledTool(tool, version, manager)
	if err != nil {
		return
	}

	if version != "" {
		fmt.Printf(constants.MsgInstallRecorded, tool, version)
	}
}
