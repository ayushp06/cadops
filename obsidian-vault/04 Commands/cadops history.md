# cadops history

Shows recent commit history in a CAD-aware terminal format.

- Prints short commit hash, commit date, and commit message.
- Lists changed CAD files for each commit when present.
- Defaults to a recent commit window and supports a simple `--limit` flag.
- Uses Git log output rather than implementing custom history storage or semantic analysis.
- Reads `.cadops/metadata/manifest.json` from each commit when available and uses it for compact CAD file annotations.
- Compares the commit manifest to the first parent manifest when both exist, which allows simple checksum-changed and size-delta indicators.
- Falls back cleanly to the standard CAD file list with `metadata unavailable` when commit-scoped metadata is absent.
