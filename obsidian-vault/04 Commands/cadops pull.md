# cadops pull

Runs light CAD-aware pre-pull checks before delegating to `git pull`.

- Warns when the working tree is dirty.
- Warns when local CAD files are modified.
- Requires Git LFS to be available before pull proceeds.
- Keeps preflight assessment separate from Git command execution.
