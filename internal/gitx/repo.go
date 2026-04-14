package gitx

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// StatusEntry represents a single line from `git status --porcelain`.
type StatusEntry struct {
	Code string
	Path string
}

// RepoState captures basic repository state used by doctor and init.
type RepoState struct {
	IsRepo        bool
	HasRemote     bool
	TrackedFiles  []string
	StatusEntries []StatusEntry
}

// IsRepo reports whether the working directory is inside a Git repository.
func IsRepo(runner Runner, dir string) bool {
	_, err := runner.Run(dir, "git", "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// HasCommits reports whether the repository has at least one commit.
func HasCommits(runner Runner, dir string) bool {
	_, err := runner.Run(dir, "git", "rev-parse", "--verify", "HEAD")
	return err == nil
}

// InitRepo runs `git init` in the target directory.
func InitRepo(runner Runner, dir string) error {
	_, err := runner.Run(dir, "git", "init")
	return err
}

// StatusPorcelain returns parsed porcelain status entries.
func StatusPorcelain(runner Runner, dir string) ([]StatusEntry, error) {
	result, err := runner.Run(dir, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	return ParseStatusPorcelain(result.Stdout), nil
}

// ParseStatusPorcelain parses the standard two-column porcelain format.
func ParseStatusPorcelain(out string) []StatusEntry {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	entries := make([]StatusEntry, 0, len(lines))
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		code := line[:2]
		path := strings.TrimSpace(line[3:])
		if strings.Contains(path, " -> ") {
			parts := strings.Split(path, " -> ")
			path = strings.TrimSpace(parts[len(parts)-1])
		}
		entries = append(entries, StatusEntry{
			Code: code,
			Path: filepath.ToSlash(path),
		})
	}
	return entries
}

// HasRemote reports whether the repository has any configured remotes.
func HasRemote(runner Runner, dir string) bool {
	result, err := runner.Run(dir, "git", "remote")
	return err == nil && strings.TrimSpace(result.Stdout) != ""
}

// Push runs `git push` in the repository.
func Push(runner Runner, dir string) error {
	_, err := runner.Run(dir, "git", "push")
	return err
}

// Pull runs `git pull` in the repository.
func Pull(runner Runner, dir string) error {
	_, err := runner.Run(dir, "git", "pull")
	return err
}

// RecentHistory returns constrained git log output for recent commits.
func RecentHistory(runner Runner, dir string, limit int) (string, error) {
	result, err := runner.Run(
		dir,
		"git",
		"log",
		"-n", strconv.Itoa(limit),
		"--date=short",
		"--pretty=format:%H%x1f%ad%x1f%s%x1e",
		"--name-only",
	)
	if err != nil {
		return "", err
	}
	return result.Stdout, nil
}

// ListTrackedFiles returns the tracked repository files.
func ListTrackedFiles(runner Runner, dir string) ([]string, error) {
	result, err := runner.Run(dir, "git", "ls-files")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(result.Stdout, "\n")
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		files = append(files, filepath.ToSlash(line))
	}
	sort.Strings(files)
	return files, nil
}

// MergeAttributes appends missing LFS attribute lines without removing unrelated content.
func MergeAttributes(existing string, extensions []string) string {
	normalized := strings.ReplaceAll(existing, "\r\n", "\n")
	normalized = strings.TrimRight(normalized, "\n")
	lines := []string{}
	if normalized != "" {
		lines = strings.Split(normalized, "\n")
	}
	seen := make(map[string]bool, len(lines))
	for _, line := range lines {
		seen[strings.TrimSpace(line)] = true
	}

	for _, extension := range extensions {
		entry := AttributeLine(extension)
		if seen[entry] {
			continue
		}
		lines = append(lines, entry)
	}

	return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
}

// AttributeLine returns the canonical LFS entry for an extension.
func AttributeLine(extension string) string {
	return fmt.Sprintf("*%s filter=lfs diff=lfs merge=lfs -text", extension)
}
