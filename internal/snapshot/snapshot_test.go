package snapshot

import (
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
