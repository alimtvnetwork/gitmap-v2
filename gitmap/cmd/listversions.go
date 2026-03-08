package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
)

// versionEntry pairs a parsed version with its changelog notes and source.
type versionEntry struct {
	Version release.Version
	Notes   []string
	Source  string
}

// runListVersions handles the "list-versions" command.
func runListVersions(args []string) {
	asJSON := hasListVersionsJSONFlag(args)
	limit := parseListVersionsLimit(args)
	source := parseListVersionsSource(args)
	entries := collectVersionEntries()
	entries = filterVersionsBySource(entries, source)
	entries = applyVersionLimit(entries, limit)

	if asJSON {
		printVersionEntriesJSON(entries)

		return
	}

	printVersionEntriesTerminal(entries)
}

// parseListVersionsSource extracts the --source value from args.
func parseListVersionsSource(args []string) string {
	for i, arg := range args {
		if arg == constants.FlagSource && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// filterVersionsBySource keeps only entries matching the given source (empty = all).
func filterVersionsBySource(entries []versionEntry, source string) []versionEntry {
	if source == "" {
		return entries
	}

	var filtered []versionEntry
	for _, e := range entries {
		if e.Source == source {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// hasListVersionsJSONFlag checks if --json is present in args.
func hasListVersionsJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == constants.FlagJSON {
			return true
		}
	}

	return false
}

// parseListVersionsLimit extracts the --limit N value from args.
func parseListVersionsLimit(args []string) int {
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

// applyVersionLimit trims entries to at most n items (0 means no limit).
func applyVersionLimit(entries []versionEntry, n int) []versionEntry {
	if n <= 0 || n >= len(entries) {
		return entries
	}

	return entries[:n]
}

// collectVersionEntries reads tags, parses, sorts, and attaches changelog + source.
func collectVersionEntries() []versionEntry {
	versions := collectVersionTags()
	changelog := loadChangelogMap()
	sources := loadVersionSourceMap()

	entries := make([]versionEntry, len(versions))
	for i, v := range versions {
		entries[i] = versionEntry{Version: v, Notes: changelog[v.String()], Source: sources[v.String()]}
	}

	return entries
}

// loadVersionSourceMap reads the Releases table to build a tag→source map.
func loadVersionSourceMap() map[string]string {
	db, err := openDB()
	if err != nil {
		return map[string]string{}
	}
	defer db.Close()

	releases, err := db.ListReleases()
	if err != nil {
		return map[string]string{}
	}

	m := make(map[string]string, len(releases))
	for _, r := range releases {
		m[r.Tag] = r.Source
	}

	return m
}

// collectVersionTags reads git tags, parses, sorts descending.
func collectVersionTags() []release.Version {
	cmd := exec.Command(constants.GitBin, constants.GitTag,
		constants.GitTagListFlag, constants.GitTagGlob)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrListVersionsNoTags)
		os.Exit(1)
	}

	versions := parseVersionTags(strings.TrimSpace(string(out)))
	if len(versions) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrListVersionsNoTags)
		os.Exit(1)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].GreaterThan(versions[j])
	})

	return versions
}

// parseVersionTags parses lines into valid versions.
func parseVersionTags(output string) []release.Version {
	lines := strings.Split(output, "\n")
	var versions []release.Version

	for _, line := range lines {
		tag := strings.TrimSpace(line)
		if len(tag) == 0 {
			continue
		}
		v, err := release.Parse(tag)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}

	return versions
}

// loadChangelogMap reads CHANGELOG.md into a version→notes map.
func loadChangelogMap() map[string][]string {
	entries, err := release.ReadChangelog()
	if err != nil {
		return map[string][]string{}
	}

	m := make(map[string][]string, len(entries))
	for _, e := range entries {
		m[e.Version] = e.Notes
	}

	return m
}

// printVersionEntriesTerminal prints versions with source and changelog sub-points.
func printVersionEntriesTerminal(entries []versionEntry) {
	for _, e := range entries {
		if e.Source != "" {
			fmt.Printf("%s  [%s]\n", e.Version.String(), e.Source)
		} else {
			fmt.Println(e.Version.String())
		}
		for _, note := range e.Notes {
			fmt.Printf("  - %s\n", note)
		}
	}
}

// lvJSONEntry is the JSON output shape for list-versions.
type lvJSONEntry struct {
	Version   string   `json:"version"`
	Changelog []string `json:"changelog,omitempty"`
}

// printVersionEntriesJSON prints versions with changelog as JSON.
func printVersionEntriesJSON(entries []versionEntry) {
	out := make([]lvJSONEntry, len(entries))
	for i, e := range entries {
		out[i] = lvJSONEntry{Version: e.Version.String(), Changelog: e.Notes}
	}

	data, _ := json.MarshalIndent(out, "", constants.JSONIndent)
	fmt.Println(string(data))
}
