package cli

import (
	"testing"

	caddiff "github.com/cadops/cadops/internal/diff"
)

func TestFormatDiffReportNoChanges(t *testing.T) {
	t.Parallel()

	got := formatDiffReport(caddiff.Summary{})
	want := "No repository changes\n"
	if got != want {
		t.Fatalf("unexpected no-change output: want %q, got %q", want, got)
	}
}

func TestFormatDiffReportGrouped(t *testing.T) {
	t.Parallel()

	got := formatDiffReport(caddiff.Summary{
		CAD: []caddiff.Entry{
			{Kind: caddiff.KindModified, Path: "bracket.sldprt", IsCAD: true},
			{Kind: caddiff.KindAdded, Path: "export.step", IsCAD: true},
		},
		Other: []caddiff.Entry{
			{Kind: caddiff.KindDeleted, Path: "old-readme.md", IsCAD: false},
		},
	})
	want := "" +
		"CAD changes:\n" +
		"  M bracket.sldprt\n" +
		"  A export.step\n" +
		"Other changes:\n" +
		"  D old-readme.md\n"
	if got != want {
		t.Fatalf("unexpected diff output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
