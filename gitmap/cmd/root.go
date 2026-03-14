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

// dispatch routes to the correct subcommand handler with audit tracking.
func dispatch(command string) {
	auditID, auditStart := recordAuditStart(command, os.Args[2:])

	if dispatchCore(command) {
		recordAuditEnd(auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchRelease(command) {
		recordAuditEnd(auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchUtility(command) {
		recordAuditEnd(auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchProjectRepos(command) {
		recordAuditEnd(auditID, auditStart, 0, "", 0)

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
		checkHelp("update", os.Args[2:])
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
	if command == constants.CmdRevert {
		runRevert(os.Args[2:])

		return true
	}
	if command == constants.CmdRevertRunner {
		runRevertRunner()

		return true
	}
	if command == constants.CmdVersion || command == constants.CmdVersionAlias {
		checkHelp("version", os.Args[2:])
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
		checkHelp("desktop-sync", os.Args[2:])
		runDesktopSync()

		return true
	}
	if command == constants.CmdRescan || command == constants.CmdRescanAlias {
		checkHelp("rescan", os.Args[2:])
		runRescan()

		return true
	}
	if command == constants.CmdSetup {
		runSetup(os.Args[2:])

		return true
	}
	if command == constants.CmdDoctor {
		checkHelp("doctor", os.Args[2:])
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
	if command == constants.CmdMultiGroup || command == constants.CmdMultiGroupAlias {
		runMultiGroup(os.Args[2:])

		return true
	}
	if command == constants.CmdDBReset {
		runDBReset(os.Args[2:])

		return true
	}
	if command == constants.CmdListVersions || command == constants.CmdListVersionsAlias {
		runListVersions(os.Args[2:])

		return true
	}
	if command == constants.CmdListReleases || command == constants.CmdListReleasesAlias {
		runListReleases(os.Args[2:])

		return true
	}
	if command == constants.CmdSEOWrite || command == constants.CmdSEOWriteAlias {
		runSEOWrite(os.Args[2:])

		return true
	}
	if command == constants.CmdAmend || command == constants.CmdAmendAlias {
		runAmend(os.Args[2:])

		return true
	}
	if command == constants.CmdAmendList || command == constants.CmdAmendListAlias {
		runAmendList(os.Args[2:])

		return true
	}
	if command == constants.CmdHistory || command == constants.CmdHistoryAlias {
		runHistory(os.Args[2:])

		return true
	}
	if command == constants.CmdHistoryReset || command == constants.CmdHistoryResetAlias {
		runHistoryReset(os.Args[2:])

		return true
	}
	if command == constants.CmdStats || command == constants.CmdStatsAlias {
		runStats(os.Args[2:])

		return true
	}
	if command == constants.CmdBookmark || command == constants.CmdBookmarkAlias {
		runBookmark(os.Args[2:])

		return true
	}
	if command == constants.CmdExport || command == constants.CmdExportAlias {
		runExport(os.Args[2:])

		return true
	}
	if command == constants.CmdImport || command == constants.CmdImportAlias {
		runImport(os.Args[2:])

		return true
	}
	if command == constants.CmdProfile || command == constants.CmdProfileAlias {
		runProfile(os.Args[2:])

		return true
	}
	if command == constants.CmdCDCmd || command == constants.CmdCDCmdAlias {
		runCD(os.Args[2:])

		return true
	}
	if command == constants.CmdDiffProfiles || command == constants.CmdDiffProfilesAlias {
		runDiffProfiles(os.Args[2:])

		return true
	}
	if command == constants.CmdWatch || command == constants.CmdWatchAlias {
		runWatch(os.Args[2:])

		return true
	}
	if command == constants.CmdGoMod || command == constants.CmdGoModAlias {
		runGoMod(os.Args[2:])

		return true
	}
	if command == constants.CmdCompletion || command == constants.CmdCompletionAlias {
		runCompletion(os.Args[2:])

		return true
	}

	return false
}

// dispatchProjectRepos routes project type query commands.
func dispatchProjectRepos(command string) bool {
	if command == constants.CmdGoRepos || command == constants.CmdGoReposAlias {
		runProjectRepos(constants.ProjectKeyGo, os.Args[2:])

		return true
	}
	if command == constants.CmdNodeRepos || command == constants.CmdNodeReposAlias {
		runProjectRepos(constants.ProjectKeyNode, os.Args[2:])

		return true
	}
	if command == constants.CmdReactRepos || command == constants.CmdReactReposAlias {
		runProjectRepos(constants.ProjectKeyReact, os.Args[2:])

		return true
	}
	if command == constants.CmdCppRepos || command == constants.CmdCppReposAlias {
		runProjectRepos(constants.ProjectKeyCpp, os.Args[2:])

		return true
	}
	if command == constants.CmdCSharpRepos || command == constants.CmdCSharpAlias {
		runProjectRepos(constants.ProjectKeyCSharp, os.Args[2:])

		return true
	}

	return false
}
