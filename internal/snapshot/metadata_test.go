package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cadops/cadops/internal/metadata"
	"github.com/cadops/cadops/internal/preview"
)

func TestRefreshMetadata(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	partPath := filepath.Join(root, "parts", "gearbox.sldprt")
	if err := os.MkdirAll(filepath.Dir(partPath), 0o755); err != nil {
		t.Fatalf("mkdir part dir: %v", err)
	}
	if err := os.WriteFile(partPath, []byte("solid-body"), 0o644); err != nil {
		t.Fatalf("write part: %v", err)
	}

	update, err := RefreshMetadata(root, []string{".sldprt"})
	if err != nil {
		t.Fatalf("refresh metadata: %v", err)
	}
	if update.Path != ".cadops/metadata/manifest.json" {
		t.Fatalf("unexpected manifest path %q", update.Path)
	}
	if update.RecordCount != 1 {
		t.Fatalf("expected 1 record, got %d", update.RecordCount)
	}

	manifest, err := metadata.Load(root)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if len(manifest.Records) != 1 || manifest.Records[0].Path != "parts/gearbox.sldprt" {
		t.Fatalf("unexpected manifest records %#v", manifest.Records)
	}
}

func TestRefreshPreviews(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	partPath := filepath.Join(root, "parts", "gearbox.sldprt")
	if err := os.MkdirAll(filepath.Dir(partPath), 0o755); err != nil {
		t.Fatalf("mkdir part dir: %v", err)
	}
	if err := os.WriteFile(partPath, []byte("solid-body"), 0o644); err != nil {
		t.Fatalf("write part: %v", err)
	}

	update, err := RefreshPreviews(root, []string{".sldprt"}, time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("refresh previews: %v", err)
	}
	if update.Path != ".cadops/previews/manifest.json" {
		t.Fatalf("unexpected manifest path %q", update.Path)
	}
	if update.RecordCount != 1 {
		t.Fatalf("expected 1 record, got %d", update.RecordCount)
	}

	manifest, err := preview.Load(root)
	if err != nil {
		t.Fatalf("load preview manifest: %v", err)
	}
	if len(manifest.Records) != 1 || manifest.Records[0].SourcePath != "parts/gearbox.sldprt" {
		t.Fatalf("unexpected preview records %#v", manifest.Records)
	}
}
