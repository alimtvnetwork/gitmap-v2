package constants

// Clone progress format strings.
const (
	ProgressBeginFmt   = "[%3d/%d]  Cloning %s ..."
	ProgressDoneFmt    = " done (%s)\n"
	ProgressFailFmt    = " FAILED\n"
	ProgressSummaryFmt = "\nClone complete: %d/%d repos in %s\n"
	ProgressDetailFmt  = "  Cloned: %d | Pulled: %d | Failed: %d\n"
)
