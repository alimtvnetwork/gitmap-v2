// Package cmd — latest-branch command handler.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/gitutil"
)

// runLatestBranch handles the 'latest-branch' / 'lb' command.
func runLatestBranch(args []string) {
	remote, allRemotes, containsFallback, top := parseLatestBranchFlags(args)

	// 1. Validate git repo.
	if !gitutil.IsInsideWorkTree() {
		fmt.Fprintln(os.Stderr, constants.ErrLatestBranchNotRepo)
		os.Exit(1)
	}

	// 2. Fetch.
	fmt.Println(constants.MsgLatestBranchFetching)
	if err := gitutil.FetchAllPrune(); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: fetch failed: %v\n", err)
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

	// 10. Display.
	shortSha := latest.Sha
	if len(shortSha) > 7 {
		shortSha = shortSha[:7]
	}
	fmt.Println()
	fmt.Printf("  Latest branch: %s\n", strings.Join(branchNames, ", "))
	fmt.Printf("  Remote:        %s\n", selectedRemote)
	fmt.Printf("  SHA:           %s\n", shortSha)
	fmt.Printf("  Commit date:   %s\n", latest.CommitDate.Format("2006-01-02T15:04:05-07:00"))
	fmt.Printf("  Subject:       %s\n", latest.Subject)
	fmt.Printf("  Ref:           %s\n", latest.RemoteRef)

	// 11. Top N.
	if top > 0 {
		count := top
		if count > len(items) {
			count = len(items)
		}
		fmt.Println()
		fmt.Printf("  Top %d most recently updated remote branches (%s):\n", count, selectedRemote)
		fmt.Printf("  %-30s %-30s %-9s %s\n", "DATE", "BRANCH", "SHA", "SUBJECT")
		for _, item := range items[:count] {
			rName := item.RemoteRef
			if idx := strings.Index(rName, "/"); idx >= 0 {
				rName = rName[idx+1:]
			}
			sha := item.Sha
			if len(sha) > 7 {
				sha = sha[:7]
			}
			fmt.Printf("  %-30s %-30s %-9s %s\n",
				item.CommitDate.Format("2006-01-02T15:04:05-07:00"),
				rName, sha, item.Subject)
		}
	}
	fmt.Println()
}

// parseLatestBranchFlags parses flags for the latest-branch command.
func parseLatestBranchFlags(args []string) (remote string, allRemotes, containsFallback bool, top int) {
	fs := flag.NewFlagSet(constants.CmdLatestBranch, flag.ExitOnError)
	remoteFlag := fs.String("remote", "origin", constants.FlagDescLBRemote)
	allRemotesFlag := fs.Bool("all-remotes", false, constants.FlagDescLBAllRemotes)
	containsFlag := fs.Bool("contains-fallback", false, constants.FlagDescLBContains)
	topFlag := fs.Int("top", 0, constants.FlagDescLBTop)
	fs.Parse(args)

	return *remoteFlag, *allRemotesFlag, *containsFlag, *topFlag
}
