# cadops doctor

Runs repository health checks for Git, Git LFS, config, attributes, metadata manifest state, preview manifest state, CAD coverage, and remotes.

- Uses `PASS`, `WARN`, and `FAIL` levels.
- Missing metadata and preview manifests are warnings, not hard failures.
