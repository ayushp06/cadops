package watch

import (
	"path/filepath"
	"strings"
)

// Filter matches repository paths against configured CAD extensions.
type Filter struct {
	allowed map[string]struct{}
}

// NewFilter builds a case-insensitive extension filter.
func NewFilter(extensions []string) Filter {
	allowed := make(map[string]struct{}, len(extensions))
	for _, extension := range extensions {
		normalized := strings.ToLower(strings.TrimSpace(extension))
		if normalized == "" {
			continue
		}
		if !strings.HasPrefix(normalized, ".") {
			normalized = "." + normalized
		}
		allowed[normalized] = struct{}{}
	}
	return Filter{allowed: allowed}
}

// Match reports whether the path has a configured extension.
func (f Filter) Match(path string) bool {
	_, ok := f.allowed[strings.ToLower(filepath.Ext(path))]
	return ok
}
