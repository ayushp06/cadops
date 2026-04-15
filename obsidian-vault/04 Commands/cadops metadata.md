# cadops metadata

Generates and inspects stored filesystem-level metadata for CAD files that match configured extensions from `.cadops.yaml`.

- `cadops metadata generate` writes a single JSON manifest to `.cadops/metadata/manifest.json`.
- Each record includes relative path, CAD type, extension, file size, modified time, SHA-256, expected Git LFS usage, and locking recommendation.
- `cadops metadata show <file>` prints a terminal-friendly view of one stored record and fails clearly when the file is missing, outside the configured CAD set, or absent from the metadata store.
- Geometry parsing, semantic CAD analysis, and preview generation remain out of scope.
