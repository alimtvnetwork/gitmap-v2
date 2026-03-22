package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// runSSHList displays all stored SSH keys as an aligned table.
func runSSHList() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	keys, err := db.ListSSHKeys()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}

	if len(keys) == 0 {
		fmt.Println("  No SSH keys stored. Run 'gitmap ssh' to generate one.")

		return
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHListHeader, len(keys))
	fmt.Fprintf(os.Stdout, constants.MsgSSHListColumns, "Name", "Path", "Fingerprint", "Created")
	fmt.Fprintf(os.Stdout, constants.MsgSSHListColumns,
		"───────────────", "──────────────────────────────",
		"─────────────────────────", "──────────")

	for _, k := range keys {
		created := k.CreatedAt
		if len(created) > 10 {
			created = created[:10]
		}

		fmt.Fprintf(os.Stdout, constants.MsgSSHListRow, k.Name, k.PrivatePath, k.Fingerprint, created)
	}
}
