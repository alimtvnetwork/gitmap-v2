package release

import (
	"fmt"
	"strings"

	"github.com/user/gitmap/constants"
)

// printInstallHint prints install one-liner commands if the current repo
// matches the gitmap source repository prefix.
func printInstallHint(v Version) {
	url := getRemoteURL()
	if ShouldPrintInstallHint(url) {
		fmt.Printf(constants.MsgInstallHintHeader, v.String())
		fmt.Print(constants.MsgInstallHintWindows)
		fmt.Print(constants.MsgInstallHintUnix)
	}
}

// ShouldPrintInstallHint returns true if the remote URL matches the
// gitmap source repository prefix.
func ShouldPrintInstallHint(remoteURL string) bool {
	if len(remoteURL) == 0 {
		return false
	}

	normalized := strings.TrimSuffix(strings.ToLower(remoteURL), ".git")
	prefix := strings.TrimSuffix(strings.ToLower(constants.GitmapRepoPrefix), ".git")

	return strings.Contains(normalized, prefix)
}
