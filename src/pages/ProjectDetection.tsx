import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Search, FolderGit2, Code2, FileCode, Braces, Cpu } from "lucide-react";

const projectTypes = [
  { icon: Code2, type: "Go", key: "go", indicator: "go.mod", color: "text-cyan-400", commands: "go-repos (gr)" },
  { icon: Braces, type: "Node.js", key: "node", indicator: "package.json", color: "text-green-400", commands: "node-repos (nr)" },
  { icon: FileCode, type: "React", key: "react", indicator: "package.json + react dep", color: "text-blue-400", commands: "react-repos (rr)" },
  { icon: Cpu, type: "C++", key: "cpp", indicator: "CMakeLists.txt / *.vcxproj", color: "text-orange-400", commands: "cpp-repos (cr)" },
  { icon: FolderGit2, type: "C#", key: "csharp", indicator: "*.csproj / *.sln", color: "text-purple-400", commands: "csharp-repos (csr)" },
];

const detectionRules = [
  { type: "Go", primary: "go.mod exists", secondary: "go.sum, *.go files", falsePositive: "Ignores vendor/, testdata/" },
  { type: "Node.js", primary: "package.json exists", secondary: "package-lock.json, yarn.lock, pnpm-lock.yaml, bun.lock", falsePositive: "Ignores node_modules/, vendor/" },
  { type: "React", primary: "package.json with react dependency", secondary: "@types/react, react-scripts, next, gatsby, remix", falsePositive: "Reclassified from Node.js — exclusive" },
  { type: "C++", primary: "CMakeLists.txt, *.vcxproj, meson.build", secondary: "Makefile + C++ sources, conanfile, vcpkg.json", falsePositive: "Ignores build/, cmake-build-*/, out/" },
  { type: "C#", primary: "*.csproj, *.sln", secondary: "*.fsproj, global.json, *.cs files", falsePositive: "Ignores bin/, obj/, packages/. .sln takes precedence" },
];

const excludeDirs = [
  "node_modules", "vendor", ".git", "dist", "build", "target",
  "bin", "obj", "out", "cmake-build-*", "testdata", "packages", ".venv", ".cache",
];

const ProjectDetectionPage = () => (
  <DocsLayout>
    <h1 className="text-3xl font-mono font-bold mb-2">Project Detection</h1>
    <p className="text-muted-foreground mb-6">
      Automatic technology stack detection during scan and rescan, with dedicated query commands per project type.
    </p>

    {/* Type cards */}
    <h2 className="text-xl font-mono font-semibold mt-8 mb-4">Supported Project Types</h2>
    <div className="grid grid-cols-2 md:grid-cols-5 gap-3 mb-8">
      {projectTypes.map((p) => (
        <div key={p.key} className="rounded-lg border border-border bg-card p-4 text-center">
          <p.icon className={`h-8 w-8 mx-auto mb-2 ${p.color}`} />
          <h3 className="font-mono font-semibold text-sm">{p.type}</h3>
          <p className="text-xs text-muted-foreground mt-1 font-mono">{p.indicator}</p>
        </div>
      ))}
    </div>

    {/* Detection rules */}
    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Detection Rules</h2>
    <div className="rounded-lg border border-border overflow-hidden overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-muted/50">
            <th className="text-left font-mono font-semibold px-4 py-2">Type</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Primary Indicator</th>
            <th className="text-left font-mono font-semibold px-4 py-2">Secondary</th>
            <th className="text-left font-mono font-semibold px-4 py-2">False Positive Prevention</th>
          </tr>
        </thead>
        <tbody>
          {detectionRules.map((r) => (
            <tr key={r.type} className="border-t border-border">
              <td className="px-4 py-2 font-mono text-primary font-semibold">{r.type}</td>
              <td className="px-4 py-2 text-foreground">{r.primary}</td>
              <td className="px-4 py-2 text-muted-foreground">{r.secondary}</td>
              <td className="px-4 py-2 text-muted-foreground">{r.falsePositive}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>

    {/* Query commands */}
    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Query Commands</h2>
    <p className="text-sm text-muted-foreground mb-4">
      Each project type has a dedicated command for instant filtered access:
    </p>
    <div className="grid md:grid-cols-2 gap-3 mb-6">
      {projectTypes.map((p) => (
        <div key={p.key} className="flex items-center gap-3 rounded-lg border border-border bg-card px-4 py-3">
          <p.icon className={`h-4 w-4 ${p.color} shrink-0`} />
          <code className="font-mono text-sm text-primary">{p.commands}</code>
          <span className="text-xs text-muted-foreground">List detected {p.type} projects</span>
        </div>
      ))}
    </div>

    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Usage</h2>
    <CodeBlock code="gitmap go-repos" title="List all Go projects" />
    <CodeBlock code="gitmap node-repos --json" title="Node.js projects as JSON" />
    <CodeBlock code="gitmap react-repos" title="List React projects" />
    <CodeBlock code="gitmap csharp-repos --json" title="C# projects as JSON" />
    <CodeBlock code="gitmap cpp-repos" title="List C++ projects" />

    {/* Monorepo handling */}
    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Monorepo & Nesting</h2>
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      <div className="rounded-lg border border-border bg-card p-4">
        <Search className="h-5 w-5 text-primary mb-2" />
        <h3 className="font-mono font-semibold text-sm mb-1">Monorepo Support</h3>
        <p className="text-xs text-muted-foreground">
          A single repo with <code className="text-primary">backend/</code> (Go) and <code className="text-primary">frontend/</code> (React)
          produces two separate detection records linked to the same repo.
        </p>
      </div>
      <div className="rounded-lg border border-border bg-card p-4">
        <FolderGit2 className="h-5 w-5 text-primary mb-2" />
        <h3 className="font-mono font-semibold text-sm mb-1">Nested Projects</h3>
        <p className="text-xs text-muted-foreground">
          A Node.js project at root with a React project at <code className="text-primary">web/</code> records both.
          The more specific classification wins at each path level.
        </p>
      </div>
    </div>

    {/* C# precedence */}
    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">C# Solution Precedence</h2>
    <p className="text-sm text-muted-foreground mb-4">
      When a <code className="text-primary font-mono">.sln</code> file is found, it defines a single project entry.
      Individual <code className="text-primary font-mono">.csproj</code> files beneath it are stored as child records, not separate projects.
      Standalone <code className="text-primary font-mono">.csproj</code> files (no parent .sln) become their own project entries.
    </p>

    {/* Excluded dirs */}
    <h2 className="text-xl font-mono font-semibold mt-10 mb-3">Excluded Directories</h2>
    <div className="flex flex-wrap gap-2 mb-8">
      {excludeDirs.map((dir) => (
        <span key={dir} className="font-mono text-xs bg-muted text-muted-foreground px-2 py-1 rounded border border-border">
          {dir}
        </span>
      ))}
    </div>

    {/* File layout */}
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
            ["detector/detector.go", "Detection orchestration"],
            ["detector/rules.go", "Per-type detection heuristics"],
            ["detector/goparser.go", "Go module parsing"],
            ["detector/csharpparser.go", "C# solution/project parsing"],
            ["detector/parser.go", "Shared parsing utilities"],
            ["model/project.go", "DetectedProject struct"],
            ["model/projecttype.go", "ProjectType enum/struct"],
            ["store/project.go", "Project CRUD operations"],
            ["cmd/projectrepos.go", "Query command handlers"],
            ["constants/constants_project.go", "Project detection constants"],
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

export default ProjectDetectionPage;
