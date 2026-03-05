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
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	cfg = config.MergeWithFlags(cfg, mode, output, outputPath)
	executeScan(dir, cfg, outFile)
}

// executeScan performs the directory scan and outputs results.
func executeScan(dir string, cfg model.Config, outFile string) {
	repos, err := scanner.ScanDir(dir, cfg.ExcludeDirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Scan error: %v\n", err)
		os.Exit(1)
	}
	records := mapper.BuildRecords(repos, cfg.DefaultMode, cfg.Notes)
	fmt.Printf("Found %d repositories.\n", len(records))
	outputRecords(records, cfg, outFile)
}

// outputRecords routes records to the correct formatter.
func outputRecords(records []model.ScanRecord, cfg model.Config, outFile string) {
	if cfg.DefaultOutput == constants.OutputCSV {
		writeCSVOutput(records, cfg, outFile)
		return
	}
	if cfg.DefaultOutput == constants.OutputJSON {
		writeJSONOutput(records, cfg, outFile)
		return
	}
	writeTerminalOutput(records)
}

// writeTerminalOutput renders records to stdout.
func writeTerminalOutput(records []model.ScanRecord) {
	err := formatter.Terminal(os.Stdout, records)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Output error: %v\n", err)
	}
}

// writeCSVOutput writes records to a CSV file.
func writeCSVOutput(records []model.ScanRecord, cfg model.Config, outFile string) {
	path := resolveOutFile(outFile, cfg.OutputDir, constants.DefaultCSVFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteCSV(file, records)
	fmt.Printf("CSV written to %s\n", path)
}

// writeJSONOutput writes records to a JSON file.
func writeJSONOutput(records []model.ScanRecord, cfg model.Config, outFile string) {
	path := resolveOutFile(outFile, cfg.OutputDir, constants.DefaultJSONFile)
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	formatter.WriteJSON(file, records)
	fmt.Printf("JSON written to %s\n", path)
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
		fmt.Fprintf(os.Stderr, "Cannot create directory: %v\n", err)

		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create file: %v\n", err)

		return nil, err
	}

	return file, nil
}
