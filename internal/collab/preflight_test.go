package collab

import (
	"testing"

	"github.com/cadops/cadops/internal/gitx"
)

func TestAssessPush(t *testing.T) {
	t.Parallel()

	report := AssessPush(
		[]gitx.StatusEntry{{Code: " M", Path: "parts/widget.sldprt"}},
		[]string{"parts/widget.sldprt", "docs/notes.txt"},
		"",
		false,
	)

	if report.CanPush {
		t.Fatal("expected push to be blocked without a remote")
	}
	if len(report.Warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d", len(report.Warnings))
	}
}

func TestAssessPushAllowsSoftWarnings(t *testing.T) {
	t.Parallel()

	report := AssessPush(
		[]gitx.StatusEntry{{Code: " M", Path: "parts/widget.sldprt"}},
		[]string{"parts/widget.sldprt"},
		"*.sldprt filter=lfs diff=lfs merge=lfs -text\n",
		true,
	)

	if !report.CanPush {
		t.Fatal("expected push to proceed when only soft warnings are present")
	}
	if len(report.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(report.Warnings))
	}
}

func TestAssessPull(t *testing.T) {
	t.Parallel()

	report := AssessPull(
		[]gitx.StatusEntry{
			{Code: " M", Path: "parts/widget.sldprt"},
			{Code: " M", Path: "README.md"},
		},
		false,
	)

	if report.CanPull {
		t.Fatal("expected pull to be blocked without Git LFS")
	}
	if len(report.Warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d", len(report.Warnings))
	}
}

func TestFindUncoveredCADFiles(t *testing.T) {
	t.Parallel()

	uncovered := FindUncoveredCADFiles(
		[]string{"parts/widget.sldprt", "assemblies/top.sldasm", "README.md"},
		"*.sldasm filter=lfs diff=lfs merge=lfs -text\n",
	)

	if len(uncovered) != 1 {
		t.Fatalf("expected 1 uncovered CAD file, got %d", len(uncovered))
	}
	if uncovered[0] != "parts/widget.sldprt" {
		t.Fatalf("unexpected uncovered file %q", uncovered[0])
	}
}
