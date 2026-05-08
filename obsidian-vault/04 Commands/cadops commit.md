# cadops commit

Wraps `git commit -m` with CAD-aware pre-commit checks.

- Requires `-m` in v1 to keep the flow explicit and non-interactive.
- Warns about unstaged changes that will be left out of the commit.
- Warns when changed CAD files are missing matching Git LFS rules.
- Warns when locking is enabled and recommended-lock CAD files do not have a local Git LFS lock.
- Warns when metadata or preview records are missing or stale for changed CAD files.
- Does not auto-stage files or generate semantic commit messages.
