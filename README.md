# CadOps

CadOps is a CAD-aware command-line workflow layer over Git and Git LFS. The MVP provides safe repository initialization, CAD-aware status output, repository health checks, repository watching, snapshot commits, and Git LFS lock helpers for file-based CAD workflows.

CadOps requires both `git` and `git lfs` to be installed on the user machine.

## Installation

### Build From Source

Prerequisites:

- Go 1.26+
- Git
- Git LFS

Build the binary locally:

```bash
make build
```

This writes the binary to:

- `bin/cadops` on Linux and macOS
- `bin/cadops.exe` on Windows

To install it to your Go bin directory:

```bash
make install
```

Make sure your Go bin directory is on `PATH`.

### Install From GitHub Releases

1. Open the latest release on GitHub.
2. Download the archive for your platform:
   - `windows-amd64`
   - `windows-arm64`
   - `linux-amd64`
   - `linux-arm64`
   - `darwin-amd64`
   - `darwin-arm64`
3. Extract the archive.
4. Put `cadops` or `cadops.exe` somewhere on your `PATH`.
5. Ensure `git` and `git lfs` are installed on the machine.

Release archives contain the binary as `cadops` or `cadops.exe`.

### Verify Installation

Create or enter a Git repository, initialize CadOps if needed, then run:

```bash
cadops doctor
```

If CadOps, Git, Git LFS, and the repository setup are available, `cadops doctor` will report repository health checks instead of failing due to a missing command.

## Commands

- `cadops init`
- `cadops status`
- `cadops doctor`
- `cadops watch`
- `cadops snapshot`
- `cadops lock <file>`
- `cadops unlock <file>`

`cadops watch` monitors the current repository recursively, reacts only to CAD extensions configured in `.cadops.yaml`, prints concise change lines, and can auto-stage changed CAD files when `auto_stage: true`.

`cadops snapshot` stages changed CAD files and creates a timestamped snapshot commit like `snapshot: 2026-04-14 15:42`. It fails if there are no relevant CAD changes.

`cadops lock` and `cadops unlock` wrap `git lfs lock` and `git lfs unlock`, validate that the target file exists, and warn when locking is recommended for the file type but Git LFS is not configured correctly for that type.

CadOps does not auto-commit from `watch` and does not generate previews yet.

## Development

```bash
make test
```
