package constants

// CD CLI commands.
const (
	CmdCD      = "cd"
	CmdCDAlias = "go"
)

// CD subcommands.
const (
	CmdCDRepos        = "repos"
	CmdCDSetDefault   = "set-default"
	CmdCDClearDefault = "clear-default"
)

// CD help text.
const HelpCD = "  cd (go) <name>      Navigate to a tracked repo directory"

// CD file.
const CDDefaultsFile = "cd-defaults.json"

// CD messages.
const (
	MsgCDMultipleHeader  = "Multiple locations found for \"%s\":\n"
	MsgCDMultipleRowFmt  = "  %d  %s\n"
	MsgCDPickPrompt      = "\nPick [1-%d]: "
	MsgCDReposHeader     = "TRACKED REPOS\n"
	MsgCDReposRowFmt     = "  %d  %s\n"
	MsgCDDefaultSet      = "Default set for %s: %s\n"
	MsgCDDefaultCleared  = "Default cleared for %s\n"
	ErrCDUsage           = "usage: gitmap cd <repo-name|repos> [--group <name>] [--pick]\n"
	ErrCDNotFound        = "no repo found matching '%s'\n"
	ErrCDInvalidPick     = "invalid selection\n"
	ErrCDSetDefaultUsage = "usage: gitmap cd set-default <name> <path>\n"
	ErrCDClearDefaultUsage = "usage: gitmap cd clear-default <name>\n"
	ErrCDDefaultNotFound = "no default set for '%s'\n"
)

// CD flag descriptions.
const (
	FlagDescCDGroup = "Filter repos list by group"
	FlagDescCDPick  = "Force interactive picker even if a default is set"
)

// CD shell wrapper functions — installed by setup/completion.
const CDFuncMarker = "# gitmap shell wrapper v2"

// CDFuncBash installs gitmap and gcd wrappers for Bash.
const CDFuncBash = `gcd() {
  local dest status
  dest="$(command gitmap cd "$@")"
  status=$?
  if [ $status -ne 0 ]; then
    return $status
  fi
  if [ -n "$dest" ] && [ -d "$dest" ]; then
    builtin cd "$dest" || return $?
  fi
}

gitmap() {
  if [ "$1" = "cd" ] || [ "$1" = "go" ]; then
    local dest status
    dest="$(command gitmap "$@")"
    status=$?
    if [ $status -ne 0 ]; then
      return $status
    fi
    if [ -n "$dest" ] && [ -d "$dest" ]; then
      builtin cd "$dest" || return $?
    fi
    return 0
  fi
  command gitmap "$@"
}`

// CDFuncZsh installs gitmap and gcd wrappers for Zsh.
const CDFuncZsh = `gcd() {
  local dest
  local status
  dest="$(command gitmap cd "$@")"
  status=$?
  if (( status != 0 )); then
    return $status
  fi
  if [[ -n "$dest" && -d "$dest" ]]; then
    builtin cd "$dest" || return $?
  fi
}

gitmap() {
  if [[ "$1" == "cd" || "$1" == "go" ]]; then
    local dest
    local status
    dest="$(command gitmap "$@")"
    status=$?
    if (( status != 0 )); then
      return $status
    fi
    if [[ -n "$dest" && -d "$dest" ]]; then
      builtin cd "$dest" || return $?
    fi
    return 0
  fi
  command gitmap "$@"
}`

// CDFuncPowerShell installs gitmap and gcd wrappers for PowerShell.
const CDFuncPowerShell = `function gcd {
  $real = Get-GitmapCommand
  if (-not $real) {
    Write-Error "gitmap executable not found"
    return
  }
  $dest = & $real cd @args
  if ($LASTEXITCODE -ne 0) {
    return
  }
  if ($dest -and (Test-Path -LiteralPath $dest)) {
    Set-Location -LiteralPath $dest
  }
}

function Get-GitmapCommand {
  $cmd = Get-Command gitmap.exe -CommandType Application -ErrorAction SilentlyContinue
  if ($cmd) {
    return $cmd.Source
  }
  $cmd = Get-Command gitmap -CommandType Application -ErrorAction SilentlyContinue
  if ($cmd) {
    return $cmd.Source
  }
  return $null
}

function gitmap {
  $real = Get-GitmapCommand
  if (-not $real) {
    Write-Error "gitmap executable not found"
    return
  }
  if ($args.Count -gt 0 -and ($args[0] -eq 'cd' -or $args[0] -eq 'go')) {
    $dest = & $real @args
    if ($LASTEXITCODE -ne 0) {
      return
    }
    if ($dest -and (Test-Path -LiteralPath $dest)) {
      Set-Location -LiteralPath $dest
    }
    return
  }
  & $real @args
}`

// CD function messages.
const (
	MsgCDFuncInstalled = "Installed 'gitmap'/'gcd' shell wrappers — restart your terminal or source your profile\n"
	MsgCDFuncAlready   = "Shell wrappers for 'gitmap'/'gcd' already installed\n"
)
