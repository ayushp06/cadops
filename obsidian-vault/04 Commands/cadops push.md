# cadops push

Runs light CAD-aware pre-push checks before delegating to `git push`.

- Warns when unstaged or uncommitted CAD changes are present locally.
- Warns when tracked CAD files are missing matching Git LFS rules in `.gitattributes`.
- Fails early when no git remote is configured.
- Keeps preflight assessment separate from Git command execution.
