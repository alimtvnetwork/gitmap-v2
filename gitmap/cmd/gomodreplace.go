package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// readModulePath reads the module path from go.mod in the current directory.
func readModulePath() string {
	data, err := os.ReadFile(constants.GoModFile)
	if err != nil {
		fmt.Fprint(os.Stderr, constants.ErrGoModNoFile)
		os.Exit(1)
	}

	return parseModuleLine(string(data))
}

// parseModuleLine extracts the module path from go.mod content.
func parseModuleLine(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, constants.GoModModuleLine) {
			return strings.TrimSpace(strings.TrimPrefix(line, constants.GoModModuleLine))
		}
	}

	fmt.Fprint(os.Stderr, constants.ErrGoModNoModule)
	os.Exit(1)

	return ""
}

// replaceModulePath replaces all occurrences of oldPath with newPath across the repo.
func replaceModulePath(oldPath, newPath string, verbose bool) int {
	replaceInGoMod(oldPath, newPath)
	goFiles := findGoFilesWithPath(oldPath)

	if len(goFiles) == 0 {
		fmt.Print(constants.MsgGoModNoImports)

		return 0
	}

	for _, f := range goFiles {
		replaceInFile(f, oldPath, newPath)
		if verbose {
			fmt.Printf(constants.MsgGoModVerboseFile, f)
		}
	}

	return len(goFiles)
}

// replaceInGoMod replaces the module line in go.mod.
func replaceInGoMod(oldPath, newPath string) {
	data, err := os.ReadFile(constants.GoModFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModReadFailed, constants.GoModFile, err)
		os.Exit(1)
	}

	updated := strings.Replace(string(data), oldPath, newPath, -1)
	writeFileContent(constants.GoModFile, updated)
}

// findGoFilesWithPath walks the repo and returns .go files containing oldPath.
func findGoFilesWithPath(oldPath string) []string {
	var matches []string

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && isExcludedDir(info.Name()) {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == constants.GoFileExt {
			if fileContains(path, oldPath) {
				matches = append(matches, path)
			}
		}

		return nil
	})

	return matches
}

// isExcludedDir checks if a directory should be skipped.
func isExcludedDir(name string) bool {
	for _, d := range constants.GoModExcludeDirs {
		if name == d {
			return true
		}
	}

	return false
}

// fileContains checks if a file contains the given substring.
func fileContains(path, substr string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(string(data), substr)
}

// replaceInFile replaces all occurrences of oldPath with newPath in a file.
func replaceInFile(path, oldPath, newPath string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModReadFailed, path, err)

		return
	}

	updated := strings.Replace(string(data), oldPath, newPath, -1)
	writeFileContent(path, updated)
}

// writeFileContent writes content to a file preserving its permissions.
func writeFileContent(path, content string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModWriteFailed, path, err)
		os.Exit(1)
	}

	err = os.WriteFile(path, []byte(content), info.Mode())
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModWriteFailed, path, err)
		os.Exit(1)
	}
}
