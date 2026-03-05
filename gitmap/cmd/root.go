// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
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
	if command == constants.CmdScan || command == constants.CmdScanAlias {
		runScan(os.Args[2:])
		return
	}
	if command == constants.CmdClone || command == constants.CmdCloneAlias {
		runClone(os.Args[2:])
		return
	}
	if command == constants.CmdUpdate {
		runUpdate()
		return
	}
	if command == constants.CmdVersion || command == constants.CmdVersionAlias {
		fmt.Printf("gitmap v%s\n", constants.Version)
		return
	}
	if command == constants.CmdDesktopSync || command == constants.CmdDesktopSyncAlias {
		runDesktopSync()
		return
	}
	if command == constants.CmdPull || command == constants.CmdPullAlias {
		runPull(os.Args[2:])
		return
	}
	if command == constants.CmdRescan || command == constants.CmdRescanAlias {
		runRescan()
		return
	}
	if command == constants.CmdHelp {
		printUsage()
		return
	}
	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	printUsage()
	os.Exit(1)
}

// printUsage displays help text for all commands.
func printUsage() {
	fmt.Printf("gitmap v%s\n\n", constants.Version)
	fmt.Println(constants.HelpUsage)
	fmt.Println()
	fmt.Println(constants.HelpCommands)
	fmt.Println(constants.HelpScan)
	fmt.Println(constants.HelpClone)
	fmt.Println(constants.HelpUpdate)
	fmt.Println(constants.HelpVersion)
	fmt.Println(constants.HelpDesktopSync)
	fmt.Println(constants.HelpPull)
	fmt.Println(constants.HelpRescan)
	fmt.Println(constants.HelpHelp)
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
	fmt.Println()
	fmt.Println(constants.HelpCloneFlags)
	fmt.Println(constants.HelpTargetDir)
	fmt.Println(constants.HelpSafePull)
	fmt.Println(constants.HelpVerbose)
}

// parseScanFlags parses flags for the scan command.
func parseScanFlags(args []string) (dir, configPath, mode, output, outFile, outputPath string, ghDesktop, openFolder, quiet bool) {
	fs := flag.NewFlagSet(constants.CmdScan, flag.ExitOnError)
	cfgFlag := fs.String("config", constants.DefaultConfigPath, constants.FlagDescConfig)
	modeFlag := fs.String("mode", "", constants.FlagDescMode)
	outputFlag := fs.String("output", "", constants.FlagDescOutput)
	outFileFlag := fs.String("out-file", "", constants.FlagDescOutFile)
	outputPathFlag := fs.String("output-path", "", constants.FlagDescOutputPath)
	ghDesktopFlag := fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	openFlag := fs.Bool("open", false, constants.FlagDescOpen)
	quietFlag := fs.Bool("quiet", false, constants.FlagDescQuiet)
	fs.Parse(args)

	dir = constants.DefaultDir
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
	}

	return dir, *cfgFlag, *modeFlag, *outputFlag, *outFileFlag, *outputPathFlag, *ghDesktopFlag, *openFlag, *quietFlag
}

// parseCloneFlags parses flags for the clone command.
func parseCloneFlags(args []string) (source, targetDir string, safePull, ghDesktop, verbose bool) {
	fs := flag.NewFlagSet(constants.CmdClone, flag.ExitOnError)
	targetFlag := fs.String("target-dir", constants.DefaultDir, constants.FlagDescTargetDir)
	safePullFlag := fs.Bool("safe-pull", false, constants.FlagDescSafePull)
	ghDesktopFlag := fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	fs.Parse(args)

	source = ""
	if fs.NArg() > 0 {
		source = fs.Arg(0)
	}

	return source, *targetFlag, *safePullFlag, *ghDesktopFlag, *verboseFlag
}
