# cadops unlock

Wraps `git lfs unlock` for a target repository file.

- Validates that the target file exists before invoking Git LFS.
- Reuses the same lock-readiness warnings as `cadops lock` for recommended locking file types.
- Keeps unlock execution separate from user-facing messaging so the command surface stays thin.
