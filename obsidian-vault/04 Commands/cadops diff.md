# cadops diff

Shows a concise Git-backed summary of current repository changes, grouped into CAD and non-CAD files using configured extensions from `.cadops.yaml`.

- Uses Git working tree state rather than semantic CAD analysis.
- Shows readable change markers such as modified, added, deleted, and renamed.
- Keeps previews, metadata comparison, and geometry-aware diffing out of scope.
