package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cadops/cadops/internal/config"
)

func TestRunConfigGetWorksFromSubdirectory(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)

	if err := config.Save(filepath.Join(root, config.FileName), config.Config{
		Version:           1,
		TrackedExtensions: []string{".sldprt", ".f3d"},
		AutoStage:         false,
		RequireLFS:        true,
		LockingEnabled:    true,
	}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	subdir := filepath.Join(root, "nested")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}

	output := captureStdout(t, func() {
		if err := runConfigGet(subdir, "tracked_extensions"); err != nil {
			t.Fatalf("run config get: %v", err)
		}
	})

	if !strings.Contains(output, ".sldprt, .f3d") {
		t.Fatalf("expected config value from repo root, got:\n%s", output)
	}
}
