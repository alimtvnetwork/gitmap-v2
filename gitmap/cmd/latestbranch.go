// Package cmd — latest-branch command handler.
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/gitutil"
)

// latestBranchJSON is the JSON output structure.
type latestBranchJSON struct {
	Branch     []string              `json:"branch"`
	Remote     string                `json:"remote"`
	Sha        string                `json:"sha"`
	CommitDate string                `json:"commitDate"`
	Subject    string                `json:"subject"`
	Ref        string                `json:"ref"`
	Top        []latestBranchTopItem `json:"top,omitempty"`
}

type latestBranchTopItem struct {
	Branch     string `json:"branch"`
	Sha        string `json:"sha"`
	CommitDate string `json:"commitDate"`
	Subject    string `json:"subject"`
}

// runLatestBranch handles the 'latest-branch' / 'lb' command.
func runLatestBranch(args []string) {
	remote, allRemotes, containsFallback, top, format, noFetch := parseLatestBranchFlags(args)
	isMachine := format == constants.OutputJSON || format == constants.OutputCSV

	// 1. Validate git repo.
	if !gitutil.IsInsideWorkTree() {
		fmt.Fprintln(os.Stderr, constants.ErrLatestBranchNotRepo)
		os.Exit(1)
	}

	// 2. Fetch (unless --no-fetch).
	if !noFetch {
		if !isMachine {
			fmt.Println(constants.MsgLatestBranchFetching)
		}
		if err := gitutil.FetchAllPrune(); err != nil && !isMachine {
			fmt.Fprintf(os.Stderr, "  Warning: fetch failed: %v\n", err)
		}
	}

	// 3. List remote branches.
	refs, err := gitutil.ListRemoteBranches()
	if err != nil || len(refs) == 0 {
		if allRemotes {
			fmt.Fprintln(os.Stderr, constants.ErrLatestBranchNoRefsAll)
		} else {
			fmt.Fprintf(os.Stderr, constants.ErrLatestBranchNoRefs, remote)
		}
		os.Exit(1)
	}

	// 4. Filter by remote.
	if !allRemotes {
		refs = gitutil.FilterByRemote(refs, remote)
		if len(refs) == 0 {
			fmt.Fprintf(os.Stderr, constants.ErrLatestBranchNoRefs, remote)
			os.Exit(1)
		}
	}

	// 5. Read tip commits.
	items, err := gitutil.ReadBranchTips(refs)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLatestBranchNoCommits+"\n")
		os.Exit(1)
	}

	// 6. Sort by date descending.
	gitutil.SortByDateDesc(items)

	// 7. Pick latest.
	latest := items[0]
	selectedRemote := latest.RemoteRef
	if idx := strings.Index(selectedRemote, "/"); idx >= 0 {
		selectedRemote = selectedRemote[:idx]
	}

	// 8. Resolve branch name(s).
	branchNames := gitutil.ResolvePointsAt(latest.Sha, selectedRemote)

	// 9. Contains fallback.
	if len(branchNames) == 0 && containsFallback {
		branchNames = gitutil.ResolveContains(latest.Sha, selectedRemote)
	}
	if len(branchNames) == 0 {
		branchNames = []string{"<unknown>"}
	}

	shortSha := latest.Sha
	if len(shortSha) > 7 {
		shortSha = shortSha[:7]
	}
	commitDate := latest.CommitDate.Format("2006-01-02T15:04:05-07:00")

	switch format {
	case constants.OutputJSON:
		printLatestBranchJSON(branchNames, selectedRemote, shortSha, commitDate, latest, items, top)
	case constants.OutputCSV:
		printLatestBranchCSV(items, selectedRemote, top)
	default:
		printLatestBranchTerminal(branchNames, selectedRemote, shortSha, commitDate, latest, items, top)
	}
}

// printLatestBranchJSON outputs JSON to stdout.
func printLatestBranchJSON(branchNames []string, remote, sha, commitDate string, latest gitutil.RemoteBranchInfo, items []gitutil.RemoteBranchInfo, top int) {
	out := latestBranchJSON{
		Branch:     branchNames,
		Remote:     remote,
		Sha:        sha,
		CommitDate: commitDate,
		Subject:    latest.Subject,
		Ref:        latest.RemoteRef,
	}
	if top > 0 {
		count := top
		if count > len(items) {
			count = len(items)
		}
		for _, item := range items[:count] {
			rName := stripRemotePrefix(item.RemoteRef)
			out.Top = append(out.Top, latestBranchTopItem{
				Branch:     rName,
				Sha:        truncSha(item.Sha),
				CommitDate: item.CommitDate.Format("2006-01-02T15:04:05-07:00"),
				Subject:    item.Subject,
			})
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", constants.JSONIndent)
	enc.Encode(out)
}

// printLatestBranchCSV outputs CSV to stdout.
func printLatestBranchCSV(items []gitutil.RemoteBranchInfo, remote string, top int) {
	count := 1
	if top > 0 {
		count = top
	}
	if count > len(items) {
		count = len(items)
	}

	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"branch", "remote", "sha", "commitDate", "subject", "ref"})
	for _, item := range items[:count] {
		rName := stripRemotePrefix(item.RemoteRef)
		w.Write([]string{
			rName,
			remote,
			truncSha(item.Sha),
			item.CommitDate.Format("2006-01-02T15:04:05-07:00"),
			item.Subject,
			item.RemoteRef,
		})
	}
	w.Flush()
}

// printLatestBranchTerminal outputs human-readable text to stdout.
func printLatestBranchTerminal(branchNames []string, remote, sha, commitDate string, latest gitutil.RemoteBranchInfo, items []gitutil.RemoteBranchInfo, top int) {
	fmt.Println()
	fmt.Printf("  Latest branch: %s\n", strings.Join(branchNames, ", "))
	fmt.Printf("  Remote:        %s\n", remote)
	fmt.Printf("  SHA:           %s\n", sha)
	fmt.Printf("  Commit date:   %s\n", commitDate)
	fmt.Printf("  Subject:       %s\n", latest.Subject)
	fmt.Printf("  Ref:           %s\n", latest.RemoteRef)

	if top > 0 {
		count := top
		if count > len(items) {
			count = len(items)
		}
		fmt.Println()
		fmt.Printf("  Top %d most recently updated remote branches (%s):\n", count, remote)
		fmt.Printf("  %-30s %-30s %-9s %s\n", "DATE", "BRANCH", "SHA", "SUBJECT")
		for _, item := range items[:count] {
			fmt.Printf("  %-30s %-30s %-9s %s\n",
				item.CommitDate.Format("2006-01-02T15:04:05-07:00"),
				stripRemotePrefix(item.RemoteRef), truncSha(item.Sha), item.Subject)
		}
	}
	fmt.Println()
}

// stripRemotePrefix removes the "<remote>/" prefix from a ref name.
func stripRemotePrefix(ref string) string {
	if idx := strings.Index(ref, "/"); idx >= 0 {
		return ref[idx+1:]
	}
	return ref
}

// truncSha returns the first 7 characters of a SHA.
func truncSha(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

// parseLatestBranchFlags parses flags for the latest-branch command.
// Supports positional integer shorthand: `gitmap lb 3` == `gitmap lb --top 3`.
// Supports --format (terminal|json|csv) and --json as shorthand for --format json.
func parseLatestBranchFlags(args []string) (remote string, allRemotes, containsFallback bool, top int, format string) {
	fs := flag.NewFlagSet(constants.CmdLatestBranch, flag.ExitOnError)
	remoteFlag := fs.String("remote", "origin", constants.FlagDescLBRemote)
	allRemotesFlag := fs.Bool("all-remotes", false, constants.FlagDescLBAllRemotes)
	containsFlag := fs.Bool("contains-fallback", false, constants.FlagDescLBContains)
	topFlag := fs.Int("top", 0, constants.FlagDescLBTop)
	formatFlag := fs.String("format", constants.OutputTerminal, constants.FlagDescLBFormat)
	jsonFlag := fs.Bool("json", false, constants.FlagDescLBJSON)
	fs.Parse(args)

	// Positional integer shorthand for --top.
	if *topFlag == 0 && fs.NArg() > 0 {
		if n, err := strconv.Atoi(fs.Arg(0)); err == nil && n > 0 {
			*topFlag = n
		}
	}

	// --json is shorthand for --format json.
	outFormat := *formatFlag
	if *jsonFlag {
		outFormat = constants.OutputJSON
	}

	return *remoteFlag, *allRemotesFlag, *containsFlag, *topFlag, outFormat
}
