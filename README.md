# GitMap

A command-line tool that scans directory trees for Git repositories, extracts clone URLs and branch info, and outputs structured data. Every scan produces **all outputs** automatically:

- **Terminal** — formatted table to stdout
- **CSV** — `gitmap.csv`
- **JSON** — `gitmap.json`
- **Folder Structure** — `folder-structure.md` (tree view of discovered repos)

All files are written to a `gitmap-output/` folder at the root of the scanned directory.

**→ [GitMap Documentation](./gitmap/README.md)**
**→ [Specifications](./spec/01-app/)**

## Quick Start

```powershell
# From the repo root:
.\run.ps1                # Pull, build, deploy to E:\bin-run
.\run.ps1 -R scan        # Build + scan parent folder
.\run.ps1 -R scan D:\repos --mode ssh   # Build + scan with SSH mode
```

Or manually:

```bash
cd gitmap
go build -o ../bin/gitmap.exe .
```

## Usage

```bash
gitmap scan ./projects                          # Outputs terminal + CSV + JSON + folder-structure.md
gitmap scan ./projects --mode ssh               # Use SSH URLs
gitmap scan ./projects --github-desktop         # Also add repos to GitHub Desktop
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored
gitmap update                                   # Self-update from source repo
gitmap doctor                                   # Diagnose PATH, version, and deployment issues
gitmap changelog --latest                       # Show latest release notes
gitmap lb                                       # Show the most recently updated remote branch
gitmap lb 5                                     # Show the 5 most recently updated remote branches
gitmap lb --json                                # Latest branch as structured JSON
```

## Build & Deploy

The project uses a single PowerShell script (`run.ps1`) at the repo root for the full lifecycle: pull, build, deploy, and optionally run.

| Flag | Description |
|------|-------------|
| `-NoPull` | Skip `git pull` |
| `-NoDeploy` | Skip deploy step |
| `-Update` | Update mode: full pipeline with post-update validation |
| `-R` | Run gitmap after build (trailing args forwarded) |

Configuration lives in `gitmap/powershell.json`.

## Web Frontend (React + Vite)

A React + TypeScript + Tailwind CSS web application scaffold.

```sh
npm i
npm run dev
```

### Tech Stack

- Vite
- TypeScript
- React
- shadcn-ui
- Tailwind CSS

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)

## License

This project is proprietary software. All rights reserved.
