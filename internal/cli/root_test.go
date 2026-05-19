package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommandPrintsBuildInfo(t *testing.T) {
	oldVersion, oldCommit, oldDate := version, commit, date
	version, commit, date = "v1.2.3", "abc1234", "2026-05-18T00:00:00Z"
	t.Cleanup(func() {
		version, commit, date = oldVersion, oldCommit, oldDate
	})

	cmd := NewRootCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute version: %v", err)
	}

	got := output.String()
	if !strings.Contains(got, "v1.2.3 (commit abc1234, built 2026-05-18T00:00:00Z)") {
		t.Fatalf("unexpected version output:\n%s", got)
	}
}
