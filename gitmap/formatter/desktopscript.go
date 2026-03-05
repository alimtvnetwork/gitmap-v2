// Package formatter — desktopscript.go generates a register-desktop.ps1 script.
package formatter

import (
	"io"

	"github.com/user/gitmap/model"
)

// WriteDesktopScript writes a PowerShell script that registers repos
// with GitHub Desktop using the embedded desktop.ps1.tmpl template.
func WriteDesktopScript(w io.Writer, records []model.ScanRecord) error {
	tmpl, err := loadTemplate("desktop.ps1.tmpl")
	if err != nil {
		return err
	}

	data := DesktopData{
		Repos: buildDesktopEntries(records),
	}

	return tmpl.Execute(w, data)
}

// buildDesktopEntries converts ScanRecords into template-friendly entries.
func buildDesktopEntries(records []model.ScanRecord) []RepoEntry {
	entries := make([]RepoEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, RepoEntry{
			Name: r.RepoName,
			Path: backslashPath(r.RelativePath),
		})
	}

	return entries
}
