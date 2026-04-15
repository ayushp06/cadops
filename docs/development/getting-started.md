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
