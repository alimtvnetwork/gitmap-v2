package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/release"
)

// runListVersions handles the "list-versions" command.
func runListVersions(args []string) {
	asJSON := hasListVersionsJSONFlag(args)
	versions := collectVersionTags()

	if asJSON {
		printVersionsJSON(versions)

		return
	}

	printVersionsTerminal(versions)
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

// collectVersionTags reads git tags, parses, sorts descending.
func collectVersionTags() []release.Version {
	cmd := exec.Command(constants.GitBin, constants.GitTag, constants.GitTagListFlag, constants.GitTagGlob)
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

// printVersionsTerminal prints versions one per line.
func printVersionsTerminal(versions []release.Version) {
	for _, v := range versions {
		fmt.Println(v.String())
	}
}

// printVersionsJSON prints versions as a JSON array.
func printVersionsJSON(versions []release.Version) {
	strs := make([]string, len(versions))
	for i, v := range versions {
		strs[i] = v.String()
	}

	data, _ := json.Marshal(strs)
	fmt.Println(string(data))
}
