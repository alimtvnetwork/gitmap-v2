# gitmap changelog

Display release notes for a specific version or the latest release.

## Alias

cl

## Usage

    gitmap changelog [version] [--open] [--source]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --latest | false | Show latest release notes |
| --open | false | Open changelog in browser |
| --source | false | Show source file path |

## Prerequisites

- Must be inside a Git repository with release metadata

## Examples

### Example 1: Show latest changelog

    gitmap changelog --latest

**Output:**

    v2.3.10 — 2025-03-10
    - Add user auth endpoint
    - Fix pagination bug
    - Update README

### Example 2: Show changelog for a specific version

    gitmap cl v2.3.7

**Output:**

    v2.3.7 — 2025-02-28
    - Initial release branch support
    - Add stats command
