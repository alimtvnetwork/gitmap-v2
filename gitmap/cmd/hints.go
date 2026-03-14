package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// printHints prints contextual helper hints to stderr.
func printHints(hints []hintEntry) {
	fmt.Fprint(os.Stderr, constants.MsgHintHeader)
	for _, h := range hints {
		fmt.Fprintf(os.Stderr, constants.MsgHintRowFmt, h.command, h.description)
	}
}

// hintEntry holds a command example and its description.
type hintEntry struct {
	command     string
	description string
}

// projectReposHints returns hints shown after go-repos, node-repos, etc.
func projectReposHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupAdd, constants.HintGroupAddDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
		{constants.HintPullGroup, constants.HintPullGroupDesc},
	}
}

// listHints returns hints shown after gitmap list.
func listHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintLsType, constants.HintLsTypeDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
	}
}

// listTypeHints returns hints shown after gitmap ls <type>.
func listTypeHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupAdd, constants.HintGroupAddDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
		{constants.HintPullGroup, constants.HintPullGroupDesc},
	}
}

// listGroupsHints returns hints shown after gitmap ls groups.
func listGroupsHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintGroupShow, constants.HintGroupShowDesc},
	}
}

// activeGroupHints returns hints shown after gitmap g (active group).
func activeGroupHints() []hintEntry {
	return []hintEntry{
		{constants.HintGPull, constants.HintGPullDesc},
		{constants.HintGStatus, constants.HintGStatusDesc},
		{constants.HintGClear, constants.HintGClearDesc},
	}
}

// groupListHints returns hints shown after gitmap group list.
func groupListHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintGroupShow, constants.HintGroupShowDesc},
		{constants.HintGroupDelete, constants.HintGroupDeleteDesc},
	}
}

// cdHints returns hints shown after gitmap cd <name>.
func cdHints() []hintEntry {
	return []hintEntry{
		{constants.HintCDSetDefault, constants.HintCDSetDefaultDesc},
		{constants.HintCDRepos, constants.HintCDReposDesc},
	}
}
