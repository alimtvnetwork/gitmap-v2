package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/setup"
)

// runSetup handles the "setup" subcommand.
func runSetup(args []string) {
	configPath, dryRun := parseSetupFlags(args)
	cfg, err := setup.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSetupLoadFailed, err)
		os.Exit(1)
	}

	printSetupBanner(dryRun)
	result := setup.Apply(cfg, dryRun)
	printSetupSummary(result)
}

// parseSetupFlags parses flags for the setup command.
func parseSetupFlags(args []string) (configPath string, dryRun bool) {
	fs := flag.NewFlagSet(constants.CmdSetup, flag.ExitOnError)
	cfgFlag := fs.String("config", constants.DefaultSetupConfigPath, constants.FlagDescSetupConfig)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	fs.Parse(args)

	return *cfgFlag, *dryRunFlag
}

// printSetupBanner shows the setup header.
func printSetupBanner(dryRun bool) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.SetupBannerTop, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.SetupBannerTitle, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.SetupBannerBottom, constants.ColorReset)
	if dryRun {
		fmt.Printf("\n  %s%s%s\n", constants.ColorYellow, constants.SetupDryRunFmt, constants.ColorReset)
	}
}

// printSetupSummary shows the final results.
func printSetupSummary(r setup.SetupResult) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermSeparator, constants.ColorReset)
	_ = filepath.Abs(".")
	printSetupCounts(r)
	printSetupErrors(r)
	fmt.Println()
}

// printSetupCounts prints applied/skipped/failed counts.
func printSetupCounts(r setup.SetupResult) {
	if r.Applied > 0 {
		fmt.Printf("  %s"+constants.SetupAppliedFmt+"%s\n", constants.ColorGreen, r.Applied, constants.ColorReset)
	}
	if r.Skipped > 0 {
		fmt.Printf("  %s"+constants.SetupSkippedFmt+"%s\n", constants.ColorDim, r.Skipped, constants.ColorReset)
	}
	if r.Failed > 0 {
		fmt.Printf("  %s"+constants.SetupFailedFmt+"%s\n", constants.ColorYellow, r.Failed, constants.ColorReset)
	}
}

// printSetupErrors prints each failed setting detail.
func printSetupErrors(r setup.SetupResult) {
	if r.Failed == 0 {
		return
	}
	for _, e := range r.Errors {
		fmt.Printf("    %s"+constants.SetupErrorEntryFmt+"%s\n", constants.ColorYellow, e, constants.ColorReset)
	}
}
