# cadops snapshot

Creates a timestamped Git snapshot commit for changed CAD files.

- Stages changed CAD files before commit creation.
- Uses a message like `snapshot: 2026-04-14 15:42`.
- Fails clearly when there are no relevant CAD changes.
- Does not include non-CAD files unless future configuration explicitly allows it.
- Does not implement smart commit grouping yet.
