package constants

// Go release asset constants.
const (
	AssetsStagingDir = "release-assets"
	GitHubTokenEnv   = "GITHUB_TOKEN"
)

// Asset flag descriptions.
const (
	FlagDescNoAssets = "Skip automatic Go binary cross-compilation"
	FlagDescTargets  = "Comma-separated cross-compile targets (e.g. windows/amd64,linux/arm64)"
)

// Asset help text.
const (
	HelpNoAssets    = "  --no-assets         Skip Go binary cross-compilation"
	HelpTargets     = "  --targets <list>    Cross-compile targets: windows/amd64,linux/arm64"
)

// Asset messages.
const (
	MsgAssetDetected     = "  → Detected Go project: %s\n"
	MsgAssetCrossCompile = "\n  Cross-compiling %d target(s)...\n"
	MsgAssetBuilt        = "  ✓ Built %s (%s/%s)\n"
	MsgAssetBuildSummary = "  → Built %d/%d binaries successfully\n"
	MsgAssetUploaded     = "  ✓ Uploaded %s\n"
	MsgAssetUploadStart  = "\n  Uploading %d asset(s) to GitHub...\n"
	MsgAssetSkipped      = "  → Skipping Go binary compilation (--no-assets)\n"
	MsgAssetNoMain       = "  → No buildable main package found, skipping binaries\n"
	MsgAssetNoGoProject  = ""
	MsgAssetStagingClean = "  ✓ Cleaned up staging directory\n"
)

// Asset dry-run messages.
const (
	MsgAssetDryRunHeader = "  [dry-run] Would cross-compile %d binaries:\n"
	MsgAssetDryRunBinary = "    → %s\n"
	MsgAssetDryRunUpload = "  [dry-run] Would upload %d assets\n"
)

// Asset error messages.
const (
	ErrAssetBuildFailed  = "  ✗ Build failed for %s/%s: %s\n"
	ErrAssetUploadRetry  = "  ⟳ Retrying upload for %s...\n"
	ErrAssetUploadFailed = "  ✗ Upload failed for %s: %v\n"
	ErrAssetNoToken      = "  ✗ GITHUB_TOKEN not set — skipping asset upload\n"
	ErrAssetRemoteParse  = "  ✗ Could not parse remote origin: %v\n"
)
