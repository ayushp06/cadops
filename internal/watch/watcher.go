package watch

import (
	"context"
	"io/fs"
	"path/filepath"
	"sort"
	"time"
)

const (
	defaultPollInterval   = time.Second
	defaultDebounceWindow = 750 * time.Millisecond
)

// ChangeKind identifies the filesystem change type.
type ChangeKind string

const (
	ChangeCreated  ChangeKind = "created"
	ChangeModified ChangeKind = "modified"
	ChangeRemoved  ChangeKind = "removed"
)

// Event describes one watched repository file change.
type Event struct {
	Path string
	Kind ChangeKind
	At   time.Time
}

// Handler receives emitted watch events.
type Handler func(Event)

// Watcher polls a repository tree and emits filtered file events.
type Watcher struct {
	root         string
	filter       Filter
	pollInterval time.Duration
	debouncer    *Debouncer
}

type fileState struct {
	modTime time.Time
	size    int64
}

// New creates a repository watcher for the configured extensions.
func New(root string, extensions []string) *Watcher {
	return &Watcher{
		root:         root,
		filter:       NewFilter(extensions),
		pollInterval: defaultPollInterval,
		debouncer:    NewDebouncer(defaultDebounceWindow),
	}
}

// Run blocks until the context is canceled or scanning fails.
func (w *Watcher) Run(ctx context.Context, handler Handler) error {
	previous, err := w.scan()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			current, err := w.scan()
			if err != nil {
				return err
			}
			for _, event := range diff(previous, current, time.Now()) {
				if w.debouncer.AllowAt(event.Path, event.At) {
					handler(event)
				}
			}
			previous = current
		}
	}
}

func (w *Watcher) scan() (map[string]fileState, error) {
	files := make(map[string]fileState)
	err := filepath.WalkDir(w.root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		relative, err := filepath.Rel(w.root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if !w.filter.Match(relative) {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		files[relative] = fileState{
			modTime: info.ModTime(),
			size:    info.Size(),
		}
		return nil
	})
	return files, err
}

func diff(previous, current map[string]fileState, at time.Time) []Event {
	events := make([]Event, 0)

	for path, state := range current {
		prior, ok := previous[path]
		switch {
		case !ok:
			events = append(events, Event{Path: path, Kind: ChangeCreated, At: at})
		case !state.modTime.Equal(prior.modTime) || state.size != prior.size:
			events = append(events, Event{Path: path, Kind: ChangeModified, At: at})
		}
	}

	for path := range previous {
		if _, ok := current[path]; !ok {
			events = append(events, Event{Path: path, Kind: ChangeRemoved, At: at})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].Path == events[j].Path {
			return events[i].Kind < events[j].Kind
		}
		return events[i].Path < events[j].Path
	})
	return events
}
