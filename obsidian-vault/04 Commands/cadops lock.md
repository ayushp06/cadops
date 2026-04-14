# cadops lock

Wraps `git lfs lock` for a target repository file.

- Validates that the target file exists before invoking Git LFS.
- Warns when locking is recommended for the file type but Git LFS is missing or `.gitattributes` is not tracking that extension correctly.
- Keeps lock execution separate from user-facing messaging so the command surface stays thin.
