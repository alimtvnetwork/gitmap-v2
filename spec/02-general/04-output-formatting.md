# Output & Formatting Patterns

## Overview

This document describes reusable patterns for CLI tools that produce
multiple output formats simultaneously from a single data source.

## Multi-Format Output Strategy

### Principle: All Formats, Every Time

When a command runs, produce **all** output formats in one pass.
Don't make the user choose — generate everything and let them
pick what they need.

| Format | Destination | Purpose |
|--------|-------------|---------|
| Terminal (colored) | stdout | Immediate human feedback |
| CSV | file | Spreadsheet / data import |
| JSON | file | Machine-readable, re-import |
| Markdown | file | Documentation / review |
| Scripts | file | Automation / re-execution |

### Output Directory

All file outputs go into a dedicated subfolder (e.g., `toolname-output/`)
inside the scanned/processed directory. This keeps outputs organized
and avoids polluting the working directory.

```
target-dir/
├── project-a/
├── project-b/
└── toolname-output/       ← all outputs here
    ├── data.csv
    ├── data.json
    ├── structure.md
    └── scripts/
```

## Terminal Output

### Colored, Structured Reports

Use ANSI codes for visual hierarchy:

| Element | Color | Purpose |
|---------|-------|---------|
| Banner/headers | Cyan | Visual identity |
| Success markers | Green | Confirmed items |
| Warnings | Yellow | Non-fatal issues |
| Data values | White | Primary content |
| Metadata | Dim/Gray | Secondary info |

### All color codes live in `constants`:

```go
const (
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorCyan   = "\033[36m"
    ColorDim    = "\033[90m"
)
```

### Terminal Report Sections

Structure terminal output as distinct sections:

1. **Banner** — tool name + version + repo count
2. **Item list** — each item with icon, path, and key data
3. **Tree visualization** — hierarchical folder structure
4. **Output file list** — what files were generated and where
5. **Next steps** — how to use the output

### Banner Pattern

```
╔══════════════════════════════════════╗
║            toolname v1.0.0          ║
╚══════════════════════════════════════╝
  ✓ Found 12 items
```

## Template-Based Script Generation

### Approach: `go:embed` Templates

For complex script outputs (PowerShell, Bash), use Go's embedded
templates rather than string concatenation:

```go
//go:embed templates/clone.ps1.tmpl
var cloneTemplate string

func WriteCloneScript(w io.Writer, data CloneData) error {
    tmpl := template.Must(template.New("clone").Parse(cloneTemplate))
    return tmpl.Execute(w, data)
}
```

### Template Data Structures

Define clear data structures for template rendering:

```go
type CloneData struct {
    Repos     []RepoEntry
    BaseDir   string
    TotalCount int
}

type RepoEntry struct {
    URL    string
    Branch string
    Path   string
    Name   string
}
```

### Script Categories

| Type | Content | Use Case |
|------|---------|----------|
| Logic scripts | Progress bars, error handling, summaries | Interactive restoration |
| Direct scripts | Raw commands, no logic | Quick copy-paste execution |
| Registration scripts | Tool-specific integrations | GitHub Desktop, etc. |

## CSV Output

### Conventions

- Always include a header row.
- Use consistent column ordering: name, URLs, branch, paths, metadata.
- Quote fields that may contain commas or special characters.
- Use standard Go `encoding/csv` writer.

```
repoName,httpsUrl,sshUrl,branch,relativePath,absolutePath,cloneInstruction,notes
```

## JSON Output

### Conventions

- Use 2-space indentation for readability.
- Output an array of record objects.
- Field names match the Go struct's `json` tags.
- The JSON output should be directly re-importable by the tool's
  clone/restore command.

```go
encoder := json.NewEncoder(w)
encoder.SetIndent("", constants.JSONIndent)
encoder.Encode(records)
```

## Markdown Output

### Folder Structure Visualization

Render a tree using Unicode box-drawing characters:

```
├── project-a/
│   ├── 📦 **service** (`main`) — git@github.com:user/service.git
│   └── 📦 **api** (`develop`) — git@github.com:user/api.git
└── project-b/
    └── 📦 **frontend** (`main`) — https://github.com/user/frontend.git
```

| Character | Constant | Usage |
|-----------|----------|-------|
| `├──` | `TreeBranch` | Non-last child |
| `└──` | `TreeCorner` | Last child |
| `│   ` | `TreePipe` | Vertical continuation |
| `    ` | `TreeSpace` | No continuation |

## Formatter Package Structure

```
formatter/
├── terminal.go       Terminal (colored stdout)
├── csv.go            CSV file output
├── json.go           JSON file output
├── structure.go      Markdown folder tree
├── clonescript.go    Logic-based clone script
├── directclone.go   Raw clone commands
├── desktopscript.go GitHub Desktop registration
├── template.go       Shared template loading
└── templates/        Embedded .tmpl files
    ├── clone.ps1.tmpl
    ├── direct-clone.ps1.tmpl
    ├── direct-clone-ssh.ps1.tmpl
    └── desktop.ps1.tmpl
```

### Rules

- Each format has its own file.
- All formatters accept `io.Writer` as first argument (testable).
- Templates are embedded via `go:embed`, not loaded from disk.
- No format string literals in formatter files — use `constants`.
