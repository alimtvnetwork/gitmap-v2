package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// runRescan handles the "rescan" subcommand.
func runRescan() {
	cache, err := loadScanCache()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRescanNoCache, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgRescanReplay, cache.Dir)
	runScanFromCache(cache)
}

// loadScanCache reads the last-scan.json from the output folder.
func loadScanCache() (model.ScanCache, error) {
	path := filepath.Join(constants.DefaultOutputFolder, constants.DefaultScanCacheFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return model.ScanCache{}, err
	}
	var cache model.ScanCache
	err = json.Unmarshal(data, &cache)

	return cache, err
}

// saveScanCache writes the current scan flags to last-scan.json.
func saveScanCache(outputDir string, cache model.ScanCache) {
	path := filepath.Join(outputDir, constants.DefaultScanCacheFile)
	data, err := json.MarshalIndent(cache, "", constants.JSONIndent)
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(path), constants.DirPermission)
	_ = os.WriteFile(path, data, 0o644)
	fmt.Printf(constants.MsgScanCacheSaved, path)
}

// runScanFromCache replays a scan using cached flags.
func runScanFromCache(c model.ScanCache) {
	args := buildScanArgs(c)
	runScan(args)
}

// buildScanArgs reconstructs CLI args from a ScanCache.
func buildScanArgs(c model.ScanCache) []string {
	var args []string
	if len(c.ConfigPath) > 0 {
		args = append(args, "--config", c.ConfigPath)
	}
	if len(c.Mode) > 0 {
		args = append(args, "--mode", c.Mode)
	}
	if len(c.Output) > 0 {
		args = append(args, "--output", c.Output)
	}
	if len(c.OutFile) > 0 {
		args = append(args, "--out-file", c.OutFile)
	}
	if len(c.OutputPath) > 0 {
		args = append(args, "--output-path", c.OutputPath)
	}
	if c.GithubDesktop {
		args = append(args, "--github-desktop")
	}
	if c.OpenFolder {
		args = append(args, "--open")
	}
	if c.Quiet {
		args = append(args, "--quiet")
	}
	if len(c.Dir) > 0 && c.Dir != constants.DefaultDir {
		args = append(args, c.Dir)
	}

	return args
}
