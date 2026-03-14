import DocsLayout from "@/components/docs/DocsLayout";
import { changelog } from "@/data/changelog";
import { Tag, ChevronDown, ChevronRight } from "lucide-react";
import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

const ChangelogPage = () => {
  const [expandedVersions, setExpandedVersions] = useState<Set<string>>(
    new Set(changelog.slice(0, 3).map((e) => e.version))
  );

  const toggle = (version: string) => {
    setExpandedVersions((prev) => {
      const next = new Set(prev);
      next.has(version) ? next.delete(version) : next.add(version);
      return next;
    });
  };

  const expandAll = () => setExpandedVersions(new Set(changelog.map((e) => e.version)));
  const collapseAll = () => setExpandedVersions(new Set());

  return (
    <DocsLayout>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-mono font-bold">Changelog</h1>
          <p className="text-muted-foreground text-sm mt-1">
            {changelog.length} releases tracked
          </p>
        </div>
        <div className="flex gap-2">
          <button onClick={expandAll} className="text-xs font-mono text-muted-foreground hover:text-foreground transition-colors px-2 py-1 rounded border border-border">
            Expand all
          </button>
          <button onClick={collapseAll} className="text-xs font-mono text-muted-foreground hover:text-foreground transition-colors px-2 py-1 rounded border border-border">
            Collapse all
          </button>
        </div>
      </div>

      <div className="relative">
        {/* Timeline line */}
        <div className="absolute left-[15px] top-0 bottom-0 w-px bg-border" />

        <div className="space-y-2">
          {changelog.map((entry, i) => {
            const isOpen = expandedVersions.has(entry.version);
            const isLatest = i === 0;

            return (
              <div key={entry.version} className="relative pl-10">
                {/* Timeline dot */}
                <div className={`absolute left-[10px] top-3 h-[11px] w-[11px] rounded-full border-2 ${isLatest ? "border-primary bg-primary" : "border-muted-foreground/40 bg-background"}`} />

                <button
                  onClick={() => toggle(entry.version)}
                  className="w-full flex items-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-card hover:bg-muted/50 transition-colors text-left"
                >
                  {isOpen ? (
                    <ChevronDown className="h-4 w-4 text-primary shrink-0" />
                  ) : (
                    <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
                  )}
                  <Tag className="h-3.5 w-3.5 text-primary shrink-0" />
                  <span className="font-mono font-semibold text-sm">{entry.version}</span>
                  {isLatest && (
                    <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-primary/10 text-primary">
                      latest
                    </span>
                  )}
                  <span className="text-xs text-muted-foreground ml-auto">
                    {entry.items.length} change{entry.items.length !== 1 ? "s" : ""}
                  </span>
                </button>

                <AnimatePresence initial={false}>
                  {isOpen && (
                    <motion.div
                      initial={{ height: 0, opacity: 0 }}
                      animate={{ height: "auto", opacity: 1 }}
                      exit={{ height: 0, opacity: 0 }}
                      transition={{ duration: 0.2 }}
                      className="overflow-hidden"
                    >
                      <ul className="mt-1 ml-4 space-y-1 pb-2">
                        {entry.items.map((item, j) => (
                          <li key={j} className="text-sm text-muted-foreground flex gap-2">
                            <span className="text-primary mt-1.5 shrink-0">•</span>
                            <span className="leading-relaxed">{item}</span>
                          </li>
                        ))}
                      </ul>
                    </motion.div>
                  )}
                </AnimatePresence>
              </div>
            );
          })}
        </div>
      </div>
    </DocsLayout>
  );
};

export default ChangelogPage;
