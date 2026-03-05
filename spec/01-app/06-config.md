# Config

## File Location

Default: `./data/config.json`
Override: `--config <path>` flag.

## Schema

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": "./gitmap-output",
  "excludeDirs": [".cache", "node_modules", "vendor", ".venv"],
  "notes": ""
}
```

## Fields

| Field         | Type     | Default            | Description                    |
|---------------|----------|--------------------|--------------------------------|
| defaultMode   | string   | "https"            | "https" or "ssh"               |
| defaultOutput | string   | "terminal"         | "terminal", "csv", or "json"   |
| outputDir     | string   | "./gitmap-output"  | Where output files are written |
| excludeDirs   | []string | []                 | Directory names to skip        |
| notes         | string   | ""                 | Default note for all records   |

## Merge Rules

1. Load config file (if it exists).
2. Apply CLI flags on top — flags always win.
3. If config file is missing, use built-in defaults silently.
