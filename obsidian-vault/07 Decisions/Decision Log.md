# Decision Log

- `cadops watch` is implemented as a dedicated internal watcher package rather than embedding filesystem logic in the Cobra command.
- The first watch implementation polls the repository tree recursively and ignores `.git`, which keeps the code portable and straightforward to extend.
- Auto-staging is supported when `auto_stage: true`, but auto-commit and preview generation remain deferred.
- `cadops snapshot` stages relevant CAD files first, then creates a timestamped commit scoped to those CAD paths only.
- `cadops config` keeps key lookup in the config package and leaves terminal formatting in the CLI layer.
- `cadops push` and `cadops pull` use a separate collaboration preflight package so warning logic stays testable without invoking Git.
- `cadops history` uses constrained `git log` output plus a separate parser/formatter package so commit rendering remains testable without embedding Git parsing in the Cobra command.
- `cadops history` reads `.cadops/metadata/manifest.json` directly from commit history when present and compares it only to the first parent manifest, which keeps historical metadata output honest without inventing semantic CAD history.
- `cadops diff` keeps Git status retrieval in the command layer, but moves metadata lookup, previous-versus-current comparison, and compact terminal formatting into `internal/diff` so the richer CAD output stays testable without semantic CAD analysis.
