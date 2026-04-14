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
