# Issue: Repeated update-flow spec mismatch (PATH sync handoff)

## Summary

The update flow expectation was repeated several times but remained inconsistently documented between general and app specs. This caused implementation drift and confusion during end-to-end deploy/update validation.

## Observed Symptoms

- `run.ps1` deploy succeeded, but active PATH sync stayed locked for many retries.
- Spec language did not clearly enforce a same-location handoff copy and explicit update mode sequencing.
- Version/changelog confirmation responsibilities were not explicit enough.

## Root Cause

1. The spec described safe-update behavior at a high level but left key sequencing details implicit.
2. General and app-specific specs diverged in terminology (`temp script` vs `handoff copy` lifecycle).
3. Post-update validation requirements (before/after version comparison and changelog source) were not strict acceptance criteria.

## Final Solution

Document and enforce one canonical two-phase process:

### Phase 1: Handoff and lock release

1. `gitmap update` creates a handoff copy in the active binary directory (`gitmap.exe.old` or `gitmap-update-<pid>.exe`).
2. Launch handoff copy with `update --from-copy`.
3. Parent exits immediately.

### Phase 2: Update pipeline and validation

1. Handoff copy resolves repo path.
2. Run `run.ps1 -Update` (pull, build, deploy).
3. Safe PATH sync uses retry first, then rename fallback.
4. Print executable-derived version comparison (`before` and `after`).
5. Run `gitmap changelog --latest` from the updated binary.
6. Run `gitmap update-cleanup` to remove handoff and `.old` artifacts.

## Acceptance Criteria

- Active PATH binary equals deployed binary version after update.
- If versions differ, update exits with clear failure output.
- Changelog output is executed via updated `gitmap` binary.
- Cleanup runs after successful update completion.

## Prevention

- Any update-flow change must update both:
  - `spec/02-general/02-powershell-build-deploy.md`
  - `spec/01-app/09-build-deploy.md`
- Keep one source-of-truth sequence and mirror it verbatim across specs.
