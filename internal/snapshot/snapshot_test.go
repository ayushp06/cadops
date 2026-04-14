package snapshot

import (
	"errors"
	"testing"
	"time"

	"github.com/cadops/cadops/internal/gitx"
)

func TestBuildMessage(t *testing.T) {
	t.Parallel()

	at := time.Date(2026, time.April, 14, 16, 7, 22, 0, time.UTC)
	got := BuildMessage(at)
	want := "snapshot: 2026-04-14 16:07"
	if got != want {
		t.Fatalf("BuildMessage() = %q, want %q", got, want)
	}
}

func TestBuildPlan(t *testing.T) {
	t.Parallel()

	entries := []gitx.StatusEntry{
		{Code: " M", Path: "parts/gearbox.SLDPRT"},
		{Code: "??", Path: "notes.txt"},
	}
	at := time.Date(2026, time.April, 14, 16, 7, 22, 0, time.UTC)

	plan, err := BuildPlan(entries, []string{".sldprt"}, at)
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	if plan.Message != "snapshot: 2026-04-14 16:07" {
		t.Fatalf("BuildPlan() message = %q", plan.Message)
	}
	if len(plan.Paths) != 1 || plan.Paths[0] != "parts/gearbox.SLDPRT" {
		t.Fatalf("BuildPlan() paths = %#v", plan.Paths)
	}
}

func TestBuildPlanNoRelevantChanges(t *testing.T) {
	t.Parallel()

	_, err := BuildPlan([]gitx.StatusEntry{{Code: " M", Path: "README.md"}}, []string{".sldprt"}, time.Now())
	if !errors.Is(err, ErrNoRelevantChanges) {
		t.Fatalf("BuildPlan() error = %v, want %v", err, ErrNoRelevantChanges)
	}
}

func TestSelectPaths(t *testing.T) {
	t.Parallel()

	entries := []gitx.StatusEntry{
		{Code: " M", Path: "parts/gearbox.SLDPRT"},
		{Code: "A ", Path: "exports/frame.step"},
		{Code: "??", Path: "notes.txt"},
		{Code: "R ", Path: "assy/model.fcstd"},
		{Code: " M", Path: "parts/gearbox.SLDPRT"},
		{Code: "D ", Path: "README.md"},
	}

	got := SelectPaths(entries, []string{".sldprt", ".step", ".fcstd"})
	want := []string{"assy/model.fcstd", "exports/frame.step", "parts/gearbox.SLDPRT"}

	if len(got) != len(want) {
		t.Fatalf("SelectPaths() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SelectPaths()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
