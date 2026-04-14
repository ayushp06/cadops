package history

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/cad"
)

const (
	recordSeparator = '\x1e'
	fieldSeparator  = '\x1f'
)

// Entry is a parsed commit plus its CAD-relevant changed files.
type Entry struct {
	Hash      string
	ShortHash string
	Date      string
	Message   string
	CADFiles  []string
}

// ParseLog parses the constrained git log output used by CadOps history.
func ParseLog(out string) []Entry {
	records := strings.Split(out, string(recordSeparator))
	entries := make([]Entry, 0, len(records))

	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		lines := strings.Split(strings.ReplaceAll(record, "\r\n", "\n"), "\n")
		if len(lines) == 0 {
			continue
		}

		meta := strings.Split(lines[0], string(fieldSeparator))
		if len(meta) != 3 {
			continue
		}

		entry := Entry{
			Hash:      meta[0],
			ShortHash: shortHash(meta[0]),
			Date:      meta[1],
			Message:   meta[2],
			CADFiles:  filterCADFiles(lines[1:]),
		}
		entries = append(entries, entry)
	}

	return entries
}

// Format renders entries in a compact terminal-friendly format.
func Format(entries []Entry) string {
	if len(entries) == 0 {
		return "No commits found\n"
	}

	var builder strings.Builder
	for i, entry := range entries {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(fmt.Sprintf("%s  %s  %s\n", entry.ShortHash, entry.Date, entry.Message))
		builder.WriteString("  CAD files:\n")
		if len(entry.CADFiles) == 0 {
			builder.WriteString("    - none\n")
			continue
		}
		for _, path := range entry.CADFiles {
			builder.WriteString("    - ")
			builder.WriteString(path)
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func filterCADFiles(paths []string) []string {
	files := make([]string, 0, len(paths))
	seen := make(map[string]bool, len(paths))

	for _, path := range paths {
		path = filepath.ToSlash(strings.TrimSpace(path))
		if path == "" || !cad.IsCADPath(path) || seen[path] {
			continue
		}
		seen[path] = true
		files = append(files, path)
	}

	return files
}

func shortHash(hash string) string {
	if len(hash) <= 7 {
		return hash
	}
	return hash[:7]
}
