# CadOps

CadOps is a CAD-aware command-line workflow layer over Git and Git LFS. The MVP provides safe repository initialization, CAD-aware status output, repository health checks, and repository watching for file-based CAD workflows.

## Commands

- `cadops init`
- `cadops status`
- `cadops doctor`
- `cadops watch`
- `cadops snapshot`

`cadops watch` monitors the current repository recursively, reacts only to CAD extensions configured in `.cadops.yaml`, prints concise change lines, and can auto-stage changed CAD files when `auto_stage: true`.

`cadops snapshot` stages changed CAD files and creates a timestamped snapshot commit like `snapshot: 2026-04-14 15:42`. It fails if there are no relevant CAD changes.

CadOps does not auto-commit from `watch`, implement smart commit grouping, or generate previews yet.

## Development

```bash
make test
```
