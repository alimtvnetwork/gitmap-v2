import { useState, useMemo } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { FileText, AlertTriangle, Compass, Terminal, ChevronRight, Search, X } from "lucide-react";

interface SpecEntry {
  id: string;
  title: string;
  link?: string;
}

interface SpecSection {
  folder: string;
  title: string;
  description: string;
  icon: React.ReactNode;
  color: string;
  entries: SpecEntry[];
}

const sections: SpecSection[] = [
  {
    folder: "01-app",
    title: "Application Specifications",
    description: "Feature specs, command designs, and refactor documentation for the gitmap CLI.",
    icon: <FileText className="h-5 w-5" />,
    color: "text-primary",
    entries: [
      { id: "01", title: "Overview" },
      { id: "02", title: "CLI Interface" },
      { id: "03", title: "Scanner" },
      { id: "04", title: "Formatter" },
      { id: "05", title: "Cloner" },
      { id: "06", title: "Config", link: "/config" },
      { id: "07", title: "Data Model" },
      { id: "08", title: "Acceptance Criteria" },
      { id: "09", title: "Build & Deploy" },
      { id: "10", title: "GitHub Desktop" },
      { id: "11", title: "Desktop Sync" },
      { id: "12", title: "Release Command", link: "/release" },
      { id: "13", title: "Release Data Model" },
      { id: "14", title: "Latest Branch" },
      { id: "15", title: "Date Display Format" },
      { id: "16", title: "Database" },
      { id: "17", title: "Repo Grouping" },
      { id: "18", title: "Compliance Audit" },
      { id: "19", title: "List Versions" },
      { id: "20", title: "Revert" },
      { id: "21", title: "List Releases" },
      { id: "22", title: "Scan Release Import" },
      { id: "23", title: "SEO Write" },
      { id: "24", title: "Amend Author" },
      { id: "25", title: "Command History", link: "/history" },
      { id: "26", title: "Stats", link: "/stats" },
      { id: "27", title: "Bookmarks", link: "/bookmarks" },
      { id: "28", title: "Export", link: "/export" },
      { id: "29", title: "Import", link: "/import" },
      { id: "30", title: "Profiles", link: "/profile" },
      { id: "31", title: "cd" },
      { id: "32", title: "Watch", link: "/watch" },
      { id: "33", title: "Diff Profiles", link: "/diff-profiles" },
      { id: "34", title: "Clone Progress" },
      { id: "35", title: "Docs Site" },
      { id: "36", title: "GoMod Rename", link: "/gomod" },
      { id: "37", title: "Project Detection", link: "/project-detection" },
      { id: "38", title: "Command Help" },
      { id: "39", title: "Shell Completion" },
      { id: "40", title: "Enhanced Groups & Listing" },
      { id: "41", title: "Go Release Assets" },
      { id: "42", title: "Cross-Platform" },
      { id: "43", title: "Interactive TUI", link: "/interactive" },
      { id: "44", title: "List DB Diagnostic" },
      { id: "45", title: "Release Pending Metadata" },
      { id: "46", title: "Clear Release JSON", link: "/clear-release-json" },
      { id: "47", title: "Zip Groups", link: "/zip-group" },
      { id: "48", title: "Repo Aliases", link: "/alias" },
      { id: "49", title: "Changelog Generate", link: "/changelog-generate" },
      { id: "50", title: "SSH Keys", link: "/ssh" },
      { id: "51", title: "Prune", link: "/prune" },
      { id: "52", title: "Upload Retry" },
      { id: "53", title: "Offline Detection" },
      { id: "54", title: "Process Locking" },
      { id: "55", title: "Temp Release", link: "/temp-release" },
      { id: "56", title: "Unified .gitmap Dir" },
      { id: "57", title: "Skip-Meta Integration Test" },
      { id: "58–78", title: "Refactors (workflow, dispatch, archive, autocommit, SEO, branches, assets, TUI, aliases, ops, status, exec, logs, compress)" },
    ],
  },
  {
    folder: "02-app-issues",
    title: "Issue Post-Mortems",
    description: "Root-cause analyses and resolution records for production bugs.",
    icon: <AlertTriangle className="h-5 w-5" />,
    color: "text-yellow-500",
    entries: [
      { id: "01", title: "Update File Lock" },
      { id: "02", title: "Update Flow Spec Alignment" },
      { id: "03", title: "Update Sync Lock Loop" },
      { id: "04", title: "Database Path Resolution" },
      { id: "05", title: "List Empty DB Path" },
      { id: "06", title: "Release Orphaned Metadata" },
      { id: "07", title: "Zip Group Release Silent Failure" },
      { id: "08", title: "Autocommit Push Rejection" },
      { id: "09", title: "List Releases Repo Source" },
      { id: "10", title: "Legacy UUID Detection" },
      { id: "11", title: "Auto Legacy Dir Migration" },
      { id: "12", title: "Legacy ID Migration" },
    ],
  },
  {
    folder: "03-general",
    title: "Design Guidelines",
    description: "Reusable architectural patterns and coding standards — generic and shareable across projects.",
    icon: <Compass className="h-5 w-5" />,
    color: "text-green-500",
    entries: [
      { id: "01", title: "CLI Design Patterns" },
      { id: "02", title: "PowerShell Build & Deploy" },
      { id: "03", title: "Self-Update Mechanism" },
      { id: "04", title: "Output Formatting" },
      { id: "05", title: "Configuration Pattern" },
      { id: "06", title: "Code Style Rules" },
      { id: "07", title: "Date Display Format" },
    ],
  },
  {
    folder: "04-generic-cli",
    title: "Generic CLI Blueprint",
    description: "A production-quality CLI implementation blueprint usable as a starting point for any Go CLI project.",
    icon: <Terminal className="h-5 w-5" />,
    color: "text-blue-400",
    entries: [
      { id: "01", title: "Overview", link: "/generic-cli" },
      { id: "02", title: "Project Structure" },
      { id: "03", title: "Subcommand Architecture" },
      { id: "04", title: "Flag Parsing" },
      { id: "05", title: "Configuration" },
      { id: "06", title: "Output Formatting" },
      { id: "07", title: "Error Handling" },
      { id: "08", title: "Code Style" },
      { id: "09", title: "Help System" },
      { id: "10", title: "Database" },
      { id: "11", title: "Build & Deploy" },
      { id: "12", title: "Testing" },
      { id: "13", title: "Checklist" },
      { id: "14", title: "Date Formatting" },
      { id: "15", title: "Constants Reference" },
      { id: "16", title: "Verbose Logging" },
      { id: "17", title: "Progress Tracking" },
      { id: "18", title: "Batch Execution" },
      { id: "19", title: "Shell Completion" },
    ],
  },
];

const container = { hidden: {}, show: { transition: { staggerChildren: 0.08 } } };
const item = { hidden: { opacity: 0, y: 12 }, show: { opacity: 1, y: 0 } };

const SpecIndexPage = () => {
  const [query, setQuery] = useState("");

  const filtered = useMemo(() => {
    const q = query.toLowerCase().trim();
    if (!q) return sections;
    return sections
      .map((section) => ({
        ...section,
        entries: section.entries.filter(
          (e) =>
            e.title.toLowerCase().includes(q) ||
            e.id.toLowerCase().includes(q) ||
            section.folder.toLowerCase().includes(q) ||
            section.title.toLowerCase().includes(q)
        ),
      }))
      .filter((s) => s.entries.length > 0);
  }, [query]);

  const totalResults = filtered.reduce((sum, s) => sum + s.entries.length, 0);

  return (
    <DocsLayout>
      <motion.div initial={{ opacity: 0, y: -10 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.3 }}>
        <h1 className="text-3xl font-mono font-bold mb-2">Spec Index</h1>
        <p className="text-muted-foreground mb-2">
          Complete table of contents for all specification documents, issue post-mortems, design guidelines, and the generic CLI blueprint.
        </p>

        {/* Search bar */}
        <div className="relative mb-6">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Filter specs… (e.g. release, database, TUI)"
            className="w-full pl-9 pr-9 py-2.5 text-sm font-mono bg-muted/30 border border-border rounded-lg text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary/50 transition-colors"
          />
          {query && (
            <button
              onClick={() => setQuery("")}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
            >
              <X className="h-4 w-4" />
            </button>
          )}
        </div>

        <p className="text-xs text-muted-foreground/60 font-mono mb-6">
          {query ? `${totalResults} result${totalResults !== 1 ? "s" : ""} matching "${query}"` : `${totalResults} documents across ${sections.length} sections`}
        </p>
      </motion.div>

      <motion.div variants={container} initial="hidden" animate="show" className="space-y-8">
        {sections.map((section) => (
          <motion.div key={section.folder} variants={item}>
            <div className="border border-border rounded-lg overflow-hidden">
              {/* Section header */}
              <div className="bg-muted/30 px-5 py-4 border-b border-border">
                <div className="flex items-center gap-3 mb-1">
                  <span className={section.color}>{section.icon}</span>
                  <h2 className="text-lg font-mono font-semibold text-foreground">
                    <span className="text-muted-foreground">{section.folder}/</span> {section.title}
                  </h2>
                  <span className="ml-auto text-xs font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded">
                    {section.entries.length} docs
                  </span>
                </div>
                <p className="text-sm text-muted-foreground ml-8">{section.description}</p>
              </div>

              {/* Entries */}
              <div className="divide-y divide-border">
                {section.entries.map((entry) => (
                  <div key={entry.id} className="flex items-center gap-3 px-5 py-2.5 hover:bg-muted/20 transition-colors group">
                    <span className="text-xs font-mono text-muted-foreground w-10 shrink-0">{entry.id}</span>
                    <span className="text-sm text-foreground">{entry.title}</span>
                    {entry.link && (
                      <Link
                        to={entry.link}
                        className="ml-auto flex items-center gap-1 text-xs font-mono text-primary opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                        docs <ChevronRight className="h-3 w-3" />
                      </Link>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </motion.div>
        ))}
      </motion.div>

      {/* See Also */}
      <div className="mt-10 pt-6 border-t border-border">
        <h3 className="text-sm font-mono font-semibold text-muted-foreground mb-3">See Also</h3>
        <div className="flex flex-wrap gap-2">
          {[
            { label: "Architecture", to: "/architecture" },
            { label: "Commands", to: "/commands" },
            { label: "Generic CLI", to: "/generic-cli" },
            { label: "Changelog", to: "/changelog" },
          ].map((link) => (
            <Link
              key={link.to}
              to={link.to}
              className="text-xs font-mono px-3 py-1.5 rounded border border-border bg-muted/30 text-muted-foreground hover:text-primary hover:border-primary/50 transition-colors"
            >
              {link.label}
            </Link>
          ))}
        </div>
      </div>
    </DocsLayout>
  );
};

export default SpecIndexPage;
