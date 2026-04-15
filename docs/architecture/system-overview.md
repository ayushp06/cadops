# System Overview

CadOps is a Go CLI built on Cobra.

- `cmd/cadops` holds the executable entrypoint.
- `internal/cli` contains thin command handlers.
- `internal/gitx` wraps Git and Git LFS command execution plus parsing helpers.
- `internal/cad` defines the supported CAD file registry.
- `internal/config` manages `.cadops.yaml`.
- `internal/commit` owns CAD-aware commit preflight checks.
- `internal/diff` owns Git-backed diff entry classification and grouping.
- `internal/files` owns recursive CAD file scanning and grouping.
- `internal/doctor` evaluates repository health checks.
- `internal/snapshot` owns CAD snapshot selection and commit message generation.
- `internal/watch` owns recursive repository watching, extension filtering, and event debouncing.

The implementation keeps parsing and merge behavior in small, testable functions so command handlers stay focused on orchestration and output.
