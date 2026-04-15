package cli

import "testing"

func TestNewCommitCmdRequiresMessageFlag(t *testing.T) {
	t.Parallel()

	cmd := newCommitCmd()
	flag := cmd.Flags().Lookup("message")
	if flag == nil {
		t.Fatal("expected message flag")
	}
	if flag.Shorthand != "m" {
		t.Fatalf("expected -m shorthand, got %q", flag.Shorthand)
	}
	if err := cmd.ValidateRequiredFlags(); err == nil {
		t.Fatal("expected message flag to be required")
	}
}
