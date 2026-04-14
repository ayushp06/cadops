# Decision Log

- `cadops watch` is implemented as a dedicated internal watcher package rather than embedding filesystem logic in the Cobra command.
- The first watch implementation polls the repository tree recursively and ignores `.git`, which keeps the code portable and straightforward to extend.
- Auto-staging is supported when `auto_stage: true`, but auto-commit and preview generation remain deferred.
- `cadops snapshot` stages relevant CAD files first, then creates a timestamped commit scoped to those CAD paths only.
