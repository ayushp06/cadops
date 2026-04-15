package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/metadata"
	"github.com/cadops/cadops/internal/snapshot"
)

func TestRunSnapshotCreatesAndCommitsMetadata(t *testing.T) {
	root := t.TempDir()
	setupSnapshotRepo(t, root)

	partPath := filepath.Join(root, "parts", "gearbox.sldprt")
	writeFile(t, partPath, "version-1")
	commitAll(t, root, "initial import")

	writeFile(t, partPath, "version-2")

	output := captureStdout(t, func() {
		if err := runSnapshot(root, time.Date(2026, time.April, 15, 13, 0, 0, 0, time.UTC)); err != nil {
			t.Fatalf("run snapshot: %v", err)
		}
	})

	if !strings.Contains(output, "Created snapshot commit: snapshot: 2026-04-15 13:00\n") {
		t.Fatalf("unexpected snapshot output:\n%s", output)
	}
	if !strings.Contains(output, "Metadata updated for 1 CAD files\n") {
		t.Fatalf("expected metadata update message, got:\n%s", output)
	}

	manifest, err := metadata.Load(root)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if len(manifest.Records) != 1 {
		t.Fatalf("expected 1 metadata record, got %d", len(manifest.Records))
	}

	runner := gitx.Runner{}
	result, err := runner.Run(root, "git", "show", "--pretty=format:", "--name-only", "HEAD")
	if err != nil {
		t.Fatalf("git show head: %v", err)
	}
	for _, path := range []string{".cadops/metadata/manifest.json", "parts/gearbox.sldprt"} {
		if !strings.Contains(result.Stdout, path) {
			t.Fatalf("expected %s in HEAD files:\n%s", path, result.Stdout)
		}
	}
}

func TestRunSnapshotMetadataFailureWarnsButStillCommits(t *testing.T) {
	root := t.TempDir()
	setupSnapshotRepo(t, root)

	partPath := filepath.Join(root, "parts", "gearbox.sldprt")
	writeFile(t, partPath, "version-1")
	commitAll(t, root, "initial import")
	writeFile(t, partPath, "version-2")

	output := captureStdout(t, func() {
		err := runSnapshotWithMetadataUpdater(root, time.Date(2026, time.April, 15, 13, 0, 0, 0, time.UTC),
			func(root string, extensions []string) (snapshot.MetadataUpdate, error) {
				return snapshot.MetadataUpdate{}, errors.New("simulated metadata failure")
			},
		)
		if err != nil {
			t.Fatalf("run snapshot with warning path: %v", err)
		}
	})

	if !strings.Contains(output, "Warning: metadata update failed: simulated metadata failure\n") {
		t.Fatalf("expected metadata warning, got:\n%s", output)
	}
	if !strings.Contains(output, "Created snapshot commit: snapshot: 2026-04-15 13:00\n") {
		t.Fatalf("expected snapshot commit output, got:\n%s", output)
	}
	if _, err := os.Stat(filepath.Join(root, ".cadops", "metadata", "manifest.json")); !os.IsNotExist(err) {
		t.Fatalf("expected no manifest on metadata failure, got err=%v", err)
	}
}

func TestRunSnapshotNoCADChangesDoesNotInvokeMetadata(t *testing.T) {
	root := t.TempDir()
	setupSnapshotRepo(t, root)
	writeFile(t, filepath.Join(root, "README.md"), "notes")

	called := false
	err := runSnapshotWithMetadataUpdater(root, time.Date(2026, time.April, 15, 13, 0, 0, 0, time.UTC),
		func(root string, extensions []string) (snapshot.MetadataUpdate, error) {
			called = true
			return snapshot.MetadataUpdate{}, nil
		},
	)
	if !errors.Is(err, snapshot.ErrNoRelevantChanges) {
		t.Fatalf("expected no relevant changes error, got %v", err)
	}
	if called {
		t.Fatal("metadata updater should not be called when no CAD files changed")
	}
}

func setupSnapshotRepo(t *testing.T, root string) {
	t.Helper()

	initGitRepo(t, root)
	configureGitIdentity(t, root)
	if err := config.Save(filepath.Join(root, config.FileName), config.Config{
		Version:           1,
		TrackedExtensions: []string{".sldprt"},
		AutoStage:         false,
		RequireLFS:        true,
		LockingEnabled:    true,
	}); err != nil {
		t.Fatalf("save config: %v", err)
	}
}

func configureGitIdentity(t *testing.T, dir string) {
	t.Helper()

	runner := gitx.Runner{}
	for _, args := range [][]string{
		{"config", "user.name", "CadOps Test"},
		{"config", "user.email", "cadops@example.com"},
	} {
		if _, err := runner.Run(dir, "git", args...); err != nil {
			t.Fatalf("git %s: %v", strings.Join(args, " "), err)
		}
	}
}

func commitAll(t *testing.T, dir, message string) {
	t.Helper()

	runner := gitx.Runner{}
	if _, err := runner.Run(dir, "git", "add", "--all"); err != nil {
		t.Fatalf("git add --all: %v", err)
	}
	if _, err := runner.Run(dir, "git", "commit", "-m", message); err != nil {
		t.Fatalf("git commit -m %q: %v", message, err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
