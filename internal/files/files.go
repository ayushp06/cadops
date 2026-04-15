package files

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/watch"
)

// Entry describes a CAD-relevant repository file.
type Entry struct {
	Path             string
	TypeName         string
	RecommendLocking bool
}

// Group collects files that share a CAD type label.
type Group struct {
	TypeName string
	Entries  []Entry
}

// Scan walks the repository tree and returns CAD-relevant files that match
// the configured extensions.
func Scan(root string, extensions []string) ([]Entry, error) {
	filter := watch.NewFilter(extensions)
	entries := make([]Entry, 0)

	err := filepath.WalkDir(root, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := dirEntry.Name()
		if dirEntry.IsDir() {
			if name == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		if !filter.Match(relPath) {
			return nil
		}

		fileType, ok := cad.Lookup(filepath.Ext(relPath))
		entry := Entry{
			Path: relPath,
		}
		if ok {
			entry.TypeName = fileType.Name
			entry.RecommendLocking = fileType.RecommendLocking
		} else {
			entry.TypeName = unknownTypeName(filepath.Ext(relPath))
		}

		entries = append(entries, entry)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].TypeName != entries[j].TypeName {
			return entries[i].TypeName < entries[j].TypeName
		}
		return entries[i].Path < entries[j].Path
	})

	return entries, nil
}

// GroupEntries collects entries by CAD type in stable order.
func GroupEntries(entries []Entry) []Group {
	byType := make(map[string][]Entry, len(entries))
	for _, entry := range entries {
		byType[entry.TypeName] = append(byType[entry.TypeName], entry)
	}

	typeNames := make([]string, 0, len(byType))
	for typeName := range byType {
		typeNames = append(typeNames, typeName)
	}
	sort.Strings(typeNames)

	groups := make([]Group, 0, len(typeNames))
	for _, typeName := range typeNames {
		groupEntries := byType[typeName]
		sort.Slice(groupEntries, func(i, j int) bool {
			return groupEntries[i].Path < groupEntries[j].Path
		})
		groups = append(groups, Group{
			TypeName: typeName,
			Entries:  groupEntries,
		})
	}

	return groups
}

func unknownTypeName(extension string) string {
	normalized := strings.ToLower(strings.TrimSpace(extension))
	if normalized == "" {
		return "Unknown CAD Type"
	}
	return "Unknown CAD Type (" + normalized + ")"
}
