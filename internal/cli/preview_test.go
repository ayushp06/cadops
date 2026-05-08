package cli

import (
	"strings"
	"testing"

	"github.com/cadops/cadops/internal/preview"
)

func TestFormatPreviewGenerateReport(t *testing.T) {
	t.Parallel()

	out := formatPreviewGenerateReport(preview.Manifest{
		Records: []preview.Record{
			{Status: preview.StatusGenerated},
			{Status: preview.StatusUnavailable},
			{Status: preview.StatusUnsupported},
		},
	})

	for _, want := range []string{
		"Generated preview records for 3 CAD files at .cadops/previews/manifest.json",
		"Generated: 1",
		"Unavailable: 1",
		"Unsupported: 1",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestFormatPreviewList(t *testing.T) {
	t.Parallel()

	out := formatPreviewList("", preview.Manifest{
		Records: []preview.Record{
			{
				SourcePath:   "parts/bracket.sldprt",
				CADType:      "SolidWorks Part",
				Status:       preview.StatusGenerated,
				ArtifactPath: "parts/bracket.png",
			},
		},
	})

	if !strings.Contains(out, "parts/bracket.sldprt | SolidWorks Part | generated | artifact: parts/bracket.png") {
		t.Fatalf("unexpected list output:\n%s", out)
	}
}

func TestFormatPreviewListEmpty(t *testing.T) {
	t.Parallel()

	if got := formatPreviewList("", preview.Manifest{}); got != "No preview records found\n" {
		t.Fatalf("unexpected empty output %q", got)
	}
}

func TestFormatPreviewRecord(t *testing.T) {
	t.Parallel()

	out := formatPreviewRecord("", preview.Record{
		SourcePath:   "export.step",
		CADType:      "STEP",
		SourceHash:   "abc",
		Status:       preview.StatusUnavailable,
		GeneratedAt:  "2026-05-08T12:00:00Z",
		ArtifactPath: "",
		Message:      "no CAD renderer configured; placeholder preview record only",
	})

	for _, want := range []string{
		"Source: export.step",
		"CAD Type: STEP",
		"Status: unavailable",
		"Message: no CAD renderer configured; placeholder preview record only",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in output:\n%s", want, out)
		}
	}
}
