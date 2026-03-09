# Commit Author Override & Amend Command

## Overview

Two related features for controlling Git commit authorship:

1. **SEO-write author flags** — set custom author name/email on each commit during `seo-write`.
2. **`gitmap amend` command** — rewrite author name/email on existing commits (all or from a specific commit onwards).

---

## Feature 1: SEO-Write Author Flags

### New Flags

| Flag              | Description                              | Default              |
|-------------------|------------------------------------------|----------------------|
| `--author-name`   | Git author name for commits              | (current git config) |
| `--author-email`  | Git author email for commits             | (current git config) |

### Behavior

- When provided, each `git commit` in the seo-write loop uses:
  ```
  git commit -m "message" --author="Name <email>"
  ```
- If only `--author-name` is provided without `--author-email`, use the name with the current git config email.
- If only `--author-email` is provided without `--author-name`, use the current git config name with the provided email.
- Dry-run mode (`--dry-run`) should display the author that would be used.

### Examples

```bash
# SEO-write with custom author
gitmap sw --url example.com --service Plumbing --area London \
  --author-name "John Smith" --author-email "john@example.com"

# Only override name (email stays from git config)
gitmap sw --url example.com --service SEO --area Bristol \
  --author-name "Marketing Bot"
```

---

## Feature 2: `gitmap amend` Command

### Synopsis

```
gitmap amend [commit-hash] --name <name> --email <email>
```

Alias: `gitmap am`

### Behavior

Rewrites the author name and/or email on existing commits using `git filter-branch` or `git rebase --exec`.

#### Mode 1: All Commits

```bash
gitmap amend --name "New Name" --email "new@email.com"
```

Rewrites **every commit** in the current branch to use the new author identity.

#### Mode 2: From a Specific Commit Onwards

```bash
gitmap amend abc123 --name "New Name" --email "new@email.com"
```

Rewrites all commits **from `abc123` (inclusive) to HEAD** to use the new author identity. Commits before `abc123` are left untouched.

#### Mode 3: Single Commit (HEAD only)

```bash
gitmap amend HEAD --name "New Name" --email "new@email.com"
```

Amends only the most recent commit.

### Flags

| Flag              | Description                              | Required |
|-------------------|------------------------------------------|----------|
| `--name <name>`   | New author name                          | Yes (at least one of name/email) |
| `--email <email>` | New author email                         | Yes (at least one of name/email) |
| `--dry-run`       | Preview which commits would be amended   | No       |
| `--force-push`    | Auto-run `git push --force-with-lease` after amend | No |

### Implementation Approach

The amend command generates and executes a `git filter-branch` command:

```bash
# All commits
git filter-branch -f --env-filter '
  export GIT_AUTHOR_NAME="New Name"
  export GIT_AUTHOR_EMAIL="new@email.com"
  export GIT_COMMITTER_NAME="New Name"
  export GIT_COMMITTER_EMAIL="new@email.com"
' -- HEAD

# From specific commit onwards
git filter-branch -f --env-filter '
  export GIT_AUTHOR_NAME="New Name"
  export GIT_AUTHOR_EMAIL="new@email.com"
  export GIT_COMMITTER_NAME="New Name"
  export GIT_COMMITTER_EMAIL="new@email.com"
' <commit-hash>^..HEAD
```

### Safety

- **Warning prompt**: Before executing, print a warning that this rewrites history and requires force-push. Proceed automatically (no interactive prompt — follows project convention).
- **Backup ref**: Git automatically creates `refs/original/` backup refs during filter-branch.
- **Dry-run**: List all commits that would be affected with their current author and the new author.

### Terminal Output

```
amend: rewriting 15 commits from abc1234..HEAD
  author: "Old Name <old@email.com>" -> "New Name <new@email.com>"

  [1/15] abc1234 - Fix login page
  [2/15] def5678 - Add dashboard
  ...
  [15/15] 9z8y7x6 - Update README

Done: 15 commits amended
Warning: Run 'git push --force-with-lease' to update the remote
```

### Examples

```bash
# Amend all commits in current branch
gitmap amend --name "John Smith" --email "john@company.com"

# Amend from a specific commit onwards
gitmap amend a1b2c3d --name "John Smith" --email "john@company.com"

# Amend only HEAD
gitmap amend HEAD --name "John Smith" --email "john@company.com"

# Preview what would change
gitmap amend --name "John Smith" --email "john@company.com" --dry-run

# Amend and auto force-push
gitmap amend a1b2c3d --name "John Smith" --email "john@company.com" --force-push
```

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_amend.go` | Command/flag/message constants |
| `cmd/amend.go` | Flag parsing, orchestration |
| `cmd/amendexec.go` | Git filter-branch execution and output |

SEO-write changes modify existing files:
- `constants/constants_seo.go` — add `FlagSEOAuthorName`, `FlagSEOAuthorEmail`
- `cmd/seowrite.go` — add fields to `seoWriteFlags`
- `cmd/seowriteloop.go` — pass author to `gitCommit`

---

## CLI Interface Updates

### Command Table Addition

| Command | Alias | Description |
|---------|-------|-------------|
| `amend [hash]` | `am` | Rewrite author name/email on commits |

### Dispatch

Add to `dispatchMisc` in `root.go`.

---

## Acceptance Criteria

- [ ] `gitmap sw --author-name "Bot" --author-email "bot@co.com"` sets author on each commit
- [ ] `gitmap amend --name "X" --email "x@y.com"` rewrites all commits
- [ ] `gitmap amend abc123 --name "X" --email "x@y.com"` rewrites from abc123 to HEAD
- [ ] `gitmap amend HEAD --name "X" --email "x@y.com"` amends only the latest commit
- [ ] `--dry-run` shows affected commits without modifying anything
- [ ] `--force-push` runs `git push --force-with-lease` after amend
- [ ] At least one of `--name` or `--email` is required
- [ ] Terminal output shows progress per commit
