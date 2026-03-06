package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/verbose"
)

// executeUpdate writes a temp PS1 script and runs it.
func executeUpdate(repoPath string) {
	scriptPath, err := writeUpdateScript(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
		os.Exit(1)
	}
	defer os.Remove(scriptPath)

	log := verbose.Get()
	if log != nil {
		log.Log(constants.UpdateScriptLogExec, scriptPath)
	}

	runUpdateScript(scriptPath)
}

// writeUpdateScript creates a temporary PowerShell script for self-update.
// Writes with UTF-8 BOM so PowerShell correctly handles Unicode characters.
func writeUpdateScript(repoPath string) (string, error) {
	runPS1 := filepath.Join(repoPath, "run.ps1")
	script := buildUpdateScript(repoPath, runPS1)

	return writeScriptToTemp(script)
}

// writeScriptToTemp writes script content to a temp file with UTF-8 BOM.
func writeScriptToTemp(script string) (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), constants.UpdateScriptGlob)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := tmpFile.Write(bom); err != nil {
		return "", err
	}
	if _, err := tmpFile.WriteString(script); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// buildUpdateScript generates the PowerShell script content.
func buildUpdateScript(repoPath, runPS1 string) string {
	return fmt.Sprintf(constants.UpdatePSHeader, repoPath) +
		fmt.Sprintf(constants.UpdatePSDeployDetect, repoPath) +
		constants.UpdatePSVersionBefore +
		fmt.Sprintf(constants.UpdatePSRunUpdate, runPS1) +
		constants.UpdatePSVersionAfter +
		constants.UpdatePSVerify +
		constants.UpdatePSPostActions
}

// runUpdateScript executes the PowerShell script with output piped to terminal.
func runUpdateScript(scriptPath string) {
	cmd := exec.Command(constants.PSBin, constants.PSExecPolicy, constants.PSBypass,
		constants.PSNoProfile, constants.PSNoLogo, constants.PSFile, scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()

	logScriptResult(err)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateFailed, err)
		os.Exit(1)
	}
}

// logScriptResult logs the update script exit status if verbose is active.
func logScriptResult(err error) {
	log := verbose.Get()
	if log != nil {
		log.Log(constants.UpdateScriptLogExit, err)
	}
}
