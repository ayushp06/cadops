# MVP Scope

CadOps MVP intentionally covers the core repository, collaboration, and inspection commands:

- `cadops init`
- `cadops status`
- `cadops diff`
- `cadops doctor`
- `cadops watch`
- `cadops snapshot`
- `cadops commit`
- `cadops config`
- `cadops push`
- `cadops pull`
- `cadops history`

The product goal is safe Git and Git LFS setup for CAD-heavy repositories without introducing a new version control model.

`cadops watch` is limited to change detection, concise status output, and optional auto-staging. Auto-commit and preview generation remain out of scope.

`cadops snapshot` is limited to CAD-file snapshots with an auto-generated timestamped commit message. Including other modified files and smart grouping remain out of scope.

`cadops commit` is limited to standard `git commit -m` execution with CAD-aware pre-commit warnings. Semantic message generation, preview generation, automatic staging, and metadata pipelines remain out of scope.

`cadops diff` is limited to Git-backed working tree summaries grouped into CAD and non-CAD changes, with optional enrichment from stored filesystem metadata in `.cadops/metadata/manifest.json`. Semantic CAD diffing, previews, and geometry-aware analysis remain out of scope.

`cadops history` is limited to recent Git-backed commit output with CAD file filtering plus optional enrichment from commit-scoped metadata manifests when they exist in Git history. Semantic change detection, previews, and geometry-aware analysis remain out of scope.
