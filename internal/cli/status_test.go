package cli

import (
	"testing"

	"github.com/cadops/cadops/internal/gitx"
)

func TestSummarizeStatus(t *testing.T) {
	t.Parallel()

	summary := summarizeStatus([]gitx.StatusEntry{
		{Code: " M", Path: "parts/a.sldprt"},
		{Code: "??", Path: "docs/spec.md"},
		{Code: "A ", Path: "assy/model.fcstd"},
	})

	if len(summary.CADFiles) != 2 {
		t.Fatalf("expected 2 CAD files, got %d", len(summary.CADFiles))
	}
	if len(summary.NonCADFiles) != 1 {
		t.Fatalf("expected 1 non-CAD file, got %d", len(summary.NonCADFiles))
	}
}
