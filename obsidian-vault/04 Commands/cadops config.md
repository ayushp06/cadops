# cadops config

Inspects the current repository `.cadops.yaml` without mutating it.

- `cadops config show` prints the supported keys in a concise terminal format.
- `cadops config get <key>` returns a single value for `version`, `tracked_extensions`, `auto_stage`, `require_lfs`, or `locking_enabled`.
- Configuration lookup stays separate from CLI formatting so the supported key behavior is directly testable.
