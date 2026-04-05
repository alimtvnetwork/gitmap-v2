# Tool Installer

Install a developer tool by name using the platform package manager.

## Alias

in

## Usage

    gitmap install <tool> [flags]

## Flags

| Flag      | Default | Description                                        |
|-----------|---------|----------------------------------------------------|
| --manager | (auto)  | Force package manager (choco, winget, apt, brew)   |
| --version | latest  | Install a specific version                         |
| --verbose | false   | Show full installer output                         |
| --dry-run | false   | Show install command without executing             |
| --check   | false   | Only check if tool is installed                    |
| --list    | false   | List all supported tools                           |

## Supported Tools

| Tool            | Binary         | Description                      |
|-----------------|----------------|----------------------------------|
| vscode          | code           | Visual Studio Code editor        |
| node            | node           | Node.js JavaScript runtime       |
| yarn            | yarn           | Yarn package manager             |
| bun             | bun            | Bun JavaScript runtime           |
| pnpm            | pnpm           | pnpm package manager             |
| python          | python3        | Python programming language      |
| go              | go             | Go programming language          |
| git             | git            | Git version control              |
| git-lfs         | git-lfs        | Git Large File Storage           |
| gh              | gh             | GitHub CLI                       |
| github-desktop  | —              | GitHub Desktop application       |
| cpp             | g++            | C++ compiler (MinGW/g++)         |
| php             | php            | PHP programming language         |
| powershell      | pwsh           | PowerShell shell                 |

## Prerequisites

- Windows: Chocolatey or Winget in PATH
- Linux: apt, dnf, or pacman available
- macOS: Homebrew installed

## Examples

### Install a tool end-to-end

    $ gitmap install vscode
      Checking if vscode is installed...
      Installing vscode...
      Verifying vscode installation...
      vscode installed successfully.

### Check if a tool is already installed

    $ gitmap in go --check
      Checking if go is installed...
      go is already installed (version: go version go1.22.4 linux/amd64)

    $ gitmap in bun --check
      Checking if bun is installed...
      bun is not installed.

### Preview install command with dry-run

    $ gitmap install python --dry-run
      Checking if python is installed...
      [dry-run] Would run: choco install python -y

    $ gitmap install node --manager brew --dry-run
      Checking if node is installed...
      [dry-run] Would run: brew install node

## See Also

- [env](env.md) — Manage environment variables and PATH
- [doctor](doctor.md) — Diagnose PATH and version issues
- [setup](setup.md) — Configure Git global settings
