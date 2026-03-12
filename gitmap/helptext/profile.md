# gitmap profile

Manage database profiles (separate repo databases for different contexts).

## Alias

pf

## Usage

    gitmap profile <create|list|switch|delete|show> [name]

## Flags

None.

## Prerequisites

- None

## Examples

### Example 1: Create and switch profile

    gitmap profile create work
    gitmap profile switch work

**Output:**

    ✓ Profile 'work' created
    ✓ Switched to profile 'work'

### Example 2: List profiles

    gitmap pf list

**Output:**

    default  42 repos
    work     18 repos  ✓ active
    3 profiles

### Example 3: Show current profile

    gitmap profile show

**Output:**

    Active profile: work
    Repos: 18
    Created: 2025-03-01

## See Also

- [diff-profiles](diff-profiles.md) — Compare repos across profiles
- [export](export.md) — Export current profile data
- [import](import.md) — Import data into a profile
- [db-reset](db-reset.md) — Reset the current profile database