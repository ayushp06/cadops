# Changelog

## Unreleased

- Added metadata-aware diff and history enrichment from `.cadops/metadata/manifest.json`.
- Added V1 preview record pipeline under `.cadops/previews/manifest.json`.
- Added `cadops preview generate`, `cadops preview list`, and `cadops preview show <file>`.
- Integrated preview status into scan, diff, history, and snapshot workflows.
- Added stale and missing metadata/preview warnings for status and commit flows.
- Expanded doctor checks for metadata and preview manifest state.
