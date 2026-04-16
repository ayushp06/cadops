# cadops diff

Shows a concise Git-backed summary of current repository changes, grouped into CAD and non-CAD files using configured extensions from `.cadops.yaml`.

- Uses Git working tree state rather than semantic CAD analysis.
- Shows readable change markers such as modified, added, deleted, and renamed.
- Enriches changed CAD files with stored metadata when `.cadops/metadata/manifest.json` is available.
- Uses the current manifest as the baseline and compares it to freshly derived metadata for the current working file when that file still exists locally.
- Surfaces compact context only: CAD type, lock recommendation, Git LFS expectation, checksum change, and file size delta when comparison is clean.
- Keeps previews and geometry-aware diffing out of scope.
