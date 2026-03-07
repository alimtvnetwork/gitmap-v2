package release

import (
	"fmt"
	"time"

	"github.com/user/gitmap/constants"
)

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
	meta := buildReleaseMeta(v, branchName, tag, sourceName, commit, assets, draft)

	err := WriteReleaseMeta(meta)
	if err != nil {
		return fmt.Errorf(constants.ErrReleaseMetaWrite, err)
	}
	fmt.Printf(constants.MsgReleaseMeta, constants.DefaultReleaseDir+"/"+v.String()+constants.ExtJSON)

	return updateLatestIfStable(v)
}

// buildReleaseMeta constructs the metadata struct for a release.
func buildReleaseMeta(v Version, branchName, tag, sourceName, commit string, assets []string, draft bool) ReleaseMeta {
	assetPaths := make([]string, len(assets))
	copy(assetPaths, assets)

	return ReleaseMeta{
		Version:      v.CoreString(),
		Branch:       branchName,
		SourceBranch: sourceName,
		Commit:       commit,
		Tag:          tag,
		Assets:       assetPaths,
		Changelog:    loadChangelogNotes(v.String()),
		Draft:        draft,
		PreRelease:   v.IsPreRelease(),
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		IsLatest:     false,
	}
}

// loadChangelogNotes reads changelog notes for a version, returning nil on error.
func loadChangelogNotes(version string) []string {
	entries, err := ReadChangelog()
	if err != nil {
		return nil
	}

	entry, found := FindChangelogEntry(entries, version)
	if found {
		return entry.Notes
	}

	return nil
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
	printDryRunSteps(branchName, tag, sourceName)
	printDryRunAssets(opts.Assets)
	printDryRunMeta(v)
	fmt.Printf(constants.MsgReleaseComplete, v.String())

	return nil
}

// printDryRunSteps prints branch/tag/push dry-run lines.
func printDryRunSteps(branchName, tag, sourceName string) {
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
}

// printDryRunAssets prints asset attachments in dry-run mode.
func printDryRunAssets(assetsPath string) {
	userAssets := CollectAssets(assetsPath)
	for _, a := range userAssets {
		fmt.Printf(constants.MsgReleaseDryRun, "Attach "+a)
	}
}

// printDryRunMeta prints metadata and latest marker in dry-run mode.
func printDryRunMeta(v Version) {
	fmt.Printf(constants.MsgReleaseDryRun, "Write metadata to "+constants.DefaultReleaseDir+"/"+v.String()+constants.ExtJSON)

	if len(v.PreRelease) == 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Mark "+v.String()+" as latest")
	}
}

// returnToBranch switches back to the original branch after a release.
func returnToBranch(branch string) error {
	if len(branch) == 0 {
		return nil
	}

	err := CheckoutBranch(branch)
	if err != nil {
		return fmt.Errorf("switch back to %s: %w", branch, err)
	}

	fmt.Printf(constants.MsgReleaseSwitchedBack, branch)

	return nil
}
