// Package cmd — projectreposoutput.go formats project query output.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/gitmap/model"
)

// printProjectsTerminal prints projects in terminal format.
func printProjectsTerminal(projects []model.DetectedProject) {
	for _, p := range projects {
		fmt.Printf("  %-6s %s\n", p.ProjectType, p.ProjectName)
		fmt.Printf("         Path: %s\n", p.AbsolutePath)
		fmt.Printf("         Indicator: %s\n\n", p.PrimaryIndicator)
	}
}

// printProjectsJSON prints projects as formatted JSON.
func printProjectsJSON(projects []model.DetectedProject) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(projects)
}
