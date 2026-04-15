# cadops scan

Audits the current repository for CAD assets using configured extensions from `.cadops.yaml`.

- Summarizes total CAD files, counts by CAD type, locking recommendations, and Git LFS expectations.
- Warns when CAD file types that are expected to use Git LFS do not have matching canonical `.gitattributes` rules.
- Shows the largest CAD files and top directories containing CAD files.
- Uses `.cadops/metadata/manifest.json` when available, but falls back to a live repository scan when metadata is absent.
- Keeps geometry parsing, semantic CAD analysis, and preview generation out of scope.
