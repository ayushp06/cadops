package snapshot

import (
	"path/filepath"

	"github.com/cadops/cadops/internal/metadata"
)

const metadataManifestPath = ".cadops/metadata/manifest.json"

// MetadataUpdate describes the repo-level metadata manifest prepared for a
// snapshot commit. Snapshot refreshes the full manifest before commit so the
// stored metadata and CAD changes land atomically in the same revision.
type MetadataUpdate struct {
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
