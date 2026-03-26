package release

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// LastMeta holds the most recent release metadata after Execute completes.
var LastMeta *ReleaseMeta

// lastZipChecksums stores SHA-1 hashes of zip group archives built during the
// current release, keyed by archive filename. Populated by buildZipGroupAssets
// and buildAdHocZipAssets, consumed by buildReleaseMeta.
var lastZipChecksums map[string]string

// pushAndFinalize pushes to remote and writes metadata.
func pushAndFinalize(v Version, branchName, tag, sourceName string, opts Options) error {
	lastZipChecksums = nil

	err := PushBranchAndTag(branchName, tag)
	if err != nil {
		return fmt.Errorf(constants.ErrReleasePushFailed, err)
	}
	fmt.Print(constants.MsgReleasePushed)

	assets := CollectAssets(opts.Assets)

	// Cross-compile Go binaries if applicable.
	goAssets := buildGoAssetsIfApplicable(v, opts)
	assets = append(assets, goAssets...)

	// Build zip group archives (persistent groups from DB).
	zipGroupAssets := buildZipGroupAssets(opts)
	assets = append(assets, zipGroupAssets...)

	// Build ad-hoc zip archives (-Z / --bundle).
	adHocAssets := buildAdHocZipAssets(opts)
	assets = append(assets, adHocAssets...)

	if opts.Compress && len(assets) > 0 {
		compressed, compErr := CompressAssets(assets)
		if compErr == nil && len(compressed) > 0 {
			for _, a := range compressed {
				fmt.Printf(constants.MsgCompressArchive, filepath.Base(a), filepath.Base(a))
			}

			assets = compressed
		}
	}

	if opts.Checksums && len(assets) > 0 {
		checksumPath, csErr := GenerateChecksums(assets)
		if csErr == nil && len(checksumPath) > 0 {
			fmt.Printf(constants.MsgChecksumGenerated, constants.ChecksumsFile)
			assets = append(assets, checksumPath)
		}
	}

	for _, a := range assets {
		fmt.Printf(constants.MsgReleaseAttach, a)
	}

	// Upload to GitHub if token is available.
	uploadToGitHub(v, assets, opts)

	fmt.Printf(constants.MsgReleaseComplete, v.String())

	return nil
}

// buildGoAssetsIfApplicable cross-compiles Go binaries when a Go project is detected.
func buildGoAssetsIfApplicable(v Version, opts Options) []string {
	if opts.NoAssets {
		fmt.Print(constants.MsgAssetSkipped)

		return nil
	}

	if !DetectGoProject() {
		return nil
	}

	modName, err := ReadModuleName()
	if err != nil {
		return nil
	}

	fmt.Printf(constants.MsgAssetDetected, modName)

	packages := FindMainPackages()
	if len(packages) == 0 {
		fmt.Print(constants.MsgAssetNoMain)

		return nil
	}

	targets, err := ResolveTargets(opts.Targets, opts.ConfigTargets)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Invalid targets: %v\n", err)

		return nil
	}

	stagingDir, err := EnsureStagingDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Staging dir: %v\n", err)

		return nil
	}

	fmt.Printf(constants.MsgAssetCrossCompile, len(targets)*len(packages))

	results := CrossCompile(v.String(), targets, packages, stagingDir)
	successful := CollectSuccessfulBuilds(results)

	fmt.Printf(constants.MsgAssetBuildSummary, len(successful), len(results))

	return successful
}

// uploadToGitHub creates a GitHub release and uploads assets.
func uploadToGitHub(v Version, assets []string, opts Options) {
	token := os.Getenv(constants.GitHubTokenEnv)
	if len(token) == 0 {
		if len(assets) > 0 {
			fmt.Fprint(os.Stderr, constants.ErrAssetNoToken)
		}

		return
	}

	owner, repo, err := ParseRemoteOrigin()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAssetRemoteParse, err)

		return
	}

	name := constants.ReleaseTagPrefix + v.String()
	if len(opts.Notes) > 0 {
		name = opts.Notes
	}

	body := DetectChangelog()
	ghRelease, err := CreateGitHubRelease(owner, repo, v.String(), name, body, token, opts.Draft)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ GitHub release creation failed: %v\n", err)

		return
	}

	if len(assets) > 0 {
		fmt.Printf(constants.MsgAssetUploadStart, len(assets))
		UploadAllAssets(owner, repo, ghRelease.ID, assets, token)
	}
}

// writeMetadata persists release info and updates latest.
func writeMetadata(v Version, branchName, tag, sourceName string, assets []string, opts Options) error {
	commit, _ := CurrentCommitSHA()
	meta := buildReleaseMeta(v, branchName, tag, sourceName, commit, assets, opts)

	err := WriteReleaseMeta(meta)
	if err != nil {
		return fmt.Errorf(constants.ErrReleaseMetaWrite, err)
	}
	fmt.Printf(constants.MsgReleaseMeta, constants.DefaultReleaseDir+"/"+v.String()+constants.ExtJSON)

	LastMeta = &meta

	return updateLatestIfStable(v)
}

// buildReleaseMeta constructs the metadata struct for a release.
func buildReleaseMeta(v Version, branchName, tag, sourceName, commit string, assets []string, opts Options) ReleaseMeta {
	assetPaths := make([]string, len(assets))
	copy(assetPaths, assets)

	zipGroups := collectZipGroupNames(opts)

	var checksums map[string]string
	if len(lastZipChecksums) > 0 {
		checksums = make(map[string]string, len(lastZipChecksums))
		for k, v := range lastZipChecksums {
			checksums[k] = v
		}
	}

	return ReleaseMeta{
		Version:           v.CoreString(),
		Branch:            branchName,
		SourceBranch:      sourceName,
		Commit:            commit,
		Tag:               tag,
		Assets:            assetPaths,
		Changelog:         loadChangelogNotes(v.String()),
		Notes:             opts.Notes,
		ZipGroups:         zipGroups,
		ZipGroupChecksums: checksums,
		Draft:             opts.Draft,
		PreRelease:        v.IsPreRelease(),
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
		IsLatest:          false,
	}
}

// collectZipGroupNames merges persistent group names and ad-hoc bundle name.
func collectZipGroupNames(opts Options) []string {
	var names []string
	names = append(names, opts.ZipGroups...)

	if len(opts.BundleName) > 0 {
		names = append(names, opts.BundleName)
	}

	if len(names) == 0 {
		return nil
	}

	return names
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

	if LastMeta != nil {
		LastMeta.IsLatest = true
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
	printDryRunGoAssets(v, opts)
	printDryRunZipGroups(opts)
	printDryRunAssets(opts.Assets, opts.Compress, opts.Checksums)
	printDryRunMeta(v)
	fmt.Printf(constants.MsgReleaseComplete, v.String())

	return nil
}

// printDryRunGoAssets shows Go cross-compile plan in dry-run mode.
func printDryRunGoAssets(v Version, opts Options) {
	if opts.NoAssets {
		fmt.Printf(constants.MsgReleaseDryRun, "Skip Go binary compilation (--no-assets)")

		return
	}

	if !DetectGoProject() {
		return
	}

	binName := resolveBinName()
	targets, err := ResolveTargets(opts.Targets, opts.ConfigTargets)
	if err != nil {
		return
	}

	names := DescribeTargets(binName, v.String(), targets)
	fmt.Printf(constants.MsgAssetDryRunHeader, len(names))

	for _, name := range names {
		fmt.Printf(constants.MsgAssetDryRunBinary, name)
	}

	fmt.Printf(constants.MsgAssetDryRunUpload, len(names))
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
func printDryRunAssets(assetsPath string, compress bool, checksums bool) {
	userAssets := CollectAssets(assetsPath)

	if compress && len(userAssets) > 0 {
		archiveNames := DescribeCompression(userAssets)
		for _, name := range archiveNames {
			fmt.Printf(constants.MsgReleaseDryRun, "Compress → "+name)
		}
	}

	for _, a := range userAssets {
		fmt.Printf(constants.MsgReleaseDryRun, "Attach "+a)
	}

	if checksums && len(userAssets) > 0 {
		fmt.Printf(constants.MsgReleaseDryRun, "Generate "+constants.ChecksumsFile+" (SHA256)")
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

// buildZipGroupAssets creates archives from persistent zip groups.
func buildZipGroupAssets(opts Options) []string {
	if len(opts.ZipGroups) == 0 {
		return nil
	}

	fmt.Printf(constants.MsgZGProcessing, len(opts.ZipGroups))

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Cannot open DB for zip groups: %v\n", err)

		return nil
	}
	defer db.Close()

	stagingDir, err := EnsureStagingDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrZGStagingDir, err)

		return nil
	}

	archives := BuildZipGroupArchives(db, opts.ZipGroups, stagingDir)

	if len(archives) == 0 {
		fmt.Printf(constants.MsgZGNoArchives, len(opts.ZipGroups))
	}

	collectZipChecksums(archives)

	return archives
}

// buildAdHocZipAssets creates archives from ad-hoc -Z paths.
func buildAdHocZipAssets(opts Options) []string {
	if len(opts.ZipItems) == 0 {
		return nil
	}

	stagingDir, err := EnsureStagingDir()
	if err != nil {
		return nil
	}

	archives := BuildAdHocArchive(opts.ZipItems, opts.BundleName, stagingDir)
	collectZipChecksums(archives)

	return archives
}

// collectZipChecksums computes SHA-1 for each archive and stores in lastZipChecksums.
func collectZipChecksums(archives []string) {
	if len(archives) == 0 {
		return
	}

	if lastZipChecksums == nil {
		lastZipChecksums = make(map[string]string)
	}

	for _, archivePath := range archives {
		hash, err := sha1File(archivePath)
		if err != nil {
			continue
		}

		lastZipChecksums[filepath.Base(archivePath)] = hash
	}
}

// printDryRunZipGroups shows zip group plan in dry-run mode.
func printDryRunZipGroups(opts Options) {
	if len(opts.ZipGroups) > 0 {
		db, err := store.OpenDefault()
		if err == nil {
			defer db.Close()

			DryRunZipGroups(db, opts.ZipGroups)
		}
	}

	DryRunAdHoc(opts.ZipItems, opts.BundleName)
}
