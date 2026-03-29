# Logging & Observability — Diagnostic Output Standards

## Overview

Standards for structured logging, verbose mode, and diagnostic output
across CLI and application code. Logs are for operators — make them
useful, consistent, and actionable.

---

## 1. Log Levels

Use a clear hierarchy. Never log everything at the same level.

| Level | Purpose | Example |
|-------|---------|---------|
| Error | Operation failed, action required | `"failed to write metadata: permission denied"` |
| Warn | Degraded but recoverable | `"legacy directory found, migrating"` |
| Info | Key lifecycle events | `"scan complete: 42 repos found"` |
| Debug | Detailed internals (verbose only) | `"checking path: /home/user/projects/api"` |

---

## 2. Structured Log Format

Every log entry must include context. Never log bare strings.

### Go

```go
verbose.Log("scan", "discovered repo: %s (branch: %s)", repoName, branch)
verbose.Log("release", "version resolved: %s (source: %s)", version, source)
verbose.Log("clone", "cloned %d/%d repos in %s", completed, total, elapsed)
```

### TypeScript

```ts
console.info("[scan] discovered repo:", { name: repoName, branch });
console.error("[release] upload failed:", { version, statusCode, error });
```

### Required Fields

| Field | When |
|-------|------|
| Component/stage | Always — identifies the subsystem |
| Entity name | When operating on a specific item |
| Count/progress | When processing a batch |
| Duration | When timing an operation |
| Error detail | On failures — include the original error |

---

## 3. Verbose Mode Pattern

Verbose output is opt-in via a global `--verbose` flag and writes
to timestamped files while showing summaries on stderr.

### Architecture

```
User runs: gitmap scan --verbose
  │
  ├─ stderr: colored summary lines (always visible)
  │    "✓ Scanned 42 repos in 1.2s"
  │
  └─ file: .gitmap/output/scan-2025-01-15-143022.log
       [14:30:22] [scan] starting directory walk: /home/user
       [14:30:22] [scan] found .git: /home/user/api/.git
       [14:30:23] [scan] remote: https://github.com/user/api.git
       ...
```

### Rules

- Verbose logs go to files, not stdout.
- Stderr shows colored summaries regardless of verbose flag.
- Stdout is reserved for machine-readable output (JSON, CSV).
- Log file names include timestamp to prevent overwrites.
- Each log line has a timestamp and stage prefix.

---

## 4. Pipeline Stage Logging

For multi-stage operations, log entry and exit of each stage.

```go
verbose.Log("stage", "starting: %s", stageName)
// ... work ...
verbose.Log("stage", "completed: %s (%d items, %s)", stageName, count, elapsed)
```

### Standard Stages (Release Example)

| # | Stage | Key Log Fields |
|---|-------|----------------|
| 1 | Version Resolution | source, resolved version |
| 2 | Git Operations | branch, commit, tag |
| 3 | Asset Collection | file paths, sizes |
| 4 | Cross-Compilation | GOOS, GOARCH, output path |
| 5 | Compression | algorithm, ratio, bytes |
| 6 | Upload | URL, status code, retry count |
| 7 | Metadata Persistence | file path, JSON size |

---

## 5. Diagnostic Output (Doctor Pattern)

Health checks follow a consistent pass/fail format.

```
✓ Config file       OK — parsed 6 fields
✓ Database          OK — 42 repos, 12 releases
✓ Legacy dirs       OK — no migration needed
✗ Git binary        FAIL — git not found in PATH
✓ Deploy path       OK — E:\bin-run\gitmap\
```

### Rules

- One line per check — icon, name, status, detail.
- Use `✓` for pass, `✗` for fail, `⚠` for warning.
- Include the `--fix-path` pattern for auto-remediation.
- Exit code reflects overall health (0 = all pass, 1 = any fail).
- Constants for all check names and messages — no magic strings.

---

## 6. Progress Reporting

Use `[current/total]` counters on stderr for batch operations.

```
[1/42] Scanning api-gateway...
[2/42] Scanning frontend-app...
...
[42/42] Scanning legacy-tool...
✓ Scanned 42 repos in 3.4s
```

### Rules

- Counter format: `[current/total]` — always padded.
- Show item name being processed.
- Print summary line at completion.
- Support `--quiet` to suppress progress (keep summary).
- Never mix progress output with stdout data.

---

## 7. Error Logging

### Context Wrapping

Always wrap errors with the operation that failed.

```go
// ✅ Correct
return fmt.Errorf("writing release metadata v%s: %w", version, err)

// ❌ Wrong
return err
return fmt.Errorf("error: %w", err)
```

### User-Facing vs Internal

| Audience | Format | Example |
|----------|--------|---------|
| User (stderr) | Plain, actionable | `"Error: config.json not found. Run 'gitmap doctor' to diagnose."` |
| Log file | Detailed, technical | `"[config] ReadFile failed: /home/user/data/config.json: ENOENT"` |

---

## 8. What NOT to Log

- Secrets, tokens, or credentials — never.
- Redundant success messages for trivial operations.
- Stack traces in user-facing output (log file only).
- Raw HTTP response bodies (log truncated previews).

---

## References

- Verbose Logging Spec: `spec/04-generic-cli/16-verbose-logging.md`
- Progress Tracking: `spec/04-generic-cli/17-progress-tracking.md`
- Error Handling: `spec/05-coding-guidelines/04-error-handling.md`
- Code Quality: `spec/05-coding-guidelines/01-code-quality-improvement.md`
