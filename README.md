# CadOps

CadOps is a CAD-aware command-line workflow layer over Git and Git LFS. The current CLI provides safe repository initialization, CAD-aware status output, repository health checks, CAD repository scanning, repository watching, snapshot commits, guarded CAD-aware commits, config inspection, metadata and preview record generation, guarded push and pull flows, CAD-aware history output, and Git LFS lock helpers for file-based CAD workflows.

CadOps requires both `git` and `git lfs` to be installed on the user machine.

## Supported CAD Files

CadOps detects and applies Git LFS policy to these CAD extensions by default:

- SolidWorks: `.sldprt`, `.sldasm`
- Generic part and exchange formats: `.prt`, `.step`, `.stp`, `.iges`, `.igs`
- Mesh and archive formats: `.stl`
- Fusion 360: `.f3d`, `.f3z`
- Inventor: `.ipt`, `.iam`
- FreeCAD: `.fcstd`

The generated `.cadops.yaml` stores the active extension list as `tracked_extensions`, and `cadops init` writes matching `.gitattributes` rules so these files are tracked through Git LFS.

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

## Quickstart

```bash
mkdir cadops-demo
cd cadops-demo
git init
cadops init
mkdir parts exports docs
printf "fake solidworks part\n" > parts/part.sldprt
printf "fake solidworks assembly\n" > parts/assembly.sldasm
printf "fake fusion archive\n" > parts/archive.f3z
printf "fake mesh\n" > exports/mesh.stl
printf "ISO-10303-21;\n" > exports/export.step
printf "notes\n" > docs/notes.txt
cadops metadata generate
cadops preview generate
cadops scan
cadops status
```

The fake CAD files are enough to exercise CadOps classification, metadata, preview records, LFS policy checks, and status output without requiring CAD software.

For a real remote workflow, create a GitHub repository with Git LFS enabled, add it as `origin`, then use the guarded CadOps wrappers:

```bash
git remote add origin https://github.com/<owner>/<repo>.git
git add .
cadops commit -m "Initialize CadOps tracking"
cadops push
cadops pull
cadops lock parts/assembly.sldasm
cadops unlock parts/assembly.sldasm
```

Git LFS locking requires a remote that supports the LFS locking API, such as GitHub. Local `file://` bare repositories are useful for push and pull smoke tests, but they do not support LFS locks.

## Command Reference

- `cadops init`
- `cadops status`
- `cadops diff`
- `cadops files`
- `cadops scan`
- `cadops metadata generate`
- `cadops metadata show <file>`
- `cadops preview generate`
- `cadops preview show <file>`
- `cadops preview list`
- `cadops doctor`
- `cadops watch`
- `cadops snapshot`
- `cadops commit -m <message>`
- `cadops config show`
- `cadops config get <key>`
- `cadops push`
- `cadops pull`
- `cadops history`
- `cadops lock <file>`
- `cadops unlock <file>`

`cadops watch` monitors the Git repository recursively, reacts only to CAD extensions configured in `.cadops.yaml`, prints concise change lines, and can auto-stage changed CAD files when `auto_stage: true`.

`cadops status` groups working tree changes into CAD and non-CAD files, reports missing or stale metadata records for changed CAD files, reports missing or stale preview records, and warns when changed CAD extensions lack matching Git LFS attributes.

`cadops snapshot` stages changed CAD files, regenerates the repo metadata and lightweight preview manifests before commit, includes `.cadops/metadata/manifest.json` and `.cadops/previews/manifest.json` in the same snapshot when refresh succeeds, and creates a timestamped snapshot commit like `snapshot: 2026-04-14 15:42`. Metadata and preview refresh warnings do not block snapshot creation.

`cadops commit -m "message"` runs CAD-aware pre-commit checks, warns about unstaged changes, missing LFS coverage, missing local locks for recommended-lock CAD files when lock inspection is available, and stale or missing metadata/preview records, then delegates to `git commit -m`. It does not auto-stage unrelated files.

`cadops files` scans the current Git repository recursively for configured CAD extensions, groups matches by CAD type, and shows each file path with its CAD type and lock recommendation.

`cadops diff` summarizes current Git working tree changes, separates CAD and non-CAD files using configured extensions, and prints concise Git-backed change statuses such as modified, added, deleted, and renamed. For changed CAD files it also uses the current metadata manifest as the baseline when available, compares that stored record against freshly derived file metadata when the file still exists locally, and adds compact context such as CAD type, `locking recommended yes/no`, `LFS expected yes/no`, `checksum changed yes/no`, and file size delta. It also reports preview status as generated, unavailable, unsupported, stale, or missing when preview records exist. When metadata cannot be read, the changed file still appears with an explicit fallback annotation instead of pretending to know more.

`cadops scan` audits the current repository for configured CAD files, summarizes counts by CAD type, highlights locking and Git LFS expectations, reports preview coverage, warns about missing `.gitattributes` LFS rules for expected CAD file types, and shows the largest CAD files plus top CAD-heavy directories. It uses stored metadata when available and falls back to a live scan when metadata is absent.

`cadops metadata generate` scans the current Git repository for configured CAD extensions, classifies matching files with the built-in CAD registry, computes filesystem-level metadata including size, modified time, SHA-256, LFS expectation, and lock recommendation, and writes a single manifest to `.cadops/metadata/manifest.json`. `cadops metadata show <file>` reads that manifest and prints the stored record for one CAD file. Snapshot commits also refresh this manifest automatically before commit for consistency.

`cadops preview generate` creates V1 preview records under `.cadops/previews/manifest.json` for configured CAD files. This does not render proprietary CAD geometry. If a same-name sidecar image such as `part.png`, `part.jpg`, or `part.jpeg` exists next to the CAD file, CadOps records it as a referenced artifact; otherwise it stores an honest unavailable or unsupported placeholder record. `cadops preview show <file>` prints one record and marks it stale when the source hash has changed. `cadops preview list` shows all stored preview records.

`cadops config show` prints the supported `.cadops.yaml` keys in a concise terminal format. `cadops config get <key>` returns a single value for `version`, `tracked_extensions`, `auto_stage`, `require_lfs`, or `locking_enabled`. Both commands can be run from the repository root or a subdirectory.

`cadops doctor` validates Git, Git LFS, repository state, `.cadops.yaml`, `.gitattributes`, metadata and preview manifest state, CAD tracking coverage, and remote configuration. Checks are reported as `PASS`, `WARN`, or `FAIL`.

`cadops push` runs light CAD-aware pre-push checks, warns about local CAD changes or missing LFS coverage, and stops early when no remote is configured before delegating to `git push`.

`cadops pull` verifies Git LFS availability, warns on a dirty working tree and modified CAD files, and then delegates to `git pull`.

`cadops history` shows a compact recent commit view with short hash, date, message, and changed CAD files for each commit. When `.cadops/metadata/manifest.json` is present in commit history, it also adds compact CAD metadata context per changed file such as CAD type, stored file size, and checksum-changed or size-delta indicators when both the commit and its first parent contain usable metadata manifests. When `.cadops/previews/manifest.json` is present in commit history, it adds preview availability for changed CAD files. Missing metadata or preview records fall back cleanly to concise unavailable annotations.

`cadops lock` and `cadops unlock` wrap `git lfs lock` and `git lfs unlock`, validate that the target file exists, and warn when locking is recommended for the file type but Git LFS is not configured correctly for that type. Locking is recommended for single-writer binary formats such as SolidWorks assemblies, Fusion 360 archives, and generic part files.

## Limitations

CadOps does not auto-commit from `watch`, does not implement semantic CAD diffing, does not implement geometry-aware metadata, and does not implement semantic CAD history yet. Preview V1 stores preview manifests and optional references to real sidecar images only; it does not render SolidWorks, NX, Inventor, or other proprietary CAD geometry.

CadOps relies on Git, Git LFS, filesystem metadata, SHA-256 hashes, configured CAD extensions, and stored CadOps manifests. It does not inspect proprietary CAD feature trees, assemblies, constraints, drawings, or model geometry.

## Development

```bash
make test
```
