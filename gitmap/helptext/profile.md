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

    Profile     Repos   Active
    default     42      
    work        18      ✓
    personal    7       
    3 profiles

### Example 3: Show current profile

    gitmap profile show

**Output:**

    Active profile: work
    Repos: 18
    Created: 2025-03-01
