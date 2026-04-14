package watch

import "fmt"

// StageFunc stages a repository path after a watch event.
type StageFunc func(path string) error

// Status captures the user-facing outcome for one CAD file event.
type Status struct {
	Event  Event
	Staged bool
}

// Line formats a concise watch status line.
func (s Status) Line() string {
	line := fmt.Sprintf("%s %s", s.Event.Kind, s.Event.Path)
	if s.Staged {
		line += " staged"
	}
	return line
}

// ProcessEvent handles one filtered filesystem event.
func ProcessEvent(event Event, autoStage bool, stage StageFunc) (Status, error) {
	status := Status{Event: event}
	if !autoStage {
		return status, nil
	}
	if stage == nil {
		return status, fmt.Errorf("auto-stage is enabled but no stage function was provided")
	}
	if err := stage(event.Path); err != nil {
		return status, err
	}
	status.Staged = true
	return status, nil
}
