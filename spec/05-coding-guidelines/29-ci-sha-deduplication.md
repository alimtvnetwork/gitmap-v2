# CI SHA-Based Build Deduplication

## Overview

When the same commit SHA is pushed multiple times (e.g., force-push to
the same branch, re-run of a merge, or tag pointing at an existing
commit), the CI pipeline should detect that the SHA has already been
successfully validated and skip redundant work.

---

## Mechanism

Use the CI provider's **cache** system with the commit SHA as the key.
A successful pipeline writes a marker file to the cache. Subsequent
runs check for the marker before executing any work.

### How It Works

1. A **gate job** runs first and probes the cache for a key derived
   from the commit SHA (e.g., `ci-passed-<SHA>`).
2. If the cache **hits**, the gate job sets an output flag
   (`already-built = true`) and all downstream jobs are skipped via
   a condition (`if: needs.gate.outputs.already-built != 'true'`).
3. If the cache **misses**, the pipeline proceeds normally.
4. A **finalize job** at the end of the pipeline writes a trivial
   marker file into the cache under the same key, recording success.

### Cache Key Design

```
ci-passed-<full-commit-sha>
```

- Full 40-character SHA — no short hashes.
- No branch or ref component — the SHA is globally unique.
- The key must be identical between the check and the write steps.

---

## Implementation (GitHub Actions)

### Gate Job

```yaml
jobs:
  sha-check:
    name: Check if already built
    runs-on: ubuntu-latest
    outputs:
      already-built: ${{ steps.cache-check.outputs.cache-hit }}
    steps:
      - name: Check SHA cache
        id: cache-check
        uses: actions/cache@v4
        with:
          path: /tmp/ci-passed
          key: ci-passed-${{ github.sha }}
          lookup-only: true
```

`lookup-only: true` avoids downloading the cache content — only the
existence check matters.

### Conditional Downstream Jobs

Every subsequent job adds:

```yaml
  lint:
    needs: sha-check
    if: needs.sha-check.outputs.already-built != 'true'
```

Jobs that depend on other jobs (e.g., `test` depends on `lint`) must
also include the gate in their `needs` array or inherit the skip
transitively.

### Finalize Job (Mark Success)

```yaml
  mark-success:
    name: Mark SHA as built
    runs-on: ubuntu-latest
    needs: [test-summary]
    if: success()
    steps:
      - name: Write marker
        run: mkdir -p /tmp/ci-passed && echo "${{ github.sha }}" > /tmp/ci-passed/sha.txt

      - name: Save to cache
        uses: actions/cache/save@v4
        with:
          path: /tmp/ci-passed
          key: ci-passed-${{ github.sha }}
```

`if: success()` ensures the marker is only written when **all**
upstream jobs pass. A failed pipeline never caches.

---

## Behavior Matrix

| Scenario | Cache State | Result |
|----------|-------------|--------|
| First push of commit ABC | Miss | Full pipeline runs |
| Second push of same ABC | Hit | All jobs skipped |
| New commit DEF | Miss | Full pipeline runs |
| Failed pipeline for GHI | Never written | Re-run executes fully |
| Force-push (new SHA) | Miss | Full pipeline runs |

---

## Edge Cases

### Re-running a failed build

Because the marker is only written on success, a failed pipeline
leaves no cache entry. Re-running the same SHA executes the full
pipeline again.

### Cache eviction

GitHub Actions evicts caches after 7 days of no access. If a SHA
is re-pushed after eviction, the pipeline runs again. This is safe
because the worst case is an unnecessary but correct re-run.

### Pull requests from forks

Forked PRs use a different cache scope. The deduplication only
applies within the same repository's cache namespace.

---

## Constraints

- The marker file content is irrelevant — only the cache key matters.
- Never use short SHAs — collisions would cause false cache hits.
- The finalize job must depend on **all** validation jobs to avoid
  marking a partially-successful run as complete.
- Do not use this pattern for release pipelines — artifact production
  must always run to ensure reproducibility.

---

## Acceptance Criteria

1. Pushing the same commit SHA twice results in the second run
   skipping all lint, test, and summary jobs.
2. A failed pipeline does not cache — re-running the same SHA
   executes the full pipeline.
3. A new commit (different SHA) always runs the full pipeline.
4. The gate job completes in under 10 seconds.
5. Skipped jobs show as "skipped" in the CI UI, not "failed".

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
