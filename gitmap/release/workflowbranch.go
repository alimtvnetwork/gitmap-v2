package release

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
)

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
	if BranchExists(branchName) {
		return checkDuplicate(v)
	}

	return fmt.Errorf(constants.ErrReleaseBranchNotFound, branchName)
}

// completeBranchRelease checks out the branch and runs tag/push/release.
func completeBranchRelease(v Version, branchName, assetsPath string, draft bool) error {
	originalBranch, _ := CurrentBranchName()

	err := CheckoutBranch(branchName)
	if err != nil {
		return fmt.Errorf("checkout branch: %w", err)
	}

	tag := v.String()
	err = CreateTag(tag, constants.ReleaseTagPrefix+tag)
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

// ExecutePending finds all release branches without tags and releases them.
func ExecutePending(assetsPath string, draft bool, dryRun bool) error {
	branches, err := listReleaseBranches()
	if err != nil {
		return fmt.Errorf("could not list release branches: %w", err)
	}

	pending := filterPendingBranches(branches)
	if len(pending) == 0 {
		fmt.Println(constants.MsgReleasePendingNone)

		return nil
	}

	fmt.Printf(constants.MsgReleasePendingFound, len(pending))

	return releasePendingBranches(pending, assetsPath, draft, dryRun)
}

// releasePendingBranches iterates and releases each pending branch.
func releasePendingBranches(pending []string, assetsPath string, draft bool, dryRun bool) error {
	for _, branchName := range pending {
		err := ExecuteFromBranch(branchName, assetsPath, draft, dryRun)
		if err != nil {
			fmt.Printf(constants.MsgReleasePendingFailed, branchName, err)
			continue
		}
	}

	return nil
}

// listReleaseBranches returns all local branches matching release/v*.
func listReleaseBranches() ([]string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, constants.GitBranchListFlag, constants.ReleaseBranchPrefix+constants.GitTagGlob)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseBranchLines(string(out)), nil
}

// parseBranchLines extracts branch names from git branch output.
func parseBranchLines(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var branches []string

	for _, line := range lines {
		branch := strings.TrimSpace(line)
		branch = strings.TrimPrefix(branch, "* ")
		if len(branch) > 0 {
			branches = append(branches, branch)
		}
	}

	return branches
}

// filterPendingBranches returns branches whose version tag does not exist.
func filterPendingBranches(branches []string) []string {
	var pending []string

	for _, branch := range branches {
		if isPendingBranch(branch) {
			pending = append(pending, branch)
		}
	}

	return pending
}

// isPendingBranch returns true when the branch has no released tag.
func isPendingBranch(branch string) bool {
	v, err := extractVersionFromBranch(branch)
	if err != nil {
		return false
	}

	tag := v.String()

	return tagIsMissing(tag)
}

// tagIsMissing returns true when a tag does not exist locally or remotely.
func tagIsMissing(tag string) bool {
	if TagExistsLocally(tag) {
		return false
	}
	if TagExistsRemote(tag) {
		return false
	}

	return true
}
