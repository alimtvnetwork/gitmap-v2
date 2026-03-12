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

    Checking for updates...
    Current: v2.5.0
    Latest:  v2.8.0
    Pulling latest... done
    Building... done
    Deploying to E:\bin-run... done
    ✓ Updated to v2.8.0

### Example 2: Already up to date

    gitmap update

**Output:**

    Checking for updates...
    Current: v2.8.0
    ✓ Already up to date.
