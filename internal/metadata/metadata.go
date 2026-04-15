package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/watch"
)

const (
	SchemaVersion    = 1
	dirName          = ".cadops/metadata"
	manifestFileName = "manifest.json"
)

// Record describes filesystem-level metadata for a CAD file.
type Record struct {
	Path               string `json:"path"`
	TypeName           string `json:"type_name"`
	Extension          string `json:"extension"`
	SizeBytes          int64  `json:"size_bytes"`
	ModifiedTime       string `json:"modified_time"`
	SHA256             string `json:"sha256"`
	GitLFSExpected     bool   `json:"git_lfs_expected"`
	LockingRecommended bool   `json:"locking_recommended"`
}

// Manifest stores repository-wide CAD metadata in a single JSON file.
type Manifest struct {
	Version     int      `json:"version"`
	GeneratedAt string   `json:"generated_at"`
	Records     []Record `json:"records"`
}

// ManifestPath returns the repository-local metadata manifest path.
func ManifestPath(root string) string {
	return filepath.Join(root, dirName, manifestFileName)
}

// Generate scans the repository for configured CAD extensions and builds a
// stable metadata manifest.
func Generate(root string, extensions []string) (Manifest, error) {
	records, err := Scan(root, extensions)
	if err != nil {
		return Manifest{}, err
	}

	return Manifest{
		Version:     SchemaVersion,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Records:     records,
	}, nil
}

// Scan walks the repository tree and builds metadata records for matching
// files. Git and CadOps metadata directories are skipped.
func Scan(root string, extensions []string) ([]Record, error) {
	filter := watch.NewFilter(extensions)
	records := make([]Record, 0)

	err := filepath.WalkDir(root, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := dirEntry.Name()
		if dirEntry.IsDir() {
			if name == ".git" || path == filepath.Join(root, ".cadops") {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		if !filter.Match(relPath) {
			return nil
		}

		info, err := dirEntry.Info()
		if err != nil {
			return err
		}

		record, err := BuildRecord(root, relPath, info)
		if err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Path < records[j].Path
	})
	return records, nil
}

// BuildRecord converts a repository-relative file path plus file info into a
// metadata record.
func BuildRecord(root, relPath string, info fs.FileInfo) (Record, error) {
	extension := strings.ToLower(filepath.Ext(relPath))
	fileType, ok := cad.Lookup(extension)
	typeName := unknownTypeName(extension)
	useLFS := false
	recommendLocking := false
	if ok {
		typeName = fileType.Name
		useLFS = fileType.UseLFS
		recommendLocking = fileType.RecommendLocking
	}

	hash, err := HashFile(filepath.Join(root, filepath.FromSlash(relPath)))
	if err != nil {
		return Record{}, err
	}

	return Record{
		Path:               filepath.ToSlash(relPath),
		TypeName:           typeName,
		Extension:          extension,
		SizeBytes:          info.Size(),
		ModifiedTime:       info.ModTime().UTC().Format(time.RFC3339),
		SHA256:             hash,
		GitLFSExpected:     useLFS,
		LockingRecommended: recommendLocking,
	}, nil
}

// HashFile returns the SHA-256 checksum for the target file.
func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sum := sha256.New()
	if _, err := io.Copy(sum, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(sum.Sum(nil)), nil
}

// Save persists the manifest to the repo-local metadata directory.
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

// Load reads the repo-local metadata manifest.
func Load(root string) (Manifest, error) {
	data, err := os.ReadFile(ManifestPath(root))
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// Lookup finds a record by repository-relative path.
func Lookup(manifest Manifest, relPath string) (Record, bool) {
	normalized := filepath.ToSlash(filepath.Clean(relPath))
	for _, record := range manifest.Records {
		if record.Path == normalized {
			return record, true
		}
	}
	return Record{}, false
}

func unknownTypeName(extension string) string {
	if extension == "" {
		return "Unknown CAD Type"
	}
	return "Unknown CAD Type (" + extension + ")"
}
