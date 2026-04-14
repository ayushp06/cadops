package config

import (
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	cfg := Default()
	if cfg.Version != 1 {
		t.Fatalf("expected version 1, got %d", cfg.Version)
	}
	if len(cfg.TrackedExtensions) == 0 {
		t.Fatal("expected tracked extensions")
	}
	if !cfg.RequireLFS {
		t.Fatal("expected RequireLFS to default to true")
	}
}

func TestLoadAndSave(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	expected := Default()

	if err := Save(path, expected); err != nil {
		t.Fatalf("save config: %v", err)
	}

	actual, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if actual.Version != expected.Version {
		t.Fatalf("expected version %d, got %d", expected.Version, actual.Version)
	}
	if len(actual.TrackedExtensions) != len(expected.TrackedExtensions) {
		t.Fatalf("expected %d tracked extensions, got %d", len(expected.TrackedExtensions), len(actual.TrackedExtensions))
	}
	if actual.LockingEnabled != expected.LockingEnabled {
		t.Fatalf("expected locking_enabled %v, got %v", expected.LockingEnabled, actual.LockingEnabled)
	}
}
