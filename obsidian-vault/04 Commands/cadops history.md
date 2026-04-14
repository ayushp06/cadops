# cadops history

Shows recent commit history in a CAD-aware terminal format.

- Prints short commit hash, commit date, and commit message.
- Lists changed CAD files for each commit when present.
- Defaults to a recent commit window and supports a simple `--limit` flag.
- Uses Git log output rather than implementing custom history storage or semantic analysis.
