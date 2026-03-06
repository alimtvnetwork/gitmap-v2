# Issue: Repeated update-flow spec mismatch (PATH sync handoff)

## Summary

The update flow expectation was repeated several times but remained inconsistently documented between general and app specs. This caused implementation drift and confusion during end-to-end deploy/update validation.

## Observed Symptoms

- `run.ps1` deploy succeeded, but active PATH sync stayed locked for many retries.
- Spec language did not clearly enforce a same-location handoff copy and explicit update mode sequencing.
- Version/changelog confirmation responsibilities were not explicit enough.
- Implementation oscillated between `cmd.Run()` (sync) and `cmd.Start()` (async) because specs did not explicitly forbid synchronous wait.

## Root Cause

1. The spec described safe-update behavior at a high level but left key sequencing details implicit.
2. General and app-specific specs diverged in terminology (`temp script` vs `handoff copy` lifecycle, `update --from-copy` vs `update-runner`).
3. Post-update validation requirements (before/after version comparison and changelog source) were not strict acceptance criteria.
4. **Specs did not explicitly state that `cmd.Run()` is forbidden in `runUpdate()`** — this allowed repeated reintroduction of the synchronous-wait bug.
5. **Specs did not mandate rename-first for PATH sync** — only mentioned rename as a "fallback", leading to copy-first implementations that always fail under lock.

## Final Solution

Document and enforce one canonical two-phase process:

### Phase 1: Handoff and lock release

1. `gitmap update` creates a handoff copy in the active binary directory (`gitmap-update-<pid>.exe`, fallback to `%TEMP%`).
2. Launch handoff copy with hidden `update-runner` command.
3. Parent exits immediately via `cmd.Start()` + `os.Exit(0)`. **Never `cmd.Run()`.**

### Phase 2: Update pipeline and validation

1. Handoff copy resolves repo path.
2. Run `run.ps1 -Update` (full pipeline: pull, build, deploy).
3. PATH sync uses **rename-first** in update mode, then copy-retry as fallback.
4. Print executable-derived version comparison (`before` and `after`).
5. Run `gitmap changelog --latest` from the updated binary.
6. Run `gitmap update-cleanup` to remove handoff and `.old` artifacts.

## Acceptance Criteria

- Active PATH binary equals deployed binary version after update.
- If versions differ, update exits with clear failure output.
- Changelog output is executed via updated `gitmap` binary.
- Cleanup runs after successful update completion.
- Zero lock-retry loops during normal update (rename-first should succeed immediately).

## Prevention

- Any update-flow change must update ALL of:
  - `spec/02-general/02-powershell-build-deploy.md`
  - `spec/02-general/03-self-update-mechanism.md`
  - `spec/01-app/09-build-deploy.md`
  - `spec/01-app/02-cli-interface.md`
  - `.lovable/memory/issues/` (if new failure mode)
- Keep one source-of-truth sequence and mirror it verbatim across specs.
- **Explicit prohibitions** must be documented (e.g., "never use `cmd.Run()`", "never add `Read-Host`").
