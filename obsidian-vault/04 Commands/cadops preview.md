# cadops preview

Creates and inspects V1 preview records for configured CAD files.

- `cadops preview generate` writes `.cadops/previews/manifest.json`.
- `cadops preview list` shows all stored preview records and their current status.
- `cadops preview show <file>` shows one preview record and reports stale when the source hash no longer matches.
- V1 does not render proprietary CAD geometry.
- Same-name sidecar images such as `part.png`, `part.jpg`, or `part.jpeg` are referenced when they already exist.
- Files without renderable sidecar images get honest unavailable or unsupported placeholder records.
