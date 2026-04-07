# CI/CD Pipeline Architecture

## Overview

The project uses three GitHub Actions workflows to enforce code quality,
produce release artifacts, and scan for vulnerabilities. All workflows
implement concurrency controls to cancel superseded runs.

---

## Workflows

### 1. CI (`ci.yml`)

**Triggers:** Push to `main`, pull requests targeting `main`.

| Stage          | Purpose                                           |
|----------------|---------------------------------------------------|
| Lint           | `go vet` + `golangci-lint` (28 linters, 5m timeout) |
| Vulnerability  | `govulncheck` — fails only on third-party issues  |
| Test           | Parallel matrix: unit, store, integration, tui    |
| Test Summary   | Aggregates failures, generates coverage breakdown |

**No binary builds on main.** Artifact production is delegated
exclusively to the release pipeline.

### 2. Release (`release.yml`)

**Triggers:** Push to `release/**` branches or `v*` tags.

Produces 6 cross-compiled binaries (windows/linux/darwin ×
amd64/arm64) for both `gitmap` and `gitmap-updater`. Generates
versioned artifacts, SHA256 checksums, and changelog excerpts.

### 3. Vulnerability Scan (`vulncheck.yml`)

**Triggers:** Weekly schedule (Mondays 09:00 UTC), manual dispatch.

Runs `govulncheck` independently of the CI pipeline for proactive
dependency monitoring.

---

## Concurrency Control

All three workflows use GitHub Actions `concurrency` groups with
`cancel-in-progress: true`. When a new commit is pushed while a
previous run is still executing, the older run is automatically
cancelled and only the latest commit is built.

### Group Keys

| Workflow   | Concurrency Group             |
|------------|-------------------------------|
| CI         | `ci-${{ github.ref }}`        |
| Release    | `release-${{ github.ref }}`   |
| Vulncheck  | `vulncheck-${{ github.ref }}` |

The `github.ref` suffix ensures that:

- Different branches run independently (a push to `feature/a` does
  not cancel a run on `main`).
- Multiple pushes to the **same** branch cancel each other, keeping
  only the latest.
- Pull request runs are scoped to the PR ref, so updating a PR
  cancels its previous CI run without affecting `main`.

### Why cancel-in-progress?

| Problem                                  | Solution                          |
|------------------------------------------|-----------------------------------|
| Wasted CI minutes on outdated commits    | Auto-cancel superseded runs       |
| Queue buildup during rapid iteration     | Only latest commit matters        |
| Stale results reported on merged PRs     | Cancelled runs produce no output  |
| Release branch rapid-fire pushes         | Only final push produces artifacts|

---

## Test Architecture

Tests run in a parallel matrix (`fail-fast: false`) across four
suites. Each suite produces:

- `test-output.txt` — full verbose output for failure analysis.
- `coverage-<suite>.out` — atomic coverage profile.

The `test-summary` job downloads all artifacts, aggregates failures
into a single report, and merges coverage profiles for a per-package
breakdown using `go tool cover`.

---

## Build Strategy

| Context         | Binary Production | Rationale                        |
|-----------------|-------------------|----------------------------------|
| `main` branch   | None              | CI validates, doesn't produce    |
| Pull requests   | None              | Same as main — validation only   |
| `release/**`    | 6 targets         | Official artifacts for release   |
| `v*` tags       | 6 targets         | Tagged release artifacts         |

This separation follows the **"build once, promote often"** principle:
binaries are only compiled in the release pipeline where they are
immediately attached to a GitHub Release with checksums.

---

## File Layout

| File                           | Purpose                          |
|--------------------------------|----------------------------------|
| `.github/workflows/ci.yml`     | Lint, test, coverage on main/PRs |
| `.github/workflows/release.yml`| Cross-compile and publish        |
| `.github/workflows/vulncheck.yml`| Weekly vulnerability scan      |
| `.golangci.yml`                | Linter configuration (28 rules)  |

---

## Acceptance Criteria

1. Pushing two commits in quick succession to the same branch
   cancels the first CI run and only completes the second.
2. Pushing to `main` runs lint, vulncheck, and tests but does
   **not** produce binary artifacts.
3. Pushing to `release/**` or a `v*` tag produces 6 binaries
   and uploads them as GitHub Release assets.
4. The weekly vulncheck runs independently and does not block
   or interfere with CI or release workflows.
5. Pull request CI runs are scoped to the PR ref and do not
   cancel runs on `main` or other PRs.
6. Test failures in one suite do not prevent other suites from
   completing (`fail-fast: false`).