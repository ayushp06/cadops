package files

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanClassifiesConfiguredCADFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "parts", "gearbox.sldprt"))
	mustWriteFile(t, filepath.Join(root, "exports", "frame.step"))
	mustWriteFile(t, filepath.Join(root, "notes.txt"))
	mustWriteFile(t, filepath.Join(root, ".git", "objects", "ignored.fcstd"))

	entries, err := Scan(root, []string{".sldprt", ".step", ".fcstd"})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}

	expected := []Entry{
		{Path: "exports/frame.step", TypeName: "STEP", RecommendLocking: false},
		{Path: "parts/gearbox.sldprt", TypeName: "SolidWorks Part", RecommendLocking: true},
	}
	if !reflect.DeepEqual(entries, expected) {
		t.Fatalf("unexpected scan results:\nwant: %#v\ngot:  %#v", expected, entries)
	}
}

func TestGroupEntries(t *testing.T) {
	t.Parallel()

	groups := GroupEntries([]Entry{
		{Path: "b/assy.iam", TypeName: "Inventor Assembly", RecommendLocking: true},
		{Path: "a/part.ipt", TypeName: "Inventor Part", RecommendLocking: true},
		{Path: "c/assy2.iam", TypeName: "Inventor Assembly", RecommendLocking: true},
	})

	expected := []Group{
		{
			TypeName: "Inventor Assembly",
			Entries: []Entry{
				{Path: "b/assy.iam", TypeName: "Inventor Assembly", RecommendLocking: true},
				{Path: "c/assy2.iam", TypeName: "Inventor Assembly", RecommendLocking: true},
			},
		},
		{
			TypeName: "Inventor Part",
			Entries: []Entry{
				{Path: "a/part.ipt", TypeName: "Inventor Part", RecommendLocking: true},
			},
		},
	}

	if !reflect.DeepEqual(groups, expected) {
		t.Fatalf("unexpected groups:\nwant: %#v\ngot:  %#v", expected, groups)
	}
}

func mustWriteFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
