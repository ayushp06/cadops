# cadops commit

Wraps `git commit -m` with CAD-aware pre-commit checks.

- Requires `-m` in v1 to keep the flow explicit and non-interactive.
- Warns about unstaged changes that will be left out of the commit.
- Warns when changed CAD files are missing matching Git LFS rules.
- Warns when locking is enabled and recommended-lock CAD files do not have a local Git LFS lock.
- Does not auto-stage files, generate semantic commit messages, or build previews.
