package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/user/gitmap/config"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/desktop"
	"github.com/user/gitmap/formatter"
	"github.com/user/gitmap/mapper"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/scanner"
)

// runScan handles the "scan" subcommand.
func runScan(args []string) {
	dir, cfgPath, mode, output, outFile, outputPath, ghDesktop, openFolder, quiet := parseScanFlags(args)
	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrConfigLoad, err)
		os.Exit(1)
	}
	cfg = config.MergeWithFlags(cfg, mode, output, outputPath)
	cache := model.ScanCache{
		Dir: dir, ConfigPath: cfgPath, Mode: mode, Output: output,
		OutFile: outFile, OutputPath: outputPath,
		GithubDesktop: ghDesktop, OpenFolder: openFolder, Quiet: quiet,
	}
	executeScan(dir, cfg, outFile, ghDesktop, openFolder, quiet, cache)
}

// executeScan performs the directory scan and outputs results.
func executeScan(dir string, cfg model.Config, outFile string, ghDesktop, openFolder bool, cache model.ScanCache) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScanFailed, err)
		os.Exit(1)
	}
	repos, err := scanner.ScanDir(absDir, cfg.ExcludeDirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScanFailed, err)
		os.Exit(1)
	}
	records := mapper.BuildRecords(repos, cfg.DefaultMode, cfg.Notes)
	outputDir := resolveOutputDir(cfg.OutputDir, absDir)
	writeAllOutputs(records, outputDir, outFile)
	saveScanCache(outputDir, cache)
	addToDesktop(records, ghDesktop)
	openOutputFolder(outputDir, openFolder)
}

// addToDesktop registers repos with GitHub Desktop if requested.
func addToDesktop(records []model.ScanRecord, enabled bool) {
	if enabled {
		summary := desktop.AddRepos(records)
		fmt.Printf(constants.MsgDesktopSummary, summary.Added, summary.Failed)
	}
}

// openOutputFolder opens the output directory in the OS file explorer.
func openOutputFolder(outputDir string, enabled bool) {
	if !enabled {
		return
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case constants.OSWindows:
		cmd = exec.Command("explorer", outputDir)
	case "darwin":
		cmd = exec.Command("open", outputDir)
	default:
		cmd = exec.Command("xdg-open", outputDir)
	}
	_ = cmd.Start()
	fmt.Printf(constants.MsgOpenedFolder, outputDir)
}


// resolveOutputDir determines the output directory relative to scan root.
func resolveOutputDir(cfgDir, scanDir string) string {
	if filepath.IsAbs(cfgDir) {
		return cfgDir
	}

	return filepath.Join(scanDir, constants.DefaultOutputFolder)
}

// writeAllOutputs writes terminal, CSV, JSON, text, folder structure, and clone scripts.
func writeAllOutputs(records []model.ScanRecord, outputDir, outFile string) {
	writeTerminalOutput(records, outputDir)
	writeCSVOutput(records, outputDir, outFile)
	writeJSONOutput(records, outputDir)
	writeTextOutput(records, outputDir)
	writeFolderStructure(records, outputDir)
	writeCloneScript(records, outputDir)
	writeDirectCloneScript(records, outputDir)
	writeDirectCloneSSHScript(records, outputDir)
	writeDesktopScript(records, outputDir)
}

// writeTextOutput writes records as plain text clone commands.
func writeTextOutput(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultTextFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteText(file, records)
	fmt.Printf(constants.MsgTextWritten, path)
}

// writeTerminalOutput renders records to stdout.
func writeTerminalOutput(records []model.ScanRecord, outputDir string) {
	err := formatter.Terminal(os.Stdout, records, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrOutputFailed, err)
	}
}

// writeCSVOutput writes records to a CSV file.
func writeCSVOutput(records []model.ScanRecord, outputDir, outFile string) {
	path := resolveOutFile(outFile, outputDir, constants.DefaultCSVFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteCSV(file, records)
	fmt.Printf(constants.MsgCSVWritten, path)
}

// writeJSONOutput writes records to a JSON file.
func writeJSONOutput(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultJSONFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteJSON(file, records)
	fmt.Printf(constants.MsgJSONWritten, path)
}

// writeFolderStructure writes a Markdown file showing the repo tree.
func writeFolderStructure(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultStructureFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteStructure(file, records)
	fmt.Printf(constants.MsgStructureWritten, path)
}

// writeCloneScript writes a PowerShell clone script.
func writeCloneScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultCloneScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteCloneScript(file, records)
	fmt.Printf(constants.MsgCloneScript, path)
}

// writeDirectCloneScript writes a plain PS1 with one git clone per line.
func writeDirectCloneScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDirectCloneScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteDirectCloneScript(file, records)
	fmt.Printf(constants.MsgDirectClone, path)
}

// writeDirectCloneSSHScript writes a plain SSH PS1 with one git clone per line.
func writeDirectCloneSSHScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDirectCloneSSHScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteDirectCloneSSHScript(file, records)
	fmt.Printf(constants.MsgDirectCloneSSH, path)
}

// writeDesktopScript writes a PowerShell script to register repos with GitHub Desktop.
func writeDesktopScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDesktopScript)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteDesktopScript(file, records)
	fmt.Printf(constants.MsgDesktopScript, path)
}

// resolveOutFile determines the output file path.
func resolveOutFile(outFile, outputDir, defaultName string) string {
	if len(outFile) > 0 {
		return outFile
	}

	return filepath.Join(outputDir, defaultName)
}

// createOutputFile ensures the directory exists and creates the file.
func createOutputFile(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), constants.DirPermission)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCreateDir, err)

		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCreateFile, err)

		return nil, err
	}

	return file, nil
}
