# gitmap seo-write

Auto-generate and commit SEO-optimized messages to a repository on a schedule.

## Alias

sw

## Usage

    gitmap seo-write [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --csv \<path\> | — | CSV file with SEO data |
| --url \<url\> | — | Target website URL |
| --service \<name\> | — | Service name |
| --area \<name\> | — | Geographic area |
| --company \<name\> | — | Company name |
| --phone \<number\> | — | Phone number |
| --email \<addr\> | — | Contact email |
| --address \<addr\> | — | Physical address |
| --max-commits \<N\> | 10 | Maximum commits per run |
| --interval \<secs\> | 60 | Seconds between commits |
| --files \<list\> | — | Files to modify |
| --rotate | false | Rotate through templates |
| --dry-run | false | Preview without committing |
| --template \<path\> | — | Custom template file |
| --create-template | false | Generate a starter template |
| --author-name \<n\> | — | Commit author name |
| --author-email \<e\> | — | Commit author email |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Run SEO writes from CSV

    gitmap seo-write --csv data.csv --max-commits 5

**Output:**

    Loading SEO data from data.csv...
    [1/5] "Best plumber in Seattle"... done
    ✓ 5 commits created

### Example 2: Dry-run preview

    gitmap sw --csv data.csv --dry-run

**Output:**

    [DRY RUN] "Best plumber in Seattle"
    [DRY RUN] "Emergency plumbing 24/7"
    No changes made.