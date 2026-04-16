# System Overview

CadOps is a Go CLI built on Cobra with thin command handlers and small internal packages.

- `internal/config` loads `.cadops.yaml`.
- `internal/gitx` wraps Git and Git LFS commands.
- `internal/metadata` builds and stores filesystem-level CAD metadata in `.cadops/metadata/manifest.json`.
- `internal/scan` aggregates repo-level CAD audit summaries, LFS checks, and largest-file views.
- `internal/diff` isolates Git-backed diff classification, metadata-aware comparison, and terminal formatting.
- `internal/history` isolates constrained Git history parsing, commit-scoped metadata enrichment, and terminal formatting.
- `internal/snapshot` isolates snapshot message generation and CAD file selection.
- `internal/watch` handles recursive repository watching and isolates pure extension filtering and debounce logic for tests.
