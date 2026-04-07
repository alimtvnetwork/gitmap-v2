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
	if len(url) == 0 {
		return
	}

	normalized := strings.TrimSuffix(strings.ToLower(url), ".git")
	prefix := strings.TrimSuffix(strings.ToLower(constants.GitmapRepoPrefix), ".git")

	if !strings.Contains(normalized, prefix) {
		return
	}

	fmt.Printf(constants.MsgInstallHintHeader, v.String())
	fmt.Print(constants.MsgInstallHintWindows)
	fmt.Print(constants.MsgInstallHintUnix)
}
