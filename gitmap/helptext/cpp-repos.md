# gitmap cpp-repos

List all detected C++ projects across tracked repositories.

## Alias

cr

## Usage

    gitmap cpp-repos [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap scan` first to detect projects (see scan.md)

## Examples

### Example 1: List C++ projects

    gitmap cpp-repos

**Output:**

    game-engine  CMake     ~/projects/game-engine
    codec-lib    Makefile  ~/work/codec-lib
    2 C++ projects detected

### Example 2: JSON output

    gitmap cr --json

**Output:**

    [{"repo":"game-engine","build_system":"CMake"}]