import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import CommandCard from "@/components/docs/CommandCard";
import CommandCategoryGroup from "@/components/docs/CommandCategoryGroup";
import SearchBar from "@/components/docs/SearchBar";
import { commands, categories } from "@/data/commands";

const CommandsPage = () => {
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    if (!search) return commands;
    const q = search.toLowerCase();
    return commands.filter(
      (c) =>
        c.name.includes(q) ||
        c.alias?.includes(q) ||
        c.description.toLowerCase().includes(q)
    );
  }, [search]);

  const isSearching = search.length > 0;

  return (
    <DocsLayout>
      <h1 className="text-3xl font-mono font-bold mb-2">Command Reference</h1>
      <p className="text-muted-foreground mb-6">
        All {commands.length} gitmap commands organized by category.
      </p>

      <SearchBar value={search} onChange={setSearch} />

      <div className="mt-6 space-y-3">
        {isSearching ? (
          <>
            {filtered.map((cmd) => (
              <CommandCard key={cmd.name} {...cmd} />
            ))}
            {filtered.length === 0 && (
              <p className="text-center text-muted-foreground py-8 font-mono text-sm">
                No commands matching "{search}"
              </p>
            )}
          </>
        ) : (
          categories.map((cat) => {
            const cmds = filtered.filter((c) => c.category === cat.key);
            if (cmds.length === 0) return null;
            return (
              <CommandCategoryGroup
                key={cat.key}
                label={cat.label}
                description={cat.description}
                commands={cmds}
              />
            );
          })
        )}
      </div>
    </DocsLayout>
  );
};

export default CommandsPage;
