# gitmap zip-group

Manage named collections of files and folders that are automatically
compressed into ZIP archives during a release.

## Alias

z

## Usage

    gitmap zip-group <subcommand> [arguments]

## Subcommands

| Subcommand | Description |
|------------|-------------|
| create     | Create a named zip group |
| add        | Add files or folders to a group |
| remove     | Remove an item from a group |
| list       | List all zip groups |
| show       | Show items in a group |
| delete     | Delete a zip group |
| rename     | Set a custom archive name for a group |

## Flags

| Flag | Description |
|------|-------------|
| --archive \<name\> | Custom output filename (used with create/rename) |

## Prerequisites

- Must be inside a Git repository with release workflow configured (see release.md)

## Examples

### Example 1: Create a group and add items

    gitmap z create docs-bundle
    gitmap z add docs-bundle ./README.md ./CHANGELOG.md ./docs/

**Output:**

    ✓ Zip group 'docs-bundle' created

    Adding items to 'docs-bundle'...
    ✓ Added ./README.md
    ✓ Added ./CHANGELOG.md
    ✓ Added ./docs/ (directory)
    3 items in group 'docs-bundle'

### Example 2: Create with custom archive name

    gitmap z create extras --archive extra-files.zip
    gitmap z add extras ./config/ ./scripts/deploy.sh

**Output:**

    ✓ Zip group 'extras' created (archive: extra-files.zip)

    ✓ Added ./config/ (directory)
    ✓ Added ./scripts/deploy.sh
    2 items in group 'extras'

### Example 3: List all zip groups

    gitmap z list

**Output:**

    GROUP           ITEMS   ARCHIVE NAME
    docs-bundle     3       docs-bundle.zip
    extras          2       extra-files.zip
    2 zip groups defined

### Example 4: Show items in a group

    gitmap z show docs-bundle

**Output:**

    Zip group: docs-bundle
    Archive:   docs-bundle.zip
    Items (3):
      ./README.md
      ./CHANGELOG.md
      ./docs/ (directory)

### Example 5: Use during release

    gitmap release v3.0.0 --zip-group docs-bundle

**Output:**

    Creating tag v3.0.0... done
    ✓ Compressed docs-bundle → docs-bundle_v3.0.0.zip (3 items)
    Uploading to GitHub... done
    ✓ Released v3.0.0

### Example 6: Ad-hoc zip with bundle

    gitmap release v3.0.0 -Z ./dist/report.pdf -Z ./dist/manual.pdf --bundle reports.zip

**Output:**

    Creating tag v3.0.0... done
    ✓ Compressed 2 items → reports.zip
      ./dist/report.pdf
      ./dist/manual.pdf
    Uploading to GitHub... done
    ✓ Released v3.0.0

## See Also

- [release](release.md) — Create a release with zip group assets
- [group](group.md) — Manage repository groups
