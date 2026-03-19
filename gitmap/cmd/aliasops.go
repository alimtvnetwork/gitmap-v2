package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// runAliasSet handles "alias set <alias> <slug>".
func runAliasSet(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrAliasEmpty)
		os.Exit(1)
	}

	alias := args[0]
	slug := args[1]

	executeAliasSet(alias, slug)
}

// executeAliasSet resolves the slug and creates or updates the alias.
func executeAliasSet(alias, slug string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	repos, err := db.FindBySlug(slug)
	if err != nil || len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrAliasRepoMissing, slug)
		os.Exit(1)
	}

	repoID := repos[0].ID

	if db.AliasExists(alias) {
		err = db.UpdateAlias(alias, repoID)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
			os.Exit(1)
		}

		fmt.Printf(constants.MsgAliasUpdated, alias, slug)
		printHints(aliasSetHints())

		return
	}

	_, err = db.CreateAlias(alias, repoID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgAliasCreated, alias, slug)
	printHints(aliasSetHints())
}

// runAliasRemove handles "alias remove <alias>".
func runAliasRemove(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrAliasEmpty)
		os.Exit(1)
	}

	alias := args[0]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.DeleteAlias(alias)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgAliasRemoved, alias)
}

// runAliasList handles "alias list".
func runAliasList() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	aliases, err := db.ListAliasesWithRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	printAliasList(aliases)
	printHints(aliasListHints())
}

// printAliasList renders the alias table to stdout.
func printAliasList(aliases []store.AliasWithRepo) {
	if len(aliases) == 0 {
		fmt.Println("  No aliases defined.")

		return
	}

	fmt.Printf(constants.MsgAliasListHeader, len(aliases))

	for _, a := range aliases {
		fmt.Printf(constants.MsgAliasListRow, a.Alias, a.Slug)
	}
}

// runAliasShow handles "alias show <alias>".
func runAliasShow(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrAliasEmpty)
		os.Exit(1)
	}

	alias := args[0]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	resolved, err := db.ResolveAlias(alias)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgAliasResolved, resolved.Alias, resolved.AbsolutePath, resolved.Slug)
}

// runAliasSuggest handles "alias suggest [--apply]".
func runAliasSuggest(args []string) {
	apply := parseAliasSuggestFlags(args)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	repos, err := db.ListUnaliasedRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	if len(repos) == 0 {
		fmt.Println(constants.MsgAliasSuggestNone)

		return
	}

	created := suggestAliases(db, repos, apply)
	fmt.Printf(constants.MsgAliasSuggestDone, created)
}

// parseAliasSuggestFlags parses flags for alias suggest.
func parseAliasSuggestFlags(args []string) bool {
	fs := flag.NewFlagSet(constants.SubCmdAliasSug, flag.ExitOnError)
	apply := fs.Bool("apply", false, constants.FlagDescAliasApply)
	fs.Parse(args)

	return *apply
}

// suggestAliases proposes aliases for unaliased repos.
func suggestAliases(db *store.DB, repos []store.UnaliasedRepo, autoApply bool) int {
	created := 0
	reader := bufio.NewReader(os.Stdin)

	for _, r := range repos {
		suggestion := r.RepoName
		if db.AliasExists(suggestion) {
			continue
		}

		if autoApply {
			createSuggestedAlias(db, suggestion, r.ID)
			created++

			continue
		}

		if promptAliasSuggestion(reader, r.Slug, suggestion) {
			createSuggestedAlias(db, suggestion, r.ID)
			created++
		}
	}

	return created
}

// promptAliasSuggestion asks the user to accept a suggested alias.
func promptAliasSuggestion(reader *bufio.Reader, slug, suggestion string) bool {
	fmt.Printf(constants.MsgAliasSuggest, slug, suggestion)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

// createSuggestedAlias creates an alias and prints confirmation.
func createSuggestedAlias(db *store.DB, alias, repoID string) {
	_, err := db.CreateAlias(alias, repoID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return
	}

	fmt.Printf(constants.MsgAliasCreated, alias, repoID)
}
