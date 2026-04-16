# MVP Scope

The MVP includes `cadops init`, `cadops status`, `cadops diff`, `cadops doctor`, `cadops watch`, `cadops snapshot`, `cadops commit`, `cadops config`, `cadops push`, `cadops pull`, and `cadops history`.

`cadops watch` covers recursive repository watching, configured CAD extension filtering, concise status lines, and optional auto-staging only.

Auto-commit and preview generation are intentionally deferred.

`cadops snapshot` creates timestamped CAD-only commits. Smart grouping and inclusion of non-CAD files are intentionally deferred.

`cadops commit` wraps standard `git commit -m` with CAD-aware warnings about unstaged changes, LFS coverage, and recommended-lock files. Automatic staging, semantic message generation, and previews are intentionally deferred.

`cadops config` is read-only in the current phase and focuses on clear inspection of `.cadops.yaml`.

`cadops diff` adds a readable Git-backed working tree summary with CAD versus non-CAD grouping plus stored metadata-aware enrichment from `.cadops/metadata/manifest.json`, while still deferring semantic CAD diffing, previews, and geometry-aware analysis.

`cadops push` and `cadops pull` add lightweight collaboration guardrails around the underlying Git commands without attempting advanced merge or history analysis.

`cadops history` adds a readable recent commit view with CAD file filtering plus commit-scoped metadata-aware enrichment when manifests are present in Git history, but still defers semantic change detection, previews, and geometry-aware analysis.
