import { useState } from "react";
import {
  FolderGit2,
  Search,
  Filter,
  Code2,
  FileCode,
  Braces,
  Cpu,
  Hash,
  MapPin,
  FileText,
  ChevronDown,
  ChevronRight,
} from "lucide-react";
import DocsLayout from "@/components/docs/DocsLayout";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

const PROJECT_TYPES = {
  go: { label: "Go", color: "bg-cyan-500/15 text-cyan-700 dark:text-cyan-400 border-cyan-500/30", icon: Code2 },
  node: { label: "Node.js", color: "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30", icon: Braces },
  react: { label: "React", color: "bg-sky-500/15 text-sky-700 dark:text-sky-400 border-sky-500/30", icon: Braces },
  cpp: { label: "C++", color: "bg-violet-500/15 text-violet-700 dark:text-violet-400 border-violet-500/30", icon: Cpu },
  csharp: { label: "C#", color: "bg-purple-500/15 text-purple-700 dark:text-purple-400 border-purple-500/30", icon: Hash },
};

type ProjectType = keyof typeof PROJECT_TYPES;

interface DetectedProject {
  id: string;
  repoName: string;
  projectType: ProjectType;
  projectName: string;
  absolutePath: string;
  repoPath: string;
  relativePath: string;
  primaryIndicator: string;
  detectedAt: string;
  goMetadata?: {
    moduleName: string;
    goVersion: string;
    runnables: { name: string; relativePath: string }[];
  };
  csharpMetadata?: {
    slnName: string;
    sdkVersion: string;
    projectFiles: { fileName: string; targetFramework: string; outputType: string }[];
  };
}

const SAMPLE_PROJECTS: DetectedProject[] = [
  {
    id: "1", repoName: "my-api", projectType: "go", projectName: "github.com/user/my-api",
    absolutePath: "/home/user/repos/my-api", repoPath: "/home/user/repos/my-api",
    relativePath: ".", primaryIndicator: "go.mod", detectedAt: "2026-03-11T09:54:00Z",
    goMetadata: { moduleName: "github.com/user/my-api", goVersion: "1.22", runnables: [
      { name: "server", relativePath: "cmd/server/main.go" },
      { name: "worker", relativePath: "cmd/worker/main.go" },
    ]},
  },
  {
    id: "2", repoName: "my-api", projectType: "react", projectName: "admin-dashboard",
    absolutePath: "/home/user/repos/my-api/web", repoPath: "/home/user/repos/my-api",
    relativePath: "web", primaryIndicator: "package.json", detectedAt: "2026-03-11T09:54:00Z",
  },
  {
    id: "3", repoName: "infra-tools", projectType: "go", projectName: "github.com/user/infra-tools",
    absolutePath: "/home/user/repos/infra-tools", repoPath: "/home/user/repos/infra-tools",
    relativePath: ".", primaryIndicator: "go.mod", detectedAt: "2026-03-11T09:55:00Z",
    goMetadata: { moduleName: "github.com/user/infra-tools", goVersion: "1.23", runnables: [
      { name: "infra-tools", relativePath: "main.go" },
    ]},
  },
  {
    id: "4", repoName: "web-platform", projectType: "react", projectName: "@platform/frontend",
    absolutePath: "/home/user/repos/web-platform", repoPath: "/home/user/repos/web-platform",
    relativePath: ".", primaryIndicator: "package.json", detectedAt: "2026-03-11T09:56:00Z",
  },
  {
    id: "5", repoName: "web-platform", projectType: "node", projectName: "@platform/api",
    absolutePath: "/home/user/repos/web-platform/api", repoPath: "/home/user/repos/web-platform",
    relativePath: "api", primaryIndicator: "package.json", detectedAt: "2026-03-11T09:56:00Z",
  },
  {
    id: "6", repoName: "signal-engine", projectType: "cpp", projectName: "signal-engine",
    absolutePath: "/home/user/repos/signal-engine", repoPath: "/home/user/repos/signal-engine",
    relativePath: ".", primaryIndicator: "CMakeLists.txt", detectedAt: "2026-03-11T09:57:00Z",
  },
  {
    id: "7", repoName: "enterprise-app", projectType: "csharp", projectName: "EnterpriseApp",
    absolutePath: "/home/user/repos/enterprise-app", repoPath: "/home/user/repos/enterprise-app",
    relativePath: ".", primaryIndicator: "EnterpriseApp.sln", detectedAt: "2026-03-11T09:58:00Z",
    csharpMetadata: { slnName: "EnterpriseApp.sln", sdkVersion: "8.0.100", projectFiles: [
      { fileName: "EnterpriseApp.Api.csproj", targetFramework: "net8.0", outputType: "Exe" },
      { fileName: "EnterpriseApp.Core.csproj", targetFramework: "net8.0", outputType: "Library" },
    ]},
  },
];

const TypeBadge = ({ type }: { type: ProjectType }) => {
  const config = PROJECT_TYPES[type];
  const Icon = config.icon;

  return (
    <Badge variant="outline" className={`${config.color} font-mono text-xs gap-1 border`}>
      <Icon className="h-3 w-3" />
      {config.label}
    </Badge>
  );
};

const ProjectCard = ({ project }: { project: DetectedProject }) => {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="border border-border rounded-lg p-4 hover:border-primary/40 transition-colors bg-card">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 mb-1">
            <TypeBadge type={project.projectType} />
            <span className="font-mono text-sm font-semibold text-foreground truncate">
              {project.projectName}
            </span>
          </div>
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-2">
            <MapPin className="h-3 w-3 shrink-0" />
            <span className="font-mono truncate">{project.relativePath === "." ? "(root)" : project.relativePath}</span>
          </div>
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-1">
            <FileText className="h-3 w-3 shrink-0" />
            <span className="font-mono">{project.primaryIndicator}</span>
          </div>
        </div>
        {(project.goMetadata || project.csharpMetadata) && (
          <button
            onClick={() => setExpanded(!expanded)}
            className="text-muted-foreground hover:text-foreground transition-colors p-1"
          >
            {expanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
          </button>
        )}
      </div>

      {expanded && project.goMetadata && (
        <div className="mt-3 pt-3 border-t border-border space-y-2">
          <div className="flex gap-4 text-xs text-muted-foreground">
            <span>Module: <span className="font-mono text-foreground">{project.goMetadata.moduleName}</span></span>
            <span>Go: <span className="font-mono text-foreground">{project.goMetadata.goVersion}</span></span>
          </div>
          {project.goMetadata.runnables.length > 0 && (
            <div>
              <span className="text-xs text-muted-foreground">Runnables:</span>
              <div className="flex flex-wrap gap-1.5 mt-1">
                {project.goMetadata.runnables.map((r) => (
                  <span key={r.name} className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-muted text-xs font-mono">
                    <FileCode className="h-3 w-3 text-primary" />
                    {r.name}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {expanded && project.csharpMetadata && (
        <div className="mt-3 pt-3 border-t border-border space-y-2">
          <div className="flex gap-4 text-xs text-muted-foreground">
            <span>Solution: <span className="font-mono text-foreground">{project.csharpMetadata.slnName}</span></span>
            <span>SDK: <span className="font-mono text-foreground">{project.csharpMetadata.sdkVersion}</span></span>
          </div>
          {project.csharpMetadata.projectFiles.length > 0 && (
            <div>
              <span className="text-xs text-muted-foreground">Project files:</span>
              <div className="space-y-1 mt-1">
                {project.csharpMetadata.projectFiles.map((f) => (
                  <div key={f.fileName} className="flex items-center gap-2 text-xs font-mono">
                    <FileCode className="h-3 w-3 text-primary shrink-0" />
                    <span className="text-foreground">{f.fileName}</span>
                    <Badge variant="outline" className="text-[10px] px-1.5 py-0">{f.targetFramework}</Badge>
                    <span className="text-muted-foreground">{f.outputType}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

const ProjectsPage = () => {
  const [search, setSearch] = useState("");
  const [activeFilter, setActiveFilter] = useState<ProjectType | "all">("all");

  const filtered = SAMPLE_PROJECTS.filter((p) => {
    if (activeFilter !== "all" && p.projectType !== activeFilter) {
      return false;
    }
    if (search.length > 0) {
      const q = search.toLowerCase();

      return (
        p.projectName.toLowerCase().includes(q) ||
        p.repoName.toLowerCase().includes(q) ||
        p.absolutePath.toLowerCase().includes(q)
      );
    }

    return true;
  });

  const grouped = filtered.reduce<Record<string, DetectedProject[]>>((acc, p) => {
    acc[p.repoName] = acc[p.repoName] || [];
    acc[p.repoName].push(p);

    return acc;
  }, {});

  const typeCounts = SAMPLE_PROJECTS.reduce<Record<string, number>>((acc, p) => {
    acc[p.projectType] = (acc[p.projectType] || 0) + 1;

    return acc;
  }, {});

  return (
    <DocsLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-mono font-bold text-foreground flex items-center gap-3">
            <FolderGit2 className="h-8 w-8 text-primary" />
            Detected Projects
          </h1>
          <p className="text-muted-foreground mt-2">
            Projects discovered inside Git repositories during scan. Each repo can contain multiple project types.
          </p>
        </div>

        {/* Summary cards */}
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3">
          <Card
            className={`cursor-pointer transition-all ${activeFilter === "all" ? "ring-2 ring-primary" : "hover:border-primary/40"}`}
            onClick={() => setActiveFilter("all")}
          >
            <CardContent className="p-3 text-center">
              <div className="text-2xl font-mono font-bold text-foreground">{SAMPLE_PROJECTS.length}</div>
              <div className="text-xs text-muted-foreground">All</div>
            </CardContent>
          </Card>
          {(Object.entries(PROJECT_TYPES) as [ProjectType, typeof PROJECT_TYPES[ProjectType]][]).map(([key, config]) => (
            <Card
              key={key}
              className={`cursor-pointer transition-all ${activeFilter === key ? "ring-2 ring-primary" : "hover:border-primary/40"}`}
              onClick={() => setActiveFilter(key)}
            >
              <CardContent className="p-3 text-center">
                <div className="text-2xl font-mono font-bold text-foreground">{typeCounts[key] || 0}</div>
                <div className="text-xs text-muted-foreground">{config.label}</div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search by project name, repo, or path..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-10 font-mono text-sm"
          />
        </div>

        {/* Grouped project list */}
        <div className="space-y-6">
          {Object.entries(grouped).map(([repoName, projects]) => (
            <div key={repoName}>
              <div className="flex items-center gap-2 mb-3">
                <FolderGit2 className="h-4 w-4 text-primary" />
                <h2 className="font-mono font-semibold text-foreground">{repoName}</h2>
                <span className="text-xs text-muted-foreground font-mono">
                  {projects[0].repoPath}
                </span>
                <Badge variant="secondary" className="ml-auto text-xs">
                  {projects.length} project{projects.length > 1 ? "s" : ""}
                </Badge>
              </div>
              <div className="grid gap-3 md:grid-cols-2">
                {projects.map((p) => (
                  <ProjectCard key={p.id} project={p} />
                ))}
              </div>
            </div>
          ))}
        </div>

        {filtered.length === 0 && (
          <div className="text-center py-12 text-muted-foreground">
            <Filter className="h-8 w-8 mx-auto mb-3 opacity-50" />
            <p className="font-mono">No projects match your filters.</p>
          </div>
        )}
      </div>
    </DocsLayout>
  );
};

export default ProjectsPage;
