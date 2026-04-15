package metadata

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildRecord(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "parts", "gearbox.sldprt")
	mustWriteFile(t, path, "solid-body")

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}

	record, err := BuildRecord(root, "parts/gearbox.sldprt", info)
	if err != nil {
		t.Fatalf("build record: %v", err)
	}

	if record.Path != "parts/gearbox.sldprt" {
		t.Fatalf("expected path to be normalized, got %q", record.Path)
	}
	if record.TypeName != "SolidWorks Part" {
		t.Fatalf("expected CAD type, got %q", record.TypeName)
	}
	if record.Extension != ".sldprt" {
		t.Fatalf("expected extension .sldprt, got %q", record.Extension)
	}
	if record.SizeBytes != int64(len("solid-body")) {
		t.Fatalf("expected size %d, got %d", len("solid-body"), record.SizeBytes)
	}
	if _, err := time.Parse(time.RFC3339, record.ModifiedTime); err != nil {
		t.Fatalf("expected RFC3339 time, got %q", record.ModifiedTime)
	}
	if !record.GitLFSExpected {
		t.Fatal("expected Git LFS to be recommended for .sldprt")
	}
	if !record.LockingRecommended {
		t.Fatal("expected locking to be recommended for .sldprt")
	}
	if len(record.SHA256) != 64 {
		t.Fatalf("expected SHA-256 hex digest, got %q", record.SHA256)
	}
}

func TestHashFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "frame.step")
	mustWriteFile(t, path, "abc")

	got, err := HashFile(path)
	if err != nil {
		t.Fatalf("hash file: %v", err)
	}

	const want = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got != want {
		t.Fatalf("unexpected hash:\nwant: %s\ngot:  %s", want, got)
	}
}

func TestScanSkipsGitAndCadOpsMetadata(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "parts", "gearbox.sldprt"), "part-data")
	mustWriteFile(t, filepath.Join(root, ".git", "objects", "ignored.step"), "ignored")
	mustWriteFile(t, filepath.Join(root, ".cadops", "metadata", "ignored.fcstd"), "ignored")

	records, err := Scan(root, []string{".sldprt", ".step", ".fcstd"})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Path != "parts/gearbox.sldprt" {
		t.Fatalf("unexpected record path %q", records[0].Path)
	}
}

func TestGenerateEmptyManifest(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	manifest, err := Generate(root, []string{".step"})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	if manifest.Version != SchemaVersion {
		t.Fatalf("expected schema version %d, got %d", SchemaVersion, manifest.Version)
	}
	if len(manifest.Records) != 0 {
		t.Fatalf("expected empty manifest, got %d records", len(manifest.Records))
	}
	if _, err := time.Parse(time.RFC3339, manifest.GeneratedAt); err != nil {
		t.Fatalf("expected RFC3339 timestamp, got %q", manifest.GeneratedAt)
	}
}

func TestSaveLoadAndLookup(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	manifest := Manifest{
		Version:     SchemaVersion,
		GeneratedAt: "2026-04-15T16:00:00Z",
		Records: []Record{
			{
				Path:               "parts/gearbox.sldprt",
				TypeName:           "SolidWorks Part",
				Extension:          ".sldprt",
				SizeBytes:          42,
				ModifiedTime:       "2026-04-15T15:59:00Z",
				SHA256:             strings.Repeat("a", 64),
				GitLFSExpected:     true,
				LockingRecommended: true,
			},
		},
	}

	if err := Save(root, manifest); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := Load(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	record, ok := Lookup(loaded, "parts/gearbox.sldprt")
	if !ok {
		t.Fatal("expected lookup to find record")
	}
	if record.SHA256 != strings.Repeat("a", 64) {
		t.Fatalf("unexpected record hash %q", record.SHA256)
	}

	if _, ok := Lookup(loaded, "missing.step"); ok {
		t.Fatal("expected missing record lookup to fail")
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
