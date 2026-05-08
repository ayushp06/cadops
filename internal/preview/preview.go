package preview

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/metadata"
)

const (
	SchemaVersion    = 1
	dirName          = ".cadops/previews"
	manifestFileName = "manifest.json"
)

// Status describes the availability of a preview record.
type Status string

const (
	StatusGenerated   Status = "generated"
	StatusUnavailable Status = "unavailable"
	StatusUnsupported Status = "unsupported"
	StatusStale       Status = "stale"
)

// Record describes a stored preview artifact decision for one CAD source file.
type Record struct {
	SourcePath   string `json:"source_path"`
	CADType      string `json:"cad_type"`
	SourceHash   string `json:"source_hash"`
	Status       Status `json:"status"`
	ArtifactPath string `json:"artifact_path,omitempty"`
	GeneratedAt  string `json:"generated_at"`
	Message      string `json:"message,omitempty"`
}

// Manifest stores repository-wide preview records.
type Manifest struct {
	Version     int      `json:"version"`
	GeneratedAt string   `json:"generated_at"`
	Records     []Record `json:"records"`
}

// ManifestPath returns the repository-local preview manifest path.
func ManifestPath(root string) string {
	return filepath.Join(root, dirName, manifestFileName)
}

// Generate creates preview records for configured CAD files without invoking
// proprietary renderers or inventing geometry previews.
func Generate(root string, extensions []string, now time.Time) (Manifest, error) {
	records, err := metadata.Scan(root, extensions)
	if err != nil {
		return Manifest{}, err
	}

	previewRecords := make([]Record, 0, len(records))
	generatedAt := now.UTC().Format(time.RFC3339)
	for _, record := range records {
		previewRecords = append(previewRecords, BuildRecord(root, record, generatedAt))
	}
	sort.Slice(previewRecords, func(i, j int) bool {
		return previewRecords[i].SourcePath < previewRecords[j].SourcePath
	})

	return Manifest{
		Version:     SchemaVersion,
		GeneratedAt: generatedAt,
		Records:     previewRecords,
	}, nil
}

// BuildRecord creates one honest V1 preview record for a metadata record.
func BuildRecord(root string, source metadata.Record, generatedAt string) Record {
	record := Record{
		SourcePath:  source.Path,
		CADType:     source.TypeName,
		SourceHash:  source.SHA256,
		Status:      StatusUnavailable,
		GeneratedAt: generatedAt,
		Message:     "no CAD renderer configured; placeholder preview record only",
	}

	if _, ok := cad.Lookup(source.Extension); !ok {
		record.Status = StatusUnsupported
		record.Message = "preview generation unsupported for configured extension " + source.Extension
		return record
	}

	if artifact, ok := findSidecarImage(root, source.Path); ok {
		record.Status = StatusGenerated
		record.ArtifactPath = artifact
		record.Message = "existing image preview referenced"
	}

	return record
}

// Load reads the repo-local preview manifest.
func Load(root string) (Manifest, error) {
	data, err := os.ReadFile(ManifestPath(root))
	if err != nil {
		return Manifest{}, err
	}
	return Parse(data)
}

// Parse decodes preview manifest JSON data.
func Parse(data []byte) (Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// Save persists a preview manifest to the repo-local preview directory.
func Save(root string, manifest Manifest) error {
	path := ManifestPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// Lookup finds a preview record by repository-relative source path.
func Lookup(manifest Manifest, relPath string) (Record, bool) {
	normalized := filepath.ToSlash(filepath.Clean(relPath))
	for _, record := range manifest.Records {
		if record.SourcePath == normalized {
			return record, true
		}
	}
	return Record{}, false
}

// EffectiveStatus reports stale when the source hash no longer matches.
func EffectiveStatus(root string, record Record) Status {
	if root == "" {
		return record.Status
	}
	stale, err := IsStale(root, record)
	if err == nil && stale {
		return StatusStale
	}
	return record.Status
}

// IsStale compares a preview record's stored source hash with the current file.
func IsStale(root string, record Record) (bool, error) {
	if record.SourcePath == "" || record.SourceHash == "" {
		return false, nil
	}
	hash, err := metadata.HashFile(filepath.Join(root, filepath.FromSlash(record.SourcePath)))
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	return hash != record.SourceHash, nil
}

func findSidecarImage(root, sourcePath string) (string, bool) {
	relDir := filepath.ToSlash(filepath.Dir(sourcePath))
	base := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	candidates := []string{
		base + ".png",
		base + ".jpg",
		base + ".jpeg",
	}
	for _, candidate := range candidates {
		relPath := filepath.ToSlash(filepath.Join(relDir, candidate))
		if relDir == "." {
			relPath = candidate
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(relPath))); err == nil {
			return relPath, true
		}
	}
	return "", false
}
