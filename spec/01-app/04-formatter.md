# Formatter

## Responsibility

Render a list of `ScanRecord` into multiple output formats.

## Output Behavior

Every scan **always produces all outputs simultaneously**:

1. **Terminal** — colored, professional output to stdout.
2. **CSV** — `gitmap.csv` written to the output directory.
3. **JSON** — `gitmap.json` written to the output directory.
4. **Folder Structure** — `folder-structure.md` written to the output directory.

The output directory defaults to `gitmap-output/` inside the scanned directory.

## Formats

### Terminal

Colored output with ANSI codes showing:
- Banner with repo count
- Each repo: name (📦), path, and clone instruction
- Folder tree with 📁 folders and 📦 repos with branch names
- Clone help: step-by-step instructions for cloning on another machine

### CSV

Write a CSV file with headers:

```
repoName,httpsUrl,sshUrl,branch,relativePath,absolutePath,cloneInstruction,notes
```

### JSON

Write a JSON array of `ScanRecord` objects with 2-space indentation.

### Folder Structure (Markdown)

Write a tree view of discovered repos:

```markdown
# Folder Structure

Git repositories discovered by gitmap.

├── 📦 **my-app** (`main`) — https://github.com/user/my-app.git
├── libs/
│   └── 📦 **utils** (`main`) — https://github.com/user/utils.git
└── 📦 **docs** (`main`) — https://github.com/user/docs.git
```

## Output Location

- Terminal: stdout.
- CSV/JSON/Markdown: `gitmap-output/` inside the scanned directory,
  or path from `--output-path` flag, or exact path from `--out-file`.
