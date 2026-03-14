export interface ChangelogEntry {
  version: string;
  items: string[];
}

export const changelog: ChangelogEntry[] = [
  {
    version: "v2.12.0",
    items: [
      "Added global ⌘K command palette searching across commands, flags, and pages.",
    ],
  },
  {
    version: "v2.11.0",
    items: [
      "Added Changelog page with timeline view and expand/collapse controls.",
      "Added Flag Reference page with sortable, searchable table of all flags.",
      "Added Interactive Examples page with animated terminal demos.",
    ],
  },
  {
    version: "v2.28.0",
    items: [
      "Removed unused `detector` import from `cmd/scan.go` that caused build failure.",
      "Updated documentation site fonts: Ubuntu for headings, Poppins for body text, Ubuntu Mono for code blocks.",
    ],
  },
  {
    version: "v2.27.0",
    items: [
      "Added `gitmap cd` (`go`) command: jump to any tracked repo by slug or partial name.",
      "Added `gitmap watch` (`w`) command: live terminal dashboard monitoring repo status.",
      "Added `gitmap diff-profiles` (`dp`) command: compare two profiles side-by-side.",
      "Added clone progress bars with retry logic and Windows long-path warnings.",
      "Built documentation site with interactive terminal preview for the watch command.",
      "Added `gitmap/Makefile` as a thin wrapper around `run.sh` for standard `make` workflows.",
      "Added `run.sh` cross-platform build script for Linux and macOS.",
      "Added `gitmap gomod` (`gm`) command: rename Go module path across an entire repo.",
    ],
  },
  {
    version: "v2.26.0",
    items: [
      "Version bump following `gitmap profile` command addition.",
      "All profile subcommands (`create`, `list`, `switch`, `delete`, `show`) fully integrated.",
    ],
  },
  {
    version: "v2.25.0",
    items: [
      "Added `gitmap profile` (`pf`) command: manage multiple database profiles.",
      "Subcommands: `create`, `list`, `switch`, `delete`, `show`.",
      "Each profile has its own SQLite database file.",
    ],
  },
  {
    version: "v2.24.0",
    items: [
      "Added `gitmap import` (`im`) command: restore database from a backup file.",
      "Merge semantics: upserts repos/releases, INSERT OR IGNORE for history/bookmarks/groups.",
    ],
  },
  {
    version: "v2.23.0",
    items: [
      "Added `gitmap export` (`ex`) command: export the full database as portable JSON.",
      "Exports all tables: repos, groups, releases, command history, and bookmarks.",
    ],
  },
  {
    version: "v2.22.0",
    items: [
      "Added `gitmap bookmark` (`bk`) command: save and replay command+flag combinations.",
      "Subcommands: `save`, `list`, `run`, `delete`.",
    ],
  },
  {
    version: "v2.21.0",
    items: [
      "Added `gitmap stats` (`ss`) command: aggregated usage statistics from command history.",
    ],
  },
  {
    version: "v2.20.0",
    items: [
      "Added `gitmap history` (`hi`) command: queryable audit trail of CLI executions.",
      "Added `gitmap history-reset` (`hr`) command: clears audit history.",
    ],
  },
  {
    version: "v2.19.0",
    items: [
      "Added `gitmap amend` (`am`) command: rewrite author name/email on existing commits.",
      "Added `--author-name` and `--author-email` flags to `gitmap seo-write`.",
    ],
  },
  {
    version: "v2.18.0",
    items: [
      "Added `gitmap seo-write` (`sw`) command: automated SEO commit scheduler.",
      "Supports CSV input, template mode with placeholder substitution, and rotation.",
    ],
  },
  {
    version: "v2.17.0",
    items: [
      "Added `Source` column to the `Releases` table.",
      "Added `--source` flag to `list-releases`, `list-versions`, and `changelog`.",
    ],
  },
  {
    version: "v2.16.0",
    items: [
      "Added `gitmap list-releases` (`lr`) command.",
      "Enhanced `gitmap scan` to import `.release/v*.json` metadata files.",
    ],
  },
  {
    version: "v2.15.0",
    items: ["Added `--limit N` flag to `gitmap list-versions` (`lv`)."],
  },
  {
    version: "v2.14.0",
    items: [
      "Added `Releases` table to SQLite database.",
      "Converted all DB table/column names from snake_case to PascalCase.",
    ],
  },
  {
    version: "v2.13.0",
    items: [
      "Release metadata JSON now includes a `changelog` field.",
      "`gitmap list-versions` now shows changelog notes.",
    ],
  },
  {
    version: "v2.12.0",
    items: [
      "Added `gitmap list-versions` (`lv`) command.",
      "Added `gitmap revert <version>` command.",
    ],
  },
  {
    version: "v2.11.0",
    items: ["Added constants inventory audit section to compliance spec."],
  },
  {
    version: "v2.10.0",
    items: ["Version bump for next development cycle."],
  },
  {
    version: "v2.9.0",
    items: [
      "Full code style refactor of `latest-branch` command.",
      "Split handler into 3 files, all under 200 lines.",
    ],
  },
  {
    version: "v2.8.0",
    items: ["Added `--filter` flag to `latest-branch`."],
  },
  {
    version: "v2.7.0",
    items: ["Added `--sort` flag to `latest-branch`."],
  },
  {
    version: "v2.6.0",
    items: ["Centralized date display formatting with local timezone conversion."],
  },
  {
    version: "v2.5.0",
    items: ["Added `--format` flag to `latest-branch` (terminal, json, csv)."],
  },
  {
    version: "v2.4.0",
    items: [
      "Added `gitmap latest-branch` (`lb`) command.",
      "Positional integer shorthand: `gitmap lb 3` equals `gitmap lb --top 3`.",
    ],
  },
];
