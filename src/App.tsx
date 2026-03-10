import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Index from "./pages/Index";
import Commands from "./pages/Commands";
import GoMod from "./pages/GoMod";
import GettingStarted from "./pages/GettingStarted";
import Config from "./pages/Config";
import Architecture from "./pages/Architecture";
import Watch from "./pages/Watch";
import MakefilePage from "./pages/Makefile";
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
          <Route path="/gomod" element={<GoMod />} />
          <Route path="/makefile" element={<MakefilePage />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
