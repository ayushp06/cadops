package gitx

import "testing"

func TestParseStatusPorcelain(t *testing.T) {
	t.Parallel()

	output := " M parts/gearbox.sldprt\nA  notes.txt\nR  old.step -> new.step\n?? scratch.fcstd\n"
	entries := ParseStatusPorcelain(output)

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	if entries[0].Code != " M" || entries[0].Path != "parts/gearbox.sldprt" {
		t.Fatalf("unexpected first entry: %+v", entries[0])
	}
	if entries[2].OldPath != "old.step" {
		t.Fatalf("expected rename source path, got %q", entries[2].OldPath)
	}
	if entries[2].Path != "new.step" {
		t.Fatalf("expected rename target path, got %q", entries[2].Path)
	}
}

func TestMergeAttributes(t *testing.T) {
	t.Parallel()

	existing := "# keep me\n*.png binary\n*.sldprt filter=lfs diff=lfs merge=lfs -text\n"
	merged := MergeAttributes(existing, []string{".sldprt", ".fcstd"})

	expected := "# keep me\n*.png binary\n*.sldprt filter=lfs diff=lfs merge=lfs -text\n*.fcstd filter=lfs diff=lfs merge=lfs -text\n"
	if merged != expected {
		t.Fatalf("unexpected merged attributes:\n%s", merged)
	}
}

func TestParseStatusPorcelainHandlesTrimmedLeadingSpace(t *testing.T) {
	t.Parallel()

	output := "M parts/gearbox.sldprt\n?? notes.txt\n"
	entries := ParseStatusPorcelain(output)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Code != " M" {
		t.Fatalf("expected normalized modified code, got %q", entries[0].Code)
	}
	if entries[0].Path != "parts/gearbox.sldprt" {
		t.Fatalf("expected normalized path, got %q", entries[0].Path)
	}
}

func TestParseStatusPorcelainUnquotesPaths(t *testing.T) {
	t.Parallel()

	output := " M \"Edesign Final Project.SLDASM\"\nR  \"old assembly.sldasm\" -> \"new assembly.sldasm\"\n"
	entries := ParseStatusPorcelain(output)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Path != "Edesign Final Project.SLDASM" {
		t.Fatalf("expected unquoted path, got %q", entries[0].Path)
	}
	if entries[1].OldPath != "old assembly.sldasm" {
		t.Fatalf("expected unquoted old path, got %q", entries[1].OldPath)
	}
	if entries[1].Path != "new assembly.sldasm" {
		t.Fatalf("expected unquoted new path, got %q", entries[1].Path)
	}
}

func TestParseFirstParent(t *testing.T) {
	t.Parallel()

	if got := ParseFirstParent("abc123 def456"); got != "abc123" {
		t.Fatalf("unexpected first parent %q", got)
	}
	if got := ParseFirstParent(""); got != "" {
		t.Fatalf("expected empty parent, got %q", got)
	}
}
