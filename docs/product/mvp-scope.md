# MVP Scope

CadOps MVP intentionally covers the core repository, collaboration, and inspection commands:

- `cadops init`
- `cadops status`
- `cadops doctor`
- `cadops watch`
- `cadops snapshot`
- `cadops config`
- `cadops push`
- `cadops pull`
- `cadops history`

The product goal is safe Git and Git LFS setup for CAD-heavy repositories without introducing a new version control model.

`cadops watch` is limited to change detection, concise status output, and optional auto-staging. Auto-commit and preview generation remain out of scope.

`cadops snapshot` is limited to CAD-file snapshots with an auto-generated timestamped commit message. Including other modified files and smart grouping remain out of scope.

`cadops history` is limited to recent Git-backed commit output with CAD file filtering. Advanced diffing, previews, and semantic change detection remain out of scope.
