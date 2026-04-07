package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/cloner"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/desktop"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/verbose"
)

// applySSHKey sets GIT_SSH_COMMAND if an SSH key name is provided.
func applySSHKey(name string) {
	if len(name) == 0 {
		return
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)
		os.Exit(1)
	}

	sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", key.PrivatePath)
	os.Setenv("GIT_SSH_COMMAND", sshCmd)
	fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)
}

// runClone handles the "clone" subcommand.
func runClone(args []string) {
	checkHelp("clone", args)
	source, targetDir, sshKeyName, safePull, ghDesktop, verboseMode := parseCloneFlags(args)
	if len(source) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrSourceRequired)
		fmt.Fprintln(os.Stderr, constants.ErrCloneUsage)
		os.Exit(1)
	}
	initCloneVerbose(verboseMode)
	requireOnline()
	applySSHKey(sshKeyName)
	source = resolveCloneShorthand(source)
	executeClone(source, targetDir, safePull, ghDesktop)
}

// initCloneVerbose sets up verbose logging if enabled.
func initCloneVerbose(enabled bool) {
	if enabled {
		log, err := verbose.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

			return
		}
		defer log.Close()
	}
}

// resolveCloneShorthand maps "json", "csv", and "text" to default output paths.
func resolveCloneShorthand(source string) string {
	shorthandMap := map[string]string{
		constants.ShorthandJSON: filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
		constants.ShorthandCSV:  filepath.Join(constants.DefaultOutputFolder, constants.DefaultCSVFile),
		constants.ShorthandText: filepath.Join(constants.DefaultOutputFolder, constants.DefaultTextFile),
	}
	resolved, ok := shorthandMap[strings.ToLower(source)]
	if ok {
		return validateShorthandPath(resolved)
	}

	return source
}

// validateShorthandPath checks that the resolved shorthand file exists.
func validateShorthandPath(resolved string) string {
	_, err := os.Stat(resolved)
	if err == nil {
		return resolved
	}
	fmt.Fprintf(os.Stderr, constants.ErrShorthandNotFound, resolved)
	os.Exit(1)

	return ""
}

// executeClone runs the clone operation and prints the summary.
func executeClone(source, targetDir string, safePull, ghDesktop bool) {
	summary, err := cloner.CloneFromFile(source, targetDir, safePull)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFailed, source, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgCloneComplete, summary.Succeeded, summary.Failed)
	printCloneFailures(summary)
	registerCloned(summary, targetDir, ghDesktop)
}

// printCloneFailures lists any repos that failed to clone.
func printCloneFailures(s model.CloneSummary) {
	if s.Failed == 0 {
		return
	}

	fmt.Println(constants.MsgFailedClones)
	for _, e := range s.Errors {
		fmt.Printf(constants.MsgFailedEntry,
			e.Record.RepoName, e.Record.RelativePath, e.Error)
	}
}

// registerCloned adds successfully cloned repos to GitHub Desktop.
func registerCloned(s model.CloneSummary, targetDir string, enabled bool) {
	if enabled {
		absTarget, _ := filepath.Abs(targetDir)
		records := make([]model.ScanRecord, 0, s.Succeeded)
		for _, r := range s.Cloned {
			r.Record.AbsolutePath = filepath.Join(absTarget, r.Record.RelativePath)
			records = append(records, r.Record)
		}
		result := desktop.AddRepos(records)
		fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
	}
}
