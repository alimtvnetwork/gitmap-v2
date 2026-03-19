package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/completion"
	"github.com/user/gitmap/constants"
)

// runCompletion handles the "completion" subcommand.
func runCompletion(args []string) {
	checkHelp("completion", args)

	if hasListFlag(args) {
		handleCompletionList(args)

		return
	}

	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrCompUsage)
		os.Exit(1)
	}

	printCompletionScript(args[0])
}

// hasListFlag checks if any --list-* flag is present.
func hasListFlag(args []string) bool {
	for _, a := range args {
		if a == constants.CompListRepos || a == constants.CompListGroups ||
			a == constants.CompListCommands || a == constants.CompListAliases ||
			a == constants.CompListZipGroups {
			return true
		}
	}

	return false
}

// handleCompletionList routes to the appropriate list printer.
func handleCompletionList(args []string) {
	for _, a := range args {
		if a == constants.CompListRepos {
			printCompletionRepos()

			return
		}
		if a == constants.CompListGroups {
			printCompletionGroups()

			return
		}
		if a == constants.CompListCommands {
			printCompletionCommands()

			return
		}
	}
}

// printCompletionRepos prints all repo slugs, one per line.
func printCompletionRepos() {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	repos, err := db.ListRepos()
	if err != nil {
		return
	}

	for _, r := range repos {
		fmt.Println(r.Slug)
	}
}

// printCompletionGroups prints all group names, one per line.
func printCompletionGroups() {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	groups, err := db.ListGroups()
	if err != nil {
		return
	}

	for _, g := range groups {
		fmt.Println(g.Name)
	}
}

// printCompletionCommands prints all command names, one per line.
func printCompletionCommands() {
	for _, cmd := range completion.AllCommands() {
		fmt.Println(cmd)
	}
}

// printCompletionScript outputs the shell completion script.
func printCompletionScript(shell string) {
	script, err := completion.Generate(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCompUnknownShell, shell)
		os.Exit(1)
	}

	fmt.Print(script)
}
