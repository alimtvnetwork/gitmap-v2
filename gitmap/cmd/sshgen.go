package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/store"
)

// runSSHGenerate generates a new SSH key pair.
func runSSHGenerate(args []string) {
	name, keyPath, email, force := parseSSHGenFlags(args)

	if err := validateSSHKeygen(); err != nil {
		fmt.Fprint(os.Stderr, constants.ErrSSHKeygenMissing)
		os.Exit(1)
	}

	if len(email) == 0 {
		email = resolveGitEmail()
	}
	if len(email) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrSSHEmailResolve)
		os.Exit(1)
	}

	keyPath = expandHome(keyPath)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHCreate, err)
		os.Exit(1)
	}
	defer db.Close()

	if db.SSHKeyExists(name) && !force {
		if !handleExistingKey(db, name, &keyPath) {
			return
		}
	}

	generateAndStore(db, name, keyPath, email)
}

// parseSSHGenFlags parses flags for SSH key generation.
func parseSSHGenFlags(args []string) (name, keyPath, email string, force bool) {
	fs := flag.NewFlagSet(constants.CmdSSH, flag.ExitOnError)
	nameFlag := fs.String("name", constants.DefaultSSHKeyName, "Key label")
	fs.StringVar(nameFlag, "n", constants.DefaultSSHKeyName, "Key label (short)")
	pathFlag := fs.String("path", "", "Key file path")
	fs.StringVar(pathFlag, "p", "", "Key file path (short)")
	emailFlag := fs.String("email", "", "Email comment")
	fs.StringVar(emailFlag, "e", "", "Email comment (short)")
	forceFlag := fs.Bool("force", false, "Skip prompt if key exists")
	fs.BoolVar(forceFlag, "f", false, "Skip prompt (short)")
	fs.Parse(args)

	path := *pathFlag
	if len(path) == 0 {
		path = defaultSSHKeyPath(*nameFlag)
	}

	return *nameFlag, path, *emailFlag, *forceFlag
}

// handleExistingKey prompts the user when a key already exists.
// Returns true if generation should proceed, false to cancel.
func handleExistingKey(db *store.DB, name string, keyPath *string) bool {
	existing, _ := db.FindSSHKeyByName(name)
	fmt.Fprintf(os.Stdout, constants.MsgSSHExists, name, existing.PrivatePath)
	fmt.Fprintf(os.Stdout, constants.MsgSSHExistsFP, existing.Fingerprint)
	fmt.Fprint(os.Stdout, constants.MsgSSHPromptAction)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToUpper(input))

	if input == "R" {
		removeKeyFiles(existing.PrivatePath)
		*keyPath = existing.PrivatePath

		return true
	}
	if input == "N" {
		fmt.Fprint(os.Stdout, constants.MsgSSHNewPathPrompt)
		newPath, _ := reader.ReadString('\n')
		*keyPath = expandHome(strings.TrimSpace(newPath))

		return true
	}

	return false
}

// generateAndStore runs ssh-keygen and stores the result in the database.
func generateAndStore(db *store.DB, name, keyPath, email string) {
	if err := ensureSSHDir(filepath.Dir(keyPath)); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHKeygen, err)
		os.Exit(1)
	}

	cmd := exec.Command(constants.SSHKeygenBin,
		"-t", constants.SSHKeyType,
		"-b", constants.SSHKeyBits,
		"-C", email,
		"-f", keyPath,
		"-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHKeygen, err)
		os.Exit(1)
	}

	pubKey, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHReadPub, err)
		os.Exit(1)
	}

	fingerprint := readFingerprint(keyPath)

	if db.SSHKeyExists(name) {
		_ = db.UpdateSSHKey(name, keyPath, string(pubKey), fingerprint, email)
	} else {
		_, _ = db.InsertSSHKey(name, keyPath, string(pubKey), fingerprint, email)
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHGenerated, name)
	fmt.Fprintf(os.Stdout, constants.MsgSSHPath, keyPath)
	fmt.Fprintf(os.Stdout, constants.MsgSSHFingerprint, fingerprint)
	fmt.Fprint(os.Stdout, constants.MsgSSHPubLabel)
	fmt.Fprintf(os.Stdout, "  %s\n", strings.TrimSpace(string(pubKey)))
	fmt.Fprint(os.Stdout, constants.MsgSSHCopyHint)

	updateSSHConfig(db)
}

// validateSSHKeygen checks if ssh-keygen is available on PATH.
func validateSSHKeygen() error {
	_, err := exec.LookPath(constants.SSHKeygenBin)

	return err
}

// resolveGitEmail reads the global Git email config.
func resolveGitEmail() string {
	out, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// readFingerprint reads the SHA256 fingerprint of a key file.
func readFingerprint(keyPath string) string {
	out, err := exec.Command(constants.SSHKeygenBin, "-lf", keyPath+".pub").Output()
	if err != nil {
		return "unknown"
	}

	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return parts[1]
	}

	return "unknown"
}

// removeKeyFiles deletes private and public key files.
func removeKeyFiles(privatePath string) {
	_ = os.Remove(privatePath)
	_ = os.Remove(privatePath + ".pub")
}

// defaultSSHKeyPath returns the default key path based on name.
func defaultSSHKeyPath(name string) string {
	home, _ := os.UserHomeDir()
	if name == constants.DefaultSSHKeyName {
		return filepath.Join(home, ".ssh", "id_rsa")
	}

	return filepath.Join(home, ".ssh", "id_rsa_"+name)
}

// expandHome expands ~ to the user's home directory.
func expandHome(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[1:])
	}

	return path
}

// ensureSSHDir creates a directory with 0700 permissions if it doesn't exist.
func ensureSSHDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}
