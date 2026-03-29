/** Repo status values used in watch/status displays. */
export enum RepoStatus {
  Clean = "clean",
  Dirty = "dirty",
}

/** Terminal line type for animated demos. */
export enum TerminalLineType {
  Input = "input",
  Output = "output",
  Header = "header",
  Accent = "accent",
}

/** Project filter including the "all" option. */
export const FILTER_ALL = "all" as const;

/** Animation timing defaults (ms). */
export const TERMINAL_INPUT_DELAY = 600;
export const TERMINAL_OUTPUT_DELAY = 120;

/** Watch dashboard refresh interval (seconds). */
export const WATCH_REFRESH_INTERVAL = 30;

/** Status indicator symbols. */
export const STATUS_ICON_DIRTY = "●";
export const STATUS_ICON_CLEAN = "✔";

/** Root-relative path placeholder. */
export const ROOT_RELATIVE_PATH = ".";
export const ROOT_RELATIVE_LABEL = "(root)";
