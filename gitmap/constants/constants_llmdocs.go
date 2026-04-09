package constants

// llm-docs command.
const (
	CmdLLMDocs      = "llm-docs"
	CmdLLMDocsAlias = "ld"
)

// llm-docs messages.
const (
	MsgLLMDocsWritten = "  ✓ LLM.md written to %s\n"
	MsgLLMDocsGenning = "  ↻ Generating LLM.md from command registry...\n"
	ErrLLMDocsWrite   = "  ✗ Could not write LLM.md: %v\n"
	HelpLLMDocs       = "  llm-docs (ld)       Generate LLM.md reference for AI assistants"
)
