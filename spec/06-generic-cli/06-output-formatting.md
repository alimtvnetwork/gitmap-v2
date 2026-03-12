# Output Formatting

## Multi-Format Output Strategy

When a command runs, produce **all** output formats in one pass.
Generate everything and let the user pick what they need.

| Format | Destination | Purpose |
|--------|-------------|---------|
| Terminal (colored) | stdout | Immediate human feedback |
| CSV | file | Spreadsheet / data import |
| JSON | file | Machine-readable, re-import |
| Markdown | file | Documentation / review |
| Scripts | file | Automation / re-execution |

### Output Directory

All file outputs go into a dedicated subfolder:

```
target-dir/
├── project-a/
└── toolname-output/
    ├── data.csv
    ├── data.json
    ├── structure.md
    └── scripts/
```

## Terminal Output

### ANSI Color Codes

All color codes live in `constants`:

```go
const (
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorCyan   = "\033[36m"
    ColorDim    = "\033[90m"
)
```

| Element | Color | Purpose |
|---------|-------|---------|
| Banner/headers | Cyan | Visual identity |
| Success markers (✓) | Green | Confirmed items |
| Warnings (⚠) | Yellow | Non-fatal issues |
| Data values | White | Primary content |
| Metadata | Dim/Gray | Secondary info |

### Terminal Report Sections

1. **Banner** — tool name + version + item count
2. **Item list** — each item with icon, path, key data
3. **Tree visualization** — hierarchical folder structure
4. **Output file list** — what files were generated
5. **Action instructions** — step-by-step next actions
6. **Related commands** — other commands the user can run

### Banner Pattern

```
╔══════════════════════════════════════╗
║            toolname v1.0.0          ║
╚══════════════════════════════════════╝
  ✓ Found 12 items
```

## CSV Output

- Always include a header row.
- Consistent column ordering: name, identifiers, metadata, paths.
- Quote fields with commas or special characters.
- Use standard library CSV writer.

## JSON Output

- 2-space indentation for readability.
- Array of record objects.
- Directly re-importable by the tool's import command.

## Markdown Output — Tree Visualization

```
├── project-a/
│   └── 📦 **service** (`main`) — https://example.com/service.git
└── project-b/
    └── 📦 **frontend** (`main`) — https://example.com/frontend.git
```

## Template-Based Script Generation

Use `go:embed` templates for complex script outputs:

```go
//go:embed templates/restore.ps1.tmpl
var restoreTemplate string

func WriteRestoreScript(w io.Writer, data RestoreData) error {
    tmpl := template.Must(template.New("restore").Parse(restoreTemplate))

    return tmpl.Execute(w, data)
}
```

## Formatter Package Structure

```
formatter/
├── terminal.go       Colored stdout output
├── csv.go            CSV file output
├── json.go           JSON file output
├── structure.go      Markdown folder tree
├── template.go       Shared template loading
└── templates/        Embedded .tmpl files
```

### Rules

- Each format has its own file.
- All formatters accept `io.Writer` as first argument (testable).
- Templates embedded via `go:embed`, not loaded from disk.
- No format string literals — use `constants`.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
