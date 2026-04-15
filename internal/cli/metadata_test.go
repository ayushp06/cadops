package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/metadata"
)

func TestFormatMetadataGenerateReport(t *testing.T) {
	t.Parallel()

	got := formatMetadataGenerateReport(0)
	want := "Generated metadata for 0 CAD files at .cadops/metadata/manifest.json"
	if got != want {
		t.Fatalf("unexpected generate report:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestFormatMetadataRecord(t *testing.T) {
	t.Parallel()

	got := formatMetadataRecord(metadata.Record{
		Path:               "parts/gearbox.sldprt",
		TypeName:           "SolidWorks Part",
		Extension:          ".sldprt",
		SizeBytes:          42,
		ModifiedTime:       "2026-04-15T16:00:00Z",
		SHA256:             strings.Repeat("a", 64),
		GitLFSExpected:     true,
		LockingRecommended: true,
	})

	if !strings.Contains(got, "Path: parts/gearbox.sldprt\n") {
		t.Fatalf("expected path in output, got:\n%s", got)
	}
	if !strings.Contains(got, "Git LFS Expected: yes\n") {
		t.Fatalf("expected Git LFS line, got:\n%s", got)
	}
	if !strings.Contains(got, "Locking Recommended: yes\n") {
		t.Fatalf("expected locking line, got:\n%s", got)
	}
}

func TestRunMetadataShow(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)

	if err := config.Save(filepath.Join(root, config.FileName), config.Config{
		Version:           1,
		TrackedExtensions: []string{".sldprt"},
		AutoStage:         false,
		RequireLFS:        true,
		LockingEnabled:    true,
	}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	partPath := filepath.Join(root, "parts", "gearbox.sldprt")
	if err := os.MkdirAll(filepath.Dir(partPath), 0o755); err != nil {
		t.Fatalf("mkdir part dir: %v", err)
	}
	if err := os.WriteFile(partPath, []byte("solid-body"), 0o644); err != nil {
		t.Fatalf("write part: %v", err)
	}

	manifest, err := metadata.Generate(root, []string{".sldprt"})
	if err != nil {
		t.Fatalf("generate manifest: %v", err)
	}
	if err := metadata.Save(root, manifest); err != nil {
		t.Fatalf("save manifest: %v", err)
	}

	output := captureStdout(t, func() {
		if err := runMetadataShow(root, "parts/gearbox.sldprt"); err != nil {
			t.Fatalf("run metadata show: %v", err)
		}
	})

	if !strings.Contains(output, "Path: parts/gearbox.sldprt\n") {
		t.Fatalf("expected path in output, got:\n%s", output)
	}
	if !strings.Contains(output, "Type: SolidWorks Part\n") {
		t.Fatalf("expected type in output, got:\n%s", output)
	}
}

func TestRunMetadataShowReturnsErrorWhenRecordMissing(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)

	if err := config.Save(filepath.Join(root, config.FileName), config.Config{
		Version:           1,
		TrackedExtensions: []string{".sldprt"},
		AutoStage:         false,
		RequireLFS:        true,
		LockingEnabled:    true,
	}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	partPath := filepath.Join(root, "parts", "gearbox.sldprt")
	if err := os.MkdirAll(filepath.Dir(partPath), 0o755); err != nil {
		t.Fatalf("mkdir part dir: %v", err)
	}
	if err := os.WriteFile(partPath, []byte("solid-body"), 0o644); err != nil {
		t.Fatalf("write part: %v", err)
	}

	if err := metadata.Save(root, metadata.Manifest{Version: metadata.SchemaVersion}); err != nil {
		t.Fatalf("save empty manifest: %v", err)
	}

	err := runMetadataShow(root, "parts/gearbox.sldprt")
	if err == nil {
		t.Fatal("expected missing metadata error")
	}
	if !strings.Contains(err.Error(), "metadata not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	runner := gitx.Runner{}
	if _, err := runner.Run(dir, "git", "init"); err != nil {
		t.Fatalf("git init: %v", err)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = original
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close reader: %v", err)
	}

	return string(data)
}
