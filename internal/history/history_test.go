package history

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cadops/cadops/internal/metadata"
)

func TestParseLog(t *testing.T) {
	t.Parallel()

	output := "" +
		"abcdef1234567890\x1f2026-04-14\x1fAdded mount redesign\n" +
		"bracket.sldprt\n" +
		"assembly.sldasm\n" +
		"README.md\n\x1e\n" +
		"1234567890abcdef\x1f2026-04-13\x1fUpdated docs\n" +
		"docs/guide.md\n\x1e\n"

	entries := ParseLog(output)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].ShortHash != "abcdef1" {
		t.Fatalf("unexpected short hash %q", entries[0].ShortHash)
	}
	if entries[0].Date != "2026-04-14" || entries[0].Message != "Added mount redesign" {
		t.Fatalf("unexpected first entry metadata: %+v", entries[0])
	}
	if len(entries[0].CADFiles) != 2 {
		t.Fatalf("expected 2 CAD files, got %d", len(entries[0].CADFiles))
	}
	if len(entries[1].CADFiles) != 0 {
		t.Fatalf("expected no CAD files for docs commit, got %d", len(entries[1].CADFiles))
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	formatted := Format([]Entry{
		{
			ShortHash: "abcdef1",
			Date:      "2026-04-14",
			Message:   "Added mount redesign",
			CADFiles:  []string{"bracket.sldprt", "assembly.sldasm"},
		},
		{
			ShortHash: "1234567",
			Date:      "2026-04-13",
			Message:   "Updated docs",
		},
	})

	expected := "" +
		"abcdef1  2026-04-14  Added mount redesign\n" +
		"  CAD files:\n" +
		"    - bracket.sldprt\n" +
		"    - assembly.sldasm\n" +
		"\n" +
		"1234567  2026-04-13  Updated docs\n" +
		"  CAD files:\n" +
		"    - none\n"

	if formatted != expected {
		t.Fatalf("unexpected format:\n%s", formatted)
	}
}

func TestFormatEmpty(t *testing.T) {
	t.Parallel()

	if got := Format(nil); got != "No commits found\n" {
		t.Fatalf("unexpected empty output %q", got)
	}
}

func TestFormatReportEmpty(t *testing.T) {
	t.Parallel()

	if got := FormatReport(Report{}); got != "No commits found\n" {
		t.Fatalf("unexpected empty report output %q", got)
	}
}

func TestBuildReportEnrichesHistoryWithMetadata(t *testing.T) {
	t.Parallel()

	report := BuildReport(
		[]Entry{
			{
				Hash:      "abcdef1234567890",
				ShortHash: "abcdef1",
				Date:      "2026-04-14",
				Message:   "Updated bracket",
				CADFiles:  []string{"bracket.sldprt"},
			},
		},
		func(revision string) (metadata.Manifest, error) {
			switch revision {
			case "abcdef1234567890":
				return metadata.Manifest{
					Records: []metadata.Record{
						{
							Path:      "bracket.sldprt",
							TypeName:  "SolidWorks Part",
							SizeBytes: 140,
							SHA256:    strings.Repeat("b", 64),
						},
					},
				}, nil
			case "parent1234567890":
				return metadata.Manifest{
					Records: []metadata.Record{
						{
							Path:      "bracket.sldprt",
							TypeName:  "SolidWorks Part",
							SizeBytes: 100,
							SHA256:    strings.Repeat("a", 64),
						},
					},
				}, nil
			default:
				return metadata.Manifest{}, os.ErrNotExist
			}
		},
		func(commit string) (string, error) {
			if commit == "abcdef1234567890" {
				return "parent1234567890", nil
			}
			return "", nil
		},
	)

	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 detailed entry, got %d", len(report.Entries))
	}
	file := report.Entries[0].CADFiles[0]
	if !file.Metadata.MetadataAvailable {
		t.Fatal("expected metadata to be available")
	}
	if !file.Metadata.ChecksumChanged {
		t.Fatal("expected checksum change")
	}
	if !file.Metadata.HasSizeDelta || file.Metadata.SizeDeltaBytes != 40 {
		t.Fatalf("unexpected size delta %+v", file.Metadata)
	}

	output := FormatReport(report)
	if !strings.Contains(output, "SolidWorks Part; size 140 B; checksum changed; delta +40 B") {
		t.Fatalf("expected enriched output, got:\n%s", output)
	}
}

func TestBuildReportFallsBackWhenMetadataMissing(t *testing.T) {
	t.Parallel()

	report := BuildReport(
		[]Entry{
			{
				Hash:      "abcdef1234567890",
				ShortHash: "abcdef1",
				Date:      "2026-04-14",
				Message:   "Updated bracket",
				CADFiles:  []string{"bracket.sldprt"},
			},
		},
		func(revision string) (metadata.Manifest, error) {
			return metadata.Manifest{}, os.ErrNotExist
		},
		func(commit string) (string, error) {
			return "", nil
		},
	)

	output := FormatReport(report)
	if !strings.Contains(output, "bracket.sldprt [metadata unavailable]") {
		t.Fatalf("expected fallback metadata annotation, got:\n%s", output)
	}
}

func TestBuildReportWarnsAndContinuesOnMetadataLookupFailure(t *testing.T) {
	t.Parallel()

	report := BuildReport(
		[]Entry{
			{
				Hash:      "abcdef1234567890",
				ShortHash: "abcdef1",
				Date:      "2026-04-14",
				Message:   "Updated bracket",
				CADFiles:  []string{"bracket.sldprt"},
			},
		},
		func(revision string) (metadata.Manifest, error) {
			return metadata.Manifest{}, errors.New("boom")
		},
		func(commit string) (string, error) {
			return "", nil
		},
	)

	output := FormatReport(report)
	if !strings.Contains(output, "Warnings:\n  - metadata lookup failed for abcdef1\n") {
		t.Fatalf("expected warning output, got:\n%s", output)
	}
	if !strings.Contains(output, "bracket.sldprt [metadata unavailable]") {
		t.Fatalf("expected standard history output to continue, got:\n%s", output)
	}
}

func TestFilterCADFilesDeduplicatesAndSkipsNonCAD(t *testing.T) {
	t.Parallel()

	files := filterCADFiles([]string{
		" bracket.sldprt ",
		"README.md",
		"assembly.sldasm",
		"bracket.sldprt",
	})

	if len(files) != 2 {
		t.Fatalf("expected 2 CAD files, got %d", len(files))
	}
	if files[0] != "bracket.sldprt" || files[1] != "assembly.sldasm" {
		t.Fatalf("unexpected CAD files %#v", files)
	}
}
