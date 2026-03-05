package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/cloner"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/verbose"
)

// runPull handles the "pull" subcommand.
func runPull(args []string) {
	slug, verboseMode := parsePullFlags(args)
	if len(slug) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrPullSlugRequired)
		fmt.Fprintln(os.Stderr, constants.ErrPullUsage)
		os.Exit(1)
	}
	if verboseMode {
		log, err := verbose.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create verbose log: %v\n", err)
		} else {
			defer log.Close()
		}
	}
	executePull(slug)
}

// parsePullFlags parses flags for the pull command.
func parsePullFlags(args []string) (slug string, verboseFlag bool) {
	fs := flag.NewFlagSet(constants.CmdPull, flag.ExitOnError)
	vFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	fs.Parse(args)

	if fs.NArg() > 0 {
		slug = fs.Arg(0)
	}

	return slug, *vFlag
}

// executePull loads gitmap.json, finds the repo by slug, and pulls it.
func executePull(slug string) {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadJSONRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPullLoadFailed, err)
		os.Exit(1)
	}

	matches := findBySlug(records, slug)
	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrPullNotFound, slug)
		listAvailableRepos(records)
		os.Exit(1)
	}

	for _, rec := range matches {
		pullOneRepo(rec)
	}
}

// loadJSONRecords reads ScanRecords from a JSON file.
func loadJSONRecords(path string) ([]model.ScanRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []model.ScanRecord
	err = json.NewDecoder(file).Decode(&records)

	return records, err
}

// findBySlug finds records matching the slug (case-insensitive, partial match).
func findBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	slugLower := strings.ToLower(slug)
	var exact []model.ScanRecord
	var partial []model.ScanRecord

	for _, r := range records {
		nameLower := strings.ToLower(r.RepoName)
		if nameLower == slugLower {
			exact = append(exact, r)
		} else if strings.Contains(nameLower, slugLower) {
			partial = append(partial, r)
		}
	}

	if len(exact) > 0 {
		return exact
	}

	return partial
}

// pullOneRepo runs safe-pull on a single repo using its absolute path.
func pullOneRepo(rec model.ScanRecord) {
	fmt.Printf(constants.MsgPullStarting, rec.RepoName, rec.AbsolutePath)

	if !cloner.IsGitRepo(rec.AbsolutePath) {
		fmt.Fprintf(os.Stderr, constants.ErrPullNotRepo, rec.AbsolutePath)
		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	if result.Success {
		fmt.Printf(constants.MsgPullSuccess, rec.RepoName)
	} else {
		fmt.Fprintf(os.Stderr, constants.MsgPullFailed, rec.RepoName, result.Error)
	}
}

// listAvailableRepos prints all available repo names for the user.
func listAvailableRepos(records []model.ScanRecord) {
	fmt.Fprintln(os.Stderr, constants.MsgPullAvailable)
	for _, r := range records {
		fmt.Fprintf(os.Stderr, "  - %s\n", r.RepoName)
	}
}
