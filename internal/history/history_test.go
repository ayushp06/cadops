package history

import "testing"

func TestParseLog(t *testing.T) {
	t.Parallel()

	output := "" +
		"abcdef1234567890\x1f2026-04-14\x1fAdded mount redesign\n" +
		"bracket.sldprt\n" +
		"assembly.sldasm\n" +
		"README.md\n\x1e\n" +
		"1234567890abcdef\x1f2026-04-13\x1fUpdated docs\n" +
		"docs/guide.md\n\x1e\n"

	entries := ParseLog(output)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].ShortHash != "abcdef1" {
		t.Fatalf("unexpected short hash %q", entries[0].ShortHash)
	}
	if entries[0].Date != "2026-04-14" || entries[0].Message != "Added mount redesign" {
		t.Fatalf("unexpected first entry metadata: %+v", entries[0])
	}
	if len(entries[0].CADFiles) != 2 {
		t.Fatalf("expected 2 CAD files, got %d", len(entries[0].CADFiles))
	}
	if len(entries[1].CADFiles) != 0 {
		t.Fatalf("expected no CAD files for docs commit, got %d", len(entries[1].CADFiles))
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	formatted := Format([]Entry{
		{
			ShortHash: "abcdef1",
			Date:      "2026-04-14",
			Message:   "Added mount redesign",
			CADFiles:  []string{"bracket.sldprt", "assembly.sldasm"},
		},
		{
			ShortHash: "1234567",
			Date:      "2026-04-13",
			Message:   "Updated docs",
		},
	})

	expected := "" +
		"abcdef1  2026-04-14  Added mount redesign\n" +
		"  CAD files:\n" +
		"    - bracket.sldprt\n" +
		"    - assembly.sldasm\n" +
		"\n" +
		"1234567  2026-04-13  Updated docs\n" +
		"  CAD files:\n" +
		"    - none\n"

	if formatted != expected {
		t.Fatalf("unexpected format:\n%s", formatted)
	}
}

func TestFormatEmpty(t *testing.T) {
	t.Parallel()

	if got := Format(nil); got != "No commits found\n" {
		t.Fatalf("unexpected empty output %q", got)
	}
}
