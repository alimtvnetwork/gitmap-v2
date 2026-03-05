# Formatter

## Responsibility

Render a list of `ScanRecord` into multiple output formats.

## Output Behavior

Every scan **always produces all outputs simultaneously**:

1. **Terminal** — colored, professional output to stdout.
2. **CSV** — `gitmap.csv` written to the output directory.
3. **JSON** — `gitmap.json` written to the output directory.
4. **Folder Structure** — `folder-structure.md` written to the output directory.
5. **Clone Script** — `clone.ps1` — self-contained PowerShell script that clones all repos.
6. **Desktop Script** — `register-desktop.ps1` — registers cloned repos with GitHub Desktop.

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

### Clone Script (`clone.ps1`)

A self-contained PowerShell script that:
- Accepts a `-TargetDir` parameter (defaults to `.`)
- Creates the folder structure under the target directory
- Clones each repo with `git clone -b <branch> <url> <path>`
- Shows progress (`[1/N]`, `[2/N]`, …) with colored output
- Prints a summary of succeeded/failed clones

### Desktop Registration Script (`register-desktop.ps1`)

A PowerShell script that:
- Accepts a `-BaseDir` parameter (defaults to `.`)
- Checks if GitHub Desktop CLI (`github`) is available
- Registers each cloned repo with GitHub Desktop
- Shows progress with colored output
- Prints a summary of registered/failed repos

## Output Location

- Terminal: stdout.
- All files: `gitmap-output/` inside the scanned directory,
  or path from `--output-path` flag, or exact path from `--out-file`.

## Output Directory Contents

```
gitmap-output/
├── gitmap.csv
├── gitmap.json
├── folder-structure.md
├── clone.ps1
└── register-desktop.ps1
```
