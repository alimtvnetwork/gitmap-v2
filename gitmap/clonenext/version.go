package clonenext

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionSuffixRe matches a trailing -vN or -vNN suffix.
var versionSuffixRe = regexp.MustCompile(`^(.+)-v(\d+)$`)

// ParsedRepo holds the decomposed parts of a versioned repo name.
type ParsedRepo struct {
	BaseName       string
	CurrentVersion int
	HasVersion     bool
}

// ParseRepoName extracts the base name and version from a repo name.
// "macro-ahk-v11" → ("macro-ahk", 11, true)
// "macro-ahk"     → ("macro-ahk", 1, false)
func ParseRepoName(name string) ParsedRepo {
	m := versionSuffixRe.FindStringSubmatch(name)
	if m == nil {
		return ParsedRepo{BaseName: name, CurrentVersion: 1, HasVersion: false}
	}
	v, _ := strconv.Atoi(m[2])

	return ParsedRepo{BaseName: m[1], CurrentVersion: v, HasVersion: true}
}

// ResolveTarget computes the target version from a version argument.
// "v++"  → current + 1
// "v15"  → 15
func ResolveTarget(parsed ParsedRepo, arg string) (int, error) {
	lower := strings.ToLower(arg)
	if lower == "v++" {
		return parsed.CurrentVersion + 1, nil
	}
	if strings.HasPrefix(lower, "v") {
		n, err := strconv.Atoi(lower[1:])
		if err != nil || n < 1 {
			return 0, fmt.Errorf("invalid version argument: %s", arg)
		}

		return n, nil
	}

	return 0, fmt.Errorf("invalid version argument: %s (expected v++ or vN)", arg)
}

// TargetRepoName builds the full repo name for the target version.
func TargetRepoName(baseName string, version int) string {
	return fmt.Sprintf("%s-v%d", baseName, version)
}

// ReplaceRepoInURL swaps the repo name in a remote URL.
func ReplaceRepoInURL(remoteURL, currentRepo, targetRepo string) string {
	return strings.Replace(remoteURL, currentRepo, targetRepo, 1)
}
