// Package formatter renders ScanRecords to terminal output.
package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/user/gitmap/model"
)

// Terminal writes records to the given writer as a formatted table.
func Terminal(w io.Writer, records []model.ScanRecord) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	err := writeTerminalHeader(tw)
	if err != nil {
		return err
	}

	return writeTerminalRows(tw, records)
}

// writeTerminalHeader prints the column header line.
func writeTerminalHeader(tw *tabwriter.Writer) error {
	_, err := fmt.Fprintln(tw,
		"REPO\tBRANCH\tPATH\tCLONE INSTRUCTION")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(tw,
		"----\t------\t----\t-----------------")

	return err
}

// writeTerminalRows prints each record as a table row.
func writeTerminalRows(tw *tabwriter.Writer, records []model.ScanRecord) error {
	for _, r := range records {
		err := writeOneRow(tw, r)
		if err != nil {
			return err
		}
	}

	return tw.Flush()
}

// writeOneRow prints a single record row.
func writeOneRow(tw *tabwriter.Writer, r model.ScanRecord) error {
	_, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
		r.RepoName, r.Branch, r.RelativePath, r.CloneInstruction)

	return err
}
