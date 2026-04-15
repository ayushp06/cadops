package scan

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/cadops/cadops/internal/metadata"
)

func TestBuildReport(t *testing.T) {
	t.Parallel()

	files := []File{
		BuildFile("parts/gearbox.sldprt", 200),
		BuildFile("parts/housing.sldprt", 300),
		BuildFile("exports/frame.step", 150),
		BuildFile("root.igs", 100),
	}

	report := BuildReport(files, gitAttributes(".sldprt"), false)

	if report.TotalFiles != 4 {
		t.Fatalf("expected 4 files, got %d", report.TotalFiles)
	}

	wantTypes := []TypeCount{
		{TypeName: "SolidWorks Part", Count: 2},
		{TypeName: "IGES", Count: 1},
		{TypeName: "STEP", Count: 1},
	}
	if !reflect.DeepEqual(report.ByType, wantTypes) {
		t.Fatalf("unexpected type counts:\nwant: %#v\ngot:  %#v", wantTypes, report.ByType)
	}

	wantLocking := []string{"parts/gearbox.sldprt", "parts/housing.sldprt"}
	if !reflect.DeepEqual(report.LockingRecommended, wantLocking) {
		t.Fatalf("unexpected locking list:\nwant: %#v\ngot:  %#v", wantLocking, report.LockingRecommended)
	}

	wantWarnings := []LFSWarning{
		{Path: "exports/frame.step", Extension: ".step"},
		{Path: "root.igs", Extension: ".igs"},
	}
	if !reflect.DeepEqual(report.LFSWarnings, wantWarnings) {
		t.Fatalf("unexpected lfs warnings:\nwant: %#v\ngot:  %#v", wantWarnings, report.LFSWarnings)
	}
}

func TestLargestFiles(t *testing.T) {
	t.Parallel()

	files := []File{
		BuildFile("c/frame.step", 150),
		BuildFile("a/part.sldprt", 300),
		BuildFile("b/assy.sldasm", 300),
	}

	got := LargestFiles(files, 2)
	want := []SizedFile{
		{Path: "a/part.sldprt", SizeBytes: 300},
		{Path: "b/assy.sldasm", SizeBytes: 300},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected largest files:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestTopDirectories(t *testing.T) {
	t.Parallel()

	files := []File{
		BuildFile("parts/a.sldprt", 1),
		BuildFile("parts/b.sldprt", 1),
		BuildFile("exports/a.step", 1),
		BuildFile("root.igs", 1),
	}

	got := TopDirectories(files, 3)
	want := []DirectoryCount{
		{Path: "parts", Count: 2},
		{Path: ".", Count: 1},
		{Path: "exports", Count: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected top directories:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestLoadFilesUsesMetadataWhenAvailable(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	manifest := metadata.Manifest{
		Version:     metadata.SchemaVersion,
		GeneratedAt: "2026-04-15T00:00:00Z",
		Records: []metadata.Record{
			{
				Path:               "parts/gearbox.sldprt",
				TypeName:           "SolidWorks Part",
				Extension:          ".sldprt",
				SizeBytes:          123,
				GitLFSExpected:     true,
				LockingRecommended: true,
			},
		},
	}
	if err := metadata.Save(root, manifest); err != nil {
		t.Fatalf("save metadata: %v", err)
	}

	files, usedMetadata, err := LoadFiles(root, []string{".sldprt"})
	if err != nil {
		t.Fatalf("load files: %v", err)
	}
	if !usedMetadata {
		t.Fatal("expected metadata to be used")
	}
	if len(files) != 1 || files[0].Path != "parts/gearbox.sldprt" {
		t.Fatalf("unexpected files: %#v", files)
	}
}

func TestLoadFilesFallsBackToLiveScan(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "parts", "gearbox.sldprt"), strings.Repeat("x", 10))

	files, usedMetadata, err := LoadFiles(root, []string{".sldprt"})
	if err != nil {
		t.Fatalf("load files: %v", err)
	}
	if usedMetadata {
		t.Fatal("expected live scan fallback")
	}
	if len(files) != 1 || files[0].Path != "parts/gearbox.sldprt" {
		t.Fatalf("unexpected files: %#v", files)
	}
}

func TestFindLFSWarnings(t *testing.T) {
	t.Parallel()

	files := []File{
		BuildFile("parts/gearbox.sldprt", 1),
		BuildFile("exports/frame.step", 1),
	}

	got := FindLFSWarnings(files, gitAttributes(".sldprt", ".step"))
	if len(got) != 0 {
		t.Fatalf("expected no warnings, got %#v", got)
	}
}

func gitAttributes(extensions ...string) string {
	lines := make([]string, 0, len(extensions))
	for _, extension := range extensions {
		lines = append(lines, "*"+extension+" filter=lfs diff=lfs merge=lfs -text")
	}
	return strings.Join(lines, "\n")
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
