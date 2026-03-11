# Project Type Detection — Scan Integration

## Extended Scan Flow

```
 1. Parse flags
 2. Load config
 3. ScanDir (discover repos)
 4. BuildRecords (extract git metadata)
 5. DetectProjects (NEW — scan inside each repo)
 6. Write outputs (terminal, CSV, JSON, scripts)
 7. Write project JSON files (NEW)
 8. Save scan cache
 9. Upsert repos to DB
10. Upsert detected projects to DB (NEW)
11. Upsert Go metadata + runnables (NEW)
12. Upsert C# metadata + files + key files (NEW)
13. Cleanup stale projects (NEW)
14. Import releases
15. Add to desktop
16. Open folder
```

---

## Detection Flow per Repo

```
repo root
  │
  ├─ Walk directory tree (skip exclusion dirs)
  │
  ├─ For each directory:
  │   ├─ Check for go.mod         → classify as "go"
  │   │   └─ Parse go.mod/go.sum metadata
  │   │   └─ Scan cmd/ for runnables
  │   ├─ Check for package.json   → read contents
  │   │   ├─ Has react dep?       → classify as "react"
  │   │   └─ No react dep?        → classify as "node"
  │   ├─ Check for CMakeLists.txt → classify as "cpp"
  │   ├─ Check for *.vcxproj      → classify as "cpp"
  │   ├─ Check for meson.build    → classify as "cpp"
  │   ├─ Check for *.csproj       → classify as "csharp"
  │   │   └─ Parse .csproj XML for framework/output/SDK
  │   ├─ Check for *.sln          → classify as "csharp"
  │   │   └─ Scan for .csproj files + key files
  │   └─ No match                 → continue
  │
  └─ Collect all DetectedProject + metadata records
```

---

## Error Handling

| Scenario                        | Behavior                              |
|---------------------------------|---------------------------------------|
| Manifest file unreadable        | Skip project, log warning to stderr   |
| JSON/XML parse failure          | Skip project, log warning to stderr   |
| DB upsert failure               | Log error, continue with next project |
| JSON file write failure         | Log error, continue scan              |
| No projects of a type found     | Skip JSON file creation for that type |
| No database for query command   | Print message, exit 1                 |
| Go cmd/ scan failure            | Log warning, record project without runnables |
| C# .csproj parse failure        | Log warning, record project without file metadata |
