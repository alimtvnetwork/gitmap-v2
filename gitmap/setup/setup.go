// Package setup configures Git global settings from a JSON config file.
package setup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// GitSetupConfig holds the full git-setup.json structure.
type GitSetupConfig struct {
	DiffTool         *ToolConfig       `json:"diffTool"`
	MergeTool        *ToolConfig       `json:"mergeTool"`
	Aliases          map[string]string `json:"aliases"`
	CredentialHelper string            `json:"credentialHelper"`
	Core             map[string]string `json:"core"`
}

// ToolConfig holds diff/merge tool configuration.
type ToolConfig struct {
	Name          string `json:"name"`
	Cmd           string `json:"cmd"`
	TrustExitCode bool   `json:"trustExitCode"`
}

// SetupResult tracks applied and failed settings.
type SetupResult struct {
	Applied int
	Skipped int
	Failed  int
	Errors  []string
}

// LoadConfig reads and parses the git-setup.json file.
func LoadConfig(path string) (GitSetupConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GitSetupConfig{}, err
	}
	var cfg GitSetupConfig
	err = json.Unmarshal(data, &cfg)

	return cfg, err
}

// Apply applies the full git setup configuration.
func Apply(cfg GitSetupConfig, dryRun bool) SetupResult {
	result := SetupResult{}

	if cfg.DiffTool != nil {
		applyDiffTool(cfg.DiffTool, dryRun, &result)
	}
	if cfg.MergeTool != nil {
		applyMergeTool(cfg.MergeTool, dryRun, &result)
	}
	if len(cfg.Aliases) > 0 {
		applyAliases(cfg.Aliases, dryRun, &result)
	}
	if len(cfg.CredentialHelper) > 0 {
		applyCredentialHelper(cfg.CredentialHelper, dryRun, &result)
	}
	if len(cfg.Core) > 0 {
		applyCoreSettings(cfg.Core, dryRun, &result)
	}

	return result
}

// applyDiffTool configures git's global diff tool.
func applyDiffTool(tool *ToolConfig, dryRun bool, r *SetupResult) {
	settings := []gitSetting{
		{"diff.tool", tool.Name},
		{fmt.Sprintf("difftool.%s.cmd", tool.Name), tool.Cmd},
		{"difftool.prompt", "false"},
	}
	if tool.TrustExitCode {
		settings = append(settings, gitSetting{
			fmt.Sprintf("difftool.%s.trustExitCode", tool.Name), "true",
		})
	}
	applySection(constants.SetupSectionDiff, settings, dryRun, r)
}

// applyMergeTool configures git's global merge tool.
func applyMergeTool(tool *ToolConfig, dryRun bool, r *SetupResult) {
	settings := []gitSetting{
		{"merge.tool", tool.Name},
		{fmt.Sprintf("mergetool.%s.cmd", tool.Name), tool.Cmd},
		{"mergetool.prompt", "false"},
		{"mergetool.keepBackup", "false"},
	}
	if tool.TrustExitCode {
		settings = append(settings, gitSetting{
			fmt.Sprintf("mergetool.%s.trustExitCode", tool.Name), "true",
		})
	}
	applySection(constants.SetupSectionMerge, settings, dryRun, r)
}

// applyAliases configures git global aliases.
func applyAliases(aliases map[string]string, dryRun bool, r *SetupResult) {
	settings := make([]gitSetting, 0, len(aliases))
	for name, value := range aliases {
		settings = append(settings, gitSetting{
			fmt.Sprintf("alias.%s", name), value,
		})
	}
	applySection(constants.SetupSectionAlias, settings, dryRun, r)
}

// applyCredentialHelper configures git's credential helper.
func applyCredentialHelper(helper string, dryRun bool, r *SetupResult) {
	settings := []gitSetting{
		{"credential.helper", helper},
	}
	applySection(constants.SetupSectionCred, settings, dryRun, r)
}

// applyCoreSettings configures git core settings.
func applyCoreSettings(core map[string]string, dryRun bool, r *SetupResult) {
	settings := make([]gitSetting, 0, len(core))
	for key, value := range core {
		gitKey := mapCoreKey(key)
		settings = append(settings, gitSetting{gitKey, value})
	}
	applySection(constants.SetupSectionCore, settings, dryRun, r)
}

// gitSetting is a key-value pair for git config.
type gitSetting struct {
	Key   string
	Value string
}

// applySection applies a group of settings and prints results.
func applySection(section string, settings []gitSetting, dryRun bool, r *SetupResult) {
	fmt.Printf("\n  %s■ %s%s\n", constants.ColorYellow, section, constants.ColorReset)
	for _, s := range settings {
		applyOne(s, dryRun, r)
	}
}

// applyOne applies a single git config --global setting.
func applyOne(s gitSetting, dryRun bool, r *SetupResult) {
	if dryRun {
		fmt.Printf("  %s[dry-run]%s git config --global %s %q\n",
			constants.ColorDim, constants.ColorReset, s.Key, s.Value)
		r.Skipped++

		return
	}

	current := getCurrentValue(s.Key)
	if current == s.Value {
		fmt.Printf("  %s⊘ %s%s = %s (already set)\n",
			constants.ColorDim, s.Key, constants.ColorReset, s.Value)
		r.Skipped++

		return
	}

	cmd := exec.Command(constants.GitBin, constants.GitConfigCmd, "--global", s.Key, s.Value)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s — %s", s.Key, err, strings.TrimSpace(string(out)))
		fmt.Printf("  %s✗ %s%s\n", constants.ColorYellow, errMsg, constants.ColorReset)
		r.Failed++
		r.Errors = append(r.Errors, errMsg)

		return
	}

	fmt.Printf("  %s✓ %s%s = %s\n", constants.ColorGreen, s.Key, constants.ColorReset, s.Value)
	r.Applied++
}

// getCurrentValue reads the current git config value for a key.
func getCurrentValue(key string) string {
	cmd := exec.Command(constants.GitBin, constants.GitConfigCmd, "--global", constants.GitGetFlag, key)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// mapCoreKey maps JSON-friendly keys to git config keys.
func mapCoreKey(jsonKey string) string {
	coreMap := map[string]string{
		"autocrlf":      "core.autocrlf",
		"longpaths":     "core.longpaths",
		"editor":        "core.editor",
		"safecrlf":      "core.safecrlf",
		"defaultBranch": "init.defaultBranch",
	}
	if mapped, ok := coreMap[jsonKey]; ok {
		return mapped
	}

	return "core." + jsonKey
}
