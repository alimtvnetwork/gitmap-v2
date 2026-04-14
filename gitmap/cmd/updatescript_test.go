package cmd

import (
	"strings"
	"testing"
)

func TestBuildUpdateScriptUsesPathAwareDeployVerification(t *testing.T) {
	script := buildUpdateScript(`D:\src\gitmap`, `D:\src\gitmap\run.ps1`)

	checks := []string{
		`$configDeployedBinary = $null`,
		`$activeCmdForDeploy = Get-Command gitmap -ErrorAction SilentlyContinue`,
		`$deployedBinary = Join-Path $effectiveDeployTarget "gitmap\gitmap.exe"`,
		`powershell.json points to an older deploy location; using PATH-derived target for verification.`,
	}

	for _, check := range checks {
		if !strings.Contains(script, check) {
			t.Fatalf("expected script to contain %q", check)
		}
	}
}