# gitmap scan

Scan a directory tree for Git repositories and record them in the local database.

## Alias

s

## Usage

    gitmap scan [dir] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --config \<path\> | ./data/config.json | Config file path |
| --mode ssh\|https | https | Clone URL style |
| --output csv\|json\|terminal | terminal | Output format |
| --output-path \<dir\> | ./gitmap-output | Output directory |
| --github-desktop | false | Add repos to GitHub Desktop |
| --open | false | Open output folder after scan |
| --quiet | false | Suppress clone help section |

## Prerequisites

- None (this is typically the first command you run)

## Examples

### Example 1: Scan a directory

    gitmap scan ~/projects

**Output:**

    Scanning ~/projects...
    Found 42 repositories
    Output written to ./gitmap-output/
    ✓ gitmap-repos.csv
    ✓ gitmap-repos.json
    ✓ terminal summary

### Example 2: JSON output with SSH URLs

    gitmap scan ~/work --output json --mode ssh

**Output:**

    Scanning ~/work...
    Found 18 repositories
    Output written to ./gitmap-output/
    ✓ gitmap-repos.json (SSH URLs)

### Example 3: Scan and open output folder

    gitmap scan . --open --quiet

**Output:**

    Scanning current directory...
    Found 7 repositories
    Opening ./gitmap-output/
