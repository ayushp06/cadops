# Getting Started

## Prerequisites

- Go 1.26+
- Git
- Git LFS

## Common Tasks

```bash
make fmt
make test
make build
```

Run the CLI locally with:

```bash
go run ./cmd/cadops --help
```

To try the watcher in a repository initialized with `.cadops.yaml`:

```bash
go run ./cmd/cadops watch
```

The watcher only reacts to configured CAD extensions. If `auto_stage: true` is set in `.cadops.yaml`, matching changes are staged automatically.

To create a CAD-only snapshot commit:

```bash
go run ./cmd/cadops snapshot
```

`snapshot` refreshes `.cadops/metadata/manifest.json` and `.cadops/previews/manifest.json` before commit and includes them in the snapshot when generation succeeds.

To inspect repository configuration:

```bash
go run ./cmd/cadops config show
go run ./cmd/cadops config get tracked_extensions
```

To list CAD-relevant files in the current repository:

```bash
go run ./cmd/cadops files
```

To summarize current Git-backed repository changes:

```bash
go run ./cmd/cadops diff
```

When `.cadops/metadata/manifest.json` is available, `diff` uses that manifest as a stored baseline for changed CAD files and enriches the output with compact metadata context such as CAD type, locking recommendation, Git LFS expectation, explicit checksum changed yes/no, and file size delta when a clean previous-versus-current comparison is possible. When `.cadops/previews/manifest.json` exists, it also reports whether preview records are generated, unavailable, unsupported, stale, or missing. Missing metadata is reported as unavailable without blocking the standard diff summary.

To audit CAD assets and repository configuration risk:

```bash
go run ./cmd/cadops scan
```

To generate and inspect stored CAD metadata:

```bash
go run ./cmd/cadops metadata generate
go run ./cmd/cadops metadata show parts/gearbox.sldprt
```

To generate and inspect preview records:

```bash
go run ./cmd/cadops preview generate
go run ./cmd/cadops preview list
go run ./cmd/cadops preview show parts/gearbox.sldprt
```

Preview V1 does not render proprietary CAD geometry. It stores JSON records under `.cadops/previews/manifest.json` and references real same-name sidecar images only when they already exist.

To create a standard Git commit with CAD-aware preflight checks:

```bash
go run ./cmd/cadops commit -m "update bracket geometry"
```

To run guarded collaboration commands:

```bash
go run ./cmd/cadops push
go run ./cmd/cadops pull
```

`push` and `pull` keep Git execution simple, but surface CAD-aware warnings before delegating to the underlying Git command.

To view recent CAD-aware commit history:

```bash
go run ./cmd/cadops history
go run ./cmd/cadops history --limit 5
```

When snapshot or other commits include `.cadops/metadata/manifest.json`, `history` reads the manifest from each commit, compares it with the first parent manifest when available, and adds compact CAD file annotations such as type, stored size, checksum change, and size delta. When preview manifests are present in commit history, it also shows preview availability. If a commit does not carry usable metadata or preview records, the command degrades to the standard CAD file list with concise unavailable output.
