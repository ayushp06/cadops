package config

import "testing"

func TestLookup(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Version:           2,
		TrackedExtensions: []string{".step", ".fcstd"},
		AutoStage:         true,
		RequireLFS:        false,
		LockingEnabled:    true,
	}

	tests := []struct {
		key      string
		want     string
		wantList []string
	}{
		{key: "version", want: "2"},
		{key: "tracked_extensions", wantList: []string{".step", ".fcstd"}},
		{key: "auto_stage", want: "true"},
		{key: "require_lfs", want: "false"},
		{key: "locking_enabled", want: "true"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.key, func(t *testing.T) {
			t.Parallel()

			value, err := Lookup(cfg, test.key)
			if err != nil {
				t.Fatalf("lookup %s: %v", test.key, err)
			}

			if test.want != "" && value.Scalar != test.want {
				t.Fatalf("expected %q, got %q", test.want, value.Scalar)
			}
			if len(test.wantList) != len(value.List) {
				t.Fatalf("expected list length %d, got %d", len(test.wantList), len(value.List))
			}
			for i := range test.wantList {
				if value.List[i] != test.wantList[i] {
					t.Fatalf("expected list item %q at %d, got %q", test.wantList[i], i, value.List[i])
				}
			}
		})
	}
}

func TestLookupUnknownKey(t *testing.T) {
	t.Parallel()

	if _, err := Lookup(Default(), "missing"); err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestFormatValue(t *testing.T) {
	t.Parallel()

	scalar := FormatValue(Value{Scalar: "true"})
	if scalar != "true" {
		t.Fatalf("expected scalar output, got %q", scalar)
	}

	list := FormatValue(Value{List: []string{".step", ".fcstd"}})
	if list != ".step, .fcstd" {
		t.Fatalf("expected joined list, got %q", list)
	}
}
