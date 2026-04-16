package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cadops/cadops/internal/metadata"
)

func TestBuildReportEnrichesCADEntriesWithMetadataComparison(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "parts", "gearbox.sldprt")
	mustWriteFile(t, path, "updated-solid-body")

	if err := metadata.Save(root, metadata.Manifest{
		Version:     metadata.SchemaVersion,
		GeneratedAt: "2026-04-16T12:00:00Z",
		Records: []metadata.Record{
			{
				Path:               "parts/gearbox.sldprt",
				TypeName:           "SolidWorks Part",
				Extension:          ".sldprt",
				SizeBytes:          int64(len("old")),
				ModifiedTime:       "2026-04-16T11:00:00Z",
				SHA256:             strings.Repeat("a", 64),
				GitLFSExpected:     true,
				LockingRecommended: true,
			},
		},
	}); err != nil {
		t.Fatalf("save manifest: %v", err)
	}

	report := BuildReport(root, []Entry{
		{Kind: KindModified, Path: "parts/gearbox.sldprt", IsCAD: true},
	})

	if len(report.CAD) != 1 {
		t.Fatalf("expected 1 CAD entry, got %d", len(report.CAD))
	}
	if len(report.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %#v", report.Warnings)
	}

	details := report.CAD[0].Metadata
	if !details.HasCurrent || !details.HasPrevious {
		t.Fatalf("expected both current and previous metadata, got %+v", details)
	}
	if !details.Comparison.ChecksumChanged {
		t.Fatal("expected checksum change to be detected")
	}
	wantDelta := int64(len("updated-solid-body") - len("old"))
	if !details.Comparison.HasSizeDelta || details.Comparison.SizeDeltaBytes != wantDelta {
		t.Fatalf("unexpected size delta %+v", details.Comparison)
	}

	output := FormatReport(report)
	if !strings.Contains(output, "SolidWorks Part; lock yes; LFS yes; checksum changed; size +15 B") {
		t.Fatalf("expected enriched metadata details, got:\n%s", output)
	}
}

func TestBuildReportFallsBackToCurrentMetadataWithoutManifest(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "exports", "frame.step"), "step-data")

	report := BuildReport(root, []Entry{
		{Kind: KindAdded, Path: "exports/frame.step", IsCAD: true},
		{Kind: KindDeleted, Path: "notes.txt", IsCAD: false},
	})

	output := FormatReport(report)
	if !strings.Contains(output, "A exports/frame.step [STEP; lock no; LFS yes]") {
		t.Fatalf("expected current metadata context without manifest, got:\n%s", output)
	}
	if !strings.Contains(output, "D notes.txt") {
		t.Fatalf("expected non-CAD fallback output, got:\n%s", output)
	}
	if strings.Contains(output, "checksum changed") {
		t.Fatalf("did not expect checksum comparison without previous metadata, got:\n%s", output)
	}
}

func TestBuildReportWarnsAndContinuesWhenCurrentMetadataLookupFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	report := BuildReport(root, []Entry{
		{Kind: KindModified, Path: "missing.sldprt", IsCAD: true},
	})

	output := FormatReport(report)
	if !strings.Contains(output, "M missing.sldprt") {
		t.Fatalf("expected diff line to remain visible, got:\n%s", output)
	}
	if !strings.Contains(output, "Warnings:\n  - metadata unavailable for missing.sldprt\n") {
		t.Fatalf("expected concise metadata warning, got:\n%s", output)
	}
}

func TestCompareRecords(t *testing.T) {
	t.Parallel()

	comparison := CompareRecords(
		metadata.Record{SHA256: strings.Repeat("a", 64), SizeBytes: 100},
		metadata.Record{SHA256: strings.Repeat("b", 64), SizeBytes: 140},
	)

	if !comparison.ChecksumChanged {
		t.Fatal("expected checksum change")
	}
	if !comparison.HasSizeDelta || comparison.SizeDeltaBytes != 40 {
		t.Fatalf("unexpected comparison %+v", comparison)
	}
}

func TestFormatSizeDelta(t *testing.T) {
	t.Parallel()

	if got := FormatSizeDelta(42); got != "+42 B" {
		t.Fatalf("unexpected positive delta %q", got)
	}
	if got := FormatSizeDelta(-9); got != "-9 B" {
		t.Fatalf("unexpected negative delta %q", got)
	}
}

func TestFormatReportNoChanges(t *testing.T) {
	t.Parallel()

	if got := FormatReport(Report{}); got != "No repository changes\n" {
		t.Fatalf("unexpected no-change output %q", got)
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
