package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runInstall handles the "install" command.
func runInstall(args []string) {
	checkHelp("install", args)

	fs := flag.NewFlagSet("install", flag.ExitOnError)

	var manager, version string
	var verbose, dryRun, check, list bool

	fs.StringVar(&manager, constants.FlagInstallManager, "", constants.FlagDescInstallManager)
	fs.StringVar(&version, constants.FlagInstallVersion, "", constants.FlagDescInstallVersion)
	fs.BoolVar(&verbose, constants.FlagInstallVerbose, false, constants.FlagDescInstallVerbose)
	fs.BoolVar(&dryRun, constants.FlagInstallDryRun, false, constants.FlagDescInstallDryRun)
	fs.BoolVar(&check, constants.FlagInstallCheck, false, constants.FlagDescInstallCheck)
	fs.BoolVar(&list, constants.FlagInstallList, false, constants.FlagDescInstallList)
	fs.Parse(args)

	if list {
		printInstallList()

		return
	}

	tool := fs.Arg(0)
	if tool == "" {
		fmt.Fprint(os.Stderr, constants.ErrInstallToolRequired)
		os.Exit(1)
	}

	validateToolName(tool)

	opts := installOptions{
		Tool:    tool,
		Manager: manager,
		Version: version,
		Verbose: verbose,
		DryRun:  dryRun,
		Check:   check,
	}

	executeInstall(opts)
}

// installOptions holds parsed install flags.
type installOptions struct {
	Tool    string
	Manager string
	Version string
	Verbose bool
	DryRun  bool
	Check   bool
}

// printInstallList prints all supported tools.
func printInstallList() {
	fmt.Print(constants.MsgInstallListHeader)

	for tool, desc := range constants.InstallToolDescriptions {
		fmt.Printf(constants.MsgInstallListRow, tool, desc)
	}
}

// validateToolName checks if the tool is supported.
func validateToolName(tool string) {
	_, exists := constants.InstallToolDescriptions[tool]
	if exists {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrInstallUnknownTool, tool)
	os.Exit(1)
}

// executeInstall runs the install flow for a tool.
func executeInstall(opts installOptions) {
	if opts.Tool == constants.ToolNppSettings {
		runNppSettingsOnly()

		return
	}

	fmt.Printf(constants.MsgInstallChecking, opts.Tool)

	existingVersion := detectInstalledVersion(opts.Tool)
	if existingVersion != "" {
		fmt.Printf(constants.MsgInstallFound, opts.Tool, existingVersion)

		return
	}

	if opts.Check {
		fmt.Printf(constants.MsgInstallNotFound, opts.Tool)

		return
	}

	installTool(opts)

	if opts.Tool == constants.ToolNpp {
		runNppSettings()
	}
}