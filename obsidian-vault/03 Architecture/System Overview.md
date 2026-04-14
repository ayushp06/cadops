# System Overview

CadOps is a Go CLI built on Cobra with thin command handlers and small internal packages.

- `internal/config` loads `.cadops.yaml`.
- `internal/gitx` wraps Git and Git LFS commands.
- `internal/snapshot` isolates snapshot message generation and CAD file selection.
- `internal/watch` handles recursive repository watching and isolates pure extension filtering and debounce logic for tests.
