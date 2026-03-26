// Package release — assets.go orchestrates cross-compilation for Go projects.
package release

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/verbose"
)

// BuildTarget represents a single GOOS/GOARCH pair for cross-compilation.
type BuildTarget struct {
	GOOS   string `json:"goos"`
	GOARCH string `json:"goarch"`
}

// CrossCompileResult holds the outcome of a cross-compile step.
type CrossCompileResult struct {
	Target  BuildTarget
	Output  string
	Success bool
	Error   string
}

// DetectGoProject checks if the current directory contains a buildable Go project.
func DetectGoProject() bool {
	_, err := os.Stat(constants.IndicatorGoMod)

	return err == nil
}

// ReadModuleName reads the module name from go.mod.
func ReadModuleName() (string, error) {
	f, err := os.Open(constants.IndicatorGoMod)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("no module directive found in go.mod")
}

// BinaryName extracts the short name from a Go module path.
// "github.com/user/gitmap" → "gitmap"
func BinaryName(moduleName string) string {
	parts := strings.Split(moduleName, "/")

	return parts[len(parts)-1]
}

// FindMainPackages locates buildable main package directories.
// Checks root main.go first, then cmd/ subdirectories.
func FindMainPackages() []string {
	if fileExists(constants.GoMainFile) {
		return []string{"."}
	}

	cmdDir := constants.GoCmdDir
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return nil
	}

	var packages []string

	for _, entry := range entries {
		if entry.IsDir() {
			mainPath := filepath.Join(cmdDir, entry.Name(), constants.GoMainFile)
			if fileExists(mainPath) {
				packages = append(packages, filepath.Join(cmdDir, entry.Name()))
			}
		}
	}

	return packages
}

// CrossCompile builds binaries for all targets and packages.
// Returns the list of successfully built binary paths.
func CrossCompile(version string, targets []BuildTarget, packages []string, stagingDir string) []CrossCompileResult {
	var results []CrossCompileResult

	binName := resolveBinName()

	for _, pkg := range packages {
		pkgSuffix := ""
		if pkg != "." {
			pkgSuffix = "-" + filepath.Base(pkg)
		}

		for _, t := range targets {
			if verbose.IsEnabled() {
				verbose.Get().Log("build: %s/%s → %s", t.GOOS, t.GOARCH,
					filepath.Join(stagingDir, formatOutputName(binName+pkgSuffix, version, t)))
			}

			result := buildSingleTarget(binName+pkgSuffix, version, t, pkg, stagingDir)
			results = append(results, result)

			if result.Success {
				if verbose.IsEnabled() {
					info, statErr := os.Stat(result.Output)
					if statErr == nil {
						verbose.Get().Log("build: %s/%s complete (%d bytes)", t.GOOS, t.GOARCH, info.Size())
					}
				}

				fmt.Printf(constants.MsgAssetBuilt, filepath.Base(result.Output), t.GOOS, t.GOARCH)
			} else {
				if verbose.IsEnabled() {
					verbose.Get().Log("build: %s/%s failed: %s", t.GOOS, t.GOARCH, result.Error)
				}

				fmt.Fprintf(os.Stderr, constants.ErrAssetBuildFailed, t.GOOS, t.GOARCH, result.Error)
			}
		}
	}

	return results
}

// resolveBinName reads go.mod for the binary name.
func resolveBinName() string {
	mod, err := ReadModuleName()
	if err != nil {
		return "app"
	}

	return BinaryName(mod)
}

// buildSingleTarget compiles one GOOS/GOARCH combination.
func buildSingleTarget(binName, version string, target BuildTarget, pkgDir, stagingDir string) CrossCompileResult {
	outputName := formatOutputName(binName, version, target)
	outputPath := filepath.Join(stagingDir, outputName)

	ldflags := fmt.Sprintf("-s -w -X main.version=%s", version)

	cmd := exec.Command("go", "build",
		"-ldflags", ldflags,
		"-o", outputPath,
		"./"+pkgDir,
	)

	cmd.Env = buildEnv(target)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return CrossCompileResult{
			Target:  target,
			Output:  outputPath,
			Success: false,
			Error:   strings.TrimSpace(string(out)),
		}
	}

	return CrossCompileResult{
		Target:  target,
		Output:  outputPath,
		Success: true,
	}
}

// formatOutputName creates the binary filename with platform suffix.
func formatOutputName(binName, version string, target BuildTarget) string {
	name := fmt.Sprintf("%s_%s_%s_%s", binName, version, target.GOOS, target.GOARCH)
	if target.GOOS == "windows" {
		name += ".exe"
	}

	return name
}

// buildEnv returns the os.Environ() with CGO_ENABLED=0, GOOS, GOARCH set.
func buildEnv(target BuildTarget) []string {
	env := os.Environ()

	env = setEnv(env, "CGO_ENABLED", "0")
	env = setEnv(env, "GOOS", target.GOOS)
	env = setEnv(env, "GOARCH", target.GOARCH)

	return env
}

// setEnv sets or replaces an environment variable in a slice.
func setEnv(env []string, key, value string) []string {
	prefix := key + "="

	for i, e := range env {
		if strings.HasPrefix(e, prefix) {
			env[i] = prefix + value

			return env
		}
	}

	return append(env, prefix+value)
}

// fileExists checks if a file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// CollectSuccessfulBuilds returns file paths from successful results.
func CollectSuccessfulBuilds(results []CrossCompileResult) []string {
	var paths []string

	for _, r := range results {
		if r.Success {
			paths = append(paths, r.Output)
		}
	}

	return paths
}

// EnsureStagingDir creates the release-assets staging directory.
func EnsureStagingDir() (string, error) {
	dir := constants.AssetsStagingDir
	err := os.MkdirAll(dir, constants.DirPermission)
	if err != nil {
		return "", fmt.Errorf("create staging dir: %w", err)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("staging: created directory %s", dir)
	}

	return dir, nil
}

// CleanupStagingDir removes the staging directory after upload.
func CleanupStagingDir() {
	if verbose.IsEnabled() {
		verbose.Get().Log("staging: removing directory %s", constants.AssetsStagingDir)
	}

	os.RemoveAll(constants.AssetsStagingDir)
}
