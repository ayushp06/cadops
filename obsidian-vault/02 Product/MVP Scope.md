# MVP Scope

The MVP includes `cadops init`, `cadops status`, `cadops doctor`, `cadops watch`, and `cadops snapshot`.

`cadops watch` covers recursive repository watching, configured CAD extension filtering, concise status lines, and optional auto-staging only.

Auto-commit and preview generation are intentionally deferred.

`cadops snapshot` creates timestamped CAD-only commits. Smart grouping and inclusion of non-CAD files are intentionally deferred.
