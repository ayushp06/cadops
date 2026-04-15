# cadops snapshot

Creates a timestamped Git snapshot commit for changed CAD files.

- Stages changed CAD files before commit creation.
- Regenerates the full `.cadops/metadata/manifest.json` before commit so CAD changes and metadata land in the same snapshot revision.
- Includes the metadata manifest in the snapshot commit when refresh succeeds and warns without blocking when metadata refresh fails.
- Uses a message like `snapshot: 2026-04-14 15:42`.
- Fails clearly when there are no relevant CAD changes.
- Does not include non-CAD files unless future configuration explicitly allows it.
- Does not implement smart commit grouping or preview generation yet.
