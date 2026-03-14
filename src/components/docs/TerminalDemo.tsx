import { useState, useEffect, useRef } from "react";
import { Play, RotateCcw } from "lucide-react";

interface TerminalLine {
  text: string;
  type?: "input" | "output" | "header" | "accent";
  delay?: number;
}

interface TerminalDemoProps {
  title: string;
  lines: TerminalLine[];
  autoPlay?: boolean;
}

const TerminalDemo = ({ title, lines, autoPlay = false }: TerminalDemoProps) => {
  const [visibleLines, setVisibleLines] = useState<number>(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const play = () => {
    setVisibleLines(0);
    setIsPlaying(true);
  };

  const reset = () => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    setVisibleLines(0);
    setIsPlaying(false);
  };

  useEffect(() => {
    if (autoPlay) play();
    return () => { if (timeoutRef.current) clearTimeout(timeoutRef.current); };
  }, []);

  useEffect(() => {
    if (!isPlaying || visibleLines >= lines.length) {
      if (visibleLines >= lines.length) setIsPlaying(false);
      return;
    }

    const delay = lines[visibleLines]?.delay ?? (lines[visibleLines]?.type === "input" ? 600 : 120);
    timeoutRef.current = setTimeout(() => {
      setVisibleLines((v) => v + 1);
    }, delay);

    return () => { if (timeoutRef.current) clearTimeout(timeoutRef.current); };
  }, [isPlaying, visibleLines, lines]);

  const colorFor = (type?: string) => {
    switch (type) {
      case "input": return "text-[hsl(var(--terminal-foreground))]";
      case "header": return "text-primary font-bold";
      case "accent": return "text-primary";
      default: return "text-[hsl(var(--foreground))]/70";
    }
  };

  return (
    <div className="rounded-lg border border-border overflow-hidden">
      <div className="flex items-center justify-between px-4 py-2 bg-muted/30 border-b border-border">
        <div className="flex items-center gap-2">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-destructive/60" />
            <div className="w-3 h-3 rounded-full bg-[hsl(45,80%,50%)]/60" />
            <div className="w-3 h-3 rounded-full bg-primary/60" />
          </div>
          <span className="text-xs font-mono text-muted-foreground ml-2">{title}</span>
        </div>
        <div className="flex gap-1">
          <button onClick={play} className="p-1 rounded hover:bg-muted/50 text-muted-foreground hover:text-foreground transition-colors" title="Play">
            <Play className="h-3.5 w-3.5" />
          </button>
          <button onClick={reset} className="p-1 rounded hover:bg-muted/50 text-muted-foreground hover:text-foreground transition-colors" title="Reset">
            <RotateCcw className="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
      <div className="bg-[hsl(var(--terminal))] p-4 font-mono text-xs leading-relaxed min-h-[120px] max-h-[320px] overflow-y-auto">
        {lines.slice(0, visibleLines).map((line, i) => (
          <div key={i} className={colorFor(line.type)}>
            {line.type === "input" && <span className="text-primary mr-1">$</span>}
            {line.text}
          </div>
        ))}
        {isPlaying && (
          <span className="inline-block w-2 h-4 bg-primary/80 animate-pulse" />
        )}
      </div>
    </div>
  );
};

export default TerminalDemo;
