# MVP Scope

The MVP includes `cadops init`, `cadops status`, `cadops doctor`, `cadops watch`, `cadops snapshot`, `cadops config`, `cadops push`, `cadops pull`, and `cadops history`.

`cadops watch` covers recursive repository watching, configured CAD extension filtering, concise status lines, and optional auto-staging only.

Auto-commit and preview generation are intentionally deferred.

`cadops snapshot` creates timestamped CAD-only commits. Smart grouping and inclusion of non-CAD files are intentionally deferred.

`cadops config` is read-only in the current phase and focuses on clear inspection of `.cadops.yaml`.

`cadops push` and `cadops pull` add lightweight collaboration guardrails around the underlying Git commands without attempting advanced merge, diff, or history analysis yet.

`cadops history` adds a readable recent commit view with CAD file filtering, but still defers semantic change detection and advanced diff presentation.
