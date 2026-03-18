import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Index from "./pages/Index";
import Commands from "./pages/Commands";
import GoMod from "./pages/GoMod";
import Projects from "./pages/Projects";
import GettingStarted from "./pages/GettingStarted";
import Config from "./pages/Config";
import Architecture from "./pages/Architecture";
import Watch from "./pages/Watch";
import Release from "./pages/Release";
import MakefilePage from "./pages/Makefile";
import HistoryPage from "./pages/History";
import StatsPage from "./pages/Stats";
import ProjectDetectionPage from "./pages/ProjectDetection";
import GenericCLIPage from "./pages/GenericCLI";
import ChangelogPage from "./pages/Changelog";
import FlagReferencePage from "./pages/FlagReference";
import InteractiveExamplesPage from "./pages/InteractiveExamples";
import InteractiveTUIPage from "./pages/InteractiveTUI";
import BatchActionsPage from "./pages/BatchActions";
import ClearReleaseJSONPage from "./pages/ClearReleaseJSON";
import BookmarksPage from "./pages/Bookmarks";
import ExportPage from "./pages/Export";
import ImportPage from "./pages/Import";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Index />} />
          <Route path="/commands" element={<Commands />} />
          <Route path="/getting-started" element={<GettingStarted />} />
          <Route path="/config" element={<Config />} />
          <Route path="/architecture" element={<Architecture />} />
          <Route path="/watch" element={<Watch />} />
          <Route path="/release" element={<Release />} />
          <Route path="/gomod" element={<GoMod />} />
          <Route path="/projects" element={<Projects />} />
          <Route path="/makefile" element={<MakefilePage />} />
          <Route path="/history" element={<HistoryPage />} />
          <Route path="/stats" element={<StatsPage />} />
          <Route path="/project-detection" element={<ProjectDetectionPage />} />
          <Route path="/generic-cli" element={<GenericCLIPage />} />
          <Route path="/changelog" element={<ChangelogPage />} />
          <Route path="/flags" element={<FlagReferencePage />} />
          <Route path="/examples" element={<InteractiveExamplesPage />} />
          <Route path="/interactive" element={<InteractiveTUIPage />} />
          <Route path="/batch-actions" element={<BatchActionsPage />} />
          <Route path="/clear-release-json" element={<ClearReleaseJSONPage />} />
          <Route path="/bookmarks" element={<BookmarksPage />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
