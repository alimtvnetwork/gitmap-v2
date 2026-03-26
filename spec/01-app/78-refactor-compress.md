# Refactor: release/compress.go

## Problem
`compress.go` is 204 lines with two responsibilities: zip-based compression (CompressAssets orchestration, zip creation, verbose logging, dry-run descriptions) and tar.gz-based compression (tar.gz creation, tar entry writing).

## Target Layout

### compress.go (~135 lines) — Orchestration & Zip
Stays:
- `CompressAssets()`
- `logCompressedArchive()`
- `compressSingle()`
- `isWindowsBinary()`
- `createZip()`
- `addFileToZip()`
- `DescribeCompression()`

### compresstar.go (~72 lines) — Tar.gz
Moves:
- `createTarGz()`
- `addFileToTar()`

Imports: `archive/tar`, `compress/gzip`, `fmt`, `io`, `os`, `path/filepath`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `release`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.
