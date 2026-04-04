import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Link2, Zap, Database, Shield, Search, Terminal } from "lucide-react";

const MOCK_ALIASES = [
  { alias: "api", slug: "github/user/api-gateway", created: "2026-03-15" },
  { alias: "web", slug: "github/user/web-frontend", created: "2026-03-14" },
  { alias: "infra", slug: "github/user/infrastructure", created: "2026-03-12" },
  { alias: "libs", slug: "github/user/shared-libs", created: "2026-03-10" },
];

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap alias list</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-primary font-bold text-xs mb-1">
        {"  "}Aliases (4):
      </div>
      {MOCK_ALIASES.map((a) => (
        <div key={a.alias} className="text-terminal-foreground text-xs">
          {"  "}
          <span className="inline-block w-[100px] text-foreground font-semibold">{a.alias}</span>
          <span className="text-muted-foreground">→ </span>
          <span className="text-primary">{a.slug}</span>
        </div>
      ))}
    </div>
  </div>
);

const SuggestPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap alias suggest</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-muted-foreground text-xs mb-2">{"  "}Suggested aliases:</div>
      <div className="text-xs text-terminal-foreground">
        {"    "}<span className="text-foreground">api-gateway</span>
        <span className="text-muted-foreground">{"  → "}</span>
        <span className="text-primary">api</span>
        <span className="text-muted-foreground">{"       Accept? (y/N): "}</span>
        <span className="text-green-400">y</span>
      </div>
      <div className="text-xs text-terminal-foreground">
        {"    "}<span className="text-foreground">web-frontend</span>
        <span className="text-muted-foreground">{" → "}</span>
        <span className="text-primary">web</span>
        <span className="text-muted-foreground">{"       Accept? (y/N): "}</span>
        <span className="text-green-400">y</span>
      </div>
      <div className="text-xs text-terminal-foreground">
        {"    "}<span className="text-foreground">shared-libs</span>
        <span className="text-muted-foreground">{"  → "}</span>
        <span className="text-primary">libs</span>
        <span className="text-muted-foreground">{"      Accept? (y/N): "}</span>
        <span className="text-red-400">n</span>
      </div>
      <div className="mt-3 text-xs text-green-400">{"  "}✓ Created 2 alias(es).</div>
    </div>
  </div>
);

const features = [
  { icon: Link2, title: "Short Names", desc: "Replace long slugs with concise aliases like 'api', 'web', or 'infra'." },
  { icon: Zap, title: "Run From Anywhere", desc: "Execute any gitmap command against a repo using -A <alias> without changing directory." },
  { icon: Search, title: "Auto-Suggest", desc: "During scan/rescan, aliases are suggested based on repo name or slug." },
  { icon: Shield, title: "Conflict-Safe", desc: "Warns and prompts when an alias collision occurs. Cannot shadow gitmap commands." },
];

const schema = [
  ["Id", "TEXT", "PRIMARY KEY", "UUID"],
  ["Alias", "TEXT", "NOT NULL UNIQUE", "Short name"],
  ["RepoId", "TEXT", "FK → Repos(Id) CASCADE", "Target repository"],
  ["CreatedAt", "TEXT", "DEFAULT CURRENT_TIMESTAMP", ""],
];

const commandInteraction = [
  ["cd", "Navigate to aliased repo directory"],
  ["pull", "Pull in aliased repo"],
  ["exec", "Execute command in aliased repo directory"],
  ["status", "Show status of aliased repo"],
  ["watch", "Watch aliased repo for changes"],
  ["release", "Run release from aliased repo"],
  ["scan", "No effect (operates on directories)"],
  ["group", "No effect (operates on group names)"],
];

const AliasPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Repo Aliases</h1>
    <p className="text-muted-foreground mb-6">
      Assign short, memorable names to repositories for quick access from anywhere.
    </p>

    <h2 className="text-xl font-heading font-semibold mt-8 mb-2">Live Preview</h2>
    <TerminalPreview />

    <h2 className="text-xl font-heading font-semibold mt-8 mb-2">Auto-Suggestion</h2>
    <p className="text-sm text-muted-foreground mb-2">
      During scan or via <code className="text-primary font-mono">gitmap alias suggest</code>,
      aliases are proposed based on repo names.
    </p>
    <SuggestPreview />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-4">Features</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div key={f.title} className="rounded-lg border border-border bg-card p-4">
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
          <p className="text-xs text-muted-foreground">{f.desc}</p>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Subcommands</h2>
    <CodeBlock code="gitmap alias set api github/user/api-gateway" title="Create an alias" />
    <CodeBlock code="gitmap a set web github/user/web-frontend" title="Using short alias 'a'" />
    <CodeBlock code="gitmap alias list" title="List all aliases" />
    <CodeBlock code="gitmap alias show api" title="Show alias details" />
    <CodeBlock code="gitmap alias remove api" title="Remove an alias" />
    <CodeBlock code="gitmap alias suggest" title="Auto-suggest aliases for unaliased repos" />
    <CodeBlock code="gitmap alias suggest --apply" title="Auto-accept all suggestions" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Global -A Flag</h2>
    <p className="text-sm text-muted-foreground mb-4">
      Any repo-targeting command accepts <code className="text-primary font-mono">-A &lt;alias&gt;</code> to
      resolve the target by alias instead of requiring <code className="text-primary font-mono">cd</code>.
    </p>
    <CodeBlock code="gitmap pull -A api" title="Pull via alias" />
    <CodeBlock code="gitmap exec -A web -- npm test" title="Run command in aliased repo" />
    <CodeBlock code="gitmap cd -A infra" title="Navigate to aliased repo" />
    <CodeBlock code="gitmap status -A api" title="Check status via alias" />

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Command Interaction</h2>
    <div className="rounded-lg border border-border overflow-hidden mb-8">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Command</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Behavior with -A</th>
          </tr>
        </thead>
        <tbody>
          {commandInteraction.map(([cmd, behavior]) => (
            <tr key={cmd} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{cmd}</td>
              <td className="px-4 py-2 text-muted-foreground">{behavior}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Conflict Handling</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">Manual Set</h3>
        <p className="text-xs text-muted-foreground">
          If alias exists, prompts: <code className="text-primary">"Reassign to new repo? (y/N)"</code>
        </p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">Auto-Suggest</h3>
        <p className="text-xs text-muted-foreground">
          Conflicting aliases are skipped with a warning message.
        </p>
      </div>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">Table Schema</h2>
    <div className="rounded-lg border border-border overflow-hidden mb-8">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Column</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Type</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Constraints</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Notes</th>
          </tr>
        </thead>
        <tbody>
          {schema.map(([col, type, constraints, notes]) => (
            <tr key={col} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{col}</td>
              <td className="px-4 py-2 font-mono text-muted-foreground">{type}</td>
              <td className="px-4 py-2 text-muted-foreground">{constraints}</td>
              <td className="px-4 py-2 text-muted-foreground">{notes}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-heading font-semibold mt-10 mb-3">File Layout</h2>
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">File</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Purpose</th>
          </tr>
        </thead>
        <tbody>
          {[
            ["cmd/alias.go", "Subcommand dispatch"],
            ["cmd/aliasops.go", "Subcommand implementation"],
            ["cmd/aliasresolve.go", "-A flag resolution logic"],
            ["store/alias.go", "Database CRUD for Aliases"],
            ["model/alias.go", "Data struct"],
            ["constants/constants_alias.go", "Messages, SQL, flag descriptions"],
            ["helptext/alias.md", "Command help"],
          ].map(([file, purpose]) => (
            <tr key={file} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{file}</td>
              <td className="px-4 py-2 text-muted-foreground">{purpose}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </DocsLayout>
);

export default AliasPage;
