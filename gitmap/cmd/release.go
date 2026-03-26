// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/config"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/release"
	"github.com/user/gitmap/store"
)

// runRelease handles the 'release' command.
func runRelease(args []string) {
	checkHelp("release", args)
	requireOnline()
	version, assets, commit, branch, bump, notes, targets, zipGroups, zipItems, bundleName, draft, dryRun, verbose, compress, checksums, noAssets, listTargets, noCommit := parseReleaseFlags(args)
	_ = verbose

	if listTargets {
		printListTargets(targets)

		return
	}

	validateReleaseFlags(version, bump, commit, branch)
	executeRelease(version, assets, commit, branch, bump, notes, targets, zipGroups, zipItems, bundleName, draft, dryRun, verbose, compress, checksums, noAssets, noCommit)
}

// executeRelease builds options and runs the release workflow.
func executeRelease(version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, noAssets, noCommit bool) {
	cfg, _ := config.LoadFromFile(constants.DefaultConfigPath)

	opts := release.Options{
		Version: version, Assets: assets,
		Commit: commit, Branch: branch,
		Bump: bump, Notes: notes, Targets: targets,
		ConfigTargets: cfg.Release.Targets,
		ZipGroups:     zipGroups,
		ZipItems:      zipItems,
		BundleName:    bundleName,
		Draft:         draft, DryRun: dryRun,
		Verbose:       verbose,
		Compress:      compress || cfg.Release.Compress,
		Checksums:     checksums || cfg.Release.Checksums,
		NoAssets:      noAssets,
		NoCommit:      noCommit,
	}
	err := release.Execute(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	persistReleaseToDB()
}

// validateReleaseFlags checks for mutually exclusive flags.
func validateReleaseFlags(version, bump, commit, branch string) {
	if len(bump) > 0 && len(version) > 0 {
		fmt.Fprint(os.Stderr, constants.ErrReleaseBumpConflict)
		os.Exit(1)
	}
	if len(commit) > 0 && len(branch) > 0 {
		fmt.Fprint(os.Stderr, constants.ErrReleaseCommitBranch)
		os.Exit(1)
	}
}

// zipGroupFlag collects multiple --zip-group values.
type zipGroupFlag []string

func (z *zipGroupFlag) String() string { return fmt.Sprintf("%v", *z) }
func (z *zipGroupFlag) Set(val string) error {
	*z = append(*z, val)

	return nil
}

// zipItemFlag collects multiple -Z values.
type zipItemFlag []string

func (z *zipItemFlag) String() string { return fmt.Sprintf("%v", *z) }
func (z *zipItemFlag) Set(val string) error {
	*z = append(*z, val)

	return nil
}

// parseReleaseFlags parses flags for the release command.
func parseReleaseFlags(args []string) (version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, noAssets, listTargets, noCommit bool) {
	fs := flag.NewFlagSet(constants.CmdRelease, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	commitFlag := fs.String("commit", "", constants.FlagDescCommit)
	branchFlag := fs.String("branch", "", constants.FlagDescRelBranch)
	bumpFlag := fs.String("bump", "", constants.FlagDescBump)
	notesFlag := fs.String("notes", "", constants.FlagDescNotes)
	targetsFlag := fs.String("targets", "", constants.FlagDescTargets)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	compressFlag := fs.Bool("compress", false, constants.FlagDescCompress)
	checksumsFlag := fs.Bool("checksums", false, constants.FlagDescChecksums)
	noAssetsFlag := fs.Bool("no-assets", false, constants.FlagDescNoAssets)
	listTargetsFlag := fs.Bool("list-targets", false, constants.FlagDescListTargets)
	bundleFlag := fs.String("bundle", "", constants.FlagDescZGBundle)
	noCommitFlag := fs.Bool("no-commit", false, constants.FlagDescNoCommit)

	var zgGroups zipGroupFlag
	var zgItems zipItemFlag

	fs.Var(&zgGroups, "zip-group", constants.FlagDescZGZipGroup)
	fs.Var(&zgItems, "Z", constants.FlagDescZGZipItem)

	// Register -N as shorthand for --notes.
	fs.StringVar(notesFlag, "N", "", constants.FlagDescNotes)

	fs.Parse(args)

	version = ""
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}

	return version, *assetsFlag, *commitFlag, *branchFlag, *bumpFlag, *notesFlag, *targetsFlag, []string(zgGroups), []string(zgItems), *bundleFlag, *draftFlag, *dryRunFlag, *verboseFlag, *compressFlag, *checksumsFlag, *noAssetsFlag, *listTargetsFlag, *noCommitFlag
}

// printListTargets resolves and prints the target matrix, then returns.
func printListTargets(flagTargets string) {
	cfg, _ := config.LoadFromFile(constants.DefaultConfigPath)

	targets, err := release.ResolveTargets(flagTargets, cfg.Release.Targets)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	source := resolveTargetSource(flagTargets, cfg.Release.Targets)

	fmt.Printf(constants.MsgListTargetsHeader, len(targets))
	fmt.Printf(constants.MsgListTargetsSource, source)

	for _, t := range targets {
		fmt.Printf(constants.MsgListTargetsRow, t.GOOS, t.GOARCH)
	}
}

// resolveTargetSource returns a human-readable label for the active target source.
func resolveTargetSource(flagTargets string, configTargets []model.ReleaseTarget) string {
	if len(flagTargets) > 0 {
		return "--targets flag"
	}

	if len(configTargets) > 0 {
		return "config.json (release.targets)"
	}

	return "built-in defaults"
}

// persistReleaseToDB saves the release metadata to SQLite if available.
func persistReleaseToDB() {
	meta := release.LastMeta
	if meta == nil {
		return
	}

	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()

	_ = db.Migrate()
	_ = db.UpsertRelease(releaseMetaToRecord(*meta))
}

// releaseMetaToRecord converts a ReleaseMeta to a ReleaseRecord for DB storage.
func releaseMetaToRecord(m release.ReleaseMeta) model.ReleaseRecord {
	return model.ReleaseRecord{
		Version:      m.Version,
		Tag:          m.Tag,
		Branch:       m.Branch,
		SourceBranch: m.SourceBranch,
		CommitSha:    m.Commit,
		Changelog:    store.JoinChangelog(m.Changelog),
		Notes:        m.Notes,
		Draft:        m.Draft,
		PreRelease:   m.PreRelease,
		IsLatest:     m.IsLatest,
		Source:       model.SourceRelease,
		CreatedAt:    m.CreatedAt,
	}
}
