// Package formatter — directclone.go generates a plain direct-clone.ps1 script.
package formatter

import (
	"io"

	"github.com/user/gitmap/model"
)

// WriteDirectCloneScript writes a plain PS1 with one git clone per line.
func WriteDirectCloneScript(w io.Writer, records []model.ScanRecord) error {
	tmpl, err := loadTemplate("direct-clone.ps1.tmpl")
	if err != nil {
		return err
	}

	data := CloneData{
		Repos: buildRepoEntries(records),
	}

	return tmpl.Execute(w, data)
}
