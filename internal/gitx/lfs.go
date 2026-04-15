package gitx

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

// EnsureLFSInstalled runs `git lfs install` when the environment supports it.
func EnsureLFSInstalled(runner Runner, dir string) error {
	_, err := runner.Run(dir, "git", "lfs", "install")
	return err
}

// ListTrackedPatterns returns the patterns from `git lfs track --list`.
func ListTrackedPatterns(runner Runner, dir string) ([]string, error) {
	result, err := runner.Run(dir, "git", "lfs", "track", "--list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(result.Stdout, "\n")
	patterns := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		patterns = append(patterns, strings.Trim(parts[0], "\""))
	}
	return patterns, nil
}

// HasLFS reports whether Git LFS is available in the current environment.
func HasLFS(runner Runner, dir string) bool {
	_, err := runner.Run(dir, "git", "lfs", "version")
	return err == nil
}

// ListLocalLocks returns the locally cached Git LFS locks by repository path.
func ListLocalLocks(runner Runner, dir string) (map[string]bool, error) {
	result, err := runner.Run(dir, "git", "lfs", "locks", "--local", "--json")
	if err != nil {
		return nil, err
	}

	var payload struct {
		Locks []struct {
			Path string `json:"path"`
		} `json:"locks"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &payload); err != nil {
		return nil, err
	}

	lockedPaths := make(map[string]bool, len(payload.Locks))
	for _, lock := range payload.Locks {
		path := filepath.ToSlash(strings.TrimSpace(lock.Path))
		if path == "" {
			continue
		}
		lockedPaths[path] = true
	}
	return lockedPaths, nil
}

// LockPath acquires a Git LFS lock for the given repository path.
func LockPath(runner Runner, dir, path string) error {
	_, err := runner.Run(dir, "git", "lfs", "lock", filepath.ToSlash(path))
	return err
}

// UnlockPath releases a Git LFS lock for the given repository path.
func UnlockPath(runner Runner, dir, path string) error {
	_, err := runner.Run(dir, "git", "lfs", "unlock", filepath.ToSlash(path))
	return err
}
