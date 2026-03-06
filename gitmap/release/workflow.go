// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"

	"github.com/user/gitmap/constants"
)


// Options holds all parameters for a release operation.
type Options struct {
	Version    string
	Assets     string
	Commit     string
	Branch     string
	Bump       string
	Draft      bool
	DryRun     bool
	Verbose    bool
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
func checkDuplicate(v Version) error {
	if ReleaseExists(v) {
		return fmt.Errorf(constants.ErrReleaseAlreadyExists, v.String(), v.String())
	}
	if TagExistsLocally(v.String()) || TagExistsRemote(v.String()) {
		return fmt.Errorf(constants.ErrReleaseTagExists, v.String())
	}

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

	err := executeSteps(v, branchName, tag, sourceRef, sourceName, opts)
	if err != nil {
		return err
	}

	return returnToBranch(originalBranch)
}

// executeSteps runs each release step in sequence.
func executeSteps(v Version, branchName, tag, sourceRef, sourceName string, opts Options) error {
	err := CreateBranch(branchName, sourceRef)
	if err != nil {
		return fmt.Errorf("create branch: %w", err)
	}
	fmt.Printf(constants.MsgReleaseBranch, branchName)

	err = CreateTag(tag, constants.ReleaseTagPrefix+tag)
	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	return pushAndFinalize(v, branchName, tag, sourceName, opts)
}

