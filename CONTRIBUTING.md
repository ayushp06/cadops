# Contributing

CadOps is a Go CLI for CAD-aware Git and Git LFS workflows. Contributions should keep command handlers thin and put testable behavior in `internal/` packages.

## Development

Prerequisites:

- Go 1.26+
- Git
- Git LFS

Common checks:

```bash
make fmt
make test
make vet
make build
```

## Guidelines

- Do not add proprietary CAD parser or renderer dependencies.
- Do not claim semantic CAD understanding unless the implementation can prove it from stored data.
- Keep Git execution behind `internal/gitx`.
- Keep metadata and preview behavior honest: filesystem metadata, hashes, Git state, and real preview artifacts only.
- Add focused tests for new pure logic and command formatting.
