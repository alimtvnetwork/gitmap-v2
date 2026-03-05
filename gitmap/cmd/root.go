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
	if command == constants.CmdScan {
		runScan(os.Args[2:])
		return
	}
	if command == constants.CmdClone {
		runClone(os.Args[2:])
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
	fmt.Println("Usage: gitmap <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  scan [dir]          Scan directory for Git repos")
	fmt.Println("  clone <source>      Re-clone from CSV/JSON/text file")
	fmt.Println("  help                Show this help message")
	fmt.Println()
	fmt.Println("Scan flags:")
	fmt.Println("  --config <path>     Config file (default: ./data/config.json)")
	fmt.Println("  --mode ssh|https    Clone URL style (default: https)")
	fmt.Println("  --output csv|json|terminal  Output format (default: terminal)")
	fmt.Println("  --output-path <dir> Output directory (default: ./gitmap-output)")
	fmt.Println("  --out-file <path>   Exact output file path")
	fmt.Println()
	fmt.Println("Clone flags:")
	fmt.Println("  --target-dir <dir>  Base directory for clones (default: .)")
}

// parseScanFlags parses flags for the scan command.
func parseScanFlags(args []string) (dir, configPath, mode, output, outFile, outputPath string) {
	fs := flag.NewFlagSet(constants.CmdScan, flag.ExitOnError)
	cfgFlag := fs.String("config", constants.DefaultConfigPath, "Path to config file")
	modeFlag := fs.String("mode", "", "Clone URL style: https or ssh")
	outputFlag := fs.String("output", "", "Output format: terminal, csv, json")
	outFileFlag := fs.String("out-file", "", "Exact output file path")
	outputPathFlag := fs.String("output-path", "", "Output directory for CSV/JSON")
	fs.Parse(args)

	dir = constants.DefaultDir
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
	}

	return dir, *cfgFlag, *modeFlag, *outputFlag, *outFileFlag, *outputPathFlag
}

// parseCloneFlags parses flags for the clone command.
func parseCloneFlags(args []string) (source, targetDir string) {
	fs := flag.NewFlagSet(constants.CmdClone, flag.ExitOnError)
	targetFlag := fs.String("target-dir", constants.DefaultDir, "Base directory for cloned repos")
	fs.Parse(args)

	source = ""
	if fs.NArg() > 0 {
		source = fs.Arg(0)
	}

	return source, *targetFlag
}
