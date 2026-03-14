// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/release"
	"github.com/user/gitmap/store"
)

// runRelease handles the 'release' command.
func runRelease(args []string) {
	checkHelp("release", args)
	version, assets, commit, branch, bump, draft, dryRun, verbose, compress, checksums := parseReleaseFlags(args)
	_ = verbose
	validateReleaseFlags(version, bump, commit, branch)
	executeRelease(version, assets, commit, branch, bump, draft, dryRun, verbose, compress, checksums)
}

// executeRelease builds options and runs the release workflow.
func executeRelease(version, assets, commit, branch, bump string, draft, dryRun, verbose, compress bool) {
	opts := release.Options{
		Version: version, Assets: assets,
		Commit: commit, Branch: branch,
		Bump: bump, Draft: draft,
		DryRun: dryRun, Verbose: verbose,
		Compress: compress,
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

// parseReleaseFlags parses flags for the release command.
func parseReleaseFlags(args []string) (version, assets, commit, branch, bump string, draft, dryRun, verbose, compress bool) {
	fs := flag.NewFlagSet(constants.CmdRelease, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	commitFlag := fs.String("commit", "", constants.FlagDescCommit)
	branchFlag := fs.String("branch", "", constants.FlagDescRelBranch)
	bumpFlag := fs.String("bump", "", constants.FlagDescBump)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	compressFlag := fs.Bool("compress", false, constants.FlagDescCompress)
	fs.Parse(args)

	version = ""
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}

	return version, *assetsFlag, *commitFlag, *branchFlag, *bumpFlag, *draftFlag, *dryRunFlag, *verboseFlag, *compressFlag
}

// persistReleaseToDB saves the release metadata to SQLite if available.
func persistReleaseToDB() {
	meta := release.LastMeta
	if meta == nil {
		return
	}

	db, err := store.Open(constants.DefaultOutputFolder)
	if err != nil {
		return
	}
	defer db.Close()

	_ = db.Migrate()
	_ = db.UpsertRelease(metaToRecord(*meta))
}

// metaToRecord converts a ReleaseMeta to a ReleaseRecord for DB storage.
func metaToRecord(m release.ReleaseMeta) model.ReleaseRecord {
	return model.ReleaseRecord{
		Version:      m.Version,
		Tag:          m.Tag,
		Branch:       m.Branch,
		SourceBranch: m.SourceBranch,
		CommitSha:    m.Commit,
		Changelog:    store.JoinChangelog(m.Changelog),
		Draft:        m.Draft,
		PreRelease:   m.PreRelease,
		IsLatest:     m.IsLatest,
		Source:       model.SourceRelease,
		CreatedAt:    m.CreatedAt,
	}
}
