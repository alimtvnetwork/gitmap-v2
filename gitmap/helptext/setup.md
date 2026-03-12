# gitmap setup

Interactive first-time configuration wizard that applies global Git settings.

## Alias

None

## Usage

    gitmap setup

## Flags

None.

## Prerequisites

- Git must be installed

## Examples

### Example 1: Run setup wizard

    gitmap setup

**Output:**

    GitMap Setup
    ─────────────
    Applying global Git configuration...
    ✓ core.autocrlf = true
    ✓ push.default = current
    ✓ pull.rebase = true
    Setup complete.

### Example 2: Re-run setup (safe to repeat)

    gitmap setup

**Output:**

    GitMap Setup
    ─────────────
    All settings already applied.
    Setup complete.
