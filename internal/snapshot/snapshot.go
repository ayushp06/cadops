package snapshot

import (
	"errors"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/watch"
)

const messageLayout = "2006-01-02 15:04"

// ErrNoRelevantChanges indicates there are no CAD changes to include.
var ErrNoRelevantChanges = errors.New("no changed CAD files to snapshot")

// Plan describes the CAD paths and commit message for a snapshot commit.
type Plan struct {
	Message string
	Paths   []string
}

// BuildMessage returns the default snapshot commit message.
func BuildMessage(at time.Time) string {
	return "snapshot: " + at.Format(messageLayout)
}

// BuildPlan selects CAD paths and generates the default commit message.
func BuildPlan(entries []gitx.StatusEntry, extensions []string, at time.Time) (Plan, error) {
	paths := SelectPaths(entries, extensions)
	if len(paths) == 0 {
		return Plan{}, ErrNoRelevantChanges
	}

	return Plan{
		Message: BuildMessage(at),
		Paths:   paths,
	}, nil
}

// SelectPaths returns the changed CAD paths that should be included in a snapshot.
func SelectPaths(entries []gitx.StatusEntry, extensions []string) []string {
	filter := watch.NewFilter(extensions)
	seen := make(map[string]struct{}, len(entries))
	paths := make([]string, 0, len(entries))

	for _, entry := range entries {
		path := filepath.ToSlash(strings.TrimSpace(entry.Path))
		if path == "" || !filter.Match(path) {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		paths = append(paths, path)
	}

	sort.Strings(paths)
	return paths
}
