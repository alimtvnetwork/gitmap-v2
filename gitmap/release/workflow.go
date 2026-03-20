// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)


// Options holds all parameters for a release operation.
type Options struct {
	Version       string
	Assets        string
	Commit        string
	Branch        string
	Bump          string
	Notes         string
	Targets       string
	ConfigTargets []model.ReleaseTarget
	ZipGroups     []string
	ZipItems      []string
	BundleName    string
	Draft         bool
	DryRun        bool
	Verbose       bool
	Compress      bool
	Checksums     bool
	NoAssets      bool
	NoCommit      bool
	SkipMeta      bool
}

// Result holds the outcome of a release operation.
type Result struct {
	Version    Version
	BranchName string
	Tag        string
	Commit     string
	Source     string
	Assets     []string
}

// Execute runs the full release workflow.
func Execute(opts Options) error {
	version, err := resolveVersion(opts)
	if err != nil {
		return err
	}

	err = checkDuplicate(version)
	if err != nil {
		return err
	}

	sourceRef, sourceName, err := ResolveSourceRef(opts.Commit, opts.Branch)
	if err != nil {
		return err
	}

	if len(opts.Notes) > 0 {
		fmt.Printf(constants.MsgReleaseNotes, opts.Notes)
	}

	return performRelease(version, sourceRef, sourceName, opts)
}

// resolveVersion determines the version from CLI args, bump, or file.
func resolveVersion(opts Options) (Version, error) {
	if len(opts.Version) > 0 {
		return Parse(opts.Version)
	}
	if len(opts.Bump) > 0 {
		return resolveBump(opts.Bump)
	}

	return resolveFromFile()
}

// resolveBump reads latest.json or falls back to git tags, then increments.
func resolveBump(level string) (Version, error) {
	current, err := resolveLatestVersion()
	if err != nil {
		return Version{}, err
	}

	bumped, err := Bump(current, level)
	if err != nil {
		return Version{}, err
	}

	fmt.Printf(constants.MsgReleaseBumpResult, current.String(), bumped.String())

	return bumped, nil
}

// resolveLatestVersion tries latest.json first, then falls back to git tags.
func resolveLatestVersion() (Version, error) {
	latest, err := ReadLatest()
	if err == nil {
		v, parseErr := Parse(latest.Tag)
		if parseErr == nil {
			return v, nil
		}
	}

	return latestFromGitTags()
}

// resolveFromFile reads version.json.
func resolveFromFile() (Version, error) {
	raw, err := ReadVersionFile()
	if err != nil {
		return Version{}, fmt.Errorf(constants.ErrReleaseVersionRequired)
	}

	fmt.Printf(constants.MsgReleaseVersionRead, constants.DefaultVersionFile, raw)

	return Parse(raw)
}

// checkDuplicate verifies the version hasn't been released.
// If a release JSON exists but no tag or branch, prompts to remove it.
func checkDuplicate(v Version) error {
	if ReleaseExists(v) {
		tagExists := TagExistsLocally(v.String()) || TagExistsRemote(v.String())
		branchName := constants.ReleaseBranchPrefix + v.String()
		branchExists := BranchExists(branchName)

		if !tagExists && !branchExists {
			return handleOrphanedMeta(v)
		}

		return fmt.Errorf(constants.ErrReleaseAlreadyExists, v.String(), v.String())
	}
	if TagExistsLocally(v.String()) || TagExistsRemote(v.String()) {
		return fmt.Errorf(constants.ErrReleaseTagExists, v.String())
	}

	return nil
}

// handleOrphanedMeta detects a release JSON with no matching tag/branch
// and prompts the user to remove it before proceeding.
func handleOrphanedMeta(v Version) error {
	fmt.Printf(constants.MsgReleaseOrphanedMeta, v.String())
	fmt.Print(constants.MsgReleaseOrphanedPrompt)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return fmt.Errorf(constants.ErrReleaseAlreadyExists, v.String(), v.String())
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer != "y" && answer != "yes" {
		return fmt.Errorf(constants.ErrReleaseAborted)
	}

	filename := v.String() + constants.ExtJSON
	path := filepath.Join(constants.DefaultReleaseDir, filename)

	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf(constants.ErrReleaseOrphanedRemove, err)
	}

	fmt.Printf(constants.MsgReleaseOrphanedRemoved, v.String())

	return nil
}

// performRelease executes the branch/tag/push/release steps.
func performRelease(v Version, sourceRef, sourceName string, opts Options) error {
	branchName := constants.ReleaseBranchPrefix + v.String()
	tag := v.String()

	originalBranch, _ := CurrentBranchName()

	fmt.Printf(constants.MsgReleaseStart, tag)

	if opts.DryRun {
		err := printDryRun(v, branchName, tag, sourceName, opts)
		if len(originalBranch) > 0 {
			fmt.Printf(constants.MsgReleaseDryRun, "Switch back to "+originalBranch)
		}

		return err
	}

	// Step 1: Write metadata JSON on the current branch before branching.
	err := writeMetadata(v, branchName, tag, sourceName, nil, opts)
	if err != nil {
		return err
	}

	// Step 2: Create the release branch, commit metadata, tag, push.
	err = executeSteps(v, branchName, tag, sourceRef, sourceName, opts)
	if err != nil {
		return err
	}

	err = returnToBranch(originalBranch)
	if err != nil {
		return err
	}

	if !opts.NoCommit {
		AutoCommit(v.String(), false)
	} else {
		fmt.Print(constants.MsgAutoCommitSkipped)
	}

	return nil
}

// executeSteps runs each release step in sequence.
func executeSteps(v Version, branchName, tag, sourceRef, sourceName string, opts Options) error {
	err := CreateBranch(branchName, sourceRef)
	if err != nil {
		return fmt.Errorf("create branch: %w", err)
	}
	fmt.Printf(constants.MsgReleaseBranch, branchName)

	// Commit the release metadata JSON on the release branch.
	err = commitReleaseMeta(v)
	if err != nil {
		return fmt.Errorf("commit release metadata: %w", err)
	}
	fmt.Printf(constants.MsgReleaseMetaCommitted, branchName)

	err = CreateTag(tag, resolveTagMessage(tag, opts))
	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	return pushAndFinalize(v, branchName, tag, sourceName, opts)
}

// commitReleaseMeta stages and commits .release/ metadata on the current branch.
func commitReleaseMeta(v Version) error {
	metaPath := filepath.Join(constants.DefaultReleaseDir, v.String()+constants.ExtJSON)
	latestPath := filepath.Join(constants.DefaultReleaseDir, constants.DefaultLatestFile)

	files := []string{metaPath}
	if _, err := os.Stat(latestPath); err == nil {
		files = append(files, latestPath)
	}

	err := stageFiles(files)
	if err != nil {
		return err
	}

	return commitStaged(fmt.Sprintf(constants.AutoCommitMsgFmt, v.String()))
}

// resolveTagMessage returns the tag annotation message, using notes if provided.
func resolveTagMessage(tag string, opts Options) string {
	if len(opts.Notes) > 0 {
		return opts.Notes
	}

	return constants.ReleaseTagPrefix + tag
}

