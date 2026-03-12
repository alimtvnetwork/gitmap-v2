# gitmap update

Self-update gitmap from the source repository. Pulls latest, rebuilds, and deploys.

## Alias

None

## Usage

    gitmap update

## Flags

None.

## Prerequisites

- Git must be installed
- Source repository must be accessible

## Examples

### Example 1: Update to latest

    gitmap update

**Output:**

    v2.5.0 → v2.8.0
    Building and deploying... done
    ✓ Updated to v2.8.0

### Example 2: Already up to date

    gitmap update

**Output:**

    Current: v2.8.0
    ✓ Already up to date.