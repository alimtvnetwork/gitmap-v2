# CI Pipeline Architecture

The CI pipeline (GitHub Actions) uses a parallel matrix strategy ('fail-fast: false') to execute four distinct test suites: unit, store, integration, and tui. Test output and coverage profiles ('-covermode=atomic') are collected as artifacts and consolidated by a final 'test-summary' job. This summary job aggregates failures into a single report, calculates project-wide coverage using 'go tool cover', and generates a per-package breakdown. To ensure visibility, the test stage uses 'set +e' and 'grep' to filter for specific Go failure patterns (e.g., '--- FAIL', 'build failed', 'undefined') before exiting with the original code.

## SHA-Based Build Deduplication (Passthrough Gate Pattern)

A 'sha-check' gate job runs before all other jobs. It probes the GitHub Actions cache for key 'ci-passed-<SHA>' using 'lookup-only: true'. Downstream jobs always run (no job-level `if` skipping) but use **step-level conditionals**: when the SHA is already cached, each job executes only an "Already validated" echo step and exits with ✅ Success. This ensures the GitHub UI always shows green checkmarks — never grey "skipped" icons that look like failures and block required status checks. When the cache misses, steps guarded by `if: needs.sha-check.outputs.already-built != 'true'` execute normally. The cache write is **inlined as the final step of `test-summary`** (not a separate `mark-success` job) to prevent `cancel-in-progress` from cancelling the cache save while all validation jobs have already passed. Failed pipelines never cache, so re-runs of the same SHA execute fully. Documented in spec/05-coding-guidelines/29-ci-sha-deduplication.md.

## Concurrency Control

All workflows use 'concurrency: group: ci-${{ github.ref }}, cancel-in-progress: true' to cancel superseded runs on the same branch while allowing independent runs on different branches.

## Lessons Learned

1. **Never use `cd` in CI scripts** — use `working-directory` in the workflow step definition. The v2.54.0 release pipeline failed with `cd: dist: No such file or directory` because the compress step ran in `gitmap-updater/` instead of `gitmap/`. Fixed by setting explicit `working-directory: gitmap/dist`. See `spec/02-app-issues/13-release-pipeline-dist-directory.md`.
2. **Pin Go tool versions** — `go install tool@latest` is non-reproducible. All tools (e.g., `golangci-lint@v1.64.8`) must use exact version tags. Documented in `setup.sh` and `spec/05-coding-guidelines/17-cicd-patterns.md`.
3. **Validate build output directories** before operating on them: `test -d "$DIR" || exit 1`.
4. **Never use job-level `if` for SHA deduplication** — GitHub treats skipped jobs as neither success nor failure, blocking required status checks. Use the passthrough gate pattern with step-level conditionals instead.
5. **Inline cache writes into the last validation job** — a separate `mark-success` job can be cancelled by `cancel-in-progress` after all validation passes, leaving the SHA uncached. Inlining the cache save as the final step of `test-summary` prevents this.
