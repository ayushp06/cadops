package diff

import (
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/watch"
)

// Kind identifies the repository change type for a diff entry.
type Kind string

const (
	KindAdded    Kind = "A"
	KindModified Kind = "M"
	KindDeleted  Kind = "D"
	KindRenamed  Kind = "R"
	KindCopied   Kind = "C"
)

// Entry is a Git-backed repository change with CAD-aware classification.
type Entry struct {
	Code    string
	Kind    Kind
	Path    string
	OldPath string
	IsCAD   bool
}

// Summary groups changed entries into CAD and non-CAD buckets.
type Summary struct {
	CAD   []Entry
	Other []Entry
}

// BuildEntries converts porcelain status entries into structured diff entries.
func BuildEntries(statusEntries []gitx.StatusEntry, extensions []string) []Entry {
	filter := watch.NewFilter(extensions)
	entries := make([]Entry, 0, len(statusEntries))

	for _, statusEntry := range statusEntries {
		entry := Entry{
			Code:    strings.TrimSpace(statusEntry.Code),
			Kind:    classifyKind(statusEntry.Code),
			Path:    statusEntry.Path,
			OldPath: statusEntry.OldPath,
		}
		entry.IsCAD = filter.Match(entry.Path) || filter.Match(entry.OldPath)
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsCAD != entries[j].IsCAD {
			return entries[i].IsCAD
		}
		if entries[i].Kind != entries[j].Kind {
			return entries[i].Kind < entries[j].Kind
		}
		return displayPath(entries[i]) < displayPath(entries[j])
	})

	return entries
}

// Summarize groups entries into CAD and non-CAD slices in stable order.
func Summarize(entries []Entry) Summary {
	summary := Summary{
		CAD:   make([]Entry, 0),
		Other: make([]Entry, 0),
	}

	for _, entry := range entries {
		if entry.IsCAD {
			summary.CAD = append(summary.CAD, entry)
			continue
		}
		summary.Other = append(summary.Other, entry)
	}

	return summary
}

// DisplayPath returns a human-readable path string for the entry.
func DisplayPath(entry Entry) string {
	return displayPath(entry)
}

func classifyKind(code string) Kind {
	trimmed := strings.TrimSpace(code)
	if trimmed == "??" {
		return KindAdded
	}
	if strings.Contains(trimmed, "R") {
		return KindRenamed
	}
	if strings.Contains(trimmed, "D") {
		return KindDeleted
	}
	if strings.Contains(trimmed, "A") {
		return KindAdded
	}
	if strings.Contains(trimmed, "C") {
		return KindCopied
	}
	return KindModified
}

func displayPath(entry Entry) string {
	if entry.OldPath != "" && entry.OldPath != entry.Path {
		return entry.OldPath + " -> " + entry.Path
	}
	return entry.Path
}
