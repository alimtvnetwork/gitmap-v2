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
	printUsageSEOFlags()
	printUsageAmendFlags()
	printUsageGoModFlags()
	printUsageInteractiveFlags()
}

// printUsageInteractiveFlags prints the interactive flags section.
func printUsageInteractiveFlags() {
	fmt.Println()
	fmt.Println(constants.HelpInteractiveFlags)
	fmt.Println(constants.HelpRefresh)
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
	fmt.Println(constants.HelpMultiGroup)
	fmt.Println(constants.HelpDBReset)
	fmt.Println(constants.HelpUpdateCleanup)
	fmt.Println(constants.HelpListVersions)
	fmt.Println(constants.HelpListReleases)
	fmt.Println(constants.HelpRevert)
	fmt.Println(constants.HelpSEOWrite)
	fmt.Println(constants.HelpAmend)
	fmt.Println(constants.HelpAmendList)
	fmt.Println(constants.HelpHistory)
	fmt.Println(constants.HelpHistoryReset)
	fmt.Println(constants.HelpStats)
	fmt.Println(constants.HelpBookmark)
	fmt.Println(constants.HelpExport)
	fmt.Println(constants.HelpImport)
	fmt.Println(constants.HelpProfile)
	fmt.Println(constants.HelpCD)
	fmt.Println(constants.HelpDiffProfiles)
	fmt.Println(constants.HelpWatch)
	fmt.Println(constants.HelpGoMod)
	fmt.Println(constants.HelpGoRepos)
	fmt.Println(constants.HelpNodeRepos)
	fmt.Println(constants.HelpReactRepos)
	fmt.Println(constants.HelpCppRepos)
	fmt.Println(constants.HelpCSharpRepos)
	fmt.Println(constants.HelpCompletion)
	fmt.Println(constants.HelpInteractive)
	fmt.Println(constants.HelpClearReleaseJSON)
	fmt.Println(constants.HelpHasAnyUpdates)
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
	fmt.Println(constants.HelpCompressFlag)
	fmt.Println(constants.HelpChecksumsFlag)
	fmt.Println(constants.HelpNoAssets)
	fmt.Println(constants.HelpTargets)
	fmt.Println(constants.HelpListTargets)
}

// printUsageSEOFlags prints the seo-write flags section.
func printUsageSEOFlags() {
	fmt.Println()
	fmt.Println(constants.HelpSEOWriteFlags)
	fmt.Println(constants.HelpSEOCSV)
	fmt.Println(constants.HelpSEOURL)
	fmt.Println(constants.HelpSEOService)
	fmt.Println(constants.HelpSEOArea)
	fmt.Println(constants.HelpSEOCompany)
	fmt.Println(constants.HelpSEOPhone)
	fmt.Println(constants.HelpSEOEmail)
	fmt.Println(constants.HelpSEOAddress)
	fmt.Println(constants.HelpSEOMaxCommits)
	fmt.Println(constants.HelpSEOInterval)
	fmt.Println(constants.HelpSEOFilesFlag)
	fmt.Println(constants.HelpSEORotate)
	fmt.Println(constants.HelpSEODryRunFlag)
	fmt.Println(constants.HelpSEOTemplateF)
	fmt.Println(constants.HelpSEOCreateTpl)
	fmt.Println(constants.HelpSEOAuthorName)
	fmt.Println(constants.HelpSEOAuthorEmail)
}

// printUsageAmendFlags prints the amend flags section.
func printUsageAmendFlags() {
	fmt.Println()
	fmt.Println(constants.HelpAmendFlags)
	fmt.Println(constants.HelpAmendName)
	fmt.Println(constants.HelpAmendEmail)
	fmt.Println(constants.HelpAmendBr)
	fmt.Println(constants.HelpAmendDry)
	fmt.Println(constants.HelpAmendForce)
}

// printUsageGoModFlags prints the gomod flags section.
func printUsageGoModFlags() {
	fmt.Println()
	fmt.Println(constants.HelpGoModFlags)
	fmt.Println(constants.HelpGoModDry)
	fmt.Println(constants.HelpGoModNoMrg)
	fmt.Println(constants.HelpGoModNoTdy)
	fmt.Println(constants.HelpGoModVerb)
	fmt.Println(constants.HelpGoModExt)
}
