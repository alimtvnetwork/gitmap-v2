package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// runListReleases handles the "list-releases" command.
func runListReleases(args []string) {
	checkHelp("list-releases", args)
	asJSON := hasListReleasesJSONFlag(args)
	limit := parseListReleasesLimit(args)
	source := parseListReleasesSource(args)
	releases := loadReleases()
	releases = filterBySource(releases, source)
	releases = applyReleaseLimit(releases, limit)

	if asJSON {
		printReleasesJSON(releases)

		return
	}

	printReleasesTerminal(releases)
}

// parseListReleasesSource extracts the --source value from args.
func parseListReleasesSource(args []string) string {
	for i, arg := range args {
		if arg == constants.FlagSource && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// filterBySource keeps only releases matching the given source (empty = all).
func filterBySource(releases []model.ReleaseRecord, source string) []model.ReleaseRecord {
	if source == "" {
		return releases
	}

	var filtered []model.ReleaseRecord
	for _, r := range releases {
		if r.Source == source {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

// hasListReleasesJSONFlag checks if --json is present in args.
func hasListReleasesJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == constants.FlagJSON {
			return true
		}
	}

	return false
}

// parseListReleasesLimit extracts the --limit N value from args.
func parseListReleasesLimit(args []string) int {
	for i, arg := range args {
		if arg == constants.FlagLimit && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err == nil && n > 0 {
				return n
			}
		}
	}

	return 0
}

// applyReleaseLimit trims releases to at most n items (0 means no limit).
func applyReleaseLimit(releases []model.ReleaseRecord, n int) []model.ReleaseRecord {
	if n <= 0 || n >= len(releases) {
		return releases
	}

	return releases[:n]
}

// loadReleases opens the DB and fetches all releases.
func loadReleases() []model.ReleaseRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrNoDatabase)
		os.Exit(1)
	}
	defer db.Close()

	releases, err := db.ListReleases()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListReleasesFailed, err)
		os.Exit(1)
	}

	return releases
}

// printReleasesTerminal renders releases as a table to stdout.
func printReleasesTerminal(releases []model.ReleaseRecord) {
	if len(releases) == 0 {
		fmt.Println(constants.MsgListReleasesEmpty)

		return
	}

	fmt.Printf(constants.MsgListReleasesHeader, len(releases))
	fmt.Println(constants.MsgListReleasesSeparator)
	fmt.Println(constants.MsgListReleasesColumns)
	for _, r := range releases {
		printReleaseRow(r)
	}
}

// printReleaseRow prints a single release row.
func printReleaseRow(r model.ReleaseRecord) {
	draft := constants.MsgNo
	if r.Draft {
		draft = constants.MsgYes
	}
	latest := constants.MsgNo
	if r.IsLatest {
		latest = constants.MsgYes
	}

	fmt.Printf(constants.MsgListReleasesRowFmt, r.Version, r.Tag, r.Branch, draft, latest, r.Source, r.CreatedAt)
}

// printReleasesJSON renders releases as JSON to stdout.
func printReleasesJSON(releases []model.ReleaseRecord) {
	data, _ := json.MarshalIndent(releases, "", constants.JSONIndent)
	fmt.Println(string(data))
}
