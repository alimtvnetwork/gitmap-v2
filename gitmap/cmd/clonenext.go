package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/gitmap/clonenext"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/desktop"
	"github.com/user/gitmap/gitutil"
	"github.com/user/gitmap/lockcheck"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/verbose"
)

// runCloneNext handles the "clone-next" subcommand.
func runCloneNext(args []string) {
	checkHelp("clone-next", args)
	versionArg, deleteFlag, keepFlag, noDesktop, createRemote, sshKeyName, verboseMode := parseCloneNextFlags(args)
	if len(versionArg) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrCloneNextUsage)
		os.Exit(1)
	}

	if verboseMode {
		log, err := verbose.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)
		} else {
			defer log.Close()
		}
	}

	requireOnline()
	applySSHKey(sshKeyName)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextCwd, err)
		os.Exit(1)
	}

	remoteURL, err := gitutil.RemoteURL(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextNoRemote, err)
		os.Exit(1)
	}

	currentFolder := filepath.Base(cwd)
	parentDir := filepath.Dir(cwd)

	// Strip .git suffix from remote URL for repo name extraction.
	repoName := extractRepoName(remoteURL)

	parsed := clonenext.ParseRepoName(repoName)
	targetVersion, err := clonenext.ResolveTarget(parsed, versionArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextBadVersion, err)
		os.Exit(1)
	}

	targetName := clonenext.TargetRepoName(parsed.BaseName, targetVersion)
	targetURL := clonenext.ReplaceRepoInURL(remoteURL, repoName, targetName)
	targetPath := filepath.Join(parentDir, targetName)

	if _, statErr := os.Stat(targetPath); statErr == nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextExists, targetPath)
		os.Exit(1)
	}

	// Optionally check and create the target GitHub repo when --create-remote is set.
	if createRemote {
		owner, _, parseErr := clonenext.ParseOwnerRepo(remoteURL)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCloneNextRemoteParse, parseErr)
			os.Exit(1)
		}

		exists, checkErr := clonenext.RepoExists(owner, targetName)
		if checkErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCloneNextRepoCheck, checkErr)
			os.Exit(1)
		}

		if !exists {
			fmt.Printf(constants.MsgCloneNextCreating, targetName)
			createErr := clonenext.CreateRepo(owner, targetName, true)
			if createErr != nil {
				fmt.Fprintf(os.Stderr, constants.ErrCloneNextRepoCreate, targetName, createErr)
				os.Exit(1)
			}
			fmt.Printf(constants.MsgCloneNextCreated, targetName)
		}
	}

	fmt.Printf(constants.MsgCloneNextCloning, targetName, parentDir)
	cloneResult := runGitClone(targetURL, targetPath)
	if !cloneResult {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextFailed, targetName)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgCloneNextDone, targetName)

	if !noDesktop {
		registerCloneNextDesktop(targetName, targetPath)
	}

	handleCloneNextRemoval(currentFolder, cwd, targetPath, deleteFlag, keepFlag)
}

// extractRepoName extracts the repository name from a remote URL.
func extractRepoName(remoteURL string) string {
	name := remoteURL
	// Remove trailing .git
	name = strings.TrimSuffix(name, ".git")
	// Get last path segment
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}

	return name
}

// runGitClone executes git clone and returns success status.
func runGitClone(url, dest string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitClone, url, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run() == nil
}

// registerCloneNextDesktop registers the cloned repo with GitHub Desktop.
func registerCloneNextDesktop(name, absPath string) {
	records := []model.ScanRecord{{
		RepoName:     name,
		AbsolutePath: absPath,
	}}
	result := desktop.AddRepos(records)
	if result.Added > 0 {
		fmt.Printf(constants.MsgCloneNextDesktop, name)
	}
}

// handleCloneNextRemoval manages removal of the current version folder.
// It changes to the parent directory first to release file locks on Windows.
func handleCloneNextRemoval(folderName, fullPath, targetPath string, deleteFlag, keepFlag bool) {
	if keepFlag {
		return
	}

	// Move out of the folder before attempting removal to avoid Windows file locks.
	parentDir := filepath.Dir(fullPath)
	if chErr := os.Chdir(parentDir); chErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not cd to %s: %v\n", parentDir, chErr)
	}

	removed := false
	if deleteFlag {
		removeFolder(folderName, fullPath)
		removed = true
	} else {
		// Prompt
		fmt.Printf(constants.MsgCloneNextRemovePrompt, folderName)
		var answer string
		_, _ = fmt.Scanln(&answer)
		if strings.ToLower(strings.TrimSpace(answer)) == "y" {
			removeFolder(folderName, fullPath)
			removed = true
		}
	}

	// After removing the old folder, move into the newly cloned directory.
	if removed {
		if chErr := os.Chdir(targetPath); chErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not cd to %s: %v\n", targetPath, chErr)
		} else {
			fmt.Printf(constants.MsgCloneNextMovedTo, filepath.Base(targetPath))
		}
	}
}

// removeFolder deletes a directory and prints the result.
func removeFolder(name, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnCloneNextRemoveFailed, name, err)

		return
	}
	fmt.Printf(constants.MsgCloneNextRemoved, name)
}
