package diff

import (
	"reflect"
	"testing"

	"github.com/cadops/cadops/internal/gitx"
)

func TestBuildEntries(t *testing.T) {
	t.Parallel()

	entries := BuildEntries([]gitx.StatusEntry{
		{Code: " M", Path: "parts/gearbox.sldprt"},
		{Code: "A ", Path: "exports/frame.step"},
		{Code: "R ", Path: "assy/new.sldasm", OldPath: "assy/old.sldasm"},
		{Code: " D", Path: "notes.txt"},
		{Code: "??", Path: "docs/todo.md"},
	}, []string{".sldprt", ".sldasm", ".step"})

	expected := []Entry{
		{Code: "A", Kind: KindAdded, Path: "exports/frame.step", IsCAD: true},
		{Code: "M", Kind: KindModified, Path: "parts/gearbox.sldprt", IsCAD: true},
		{Code: "R", Kind: KindRenamed, Path: "assy/new.sldasm", OldPath: "assy/old.sldasm", IsCAD: true},
		{Code: "??", Kind: KindAdded, Path: "docs/todo.md", IsCAD: false},
		{Code: "D", Kind: KindDeleted, Path: "notes.txt", IsCAD: false},
	}

	if !reflect.DeepEqual(entries, expected) {
		t.Fatalf("unexpected entries:\nwant: %#v\ngot:  %#v", expected, entries)
	}
}

func TestSummarize(t *testing.T) {
	t.Parallel()

	summary := Summarize([]Entry{
		{Kind: KindModified, Path: "part.sldprt", IsCAD: true},
		{Kind: KindDeleted, Path: "notes.txt", IsCAD: false},
	})

	if len(summary.CAD) != 1 {
		t.Fatalf("expected 1 CAD entry, got %d", len(summary.CAD))
	}
	if len(summary.Other) != 1 {
		t.Fatalf("expected 1 non-CAD entry, got %d", len(summary.Other))
	}
}

func TestDisplayPathRename(t *testing.T) {
	t.Parallel()

	got := DisplayPath(Entry{Path: "new.step", OldPath: "old.step"})
	want := "old.step -> new.step"
	if got != want {
		t.Fatalf("unexpected display path: want %q, got %q", want, got)
	}
}
