package watch

import (
	"path/filepath"
	"time"
)

// Debouncer suppresses repeated file events for the same path within a time window.
type Debouncer struct {
	window   time.Duration
	lastSeen map[string]time.Time
}

// NewDebouncer creates a path-based debouncer.
func NewDebouncer(window time.Duration) *Debouncer {
	return &Debouncer{
		window:   window,
		lastSeen: make(map[string]time.Time),
	}
}

// AllowAt reports whether an event at the given time should be emitted.
func (d *Debouncer) AllowAt(path string, at time.Time) bool {
	if d == nil {
		return true
	}

	normalized := filepath.ToSlash(path)
	last, ok := d.lastSeen[normalized]
	if ok && at.Sub(last) < d.window {
		return false
	}

	d.lastSeen[normalized] = at
	return true
}
