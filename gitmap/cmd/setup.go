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
	fmt.Printf("  %s╔══════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║         gitmap setup                 ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	if dryRun {
		fmt.Printf("\n  %s[DRY RUN] No changes will be made%s\n", constants.ColorYellow, constants.ColorReset)
	}
}

// printSetupSummary shows the final results.
func printSetupSummary(r setup.SetupResult) {
	fmt.Println()
	fmt.Printf("  %s──────────────────────────────────────────%s\n", constants.ColorDim, constants.ColorReset)

	binDir, _ := filepath.Abs(".")
	_ = binDir

	if r.Applied > 0 {
		fmt.Printf("  %s✓ %d settings applied%s\n", constants.ColorGreen, r.Applied, constants.ColorReset)
	}
	if r.Skipped > 0 {
		fmt.Printf("  %s⊘ %d settings unchanged%s\n", constants.ColorDim, r.Skipped, constants.ColorReset)
	}
	if r.Failed > 0 {
		fmt.Printf("  %s✗ %d settings failed%s\n", constants.ColorYellow, r.Failed, constants.ColorReset)
		for _, e := range r.Errors {
			fmt.Printf("    %s- %s%s\n", constants.ColorYellow, e, constants.ColorReset)
		}
	}
	fmt.Println()
}
