package release

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
)

// ExecuteFromBranch runs the release workflow from an existing release branch.
func ExecuteFromBranch(branchName, assetsPath string, draft bool, dryRun bool, noCommit bool) error {
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

	return completeBranchRelease(version, branchName, assetsPath, draft, noCommit)
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
func completeBranchRelease(v Version, branchName, assetsPath string, draft bool, noCommit bool) error {
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

	err = returnToBranch(originalBranch)
	if err != nil {
		return err
	}

	if !noCommit {
		AutoCommit(v.String(), false)
	} else {
		fmt.Print(constants.MsgAutoCommitSkipped)
	}

	return nil
}

// ExecutePending finds all release branches without tags and releases them.
// Also discovers unreleased versions from .release/v*.json metadata files.
func ExecutePending(assetsPath string, draft bool, dryRun bool, noCommit bool) error {
	branches, err := listReleaseBranches()
	if err != nil {
		return fmt.Errorf("could not list release branches: %w", err)
	}

	pending := filterPendingBranches(branches)
	metaPending := discoverMetadataPending(pending)

	total := len(pending) + len(metaPending)
	if total == 0 {
		fmt.Println(constants.MsgReleasePendingNone)

		return nil
	}

	fmt.Printf(constants.MsgReleasePendingFound, total)

	if len(metaPending) > 0 {
		fmt.Printf(constants.MsgPendingMetaFound, len(metaPending))
	}

	err = releasePendingBranches(pending, assetsPath, draft, dryRun, noCommit)
	if err != nil {
		return err
	}

	return releasePendingFromMetadata(metaPending, assetsPath, draft, dryRun)
}

// discoverMetadataPending finds .release/v*.json files where neither
// the Git branch nor the tag exists. Skips versions already in pendingBranches.
func discoverMetadataPending(pendingBranches []string) []ReleaseMeta {
	metaFiles, err := ListReleaseMetaFiles()
	if err != nil {
		return nil
	}

	branchSet := make(map[string]bool, len(pendingBranches))
	for _, b := range pendingBranches {
		branchSet[b] = true
	}

	var pending []ReleaseMeta

	for _, meta := range metaFiles {
		if isMetaPending(meta, branchSet) {
			pending = append(pending, meta)
		}
	}

	return pending
}

// isMetaPending returns true when the metadata version has no branch/tag
// and is not already covered by a pending branch.
func isMetaPending(meta ReleaseMeta, branchSet map[string]bool) bool {
	if len(meta.Commit) == 0 {
		return false
	}

	branchName := constants.ReleaseBranchPrefix + meta.Tag
	if branchSet[branchName] {
		return false
	}
	if BranchExists(branchName) {
		return false
	}
	if TagExistsLocally(meta.Tag) || TagExistsRemote(meta.Tag) {
		return false
	}

	return true
}

// releasePendingFromMetadata creates branch+tag from stored commit SHA.
func releasePendingFromMetadata(pending []ReleaseMeta, assetsPath string, draft bool, dryRun bool) error {
	for _, meta := range pending {
		err := releaseFromMetadata(meta, assetsPath, draft, dryRun)
		if err != nil {
			fmt.Printf(constants.MsgReleasePendingFailed, meta.Tag, err)
			continue
		}
	}

	return nil
}

// releaseFromMetadata creates a release branch+tag from a metadata file's commit SHA.
func releaseFromMetadata(meta ReleaseMeta, assetsPath string, draft bool, dryRun bool) error {
	v, err := Parse(meta.Tag)
	if err != nil {
		return fmt.Errorf("invalid version in metadata: %s", meta.Tag)
	}

	if !CommitExists(meta.Commit) {
		fmt.Printf(constants.WarnPendingMetaNoCommit, meta.Tag, meta.Commit)

		return nil
	}

	shortSHA := meta.Commit
	if len(shortSHA) > 7 {
		shortSHA = shortSHA[:7]
	}

	fmt.Printf(constants.MsgPendingMetaRelease, meta.Tag, shortSHA)

	if dryRun {
		branchName := constants.ReleaseBranchPrefix + v.String()
		fmt.Printf(constants.MsgReleaseDryRun, "Create branch "+branchName+" from commit "+shortSHA)
		fmt.Printf(constants.MsgReleaseDryRun, "Create tag "+v.String())
		fmt.Printf(constants.MsgReleaseDryRun, "Push branch and tag to origin")

		return nil
	}

	branchName := constants.ReleaseBranchPrefix + v.String()

	err = CreateBranch(branchName, meta.Commit)
	if err != nil {
		return fmt.Errorf("create branch from metadata: %w", err)
	}
	fmt.Printf(constants.MsgReleaseBranch, branchName)

	tag := v.String()
	err = CreateTag(tag, constants.ReleaseTagPrefix+tag)
	if err != nil {
		return fmt.Errorf("create tag from metadata: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	opts := Options{Assets: assetsPath, Draft: draft}

	return pushAndFinalize(v, branchName, tag, "metadata:"+meta.Commit, opts)
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
