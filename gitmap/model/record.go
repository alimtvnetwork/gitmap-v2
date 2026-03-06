// Package model defines the core data structures for gitmap.
package model

import "github.com/user/gitmap/constants"

// ScanRecord holds all information about a discovered Git repository.
type ScanRecord struct {
	ID               string `json:"id"               csv:"id"`
	Slug             string `json:"slug"             csv:"slug"`
	RepoName         string `json:"repoName"         csv:"repoName"`
	HTTPSUrl         string `json:"httpsUrl"          csv:"httpsUrl"`
	SSHUrl           string `json:"sshUrl"            csv:"sshUrl"`
	Branch           string `json:"branch"            csv:"branch"`
	RelativePath     string `json:"relativePath"      csv:"relativePath"`
	AbsolutePath     string `json:"absolutePath"      csv:"absolutePath"`
	CloneInstruction string `json:"cloneInstruction"  csv:"cloneInstruction"`
	Notes            string `json:"notes"             csv:"notes"`
}

// Config holds application configuration loaded from JSON and CLI flags.
type Config struct {
	DefaultMode   string   `json:"defaultMode"`
	DefaultOutput string   `json:"defaultOutput"`
	OutputDir     string   `json:"outputDir"`
	ExcludeDirs   []string `json:"excludeDirs"`
	Notes         string   `json:"notes"`
}

// DefaultConfig returns a Config with sensible built-in defaults.
func DefaultConfig() Config {

	return Config{
		DefaultMode:   constants.ModeHTTPS,
		DefaultOutput: constants.OutputTerminal,
		OutputDir:     constants.DefaultOutputDir,
		ExcludeDirs:   []string{},
		Notes:         "",
	}
}

// CloneResult tracks the outcome of a single clone operation.
type CloneResult struct {
	Record  ScanRecord
	Success bool
	Error   string
}

// CloneSummary aggregates results of a batch clone operation.
type CloneSummary struct {
	Succeeded int
	Failed    int
	Cloned    []CloneResult
	Errors    []CloneResult
}

// ScanCache stores the flags used for the last scan so rescan can replay them.
type ScanCache struct {
	Dir           string `json:"dir"`
	ConfigPath    string `json:"configPath"`
	Mode          string `json:"mode"`
	Output        string `json:"output"`
	OutFile       string `json:"outFile"`
	OutputPath    string `json:"outputPath"`
	GithubDesktop bool   `json:"githubDesktop"`
	OpenFolder    bool   `json:"openFolder"`
	Quiet         bool   `json:"quiet"`
}
