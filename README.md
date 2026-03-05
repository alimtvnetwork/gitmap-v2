# Project Repository

This repository contains two main components:

## 1. GitMap CLI (Go)

A command-line tool that scans directory trees for Git repositories, extracts clone URLs and branch info, and outputs structured data. Every scan produces **all outputs** automatically:

- **Terminal** — formatted table to stdout
- **CSV** — `gitmap.csv`
- **JSON** — `gitmap.json`
- **Folder Structure** — `folder-structure.md` (tree view of discovered repos)

All files are written to a `gitmap-output/` folder at the root of the scanned directory.

**→ [GitMap Documentation](./gitmap/README.md)**  
**→ [Specifications](./spec/01-app/)**

### Quick Start

```powershell
# From the repo root:
.\run.ps1                # Pull, build, deploy to E:\bin-run
.\run.ps1 -Run           # Build + run on parent folder
.\run.ps1 -Run -RunPath "D:\repos"  # Build + run on specific path
```

Or manually:

```bash
cd gitmap
go build -o ../bin/gitmap.exe .
```

### Usage

```bash
gitmap scan ./projects                          # Outputs terminal + CSV + JSON + folder-structure.md
gitmap scan ./projects --mode ssh               # Use SSH URLs
gitmap scan ./projects --github-desktop         # Also add repos to GitHub Desktop
gitmap clone ./gitmap-output/gitmap.json --target-dir ./restored
```

---

## 2. Web Frontend (React + Vite)

A React + TypeScript + Tailwind CSS web application scaffold.

### Setup

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
