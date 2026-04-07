# Post-Mortem #16: CI Passthrough Gate Pattern

## Problem

The CI pipeline used **job-level `if` conditionals** to skip downstream
jobs when the SHA-based deduplication cache detected a previously
validated commit. While functionally correct, this caused all skipped
jobs to appear as **grey "Skipped"** in the GitHub Actions UI instead
of green ✅ Success.

### Impact

1. **Visual misrepresentation** — the CI summary looked like a partial
   failure, confusing contributors and reviewers.
2. **Required status checks blocked** — GitHub treats skipped jobs as
   neither success nor failure. If any skipped job was configured as a
   required status check, pull requests could not be merged.
3. **False alarms** — repository badges and notifications reported
   ambiguous status for perfectly valid commits.

### Root Cause

```yaml
# BEFORE — job-level skip (broken)
lint:
  needs: sha-check
  if: needs.sha-check.outputs.already-built != 'true'
  steps:
    - uses: actions/checkout@v4
    - name: Run lint
      run: golangci-lint run
```

When `already-built == 'true'`, the entire `lint` job was skipped.
GitHub rendered this as a grey icon with no success/failure verdict.

---

## Solution: Passthrough Gate Pattern

Replace job-level `if` with **step-level conditionals**. Every job
always runs, but individual steps are guarded:

```yaml
# AFTER — step-level guard (correct)
lint:
  needs: sha-check
  steps:
    - name: Already validated
      if: needs.sha-check.outputs.already-built == 'true'
      run: echo "✅ SHA ${{ github.sha }} already passed lint"

    - uses: actions/checkout@v4
      if: needs.sha-check.outputs.already-built != 'true'

    - name: Run lint
      if: needs.sha-check.outputs.already-built != 'true'
      run: golangci-lint run
```

When a cached SHA is detected:
- The "Already validated" step prints a confirmation message.
- All other steps are skipped (step-level, not job-level).
- The job itself completes with ✅ Success.

---

## Files Changed

| File | Change |
|------|--------|
| `.github/workflows/ci.yml` | Removed job-level `if`; added step-level conditionals to lint, vulncheck, test, test-summary |
| `spec/05-coding-guidelines/29-ci-sha-deduplication.md` | Documented the passthrough gate pattern and updated acceptance criteria |
| `.lovable/memory/tech/ci-pipeline-architecture.md` | Updated SHA deduplication section |

---

## Prevention Rules

1. **Never use job-level `if` for deduplication** — always use
   step-level conditionals so jobs report ✅ Success.
2. **Test with required status checks** — verify that cached SHA
   runs still satisfy branch protection rules.
3. **Check the GitHub UI** — after implementing cache-based skipping,
   manually confirm that all jobs show green checkmarks (not grey).

---

## Behavior Matrix

| Scenario | Cache | Jobs | UI Status |
|----------|-------|------|-----------|
| First push (new SHA) | Miss | Full execution | ✅ Green |
| Second push (same SHA) | Hit | "Already validated" echo | ✅ Green |
| Failed pipeline | Never cached | Re-run executes fully | ❌ Red |
| Job-level skip (old pattern) | Hit | Entire job skipped | ⚪ Grey |
| Concurrency cancellation | Hit/Miss | `mark-success` cancelled | ⚪ Grey (safe) |

---

## Known Behavior: Concurrency Cancellation

When `concurrency: cancel-in-progress: true` is active, a newer push to
the same branch cancels the in-progress run. If cancellation occurs after
all validation jobs pass but before `mark-success` completes, the SHA is
**not cached** — the `mark-success` job appears as cancelled (grey).

### Why This Is Safe

- All validation jobs (lint, vulncheck, test, test-summary) already
  completed successfully — the code **was** validated.
- The only consequence is that the SHA marker is not written to the cache.
- If the same SHA is encountered again, the pipeline runs a full
  validation — an unnecessary but correct re-run (identical to the
  "cache eviction" edge case).
- The newer commit that triggered the cancellation gets its own full
  pipeline run and will cache its own SHA on success.

### When It Happens

- Force-pushing rapid successive commits to the same branch.
- Merging multiple PRs in quick succession to `main`.
- Any scenario where two pushes target the same `concurrency` group
  before the first run's `mark-success` job finishes.

### Mitigation

No action required. The pattern is self-healing: the worst case is a
redundant but correct re-validation. Attempting to protect `mark-success`
from cancellation (e.g., via a separate concurrency group) would defeat
the purpose of cancelling superseded runs.

---

## Acceptance Criteria

1. All jobs show ✅ green in the CI UI for cached SHA runs.
2. No job appears as grey "Skipped" when SHA deduplication triggers.
3. Required status checks pass for deduplicated commits.
4. Each skipped job prints "Already validated" in its log output.
5. New commits (different SHA) still execute the full pipeline.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
