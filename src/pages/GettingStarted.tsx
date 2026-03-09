import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

const GettingStartedPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-mono font-bold mb-2">Getting Started</h1>
      <p className="text-muted-foreground mb-8">
        Get up and running with gitmap in under 5 minutes.
      </p>

      <section className="space-y-8">
        <div>
          <h2 className="text-xl font-mono font-semibold mb-3 text-foreground">1. Install gitmap</h2>
          <p className="text-muted-foreground mb-3">
            Build from source using Go 1.21+:
          </p>
          <CodeBlock
            code={`go install github.com/user/gitmap@latest`}
            title="Terminal"
          />
          <p className="text-sm text-muted-foreground mt-2">
            Or clone the repo and build with the platform-appropriate script:
          </p>
          <CodeBlock code={`# Windows (PowerShell)\ngit clone https://github.com/user/gitmap.git\ncd gitmap\n./run.ps1`} title="PowerShell" />
          <CodeBlock code={`# Linux / macOS (Bash)\ngit clone https://github.com/user/gitmap.git\ncd gitmap\nchmod +x run.sh\n./run.sh`} title="Bash" />
          <CodeBlock code={`# Or use Make (requires run.sh)\ncd gitmap\nmake build`} title="Makefile" />

          <div className="mt-4 p-4 rounded-lg border border-border bg-muted/30">
            <h3 className="text-sm font-mono font-semibold text-foreground mb-2">Build script flags</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-1 text-sm text-muted-foreground">
              <div><code className="font-mono text-primary">run.ps1 -NoPull</code> / <code className="font-mono text-primary">run.sh -n</code> — skip git pull</div>
              <div><code className="font-mono text-primary">run.ps1 -NoDeploy</code> / <code className="font-mono text-primary">run.sh -d</code> — build only</div>
              <div><code className="font-mono text-primary">run.ps1 -Update</code> / <code className="font-mono text-primary">run.sh -u</code> — full update pipeline</div>
              <div><code className="font-mono text-primary">run.ps1 -R list</code> / <code className="font-mono text-primary">run.sh -r list</code> — build &amp; run</div>
              <div><code className="font-mono text-primary">run.sh -t</code> / <code className="font-mono text-primary">make test</code> — run tests</div>
            </div>
          </div>
        </div>

        <div>
          <h2 className="text-xl font-mono font-semibold mb-3 text-foreground">2. Run your first scan</h2>
          <p className="text-muted-foreground mb-3">
            Point gitmap at a directory containing Git repositories:
          </p>
          <CodeBlock code={`gitmap scan ~/projects`} title="Terminal" />
          <p className="text-sm text-muted-foreground mt-2">
            This generates <code className="font-mono text-primary">gitmap-output/</code> containing CSV, JSON,
            folder structure, and clone scripts.
          </p>
        </div>

        <div>
          <h2 className="text-xl font-mono font-semibold mb-3 text-foreground">3. Clone on another machine</h2>
          <p className="text-muted-foreground mb-3">
            Copy the output files and restore the exact folder structure:
          </p>
          <CodeBlock
            code={`gitmap clone json --target-dir ./projects`}
            title="Terminal"
          />
          <p className="text-sm text-muted-foreground mt-2">
            Shorthands <code className="font-mono text-primary">json</code>,{" "}
            <code className="font-mono text-primary">csv</code>, and{" "}
            <code className="font-mono text-primary">text</code> auto-resolve to the default output files.
          </p>
        </div>

        <div>
          <h2 className="text-xl font-mono font-semibold mb-3 text-foreground">4. Set up shell navigation</h2>
          <p className="text-muted-foreground mb-3">
            Add a wrapper function to quickly cd into tracked repos:
          </p>
          <CodeBlock
            code={`# PowerShell ($PROFILE)\nfunction gcd { Set-Location (gitmap cd $args) }\n\n# Bash/Zsh (~/.bashrc or ~/.zshrc)\ngcd() { cd "$(gitmap cd "$@")" ; }`}
            title="Shell Profile"
          />
          <p className="text-sm text-muted-foreground mt-2">
            Then use <code className="font-mono text-primary">gcd myrepo</code> to jump to any tracked repo.
          </p>
        </div>

        <div>
          <h2 className="text-xl font-mono font-semibold mb-3 text-foreground">5. Monitor your repos</h2>
          <p className="text-muted-foreground mb-3">
            Start a live dashboard to watch all tracked repos:
          </p>
          <CodeBlock code={`gitmap watch --interval 15`} title="Terminal" />
        </div>
      </section>
    </DocsLayout>
  );
};

export default GettingStartedPage;
