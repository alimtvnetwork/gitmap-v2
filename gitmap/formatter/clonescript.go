// Package formatter — clonescript.go generates a clone.ps1 PowerShell script.
package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// WriteCloneScript writes a self-contained PowerShell clone script.
func WriteCloneScript(w io.Writer, records []model.ScanRecord) error {
	writeScriptHeader(w, len(records))
	writeScriptBody(w, records)
	writeScriptFooter(w)

	return nil
}

// writeScriptHeader writes the PowerShell script preamble.
func writeScriptHeader(w io.Writer, count int) {
	fmt.Fprintln(w, constants.ScriptHeader)
	fmt.Fprintf(w, constants.ScriptParamBlock)
	fmt.Fprintln(w, constants.ScriptErrorPref)
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ScriptBanner, count)
	fmt.Fprintln(w)
}

// writeScriptBody writes the clone commands for each repo.
func writeScriptBody(w io.Writer, records []model.ScanRecord) {
	fmt.Fprintln(w, constants.ScriptCounters)
	fmt.Fprintln(w)
	for i, r := range records {
		writeScriptEntry(w, r, i+1, len(records))
	}
}

// writeScriptEntry writes one repo clone block.
func writeScriptEntry(w io.Writer, r model.ScanRecord, idx, total int) {
	relPath := strings.ReplaceAll(r.RelativePath, "/", "\\")
	fmt.Fprintf(w, constants.ScriptRepoHeader, idx, total, r.RepoName)
	fmt.Fprintf(w, constants.ScriptMkdir, relPath)
	fmt.Fprintf(w, constants.ScriptCloneCmd, r.Branch, cloneURL(r), relPath)
	fmt.Fprintln(w, constants.ScriptCloneCheck)
	fmt.Fprintln(w)
}

// cloneURL picks the best URL from a record.
func cloneURL(r model.ScanRecord) string {
	if len(r.HTTPSUrl) > 0 {
		return r.HTTPSUrl
	}

	return r.SSHUrl
}

// writeScriptFooter writes the summary section.
func writeScriptFooter(w io.Writer) {
	fmt.Fprintln(w, constants.ScriptSummary)
}
