package constants

// Doctor command messages.
const (
	DoctorBannerFmt       = "\n  gitmap doctor (v%s)\n"
	DoctorBannerRule      = "  ──────────────────────────────────────────"
	DoctorIssuesFmt       = "  Found %d issue(s). See recommendations above.\n"
	DoctorFixPathTip      = "  Tip: run 'gitmap doctor --fix-path' to auto-sync the PATH binary.\n\n"
	DoctorAllPassed       = "  All checks passed."
	DoctorFixBannerFmt    = "\n  gitmap doctor --fix-path (v%s)\n"
	DoctorActivePathFmt   = "  Active PATH:  %s (%s)\n"
	DoctorDeployedFmt     = "  Deployed:     %s (%s)\n"
	DoctorSyncingFmt      = "  Syncing %s -> %s...\n"
	DoctorRetryFmt        = "  [%d/%d] File in use, retrying...\n"
	DoctorRenamedMsg      = "  Renamed locked binary to .old, copying fresh..."
	DoctorKillingMsg      = "  Attempting to stop stale gitmap processes..."
	DoctorKilledFmt       = "  Stopped process(es): %s\n"
	DoctorSyncFailTitle   = "Could not sync PATH binary after all fallback attempts"
	DoctorSyncFailDetail  = "The file is still locked by another process."
	DoctorSyncFailFix1    = "Close all terminals and apps using gitmap, then run:"
	DoctorSyncFailFix2Fmt = "  Copy-Item \"%s\" \"%s\" -Force"
	DoctorFixFlagDesc     = "Sync the active PATH binary from the deployed binary"
	DoctorOKPathFmt       = "PATH binary synced successfully: %s"
	DoctorWarnSyncFmt     = "Synced but version mismatch: got %s, expected %s"
	DoctorNotOnPath       = "gitmap not found on PATH"
	DoctorNoSync          = "Cannot sync — no active binary to replace."
	DoctorAddPathFix      = "Add your deploy directory to PATH first."
	DoctorCannotResolve   = "Cannot resolve deployed binary"
	DoctorAlreadySynced   = "PATH already points to the deployed binary. Nothing to sync."
	DoctorVersionsMatch   = "Versions already match (%s). No sync needed."
	DoctorRepoPathMissing = "RepoPath not embedded"
	DoctorRepoPathDetail  = "Binary was not built with run.ps1. Self-update will not work."
	DoctorRepoPathFix     = "Rebuild with: .\\run.ps1"
	DoctorRepoPathOKFmt   = "RepoPath: %s"
	DoctorPathBinaryFmt   = "PATH binary: %s (%s)"
	DoctorPathMissTitle   = "gitmap not found on PATH"
	DoctorPathMissDetail  = "The gitmap binary is not accessible from your terminal."
	DoctorPathMissFix     = "Add your deploy directory to PATH (e.g., E:\\bin-run\\gitmap)"
	DoctorDeployReadFail  = "Cannot read powershell.json"
	DoctorDeployReadDet   = "Deploy path detection unavailable."
	DoctorNoDeployPath    = "No deployPath in powershell.json"
	DoctorNoDeployDet     = "Deploy target not configured."
	DoctorDeployNotFound  = "Deployed binary not found"
	DoctorDeployRunFix    = "Run: .\\run.ps1"
	DoctorDeployOKFmt     = "Deployed binary: %s (%s)"
	DoctorGitMissTitle    = "git not found on PATH"
	DoctorGitMissDetail   = "Git is required for most gitmap commands."
	DoctorGitOKFmt        = "Git: %s (%s)"
	DoctorGitOKPathFmt    = "Git: %s (version unknown)"
	DoctorGoWarn          = "Go not found on PATH (needed only for building from source)"
	DoctorGoOKFmt         = "Go: %s"
	DoctorGoOKPathFmt     = "Go: %s (version unknown)"
	DoctorChangelogWarn   = "CHANGELOG.md not found (changelog command will not work)"
	DoctorChangelogOK     = "CHANGELOG.md present"
	DoctorVersionMismatch = "PATH binary version mismatch"
	DoctorVMismatchFmt    = "PATH: %s, Source: %s"
	DoctorVMismatchFix    = "Run: gitmap update  or  gitmap doctor --fix-path"
	DoctorDeployMismatch  = "Deployed binary version mismatch"
	DoctorDMismatchFmt    = "Deployed: %s, Source: %s"
	DoctorDMismatchFix    = "Run: .\\run.ps1 -NoPull"
	DoctorBinariesDiffer  = "PATH and deployed binaries differ"
	DoctorBDifferFmt      = "PATH: %s (%s), Deployed: %s (%s)"
	DoctorBDifferFix      = "Run: gitmap doctor --fix-path"
	DoctorSourceOKFmt     = "Source version: %s (all binaries match)"
	DoctorResolveNoRepo   = "RepoPath not embedded — rebuild with run.ps1"
	DoctorResolveNoRead   = "cannot read powershell.json: %v"
	DoctorResolveNoDeploy = "no deployPath in powershell.json"
	DoctorResolveNotFound = "deployed binary not found: %s"
	DoctorDefaultBinary   = "gitmap.exe"
)

// Doctor binary and tool lookup names.
const (
	GitMapBin            = "gitmap"
	GoBin                = "go"
	GoVersionArg         = "version"
	PowershellConfigFile = "powershell.json"
	JSONKeyDeployPath    = "deployPath"
	JSONKeyBinaryName    = "binaryName"
	BackupSuffix         = ".old"
	GitMapSubdir         = "gitmap"
)

// Doctor format markers.
const (
	DoctorOKFmt    = "  %s[OK]%s %s\n"
	DoctorIssueFmt = "  %s[!!]%s %s\n"
	DoctorFixFmt   = "       %sFix:%s %s\n"
	DoctorWarnFmt  = "  %s[--]%s %s\n"
	DoctorDetail   = "       %s\n"
)

// Doctor config validation messages.
const (
	DoctorConfigMissing = "config.json not found (using defaults)"
	DoctorConfigInvalid = "config.json is not valid JSON"
	DoctorConfigOKFmt   = "Config: %s"
)

// Doctor database validation messages.
const (
	DoctorDBOpenFail    = "Database cannot be opened"
	DoctorDBMigrateFail = "Database migration failed"
	DoctorDBOK          = "Database: %s"
)

// Doctor lock file messages.
const (
	DoctorLockNone   = "No stale lock file"
	DoctorLockExists = "Lock file exists — another gitmap may be running (or stale)"
)

// Doctor network messages.
const (
	DoctorNetworkOK      = "Network: github.com reachable"
	DoctorNetworkOffline = "Network: github.com unreachable (offline mode)"
)

// Doctor legacy directory messages.
const (
	DoctorLegacyDirsOK = "No legacy directories (.release/, gitmap-output/, .deployed/)"
)
