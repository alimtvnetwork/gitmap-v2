import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitBranch, Tag, Upload, Clock, Shield, Eye } from "lucide-react";

const features = [
  { icon: GitBranch, title: "Branch + Tag", desc: "Creates release/vX.Y.Z branch and vX.Y.Z tag in one step." },
  { icon: Tag, title: "Semver Padding", desc: "Partial versions auto-pad: v1 → v1.0.0, v1.2 → v1.2.0." },
  { icon: Upload, title: "Auto Push", desc: "Pushes branch and tag to origin after creation." },
  { icon: Clock, title: "Auto Increment", desc: "Use --bump to increment from the latest released version." },
  { icon: Shield, title: "Duplicate Detection", desc: "Aborts if the version tag or metadata file already exists." },
  { icon: Eye, title: "Dry Run", desc: "Preview all steps without executing with --dry-run." },
];

const releaseFlags = [
  { flag: "--assets <path>", description: "Directory or file to record as release assets" },
  { flag: "--commit <sha>", description: "Create release from a specific commit" },
  { flag: "--branch <name>", description: "Create release from latest commit of a branch" },
  { flag: "--bump major|minor|patch", description: "Auto-increment from the latest released version" },
  { flag: "--draft", description: "Mark release metadata as draft" },
  { flag: "--dry-run", description: "Preview release steps without executing" },
  { flag: "--verbose", description: "Write detailed debug log" },
];

const branchFlags = [
  { flag: "--assets <path>", description: "Directory or file to record" },
  { flag: "--draft", description: "Mark release metadata as draft" },
  { flag: "--dry-run", description: "Preview steps without executing" },
  { flag: "--verbose", description: "Write detailed debug log" },
];

const bumpExamples = [
  { current: "1.2.3", patch: "1.2.4", minor: "1.3.0", major: "2.0.0" },
  { current: "0.9.1", patch: "0.9.2", minor: "0.10.0", major: "1.0.0" },
];

const paddingExamples = [
  { input: "v1", resolved: "v1.0.0", branch: "release/v1.0.0", tag: "v1.0.0" },
  { input: "v1.2", resolved: "v1.2.0", branch: "release/v1.2.0", tag: "v1.2.0" },
  { input: "v1.2.3", resolved: "v1.2.3", branch: "release/v1.2.3", tag: "v1.2.3" },
];

const errorScenarios = [
  { scenario: "Invalid version string", behavior: "'abc' is not a valid version." },
  { scenario: "--commit SHA not found", behavior: "commit abc123 not found." },
  { scenario: "--branch does not exist", behavior: "branch develop does not exist." },
  { scenario: "Push to remote fails", behavior: "failed to push to remote: <detail>" },
  { scenario: "Version already released", behavior: "Version v1.2.3 is already released." },
  { scenario: "--bump + version argument", behavior: "--bump cannot be used with an explicit version argument." },
  { scenario: "--commit + --branch", behavior: "--commit and --branch are mutually exclusive." },
];

const ReleasePage = () => {
  return (
    <DocsLayout>
      <div className="space-y-10">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-mono font-bold text-foreground mb-3">Release Command</h1>
          <p className="text-muted-foreground leading-relaxed max-w-2xl">
            Automate Git release workflows: create branches, tags, push to remote, and track
            release history. Supports semver, partial versions, pre-release suffixes, draft mode,
            dry-run preview, and auto-increment.
          </p>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {features.map((f) => (
            <div key={f.title} className="p-5 rounded-lg border border-border bg-card hover:border-primary/40 transition-colors group">
              <div className="h-9 w-9 rounded-md bg-primary/10 flex items-center justify-center mb-3 group-hover:bg-primary/20 transition-colors">
                <f.icon className="h-4 w-4 text-primary" />
              </div>
              <h3 className="font-mono font-semibold text-foreground text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground leading-relaxed">{f.desc}</p>
            </div>
          ))}
        </div>

        {/* Commands */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Commands</h2>
          <div className="space-y-4">
            <div className="p-4 rounded-lg border border-border bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-1">gitmap release [version] <span className="text-muted-foreground font-normal text-sm">(alias: r)</span></h3>
              <p className="text-sm text-muted-foreground">Create a release branch, Git tag, and push to remote.</p>
            </div>
            <div className="p-4 rounded-lg border border-border bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-1">gitmap release-branch &lt;branch&gt; <span className="text-muted-foreground font-normal text-sm">(alias: rb)</span></h3>
              <p className="text-sm text-muted-foreground">Complete a release from an existing release/vX.Y.Z branch.</p>
            </div>
            <div className="p-4 rounded-lg border border-border bg-card">
              <h3 className="font-mono font-semibold text-foreground mb-1">gitmap release-pending <span className="text-muted-foreground font-normal text-sm">(alias: rp)</span></h3>
              <p className="text-sm text-muted-foreground">Release all release/v* branches that are missing tags.</p>
            </div>
          </div>
        </div>

        {/* Workflow Diagram */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Release Workflow</h2>
          <div className="p-5 rounded-lg border border-border bg-card font-mono text-sm space-y-1">
            <p className="text-muted-foreground mb-3">Steps executed by <span className="text-primary">gitmap release [version]</span>:</p>
            {[
              "1. Resolve version (CLI → --bump → version.json → error)",
              "2. Pad partial version to full semver",
              "3. Check .release/ and git tags for duplicates",
              "4. Resolve source commit (--commit / --branch / HEAD)",
              "5. Create branch release/vX.Y.Z from source",
              "6. Create git tag vX.Y.Z",
              "7. Push branch + tag to origin",
              "8. Collect --assets contents",
              "9. Write .release/vX.Y.Z.json",
              "10. Update .release/latest.json (if highest stable)",
            ].map((step) => (
              <p key={step} className="text-foreground/80 pl-2">{step}</p>
            ))}
          </div>
        </div>

        {/* Version Resolution */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Version Resolution</h2>
          <p className="text-sm text-muted-foreground mb-3">Version is resolved in priority order:</p>
          <div className="space-y-2 mb-6">
            {[
              { label: "CLI argument", example: "gitmap release v1.2.3" },
              { label: "--bump flag", example: "reads latest, increments" },
              { label: "version.json", example: '{ "version": "1.2.3" }' },
              { label: "Error", example: "no version source found" },
            ].map((item, i) => (
              <div key={item.label} className="flex items-start gap-3 text-sm">
                <span className="font-mono text-primary font-bold">{i + 1}.</span>
                <span className="font-mono font-semibold text-foreground">{item.label}</span>
                <span className="text-muted-foreground">— {item.example}</span>
              </div>
            ))}
          </div>

          <h3 className="text-base font-mono font-semibold text-foreground mb-3">Partial Version Padding</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Input</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Resolved</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Branch</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Tag</th>
                </tr>
              </thead>
              <tbody>
                {paddingExamples.map((row) => (
                  <tr key={row.input} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary">{row.input}</td>
                    <td className="py-2 px-3 font-mono text-foreground">{row.resolved}</td>
                    <td className="py-2 px-3 font-mono text-foreground/80">{row.branch}</td>
                    <td className="py-2 px-3 font-mono text-foreground/80">{row.tag}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Release Flags */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Release Flags</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Flag</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Description</th>
                </tr>
              </thead>
              <tbody>
                {releaseFlags.map((f) => (
                  <tr key={f.flag} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary whitespace-nowrap">{f.flag}</td>
                    <td className="py-2 px-3 text-muted-foreground">{f.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h3 className="text-base font-mono font-semibold text-foreground mt-6 mb-3">Release-Branch / Release-Pending Flags</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Flag</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Description</th>
                </tr>
              </thead>
              <tbody>
                {branchFlags.map((f) => (
                  <tr key={f.flag} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary whitespace-nowrap">{f.flag}</td>
                    <td className="py-2 px-3 text-muted-foreground">{f.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Auto-Increment */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Auto-Increment (--bump)</h2>
          <p className="text-sm text-muted-foreground mb-3">
            Reads the latest version from <span className="font-mono text-foreground">.release/latest.json</span> and increments.
            Falls back to scanning local Git tags.
          </p>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Current</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">--bump patch</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">--bump minor</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">--bump major</th>
                </tr>
              </thead>
              <tbody>
                {bumpExamples.map((row) => (
                  <tr key={row.current} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-foreground">{row.current}</td>
                    <td className="py-2 px-3 font-mono text-primary">{row.patch}</td>
                    <td className="py-2 px-3 font-mono text-primary">{row.minor}</td>
                    <td className="py-2 px-3 font-mono text-primary">{row.major}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Usage Examples */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Usage Examples</h2>
          <div className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground mb-2">Full semver release from HEAD</p>
              <CodeBlock code="gitmap release v1.2.3" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Partial version (padded to v1.0.0)</p>
              <CodeBlock code="gitmap release v1" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">With assets</p>
              <CodeBlock code="gitmap release v2.0.0 --assets ./dist" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">From specific commit or branch</p>
              <CodeBlock code={`gitmap release v1.2.3 --commit abc123def\ngitmap release v1.0.0 --branch develop`} />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Auto-increment</p>
              <CodeBlock code={`gitmap release --bump patch\ngitmap release --bump minor --assets ./bin`} />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Draft and dry-run</p>
              <CodeBlock code={`gitmap release v3.0.0-rc.1 --draft\ngitmap release v1.0.0 --dry-run`} />
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-2">Release from existing branch / pending</p>
              <CodeBlock code={`gitmap release-branch release/v1.2.0\ngitmap release-pending\ngitmap release-pending --dry-run`} />
            </div>
          </div>
        </div>

        {/* Dry-Run Output */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Dry-Run Preview</h2>
          <div className="bg-card border border-border rounded-lg p-5 font-mono text-sm space-y-1">
            <p className="text-muted-foreground">{`$ gitmap release v1.2.3 --dry-run`}</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Create branch release/v1.2.3 from main</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Create tag v1.2.3</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Push branch and tag to origin</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Write metadata to .release/v1.2.3.json</p>
            <p className="text-foreground/80">&nbsp;&nbsp;[dry-run] Mark v1.2.3 as latest</p>
          </div>
        </div>

        {/* Error Scenarios */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Error Scenarios</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Scenario</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Error Message</th>
                </tr>
              </thead>
              <tbody>
                {errorScenarios.map((row) => (
                  <tr key={row.scenario} className="border-b border-border/50">
                    <td className="py-2 px-3 text-foreground">{row.scenario}</td>
                    <td className="py-2 px-3 font-mono text-destructive/80">{row.behavior}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Package Layout */}
        <div>
          <h2 className="text-xl font-mono font-bold text-foreground mb-4">Package Layout</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">File</th>
                  <th className="text-left py-2 px-3 font-mono text-muted-foreground font-medium">Responsibility</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { file: "release/semver.go", desc: "Version parsing, padding, comparison, bumping" },
                  { file: "release/metadata.go", desc: "Read/write .release/*.json, latest.json, version.json" },
                  { file: "release/gitops.go", desc: "Branch, tag, push, checkout Git operations" },
                  { file: "release/github.go", desc: "Asset collection, changelog/readme detection" },
                  { file: "release/workflow.go", desc: "Orchestration: Execute(), ExecuteFromBranch()" },
                ].map((row) => (
                  <tr key={row.file} className="border-b border-border/50">
                    <td className="py-2 px-3 font-mono text-primary">{row.file}</td>
                    <td className="py-2 px-3 text-muted-foreground">{row.desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </DocsLayout>
  );
};

export default ReleasePage;
