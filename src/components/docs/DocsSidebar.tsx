import {
  Home,
  BookOpen,
  Rocket,
  Settings,
  Boxes,
  Monitor,
  Hammer,
  GitBranch,
  Tag,
  Sun,
  Moon,
  FolderGit2,
  Clock,
  BarChart3,
  Search,
  Terminal,
  FileText,
  Flag,
  PlayCircle,
} from "lucide-react";
import { NavLink } from "@/components/NavLink";
import { useLocation } from "react-router-dom";
import { useEffect, useState } from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarFooter,
  useSidebar,
} from "@/components/ui/sidebar";

const navItems = [
  { title: "Home", url: "/", icon: Home },
  { title: "Commands", url: "/commands", icon: BookOpen },
  { title: "Getting Started", url: "/getting-started", icon: Rocket },
  { title: "Configuration", url: "/config", icon: Settings },
  { title: "Architecture", url: "/architecture", icon: Boxes },
  { title: "Watch", url: "/watch", icon: Monitor },
  { title: "Release", url: "/release", icon: Tag },
  { title: "GoMod", url: "/gomod", icon: GitBranch },
  { title: "Projects", url: "/projects", icon: FolderGit2 },
  { title: "Makefile", url: "/makefile", icon: Hammer },
  { title: "History", url: "/history", icon: Clock },
  { title: "Stats", url: "/stats", icon: BarChart3 },
  { title: "Detection", url: "/project-detection", icon: Search },
  { title: "Generic CLI", url: "/generic-cli", icon: Terminal },
  { title: "Changelog", url: "/changelog", icon: FileText },
  { title: "Flags", url: "/flags", icon: Flag },
  { title: "Examples", url: "/examples", icon: PlayCircle },
];

export function DocsSidebar() {
  const { state } = useSidebar();
  const collapsed = state === "collapsed";
  const location = useLocation();
  const [dark, setDark] = useState(() => document.documentElement.classList.contains("dark"));

  useEffect(() => {
    document.documentElement.classList.toggle("dark", dark);
  }, [dark]);

  return (
    <Sidebar collapsible="icon">
      <SidebarContent>
        <SidebarGroup>
          {!collapsed && (
            <div className="px-3 py-4">
              <div className="flex items-center gap-2">
                <span className="font-mono font-bold text-lg text-primary">gitmap</span>
                <span className="text-xs font-mono text-muted-foreground">docs</span>
              </div>
            </div>
          )}
          {collapsed && (
            <div className="flex justify-center py-4">
              <span className="font-mono font-bold text-lg text-primary">g</span>
            </div>
          )}
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <NavLink
                      to={item.url}
                      end={item.url === "/"}
                      className="hover:bg-muted/50"
                      activeClassName="bg-muted text-primary font-medium"
                    >
                      <item.icon className="mr-2 h-4 w-4" />
                      {!collapsed && <span>{item.title}</span>}
                    </NavLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <button
          onClick={() => setDark(!dark)}
          className="flex items-center gap-2 px-3 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors w-full"
        >
          {dark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          {!collapsed && <span>{dark ? "Light mode" : "Dark mode"}</span>}
        </button>
      </SidebarFooter>
    </Sidebar>
  );
}
