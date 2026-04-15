package cli

import (
	"testing"

	cadfiles "github.com/cadops/cadops/internal/files"
)

func TestFormatFilesReportEmpty(t *testing.T) {
	t.Parallel()

	got := formatFilesReport(nil)
	want := "No CAD files found for configured extensions\n"
	if got != want {
		t.Fatalf("unexpected empty output:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestFormatFilesReportGrouped(t *testing.T) {
	t.Parallel()

	got := formatFilesReport([]cadfiles.Group{
		{
			TypeName: "STEP",
			Entries: []cadfiles.Entry{
				{Path: "exports/frame.step", TypeName: "STEP", RecommendLocking: false},
			},
		},
		{
			TypeName: "SolidWorks Part",
			Entries: []cadfiles.Entry{
				{Path: "parts/gearbox.sldprt", TypeName: "SolidWorks Part", RecommendLocking: true},
			},
		},
	})
	want := "" +
		"CAD files: 2\n" +
		"STEP\n" +
		"- exports/frame.step | type: STEP | lock: no\n" +
		"SolidWorks Part\n" +
		"- parts/gearbox.sldprt | type: SolidWorks Part | lock: yes\n"
	if got != want {
		t.Fatalf("unexpected grouped output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
