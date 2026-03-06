package cmd

import (
	"fmt"

	"github.com/user/gitmap/constants"
)

// printUsage displays help text for all commands and flags.
func printUsage() {
	fmt.Printf(constants.UsageHeaderFmt, constants.Version)
	fmt.Println(constants.HelpUsage)
	fmt.Println()
	printUsageCommands()
	printUsageScanFlags()
	printUsageCloneFlags()
	printUsageReleaseFlags()
}

// printUsageCommands prints the available commands section.
func printUsageCommands() {
	fmt.Println(constants.HelpCommands)
	fmt.Println(constants.HelpScan)
	fmt.Println(constants.HelpClone)
	fmt.Println(constants.HelpUpdate)
	fmt.Println(constants.HelpVersion)
	fmt.Println(constants.HelpDesktopSync)
	fmt.Println(constants.HelpPull)
	fmt.Println(constants.HelpRescan)
	fmt.Println(constants.HelpSetup)
	fmt.Println(constants.HelpStatus)
	fmt.Println(constants.HelpExec)
	fmt.Println(constants.HelpRelease)
	fmt.Println(constants.HelpReleaseBr)
	fmt.Println(constants.HelpReleasePend)
	fmt.Println(constants.HelpChangelog)
	fmt.Println(constants.HelpDoctor)
	fmt.Println(constants.HelpLatestBr)
	fmt.Println(constants.HelpList)
	fmt.Println(constants.HelpGroup)
	fmt.Println(constants.HelpDBReset)
	fmt.Println(constants.HelpUpdateCleanup)
	fmt.Println(constants.HelpHelp)
}

// printUsageScanFlags prints the scan flags section.
func printUsageScanFlags() {
	fmt.Println()
	fmt.Println(constants.HelpScanFlags)
	fmt.Println(constants.HelpConfig)
	fmt.Println(constants.HelpMode)
	fmt.Println(constants.HelpOutput)
	fmt.Println(constants.HelpOutputPath)
	fmt.Println(constants.HelpOutFile)
	fmt.Println(constants.HelpGitHubDesktop)
	fmt.Println(constants.HelpOpen)
	fmt.Println(constants.HelpQuiet)
}

// printUsageCloneFlags prints the clone flags section.
func printUsageCloneFlags() {
	fmt.Println()
	fmt.Println(constants.HelpCloneFlags)
	fmt.Println(constants.HelpTargetDir)
	fmt.Println(constants.HelpSafePull)
	fmt.Println(constants.HelpVerbose)
}

// printUsageReleaseFlags prints the release flags section.
func printUsageReleaseFlags() {
	fmt.Println()
	fmt.Println(constants.HelpReleaseFlags)
	fmt.Println(constants.HelpAssets)
	fmt.Println(constants.HelpCommit)
	fmt.Println(constants.HelpRelBranch)
	fmt.Println(constants.HelpBump)
	fmt.Println(constants.HelpDraft)
	fmt.Println(constants.HelpDryRun)
}
