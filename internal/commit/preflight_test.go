package commit

import (
	"reflect"
	"testing"

	"github.com/cadops/cadops/internal/collab"
	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
)

func TestAssessWarningsAndFailures(t *testing.T) {
	t.Parallel()

	cfg := config.Config{LockingEnabled: true}
	report := Assess(
		cfg,
		[]gitx.StatusEntry{
			{Code: "M ", Path: "parts/gearbox.sldprt"},
			{Code: " M", Path: "notes.txt"},
			{Code: "R ", Path: "assy/new.sldasm", OldPath: "assy/old.sldasm"},
		},
		"*.step filter=lfs diff=lfs merge=lfs -text\n",
		map[string]bool{"parts/gearbox.sldprt": true},
	)

	if !report.CanCommit {
		t.Fatal("expected commit to be allowed")
	}
	if !report.HasStagedChanges {
		t.Fatal("expected staged changes to be detected")
	}
	if !report.HasUnstagedChanges {
		t.Fatal("expected unstaged changes to be detected")
	}

	wantWarnings := []collab.Warning{
		{
			Title:   "Unstaged changes",
			Details: "local changes are present and will not be included in this commit",
		},
		{
			Title:   "LFS coverage",
			Details: "changed CAD files are missing matching Git LFS rules: assy/new.sldasm, parts/gearbox.sldprt",
		},
		{
			Title:   "Locking",
			Details: "recommended-lock CAD files are modified without a local Git LFS lock: assy/new.sldasm",
		},
	}
	if !reflect.DeepEqual(report.Warnings, wantWarnings) {
		t.Fatalf("unexpected warnings:\nwant: %#v\ngot:  %#v", wantWarnings, report.Warnings)
	}
}

func TestAssessNothingToCommit(t *testing.T) {
	t.Parallel()

	report := Assess(config.Config{}, []gitx.StatusEntry{
		{Code: " M", Path: "parts/gearbox.sldprt"},
		{Code: "??", Path: "scratch.fcstd"},
	}, "", nil)

	if report.CanCommit {
		t.Fatal("expected commit to be blocked")
	}
	if report.HasStagedChanges {
		t.Fatal("did not expect staged changes")
	}
	if !report.HasUnstagedChanges {
		t.Fatal("expected unstaged changes")
	}
	if len(report.Warnings) != 0 {
		t.Fatalf("expected no warnings when commit is blocked early, got %#v", report.Warnings)
	}
}

func TestFindUncoveredChangedCADFiles(t *testing.T) {
	t.Parallel()

	got := FindUncoveredChangedCADFiles([]gitx.StatusEntry{
		{Code: "M ", Path: "parts/gearbox.sldprt"},
		{Code: "R ", Path: "exports/new.step", OldPath: "exports/old.step"},
		{Code: "A ", Path: "notes.txt"},
	}, "*.step filter=lfs diff=lfs merge=lfs -text\n")
	want := []string{"parts/gearbox.sldprt"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected uncovered files: want %#v, got %#v", want, got)
	}
}

func TestFindRecommendedLockWarnings(t *testing.T) {
	t.Parallel()

	got := FindRecommendedLockWarnings([]gitx.StatusEntry{
		{Code: "M ", Path: "parts/gearbox.sldprt"},
		{Code: "R ", Path: "assy/new.sldasm", OldPath: "assy/old.sldasm"},
		{Code: "A ", Path: "exports/frame.step"},
	}, map[string]bool{"parts/gearbox.sldprt": true})
	want := []string{"assy/new.sldasm"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected lock warnings: want %#v, got %#v", want, got)
	}
}
