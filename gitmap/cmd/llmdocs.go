// Package cmd — llmdocs.go generates a consolidated LLM.md reference file.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// runLLMDocs generates LLM.md in the current working directory.
func runLLMDocs(args []string) {
	checkHelp("llm-docs", args)
	fmt.Print(constants.MsgLLMDocsGenning)

	content := buildLLMDocument()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsWrite, err)
		os.Exit(1)
	}

	outPath := filepath.Join(wd, "LLM.md")

	if writeErr := os.WriteFile(outPath, []byte(content), 0o644); writeErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsWrite, writeErr)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgLLMDocsWritten, outPath)
}

// buildLLMDocument assembles the complete LLM.md content dynamically.
func buildLLMDocument() string {
	var sb strings.Builder

	writeLLMHeader(&sb)
	writeLLMArchitecture(&sb)
	writeLLMCommands(&sb)
	writeLLMGlobalFlags(&sb)
	writeLLMCodingConventions(&sb)
	writeLLMProjectStructure(&sb)
	writeLLMDatabase(&sb)
	writeLLMInstallation(&sb)
	writeLLMPatterns(&sb)

	return sb.String()
}
