// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

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

	err = CreateTag(tag, "Release "+tag)
	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	return pushAndFinalize(v, branchName, tag, sourceName, opts)
}

// pushAndFinalize pushes to remote and writes metadata.
func pushAndFinalize(v Version, branchName, tag, sourceName string, opts Options) error {
	err := PushBranchAndTag(branchName, tag)
	if err != nil {
		return fmt.Errorf(constants.ErrReleasePushFailed, err)
	}
	fmt.Print(constants.MsgReleasePushed)

	assets := CollectAssets(opts.Assets)
	for _, a := range assets {
		fmt.Printf(constants.MsgReleaseAttach, a)
	}

	return writeMetadata(v, branchName, tag, sourceName, assets, opts.Draft)
}

// writeMetadata persists release info and updates latest.
func writeMetadata(v Version, branchName, tag, sourceName string, assets []string, draft bool) error {
	commit, _ := CurrentCommitSHA()
	assetPaths := make([]string, len(assets))
	copy(assetPaths, assets)

	meta := ReleaseMeta{
		Version:      v.CoreString(),
		Branch:       branchName,
		SourceBranch: sourceName,
		Commit:       commit,
		Tag:          tag,
		Assets:       assetPaths,
		Draft:        draft,
		PreRelease:   v.IsPreRelease(),
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		IsLatest:     false,
	}

	err := WriteReleaseMeta(meta)
	if err != nil {
		return fmt.Errorf(constants.ErrReleaseMetaWrite, err)
	}
	fmt.Printf(constants.MsgReleaseMeta, constants.DefaultReleaseDir+"/"+v.String()+constants.ExtJSON)

	return updateLatestIfStable(v)
}

// updateLatestIfStable marks the release as latest if stable.
func updateLatestIfStable(v Version) error {
	if v.IsPreRelease() {
		fmt.Printf(constants.MsgReleaseComplete, v.String())
		return nil
	}

	err := WriteLatest(v)
	if err != nil {
		return err
	}

	fmt.Printf(constants.MsgReleaseLatest, v.String())
	fmt.Printf(constants.MsgReleaseComplete, v.String())

	return nil
}

// printDryRun shows what would happen without executing.
func printDryRun(v Version, branchName, tag, sourceName string, opts Options) error {
	fmt.Printf(constants.MsgReleaseDryRun, "Create branch "+branchName+" from "+sourceName)
	fmt.Printf(constants.MsgReleaseDryRun, "Create tag "+tag)
	fmt.Printf(constants.MsgReleaseDryRun, "Push branch and tag to origin")

	body := DetectChangelog()
	if len(body) > 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Use CHANGELOG.md as release body")
	}
	readme := DetectReadme()
	if len(readme) > 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Attach README.md")
	}
	userAssets := CollectAssets(opts.Assets)
	for _, a := range userAssets {
		fmt.Printf(constants.MsgReleaseDryRun, "Attach "+a)
	}

	fmt.Printf(constants.MsgReleaseDryRun, "Write metadata to "+constants.DefaultReleaseDir+"/"+v.String()+constants.ExtJSON)

	if v.IsPreRelease() == false {
		fmt.Printf(constants.MsgReleaseDryRun, "Mark "+v.String()+" as latest")
	}

	fmt.Printf(constants.MsgReleaseComplete, v.String())

	return nil
}

// ExecuteFromBranch runs the release workflow from an existing release branch.
func ExecuteFromBranch(branchName, assetsPath string, draft bool, dryRun bool) error {
	version, err := extractVersionFromBranch(branchName)
	if err != nil {
		return err
	}

	err = validateExistingBranch(branchName, version)
	if err != nil {
		return err
	}

	fmt.Printf(constants.MsgReleaseBranchStart, branchName)

	if dryRun {
		return printDryRun(version, branchName, version.String(), branchName, Options{
			Assets: assetsPath, Draft: draft, DryRun: true,
		})
	}

	return completeBranchRelease(version, branchName, assetsPath, draft)
}

// extractVersionFromBranch parses the version from a release branch name.
func extractVersionFromBranch(branchName string) (Version, error) {
	prefix := constants.ReleaseBranchPrefix
	if len(branchName) <= len(prefix) {
		return Version{}, fmt.Errorf(constants.ErrReleaseInvalidVersion, branchName)
	}

	versionStr := branchName[len(prefix):]

	return Parse(versionStr)
}

// validateExistingBranch checks the branch exists and tag doesn't.
func validateExistingBranch(branchName string, v Version) error {
	if BranchExists(branchName) == false {
		return fmt.Errorf(constants.ErrReleaseBranchNotFound, branchName)
	}

	return checkDuplicate(v)
}

// completeBranchRelease checks out the branch and runs tag/push/release.
func completeBranchRelease(v Version, branchName, assetsPath string, draft bool) error {
	originalBranch, _ := CurrentBranchName()

	err := CheckoutBranch(branchName)
	if err != nil {
		return fmt.Errorf("checkout branch: %w", err)
	}

	tag := v.String()
	err = CreateTag(tag, "Release "+tag)
	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	opts := Options{Assets: assetsPath, Draft: draft}

	err = pushAndFinalize(v, branchName, tag, branchName, opts)
	if err != nil {
		return err
	}

	return returnToBranch(originalBranch)
}
