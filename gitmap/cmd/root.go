// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// Run is the main entry point for the CLI.
func Run() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	dispatch(command)
}

// dispatch routes to the correct subcommand handler.
func dispatch(command string) {
	if dispatchCore(command) {
		return
	}
	if dispatchRelease(command) {
		return
	}
	if dispatchUtility(command) {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrUnknownCommand, command)
	printUsage()
	os.Exit(1)
}

// dispatchCore routes scan, clone, pull, and status commands.
func dispatchCore(command string) bool {
	if command == constants.CmdScan || command == constants.CmdScanAlias {
		runScan(os.Args[2:])

		return true
	}
	if command == constants.CmdClone || command == constants.CmdCloneAlias {
		runClone(os.Args[2:])

		return true
	}
	if command == constants.CmdPull || command == constants.CmdPullAlias {
		runPull(os.Args[2:])

		return true
	}
	if command == constants.CmdStatus || command == constants.CmdStatusAlias {
		runStatus(os.Args[2:])

		return true
	}
	if command == constants.CmdExec || command == constants.CmdExecAlias {
		runExec(os.Args[2:])

		return true
	}

	return false
}

// dispatchRelease routes release-related commands.
func dispatchRelease(command string) bool {
	if command == constants.CmdRelease || command == constants.CmdReleaseAlias {
		runRelease(os.Args[2:])

		return true
	}
	if command == constants.CmdReleaseBranch || command == constants.CmdReleaseBranchAlias {
		runReleaseBranch(os.Args[2:])

		return true
	}
	if command == constants.CmdReleasePending || command == constants.CmdReleasePendingAlias {
		runReleasePending(os.Args[2:])

		return true
	}
	if command == constants.CmdChangelog || command == constants.CmdChangelogAlias {
		runChangelog(os.Args[2:])

		return true
	}
	if command == constants.CmdChangelogMD {
		runChangelog([]string{constants.FlagOpenValue})

		return true
	}

	return false
}

// dispatchUtility routes setup, update, doctor, and other utility commands.
func dispatchUtility(command string) bool {
	if command == constants.CmdUpdate {
		runUpdate()

		return true
	}
	if command == constants.CmdUpdateRunner {
		runUpdateRunner()

		return true
	}
	if command == constants.CmdUpdateCleanup {
		runUpdateCleanup()

		return true
	}
	if command == constants.CmdVersion || command == constants.CmdVersionAlias {
		fmt.Printf(constants.MsgVersionFmt, constants.Version)

		return true
	}
	if command == constants.CmdHelp {
		printUsage()

		return true
	}

	return dispatchMisc(command)
}

// dispatchMisc routes remaining miscellaneous commands.
func dispatchMisc(command string) bool {
	if command == constants.CmdDesktopSync || command == constants.CmdDesktopSyncAlias {
		runDesktopSync()

		return true
	}
	if command == constants.CmdRescan || command == constants.CmdRescanAlias {
		runRescan()

		return true
	}
	if command == constants.CmdSetup {
		runSetup(os.Args[2:])

		return true
	}
	if command == constants.CmdDoctor {
		runDoctor()

		return true
	}
	if command == constants.CmdLatestBranch || command == constants.CmdLatestBranchAlias {
		runLatestBranch(os.Args[2:])

		return true
	}
	if command == constants.CmdList || command == constants.CmdListAlias {
		runList(os.Args[2:])

		return true
	}
	if command == constants.CmdGroup || command == constants.CmdGroupAlias {
		runGroup(os.Args[2:])

		return true
	}
	if command == constants.CmdDBReset {
		runDBReset(os.Args[2:])

		return true
	}

	return false
}
