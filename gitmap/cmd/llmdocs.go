// Package cmd — llmdocs.go generates a consolidated LLM.md reference file.
package cmd

import (
	"encoding/json"
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
	format := fs.String(constants.FlagLLMDocsFormat, "markdown", constants.FlagDescLLMDocsFormat)

	reordered := reorderFlagsBeforeArgs(args)

	if err := fs.Parse(reordered); err != nil {
		fmt.Fprintf(os.Stderr, "llm-docs: %v\n", err)
		os.Exit(1)
	}

	if *format != "markdown" && *format != "json" {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsFormat, *format)
		os.Exit(1)
	}

	content := buildLLMOutput(*format)

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

	ext := ".md"
	if *format == "json" {
		ext = ".json"
	}

	outPath := filepath.Join(wd, "LLM"+ext)

	if writeErr := os.WriteFile(outPath, []byte(content), 0o644); writeErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsWrite, writeErr)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgLLMDocsWritten, outPath)
}

// buildLLMOutput returns the document in the requested format.
func buildLLMOutput(format string) string {
	if format == "json" {
		return buildLLMJSON()
	}

	return buildLLMDocument()
}

// buildLLMJSON assembles a JSON representation of the LLM reference.
func buildLLMJSON() string {
	groups := buildCommandGroups()

	type jsonCmd struct {
		Name    string `json:"name"`
		Alias   string `json:"alias"`
		Desc    string `json:"description"`
		Example string `json:"example,omitempty"`
	}

	type jsonGroup struct {
		Title    string    `json:"title"`
		Commands []jsonCmd `json:"commands"`
	}

	out := make([]jsonGroup, 0, len(groups))

	for _, g := range groups {
		jg := jsonGroup{Title: g.title}

		for _, c := range g.commands {
			jg.Commands = append(jg.Commands, jsonCmd{
				Name:    c.name,
				Alias:   c.alias,
				Desc:    c.desc,
				Example: c.example,
			})
		}

		out = append(out, jg)
	}

	data, _ := json.MarshalIndent(out, "", "  ")

	return string(data) + "\n"
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
