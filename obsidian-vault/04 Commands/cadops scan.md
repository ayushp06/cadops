# cadops scan

Audits the current repository for CAD assets using configured extensions from `.cadops.yaml`.

- Summarizes total CAD files, counts by CAD type, locking recommendations, and Git LFS expectations.
- Warns when CAD file types that are expected to use Git LFS do not have matching canonical `.gitattributes` rules.
- Shows the largest CAD files and top directories containing CAD files.
- Uses `.cadops/metadata/manifest.json` when available, but falls back to a live repository scan when metadata is absent.
- Reports preview coverage from `.cadops/previews/manifest.json`, including generated, unavailable, unsupported, stale, and missing records.
- Keeps geometry parsing, semantic CAD analysis, and preview generation out of scope.
