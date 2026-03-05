package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/config"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/formatter"
	"github.com/user/gitmap/mapper"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/scanner"
)

// runScan handles the "scan" subcommand.
func runScan(args []string) {
	dir, cfgPath, mode, output, outFile, outputPath := parseScanFlags(args)
	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrConfigLoad, err)
		os.Exit(1)
	}
	cfg = config.MergeWithFlags(cfg, mode, output, outputPath)
	executeScan(dir, cfg, outFile)
}

// executeScan performs the directory scan and outputs results.
func executeScan(dir string, cfg model.Config, outFile string) {
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
	fmt.Printf(constants.MsgFoundRepos, len(records))
	outputDir := resolveOutputDir(cfg.OutputDir, absDir)
	writeAllOutputs(records, outputDir, outFile)
}

// resolveOutputDir determines the output directory relative to scan root.
func resolveOutputDir(cfgDir, scanDir string) string {
	if filepath.IsAbs(cfgDir) {
		return cfgDir
	}

	return filepath.Join(scanDir, constants.DefaultOutputFolder)
}

// writeAllOutputs writes terminal, CSV, JSON, and folder structure.
func writeAllOutputs(records []model.ScanRecord, outputDir, outFile string) {
	writeTerminalOutput(records)
	writeCSVOutput(records, outputDir, outFile)
	writeJSONOutput(records, outputDir)
	writeFolderStructure(records, outputDir)
}

// writeTerminalOutput renders records to stdout.
func writeTerminalOutput(records []model.ScanRecord) {
	err := formatter.Terminal(os.Stdout, records)
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
