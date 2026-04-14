package locking

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveTarget(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	target := filepath.Join(root, "parts", "gearbox.sldprt")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(target, []byte("cad"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := ResolveTarget(root, filepath.Join("parts", "gearbox.sldprt"))
	if err != nil {
		t.Fatalf("ResolveTarget() error = %v", err)
	}
	if got != "parts/gearbox.sldprt" {
		t.Fatalf("ResolveTarget() = %q", got)
	}
}

func TestResolveTargetMissing(t *testing.T) {
	t.Parallel()

	_, err := ResolveTarget(t.TempDir(), "missing.sldprt")
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("ResolveTarget() error = %v", err)
	}
}

func TestAssessTargetWarnings(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitattributes"), []byte("*.step filter=lfs diff=lfs merge=lfs -text\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	assessment, err := AssessTarget(root, "parts/gearbox.sldprt", true)
	if err != nil {
		t.Fatalf("AssessTarget() error = %v", err)
	}
	if !assessment.LockRecommended {
		t.Fatal("expected locking to be recommended")
	}
	if len(assessment.ConfigurationWarn) != 1 {
		t.Fatalf("warnings = %#v", assessment.ConfigurationWarn)
	}
}

func TestAssessTargetWithoutLFS(t *testing.T) {
	t.Parallel()

	assessment, err := AssessTarget(t.TempDir(), "parts/gearbox.sldprt", false)
	if err != nil {
		t.Fatalf("AssessTarget() error = %v", err)
	}
	if len(assessment.ConfigurationWarn) != 1 {
		t.Fatalf("warnings = %#v", assessment.ConfigurationWarn)
	}
}
