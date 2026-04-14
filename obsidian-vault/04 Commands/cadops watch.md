# cadops watch

Watches the current repository recursively for configured CAD file changes from `.cadops.yaml`.

- Prints concise change lines for created, modified, and removed CAD files.
- Debounces rapid repeated events for the same file.
- Stages changed CAD files automatically when `auto_stage: true`.
- Does not auto-commit or generate previews yet.
