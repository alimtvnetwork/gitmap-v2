# Verbose Logging

## Purpose

A shared `--verbose` flag enables detailed debug logging to a timestamped
file. Normal runs produce clean user-facing output only. Verbose runs
capture full diagnostics for troubleshooting without polluting stdout.

---

## Design Rules

| Rule | Detail |
|------|--------|
| Off by default | No log file created unless `--verbose` is passed |
| File + stderr | Every verbose entry writes to both the log file and stderr |
| Timestamped entries | Each line prefixed with `[HH:MM:SS.mmm]` |
| Timestamped filenames | Log file named `toolname-verbose-YYYY-MM-DD_HH-mm-ss.log` |
| Output directory | Logs written to the tool's default output folder |
| Dim on stderr | Verbose stderr output uses dim/gray ANSI color |
| No stdout pollution | Verbose output never mixes with normal command output |
| Global singleton | One logger instance shared across all packages |

---

## Package Structure

```
verbose/
└── verbose.go     Logger type, Init, Close, Log, IsEnabled, Get
```

Single file. No sub-packages. No external dependencies beyond `constants`.

---

## Logger API

```go
// Init creates the log file and enables verbose logging.
// Call once at startup when --verbose is set.
func Init() (*Logger, error)

// Close flushes and closes the log file.
func (l *Logger) Close()

// Log writes a formatted message to the log file and stderr.
func (l *Logger) Log(format string, args ...interface{})

// IsEnabled returns true if verbose mode is active.
func IsEnabled() bool

// Get returns the global logger (may be nil).
func Get() *Logger
```

---

## Logger Type

```go
type Logger struct {
    file    *os.File
    enabled bool
}

var global *Logger
```

- `file` — open handle to the log file
- `enabled` — guards all write operations
- `global` — package-level singleton set by `Init()`

---

## Init Flow

```go
func Init() (*Logger, error) {
    logDir := constants.DefaultOutputFolder
    _ = os.MkdirAll(logDir, constants.DirPermission)

    timestamp := time.Now().Format("2006-01-02_15-04-05")
    logPath := filepath.Join(logDir, fmt.Sprintf(constants.VerboseLogFileFmt, timestamp))

    file, err := os.Create(logPath)
    if err != nil {
        return nil, err
    }

    l := &Logger{file: file, enabled: true}
    global = l
    fmt.Printf(constants.MsgVerboseLogFile, logPath)

    return l, nil
}
```

**Key points:**

- Creates the output directory if missing (no error on existing)
- Prints the log file path to stdout so the user knows where to find it
- Returns both the logger and any error — caller decides whether to abort

---

## Log Entry Format

```go
func writeLogEntry(l *Logger, format string, args ...interface{}) {
    line := fmt.Sprintf(format, args...)
    ts := time.Now().Format("15:04:05.000")
    entry := fmt.Sprintf("[%s] %s\n", ts, line)
    l.file.WriteString(entry)
    fmt.Fprint(os.Stderr, constants.ColorDim+entry+constants.ColorReset)
}
```

**Example output:**

```
[14:32:07.123] git clone https://github.com/user/repo.git
[14:32:09.456] clone completed in 2.3s
[14:32:09.460] retry attempt 1/4 for locked file
```

---

## Flag Registration

The `--verbose` flag is a **global flag** registered on the root command,
not per-subcommand.

```go
// In cmd/rootflags.go
fs.BoolVar(&verboseFlag, constants.FlagVerbose, false, constants.FlagDescVerbose)
```

---

## Command Handler Pattern

Every command that supports verbose logging follows this pattern:

```go
func runPull(args []string) {
    checkHelp("pull", args)
    slug, group, all, verboseFlag := parsePullFlags(args)

    if verboseFlag {
        initVerboseLog()       // Init + defer Close
    }

    // ... command logic ...
}

func initVerboseLog() {
    logger, err := verbose.Init()
    if err != nil {
        fmt.Fprintf(os.Stderr, constants.ErrVerboseInit, err)
        return                 // Non-fatal — continue without logging
    }
    defer logger.Close()
}
```

**Rules:**

- Verbose init failure is **non-fatal** — warn and continue
- `defer Close()` in the same function that calls `Init()`
- Never pass the logger as a parameter — use `verbose.Get()` or `verbose.IsEnabled()`

---

## What to Log

| Category | Examples |
|----------|----------|
| Git operations | Clone/pull commands, remote URLs, branch names |
| Retry attempts | Attempt number, delay, reason for retry |
| File I/O | Paths read/written, file sizes, permissions |
| External processes | Command lines, exit codes, stdout/stderr |
| Timing | Operation durations, elapsed time |
| Environment | OS, paths, config values loaded |
| Errors (detailed) | Full error chains, stack context |

**What NOT to log:**

- Secrets, tokens, or credentials
- Routine success paths that add no diagnostic value
- Data that duplicates normal stdout output

---

## Constants

All verbose-related literals live in the constants package:

```go
// constants/constants.go
const VerboseLogFileFmt = "toolname-verbose-%s.log"

// constants/constants_cli.go
const FlagVerbose    = "verbose"
const FlagDescVerbose = "Enable verbose debug logging to file"

// constants/constants_messages.go
const MsgVerboseLogFile = "Verbose log: %s\n"
const ErrVerboseInit    = "Warning: could not initialize verbose log: %v\n"
```

---

## Conditional Logging in Libraries

Domain packages (scanner, cloner, mapper) check `verbose.IsEnabled()`
before calling `verbose.Get().Log(...)`:

```go
func safePullOne(repo model.Record) error {
    logger := verbose.Get()

    if logger != nil {
        logger.Log("pulling %s at %s", repo.Name, repo.Path)
    }

    // ... pull logic ...

    if logger != nil {
        logger.Log("pull complete for %s (%.1fs)", repo.Name, elapsed.Seconds())
    }

    return nil
}
```

**Rules:**

- Always nil-check `verbose.Get()` — verbose may not be active
- Keep log calls outside hot loops to avoid performance overhead
- Use `fmt.Sprintf`-style formatting — no structured logging libraries

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
