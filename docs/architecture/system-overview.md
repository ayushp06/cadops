# System Overview

CadOps is a Go CLI built on Cobra.

- `cmd/cadops` holds the executable entrypoint.
- `internal/cli` contains thin command handlers.
- `internal/gitx` wraps Git and Git LFS command execution plus parsing helpers.
- `internal/cad` defines the supported CAD file registry.
- `internal/config` manages `.cadops.yaml`.
- `internal/doctor` evaluates repository health checks.

The implementation keeps parsing and merge behavior in small, testable functions so command handlers stay focused on orchestration and output.
