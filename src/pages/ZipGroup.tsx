import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Archive, FolderPlus, Trash2, List, Eye, Database, Tag } from "lucide-react";

const MOCK_GROUPS = [
  { name: "docs-bundle", items: 3, archive: "docs-bundle_v3.0.0.zip" },
  { name: "extras", items: 2, archive: "extra-files.zip" },
  { name: "configs", items: 4, archive: "" },
];

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap zip-group list</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-primary font-bold text-xs mb-1">
        {"  "}Zip Groups (3):
      </div>
      {MOCK_GROUPS.map((g) => (
        <div key={g.name} className="text-terminal-foreground text-xs">
          {"  "}
          <span className="inline-block w-[160px] text-foreground">{g.name}</span>
          <span className="inline-block w-[100px] text-muted-foreground">{g.items} item(s)</span>
          <span className="text-primary">{g.archive || "—"}</span>
        </div>
      ))}
    </div>
  </div>
);

const ShowPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap z show docs-bundle</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed overflow-x-auto">
      <div className="text-primary font-bold text-xs mb-1">{"  "}docs-bundle (3 item(s)):</div>
      <div className="text-xs text-terminal-foreground">{"    "}📄 README.md</div>
      <div className="text-xs text-terminal-foreground">{"    "}📄 CHANGELOG.md</div>
      <div className="text-xs text-terminal-foreground">{"    "}📁 docs/</div>
      <div className="text-xs text-muted-foreground mt-2">{"  "}Archive: docs-bundle_v3.0.0.zip</div>
    </div>
  </div>
);

const features = [
  { icon: Archive, title: "Max Compression", desc: "All archives use ZIP format with Deflate level 9 for maximum compression." },
  { icon: FolderPlus, title: "Persistent Groups", desc: "Save named file/folder collections in SQLite for reuse across releases." },
  { icon: Tag, title: "Release Integration", desc: "Use --zip-group to include groups as release assets, or -Z for ad-hoc items." },
  { icon: Database, title: "Metadata Tracking", desc: "Zip group definitions are recorded in .release/vX.Y.Z.json under zipGroups." },
];

const schema = [
  ["Id", "TEXT", "PRIMARY KEY", "UUID"],
  ["Name", "TEXT", "NOT NULL UNIQUE", "Group name"],
  ["ArchiveName", "TEXT", "DEFAULT ''", "Custom output filename"],
  ["CreatedAt", "TEXT", "DEFAULT CURRENT_TIMESTAMP", ""],
];

const itemSchema = [
  ["GroupId", "TEXT", "FK → ZipGroups(Id) CASCADE", "Parent group"],
  ["Path", "TEXT", "NOT NULL", "File or folder path"],
  ["IsFolder", "INTEGER", "DEFAULT 0", "1 = folder, 0 = file"],
];

const releaseFlags = [
  ["--zip-group <name>", "Include a persistent zip group as a release asset"],
  ["-Z <path>", "Add ad-hoc file or folder to zip as a release asset"],
  ["--bundle <name.zip>", "Bundle all -Z items into a single named archive"],
];

const ZipGroupPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-mono font-bold mb-2">Zip Groups</h1>
    <p className="text-muted-foreground mb-6">
      Define named collections of files and folders that are automatically compressed
      into ZIP archives during a release.
    </p>

    <h2 className="text-xl font-mono font-semibold mt-8 mb-2">Live Preview</h2>
    <TerminalPreview />
    <ShowPreview />

    <h2 className="text-xl font-mono font-semibold mt-10 mb-4">Features</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div key={f.title} className="rounded-lg border border-border bg-card p-4">
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
          <p className="text-xs text-muted-foreground">{f.desc}</p>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Subcommands</h2>
    <CodeBlock code="gitmap z create docs-bundle" title="Create a zip group" />
    <CodeBlock code="gitmap z add docs-bundle ./README.md ./CHANGELOG.md ./docs/" title="Add items to a group" />
    <CodeBlock code="gitmap z create extras --archive extra-files.zip" title="Create with custom archive name" />
    <CodeBlock code="gitmap z show docs-bundle" title="Show group contents" />
    <CodeBlock code="gitmap z list" title="List all zip groups" />
    <CodeBlock code="gitmap z rename docs-bundle --archive release-docs.zip" title="Set custom archive name" />
    <CodeBlock code="gitmap z remove docs-bundle ./CHANGELOG.md" title="Remove an item from a group" />
    <CodeBlock code="gitmap z delete extras" title="Delete a zip group" />

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Release Integration</h2>
    <p className="text-sm text-muted-foreground mb-4">
      Use persistent groups or ad-hoc items during a release. Each group produces a single
      <code className="text-primary font-mono"> .zip</code> archive attached as a release asset.
    </p>
    <CodeBlock code="gitmap release v3.0.0 --zip-group docs-bundle" title="Release with a persistent zip group" />
    <CodeBlock code="gitmap release v3.0.0 -Z ./dist/report.pdf -Z ./dist/manual.pdf --bundle docs.zip" title="Ad-hoc bundle" />
    <CodeBlock code="gitmap release v3.0.0 --zip-group docs-bundle -Z ./extras/notes.txt" title="Combined: group + ad-hoc" />

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Release Flags</h2>
    <div className="rounded-lg border border-border overflow-hidden mb-8">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Flag</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Description</th>
          </tr>
        </thead>
        <tbody>
          {releaseFlags.map(([flag, desc]) => (
            <tr key={flag} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary">{flag}</td>
              <td className="px-4 py-2 text-muted-foreground">{desc}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Ad-Hoc Bundling</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">Without --bundle</h3>
        <p className="text-xs text-muted-foreground">Each <code className="text-primary">-Z</code> item is zipped individually into its own archive.</p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <h3 className="font-mono font-semibold text-sm mb-1 text-primary">With --bundle name.zip</h3>
        <p className="text-xs text-muted-foreground">All <code className="text-primary">-Z</code> items are bundled into a single named archive.</p>
      </div>
    </div>

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Table Schema: ZipGroups</h2>
    <div className="rounded-lg border border-border overflow-hidden mb-6">
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

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Table Schema: ZipGroupItems</h2>
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
          {itemSchema.map(([col, type, constraints, notes]) => (
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

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">File Layout</h2>
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
            ["cmd/zipgroup.go", "Subcommand dispatch"],
            ["cmd/zipgroupops.go", "Subcommand implementation"],
            ["release/ziparchive.go", "ZIP creation with max compression"],
            ["store/zipgroup.go", "Database CRUD for ZipGroups/ZipGroupItems"],
            ["model/zipgroup.go", "Data structs"],
            ["constants/constants_zipgroup.go", "Messages, SQL, flag descriptions"],
            ["helptext/zip-group.md", "Command help"],
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

export default ZipGroupPage;
