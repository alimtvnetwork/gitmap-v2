# Compliance Audit Summary

Last updated: 2026-03-06

## Rules Checked

| # | Rule | Source |
|---|------|--------|
| 1 | No negation in `if` conditions (`!`, `!=`, `== false`) | 01-overview §Code Style |
| 2 | Functions: 8–15 lines | 01-overview §Code Style |
| 3 | Files: 100–200 lines max | 01-overview §Code Style |
| 4 | One responsibility per package | 01-overview §Code Style |
| 5 | Blank line before `return` (unless sole line in `if`) | 01-overview §Code Style |
| 6 | No magic strings — all literals in `constants` | 01-overview §Code Style |
| 7 | No `switch` statements — use `if`/`else if` chains | 02-general/06 §Conditionals |

## Package: `cmd`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `root.go` | ~60 | ✅ Pass | |
| `rootflags.go` | ~50 | ✅ Pass | |
| `rootusage.go` | ~45 | ✅ Pass | |
| `scan.go` | ~113 | ✅ Pass | Split from 257 lines |
| `scanoutput.go` | ~155 | ✅ Pass | Extracted from scan.go |
| `clone.go` | ~140 | ✅ Pass | |
| `pull.go` | ~100 | ✅ Pass | Magic strings extracted |
| `rescan.go` | ~110 | ✅ Pass | Magic strings extracted |
| `status.go` | ~187 | ✅ Pass | |
| `statusformat.go` | ~135 | ✅ Pass | |
| `exec.go` | ~120 | ✅ Pass | |
| `list.go` | ~80 | ✅ Pass | |
| `setup.go` | ~60 | ✅ Pass | |
| `update.go` | ~90 | ✅ Pass | |
| `updatescript.go` | ~120 | ✅ Pass | Magic strings extracted |
| `updatecleanup.go` | ~100 | ✅ Pass | Magic strings extracted |
| `release.go` | ~130 | ✅ Pass | |
| `releasebranch.go` | ~60 | ✅ Pass | |
| `releasepending.go` | ~40 | ✅ Pass | |
| `changelog.go` | ~80 | ✅ Pass | Magic strings extracted |
| `latestbranch.go` | ~80 | ✅ Pass | |
| `latestbranchresolve.go` | ~90 | ✅ Pass | |
| `latestbranchoutput.go` | ~100 | ✅ Pass | Magic strings extracted |
| `desktopsync.go` | ~100 | ✅ Pass | |
| `doctor.go` | ~60 | ✅ Pass | |
| `doctorchecks.go` | ~165 | ✅ Pass | Split; version logic extracted |
| `doctorversion.go` | ~120 | ✅ Pass | Extracted from doctorchecks.go |
| `doctorfixpath.go` | ~170 | ✅ Pass | Split; sync logic extracted |
| `doctorsync.go` | ~110 | ✅ Pass | Extracted from doctorfixpath.go |
| `group.go` | ~30 | ✅ Pass | |
| `groupcreate.go` | ~60 | ✅ Pass | |
| `groupdelete.go` | ~60 | ✅ Pass | |
| `groupadd.go` | ~60 | ✅ Pass | |
| `groupremove.go` | ~60 | ✅ Pass | |
| `grouplist.go` | ~50 | ✅ Pass | |
| `groupshow.go` | ~60 | ✅ Pass | |
| `flags_test.go` | ~40 | ✅ Pass | |

## Package: `constants`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `constants.go` | ~30 | ✅ Pass | |
| `constants_cli.go` | ~60 | ✅ Pass | |
| `constants_doctor.go` | ~50 | ✅ Pass | OS/binary constants added |
| `constants_git.go` | ~45 | ✅ Pass | |
| `constants_messages.go` | ~130 | ✅ Pass | OS command constants added |
| `constants_release.go` | ~35 | ✅ Pass | |
| `constants_store.go` | ~20 | ✅ Pass | |
| `constants_terminal.go` | ~180 | ✅ Pass | Format strings extracted |
| `constants_update.go` | ~40 | ✅ Pass | PS/shell constants added |

## Package: `release`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `workflow.go` | ~416 | ⚠️ Pending trim | Needs deletion of lines 165–416 (moved to split files) |
| `workflowfinalize.go` | ~130 | ✅ Pass | Extracted from workflow.go |
| `workflowbranch.go` | ~165 | ✅ Pass | Extracted from workflow.go |
| `gitops.go` | ~100 | ✅ Pass | Rewritten; query functions extracted |
| `gitopsquery.go` | ~135 | ✅ Pass | Extracted from gitops.go |
| `changelog.go` | ~120 | ⚠️ Negation | `== false` on lines 56, 73 |
| `github.go` | ~60 | ⚠️ Negation | `== false` on line 30 |
| `metadata.go` | ~130 | ⚠️ Negation | `== false` on line 78 |
| `metadata_test.go` | ~40 | ✅ Pass | |
| `semver.go` | ~120 | ⚠️ Switch | `switch` on line 90; needs if/else chain |
| `semver_test.go` | ~80 | ✅ Pass | |

## Package: `formatter`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `terminal.go` | ~223 | ⚠️ Pending trim | Needs deletion of lines 125–223 (moved to terminaltree.go) |
| `terminaltree.go` | ~110 | ✅ Pass | Extracted from terminal.go |
| `csv.go` | ~60 | ✅ Pass | |
| `json.go` | ~30 | ✅ Pass | |
| `text.go` | ~30 | ✅ Pass | |
| `structure.go` | ~100 | ✅ Pass | |
| `clonescript.go` | ~40 | ✅ Pass | |
| `directclone.go` | ~70 | ✅ Pass | |
| `desktopscript.go` | ~50 | ✅ Pass | |
| `template.go` | ~30 | ✅ Pass | |
| `formatter_test.go` | ~60 | ✅ Pass | |

## Package: `cloner`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `cloner.go` | ~90 | ✅ Pass | |
| `safe_pull.go` | ~213 | ⚠️ Pending trim | Needs deletion of lines 112–213 (moved to pulldiag.go) |
| `pulldiag.go` | ~130 | ✅ Pass | Extracted from safe_pull.go |

## Package: `setup`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `setup.go` | ~206 | ⚠️ Pending trim | Needs deletion of lines 133–206 (moved to setupapply.go) |
| `setupapply.go` | ~100 | ✅ Pass | Extracted from setup.go |

## Package: `config`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `config.go` | ~50 | ⚠️ Negation | `os.IsNotExist` on line 17 |
| `config_test.go` | ~30 | ✅ Pass | |

## Package: `scanner`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `scanner.go` | ~80 | ✅ Pass | |

## Package: `mapper`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `mapper.go` | ~110 | ✅ Pass | |
| `mapper_test.go` | ~50 | ✅ Pass | |

## Package: `model`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `record.go` | ~67 | ✅ Pass | |
| `group.go` | ~20 | ✅ Pass | |

## Package: `store`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `store.go` | ~80 | ✅ Pass | |
| `repo.go` | ~90 | ✅ Pass | |
| `group.go` | ~130 | ✅ Pass | |

## Package: `desktop`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `desktop.go` | ~60 | ✅ Pass | |

## Package: `gitutil`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `gitutil.go` | ~120 | ✅ Pass | |
| `latestbranch.go` | ~110 | ✅ Pass | |
| `latestbranchresolve.go` | ~90 | ✅ Pass | |
| `dateformat.go` | ~40 | ✅ Pass | |

## Package: `verbose`

| File | Lines | Status | Notes |
|------|-------|--------|-------|
| `verbose.go` | ~60 | ⚠️ Negation | `!l.enabled` on line 53 |

## Pending Work Summary

| Category | Count | Details |
|----------|-------|---------|
| File trims (over 200 lines) | 4 | `workflow.go`, `terminal.go`, `safe_pull.go`, `setup.go` |
| Negation fixes (`== false`, `!`, `os.IsNotExist`) | 5 | `changelog.go`, `github.go`, `metadata.go`, `config.go`, `verbose.go` |
| Switch → if/else | 1 | `semver.go` |
| Missing constants | ~10 | Git args, bump levels, release tag prefix |
