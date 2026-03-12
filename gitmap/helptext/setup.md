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

    Applying global Git configuration...
    ✓ 3 settings applied
    Setup complete.

## See Also

- [scan](scan.md) — Scan directories after setup
- [doctor](doctor.md) — Diagnose installation issues
- [update](update.md) — Update gitmap to the latest version

### Example 2: Re-run setup (safe to repeat)

    gitmap setup

**Output:**

    All settings already applied.
    Setup complete.