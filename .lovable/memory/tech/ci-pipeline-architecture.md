# CI Pipeline Architecture

The CI pipeline (GitHub Actions) uses a parallel matrix strategy ('fail-fast: false') to execute four distinct test suites: unit, store, integration, and tui. Test output and coverage profiles ('-covermode=atomic') are collected as artifacts and consolidated by a final 'test-summary' job. This summary job aggregates failures into a single report, calculates project-wide coverage using 'go tool cover', and generates a per-package breakdown. To ensure visibility, the test stage uses 'set +e' and 'grep' to filter for specific Go failure patterns (e.g., '--- FAIL', 'build failed', 'undefined') before exiting with the original code.

## SHA-Based Build Deduplication

A 'sha-check' gate job runs before all other jobs. It probes the GitHub Actions cache for key 'ci-passed-<SHA>' using 'lookup-only: true'. If the commit has already passed CI (cache hit), all downstream jobs (lint, vulncheck, test, test-summary) are skipped via 'if: needs.sha-check.outputs.already-built != true'. On full pipeline success, a 'mark-success' job writes a marker file to the cache using 'actions/cache/save@v4'. Failed pipelines never cache, so re-runs of the same SHA execute fully. Documented in spec/05-coding-guidelines/29-ci-sha-deduplication.md.

## Concurrency Control

All workflows use 'concurrency: group: ci-${{ github.ref }}, cancel-in-progress: true' to cancel superseded runs on the same branch while allowing independent runs on different branches.
