package preview

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cadops/cadops/internal/metadata"
)

func TestGenerateCreatesUnavailableRecordsWithoutRenderer(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWritePreviewFile(t, filepath.Join(root, "part.sldprt"), "solidworks")
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

	manifest, err := Generate(root, []string{".sldprt"}, now)
	if err != nil {
		t.Fatalf("generate preview: %v", err)
	}

	if len(manifest.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(manifest.Records))
	}
	record := manifest.Records[0]
	if record.SourcePath != "part.sldprt" {
		t.Fatalf("unexpected source path %q", record.SourcePath)
	}
	if record.Status != StatusUnavailable {
		t.Fatalf("expected unavailable status, got %q", record.Status)
	}
	if record.ArtifactPath != "" {
		t.Fatalf("did not expect generated artifact path, got %q", record.ArtifactPath)
	}
	if record.GeneratedAt != "2026-05-08T12:00:00Z" {
		t.Fatalf("unexpected generated timestamp %q", record.GeneratedAt)
	}
}

func TestBuildRecordReferencesExistingSidecarImage(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWritePreviewFile(t, filepath.Join(root, "parts", "bracket.sldprt"), "solidworks")
	mustWritePreviewFile(t, filepath.Join(root, "parts", "bracket.png"), "image")

	info, err := os.Stat(filepath.Join(root, "parts", "bracket.sldprt"))
	if err != nil {
		t.Fatalf("stat source: %v", err)
	}
	source, err := metadata.BuildRecord(root, "parts/bracket.sldprt", info)
	if err != nil {
		t.Fatalf("build metadata: %v", err)
	}

	record := BuildRecord(root, source, "2026-05-08T12:00:00Z")
	if record.Status != StatusGenerated {
		t.Fatalf("expected generated status, got %q", record.Status)
	}
	if record.ArtifactPath != "parts/bracket.png" {
		t.Fatalf("unexpected artifact path %q", record.ArtifactPath)
	}
}

func TestBuildRecordMarksUnknownConfiguredExtensionUnsupported(t *testing.T) {
	t.Parallel()

	record := BuildRecord(t.TempDir(), metadata.Record{
		Path:      "custom.widget",
		TypeName:  "Unknown CAD Type (.widget)",
		Extension: ".widget",
		SHA256:    "abc",
	}, "2026-05-08T12:00:00Z")

	if record.Status != StatusUnsupported {
		t.Fatalf("expected unsupported status, got %q", record.Status)
	}
}

func TestIsStaleDetectsSourceHashChange(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	sourcePath := filepath.Join(root, "export.step")
	mustWritePreviewFile(t, sourcePath, "old")
	oldHash, err := metadata.HashFile(sourcePath)
	if err != nil {
		t.Fatalf("hash old source: %v", err)
	}

	record := Record{SourcePath: "export.step", SourceHash: oldHash, Status: StatusUnavailable}
	mustWritePreviewFile(t, sourcePath, "new")

	stale, err := IsStale(root, record)
	if err != nil {
		t.Fatalf("stale check: %v", err)
	}
	if !stale {
		t.Fatal("expected stale preview record")
	}
	if got := EffectiveStatus(root, record); got != StatusStale {
		t.Fatalf("expected effective stale status, got %q", got)
	}
}

func TestLookupNormalizesPath(t *testing.T) {
	t.Parallel()

	manifest := Manifest{Records: []Record{{SourcePath: "parts/bracket.sldprt"}}}
	if _, ok := Lookup(manifest, filepath.Join("parts", "bracket.sldprt")); !ok {
		t.Fatal("expected lookup to normalize path")
	}
}

func mustWritePreviewFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
