// Package gitutil — latest-branch helpers.
package gitutil

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/user/gitmap/constants"
)

// RemoteBranchInfo holds commit metadata for a remote-tracking branch.
type RemoteBranchInfo struct {
	RemoteRef  string
	CommitDate time.Time
	Sha        string
	Subject    string
}

// LatestBranchResult holds the resolved latest branch information.
type LatestBranchResult struct {
	BranchNames []string
	Remote      string
	Sha         string
	CommitDate  string
	Subject     string
	RemoteRef   string
}

// IsInsideWorkTree checks if the current directory is inside a git repo.
func IsInsideWorkTree() bool {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// FetchAllPrune runs git fetch --all --prune.
func FetchAllPrune() error {
	cmd := exec.Command(constants.GitBin, "fetch", "--all", "--prune")
	return cmd.Run()
}

// ListRemoteBranches returns trimmed remote-tracking branch names,
// excluding HEAD pointer lines.
func ListRemoteBranches() ([]string, error) {
	cmd := exec.Command(constants.GitBin, "branch", "-r")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var refs []string
	for _, line := range strings.Split(string(out), "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if strings.Contains(trimmed, " -> ") {
			continue
		}
		refs = append(refs, trimmed)
	}
	return refs, nil
}

// FilterByRemote keeps only refs starting with "<remote>/".
func FilterByRemote(refs []string, remote string) []string {
	prefix := remote + "/"
	var filtered []string
	for _, r := range refs {
		if strings.HasPrefix(r, prefix) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

// FilterByPattern keeps only refs whose branch name (after remote prefix)
// matches the given glob or substring pattern.
func FilterByPattern(refs []string, pattern string) []string {
	var filtered []string
	for _, r := range refs {
		name := stripPrefix(r)
		if matchesPattern(name, pattern) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

// stripPrefix removes the "<remote>/" prefix from a ref.
func stripPrefix(ref string) string {
	if idx := strings.Index(ref, "/"); idx >= 0 {

		return ref[idx+1:]
	}

	return ref
}

// matchesPattern checks if name matches a glob pattern or contains
// the pattern as a substring.
func matchesPattern(name, pattern string) bool {
	matched, err := filepath.Match(pattern, name)
	if err == nil && matched {

		return true
	}

	return strings.Contains(name, pattern)
}

// ReadBranchTips reads commit date, SHA, and subject for each ref.
func ReadBranchTips(refs []string) ([]RemoteBranchInfo, error) {
	var items []RemoteBranchInfo
	for _, ref := range refs {
		cmd := exec.Command(constants.GitBin, "log", "-1", "--format=%cI|%H|%s", ref)
		out, err := cmd.Output()
		if err != nil {
			continue
		}
		line := strings.TrimSpace(string(out))
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue
		}
		t, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			continue
		}
		items = append(items, RemoteBranchInfo{
			RemoteRef:  ref,
			CommitDate: t,
			Sha:        parts[1],
			Subject:    parts[2],
		})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("could not read commit info for remote branches")
	}
	return items, nil
}

// SortByDateDesc sorts items by CommitDate descending.
func SortByDateDesc(items []RemoteBranchInfo) {
	sort.Slice(items, func(i, j int) bool {

		return items[i].CommitDate.After(items[j].CommitDate)
	})
}

// SortByNameAsc sorts items by branch name (RemoteRef) ascending.
func SortByNameAsc(items []RemoteBranchInfo) {
	sort.Slice(items, func(i, j int) bool {

		return items[i].RemoteRef < items[j].RemoteRef
	})
}

// ResolvePointsAt returns branch names that point exactly at sha.
func ResolvePointsAt(sha, remote string) []string {
	cmd := exec.Command(constants.GitBin, "for-each-ref",
		fmt.Sprintf("--points-at=%s", sha),
		fmt.Sprintf("refs/remotes/%s", remote),
		"--format=%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	prefix := remote + "/"
	var names []string
	seen := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 || trimmed == remote+"/HEAD" {
			continue
		}
		name := strings.TrimPrefix(trimmed, prefix)
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

// ResolveContains returns branch names whose history contains sha.
func ResolveContains(sha, remote string) []string {
	cmd := exec.Command(constants.GitBin, "branch", "-r", "--contains", sha)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	prefix := remote + "/"
	var names []string
	seen := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 || strings.Contains(trimmed, " -> ") {
			continue
		}
		if !strings.HasPrefix(trimmed, prefix) {
			continue
		}
		name := strings.TrimPrefix(trimmed, prefix)
		if name == "HEAD" {
			continue
		}
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}
