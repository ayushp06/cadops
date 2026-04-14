package watch

import (
	"testing"
	"time"
)

func TestDebouncerAllowAt(t *testing.T) {
	t.Parallel()

	debouncer := NewDebouncer(500 * time.Millisecond)
	start := time.Date(2026, time.April, 14, 12, 0, 0, 0, time.UTC)

	if !debouncer.AllowAt("parts/a.sldprt", start) {
		t.Fatal("expected first event to pass")
	}
	if debouncer.AllowAt("parts/a.sldprt", start.Add(300*time.Millisecond)) {
		t.Fatal("expected repeated event inside debounce window to be suppressed")
	}
	if !debouncer.AllowAt("parts/a.sldprt", start.Add(600*time.Millisecond)) {
		t.Fatal("expected event outside debounce window to pass")
	}
}

func TestDebouncerIsPerPath(t *testing.T) {
	t.Parallel()

	debouncer := NewDebouncer(time.Second)
	start := time.Date(2026, time.April, 14, 12, 0, 0, 0, time.UTC)

	if !debouncer.AllowAt("parts/a.sldprt", start) {
		t.Fatal("expected first event to pass")
	}
	if !debouncer.AllowAt("parts/b.sldprt", start.Add(100*time.Millisecond)) {
		t.Fatal("expected different path to pass")
	}
}
