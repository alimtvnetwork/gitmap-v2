package cmd

import (
	"fmt"

	"github.com/user/gitmap/constants"
)

// printUsageCompact prints a minimal command list without descriptions.
func printUsageCompact() {
	fmt.Printf(constants.UsageHeaderFmt, constants.Version)
	fmt.Println(constants.HelpUsage)
	fmt.Println()
	fmt.Println(constants.HelpGroupScanning)
	fmt.Println(constants.CompactScanning)
	fmt.Println(constants.HelpGroupCloning)
	fmt.Println(constants.CompactCloning)
	fmt.Println(constants.HelpGroupGitOps)
	fmt.Println(constants.CompactGitOps)
	fmt.Println(constants.HelpGroupNavigation)
	fmt.Println(constants.CompactNavigation)
	fmt.Println(constants.HelpGroupRelease)
	fmt.Println(constants.CompactRelease)
	fmt.Println(constants.HelpGroupReleaseInfo)
	fmt.Println(constants.CompactRelInfo)
	fmt.Println(constants.HelpGroupData)
	fmt.Println(constants.CompactData)
	fmt.Println(constants.HelpGroupHistory)
	fmt.Println(constants.CompactHistory)
	fmt.Println(constants.HelpGroupAmendGroup)
	fmt.Println(constants.CompactAmend)
	fmt.Println(constants.HelpGroupProject)
	fmt.Println(constants.CompactProject)
	fmt.Println(constants.HelpGroupSSH)
	fmt.Println(constants.CompactSSH)
	fmt.Println(constants.HelpGroupZip)
	fmt.Println(constants.CompactZip)
	fmt.Println(constants.HelpGroupEnvTools)
	fmt.Println(constants.CompactEnvTools)
	fmt.Println(constants.HelpGroupTasks)
	fmt.Println(constants.CompactTasks)
	fmt.Println(constants.HelpGroupVisualize)
	fmt.Println(constants.CompactVisualize)
	fmt.Println(constants.HelpGroupUtilities)
	fmt.Println(constants.CompactUtilities)
	fmt.Println()
	fmt.Println(constants.HelpGroupHint)
}
