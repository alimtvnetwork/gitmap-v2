package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/completion"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
	"github.com/user/gitmap/setup"
	"github.com/user/gitmap/store"
)

// runSetup handles the "setup" subcommand.
func runSetup(args []string) {
	checkHelp("setup", args)
	configPath, dryRun, hasConfig := parseSetupFlags(args)
	configPath = resolveSetupConfigPath(configPath, hasConfig)
	cfg := mustLoadSetupConfig(configPath)
	runSetupSteps(cfg, dryRun)
}

// mustLoadSetupConfig loads the resolved setup config or exits.
func mustLoadSetupConfig(configPath string) setup.GitSetupConfig {
	cfg, err := setup.LoadConfig(configPath)
	if err == nil {
		return cfg
	}

	fmt.Fprintf(os.Stderr, constants.ErrSetupLoadFailed, configPath, err)
	os.Exit(1)

	return setup.GitSetupConfig{}
}

// runSetupSteps applies setup and installs related shell helpers.
func runSetupSteps(cfg setup.GitSetupConfig, dryRun bool) {
	printSetupBanner(dryRun)
	result := setup.Apply(cfg, dryRun)
	installShellCompletion(dryRun)
	installCDFunction(dryRun)
	ensureGitignoreStep(dryRun)
	printSetupSummary(result)
}

// resolveSetupConfigPath prefers the bundled config unless overridden.
func resolveSetupConfigPath(configPath string, hasConfig bool) string {
	if hasConfig {
		return configPath
	}

	return resolveDefaultSetupConfigPath(configPath, store.BinaryDataDir(), constants.RepoPath)
}

// resolveDefaultSetupConfigPath picks the best default setup config path.
func resolveDefaultSetupConfigPath(configPath, binaryDataDir, repoPath string) string {
	name := filepath.Base(configPath)
	repoConfigPath := resolveRepoSetupConfigPath(repoPath, name)

	return firstExistingPath(
		filepath.Join(binaryDataDir, name),
		repoConfigPath,
		filepath.Join(constants.GitMapSubdir, constants.DBDir, name),
		configPath,
	)
}

// resolveRepoSetupConfigPath returns the source-repo setup config path.
func resolveRepoSetupConfigPath(repoPath, name string) string {
	if len(repoPath) == 0 {
		return ""
	}

	return filepath.Join(repoPath, constants.GitMapSubdir, constants.DBDir, name)
}

// firstExistingPath returns the first existing path or the first candidate.
func firstExistingPath(paths ...string) string {
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		_, err := os.Stat(path)
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			return path
		}
	}

	return paths[0]
}

// installShellCompletion detects the shell and installs completions.
func installShellCompletion(dryRun bool) {
	shell := completion.DetectShell()
	fmt.Printf("\n  %s%s %s%s\n", constants.ColorYellow, constants.SetupSectionComp, shell, constants.ColorReset)

	if dryRun {
		fmt.Printf("  %s[dry-run]%s would install %s completion\n",
			constants.ColorDim, constants.ColorReset, shell)

		return
	}

	err := completion.Install(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s%s\n", constants.ColorYellow, err, constants.ColorReset)

		return
	}

	fmt.Fprintf(os.Stderr, constants.MsgCompInstalled, shell)
}

// installCDFunction detects the shell and installs the gcd wrapper.
func installCDFunction(dryRun bool) {
	shell := completion.DetectShell()
	fmt.Printf("\n  %s%s %s%s\n", constants.ColorYellow, "cd function:", shell, constants.ColorReset)

	if dryRun {
		fmt.Printf("  %s[dry-run]%s would install gcd function for %s\n",
			constants.ColorDim, constants.ColorReset, shell)

		return
	}

	err := completion.InstallCDFunction(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s%s\n", constants.ColorYellow, err, constants.ColorReset)
	}
}

// ensureGitignoreStep adds release-related paths to .gitignore during setup.
func ensureGitignoreStep(dryRun bool) {
	fmt.Printf("\n  %s■ Gitignore —%s\n", constants.ColorYellow, constants.ColorReset)

	if dryRun {
		fmt.Printf("  %s[dry-run]%s would ensure release paths are in .gitignore\n",
			constants.ColorDim, constants.ColorReset)

		return
	}

	release.EnsureGitignore()
	fmt.Printf("  %s✓%s Release paths verified in .gitignore\n", constants.ColorGreen, constants.ColorReset)
}

// parseSetupFlags parses flags for the setup command.
func parseSetupFlags(args []string) (configPath string, dryRun, hasConfig bool) {
	fs := flag.NewFlagSet(constants.CmdSetup, flag.ExitOnError)
	cfgFlag := fs.String("config", constants.DefaultSetupConfigPath, constants.FlagDescSetupConfig)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	fs.Parse(args)
	fs.Visit(func(f *flag.Flag) {
		hasConfig = hasConfig || f.Name == "config"
	})

	return *cfgFlag, *dryRunFlag, hasConfig
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
	_, _ = filepath.Abs(".")
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
