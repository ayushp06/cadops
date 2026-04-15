package cli

import (
	"strings"
	"testing"

	"github.com/cadops/cadops/internal/scan"
)

func TestFormatScanReportNoCADFiles(t *testing.T) {
	t.Parallel()

	got := formatScanReport(scan.Report{})
	if !strings.Contains(got, "No CAD files found for configured extensions\n") {
		t.Fatalf("expected no-CAD-files message, got:\n%s", got)
	}
}

func TestFormatScanReport(t *testing.T) {
	t.Parallel()

	got := formatScanReport(scan.Report{
		TotalFiles:         2,
		ByType:             []scan.TypeCount{{TypeName: "SolidWorks Part", Count: 1}, {TypeName: "STEP", Count: 1}},
		LockingRecommended: []string{"parts/gearbox.sldprt"},
		GitLFSExpected:     []string{"exports/frame.step", "parts/gearbox.sldprt"},
		LFSWarnings:        []scan.LFSWarning{{Path: "exports/frame.step", Extension: ".step"}},
		LargestFiles: []scan.SizedFile{
			{Path: "parts/gearbox.sldprt", SizeBytes: 2048},
			{Path: "exports/frame.step", SizeBytes: 1024},
		},
		TopDirectories: []scan.DirectoryCount{
			{Path: "parts", Count: 1},
			{Path: "exports", Count: 1},
		},
		DisplayedLocking:     []string{"parts/gearbox.sldprt"},
		DisplayedLFSExpected: []string{"exports/frame.step", "parts/gearbox.sldprt"},
		UsedMetadata:         true,
	})

	for _, fragment := range []string{
		"Repo Summary\n",
		"- Total CAD files: 2\n",
		"- Data source: metadata manifest\n",
		"Counts By Type\n",
		"Locking Recommendations\n",
		"Git LFS Expected\n",
		"LFS Status\n",
		"Largest CAD Files\n",
		"- parts/gearbox.sldprt | 2.0 KB\n",
	} {
		if !strings.Contains(got, fragment) {
			t.Fatalf("expected fragment %q in output:\n%s", fragment, got)
		}
	}
}
