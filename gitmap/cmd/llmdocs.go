// Package cmd — llmdocs.go generates a consolidated LLM.md reference file.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// runLLMDocs generates LLM.md or prints to stdout with --stdout.
func runLLMDocs(args []string) {
	checkHelp("llm-docs", args)

	fs := flag.NewFlagSet("llm-docs", flag.ExitOnError)
	toStdout := fs.Bool(constants.FlagLLMDocsStdout, false, constants.FlagDescLLMDocsStdout)

	reordered := reorderFlagsBeforeArgs(args, nil)

	if err := fs.Parse(reordered); err != nil {
		fmt.Fprintf(os.Stderr, "llm-docs: %v\n", err)
		os.Exit(1)
	}

	content := buildLLMDocument()

	if *toStdout {
		fmt.Print(content)

		return
	}

	fmt.Print(constants.MsgLLMDocsGenning)

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
