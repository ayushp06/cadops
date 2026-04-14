package watch

import "testing"

func TestFilterMatch(t *testing.T) {
	t.Parallel()

	filter := NewFilter([]string{".sldprt", "STEP", " .FCSTD "})

	tests := []struct {
		path string
		want bool
	}{
		{path: "parts/gearbox.SLDPRT", want: true},
		{path: "exports/frame.step", want: true},
		{path: "models/layout.fcstd", want: true},
		{path: "docs/spec.md", want: false},
		{path: "README", want: false},
	}

	for _, test := range tests {
		if got := filter.Match(test.path); got != test.want {
			t.Fatalf("Match(%q) = %v, want %v", test.path, got, test.want)
		}
	}
}
