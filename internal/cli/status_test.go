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

func TestStatusEntryPathsSorts(t *testing.T) {
	t.Parallel()

	paths := statusEntryPaths([]gitx.StatusEntry{
		{Path: "b.step"},
		{Path: "a.step"},
	})
	if len(paths) != 2 || paths[0] != "a.step" || paths[1] != "b.step" {
		t.Fatalf("unexpected paths %#v", paths)
	}
}
