package constants

// Compress messages.
const (
	MsgCompressArchive = "  ✓ Compressed %s → %s\n"
	ErrCompressFailed  = "  ✗ Failed to compress %s: %v\n"
	FlagDescCompress   = "Wrap release assets in .zip (Windows) or .tar.gz (Linux/macOS)"
)

// Compress help text.
const HelpCompress = "  --compress          Wrap assets in .zip (Windows) or .tar.gz archives"
