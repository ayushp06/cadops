package snapshot

import (
	"path/filepath"
	"time"

	"github.com/cadops/cadops/internal/metadata"
	"github.com/cadops/cadops/internal/preview"
)

const metadataManifestPath = ".cadops/metadata/manifest.json"
const previewManifestPath = ".cadops/previews/manifest.json"

// MetadataUpdate describes the repo-level metadata manifest prepared for a
// snapshot commit. Snapshot refreshes the full manifest before commit so the
// stored metadata and CAD changes land atomically in the same revision.
type MetadataUpdate struct {
	Path        string
	RecordCount int
}

// PreviewUpdate describes the repo-level preview manifest prepared for a
// snapshot commit.
type PreviewUpdate struct {
	Path        string
	RecordCount int
}

// RefreshMetadata regenerates the full repository metadata manifest for the
// configured CAD extensions and writes it to the CadOps-owned metadata path.
func RefreshMetadata(root string, extensions []string) (MetadataUpdate, error) {
	manifest, err := metadata.Generate(root, extensions)
	if err != nil {
		return MetadataUpdate{}, err
	}
	if err := metadata.Save(root, manifest); err != nil {
		return MetadataUpdate{}, err
	}

	return MetadataUpdate{
		Path:        filepath.ToSlash(metadataManifestPath),
		RecordCount: len(manifest.Records),
	}, nil
}

// RefreshPreviews regenerates lightweight preview records for configured CAD
// files without requiring external CAD software.
func RefreshPreviews(root string, extensions []string, now time.Time) (PreviewUpdate, error) {
	manifest, err := preview.Generate(root, extensions, now)
	if err != nil {
		return PreviewUpdate{}, err
	}
	if err := preview.Save(root, manifest); err != nil {
		return PreviewUpdate{}, err
	}

	return PreviewUpdate{
		Path:        filepath.ToSlash(previewManifestPath),
		RecordCount: len(manifest.Records),
	}, nil
}
