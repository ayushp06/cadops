package watch

import (
	"errors"
	"testing"
	"time"
)

func TestProcessEventAutoStage(t *testing.T) {
	t.Parallel()

	event := Event{
		Path: "parts/gearbox.sldprt",
		Kind: ChangeModified,
		At:   time.Date(2026, time.April, 14, 12, 0, 0, 0, time.UTC),
	}

	var stagedPath string
	status, err := ProcessEvent(event, true, func(path string) error {
		stagedPath = path
		return nil
	})
	if err != nil {
		t.Fatalf("ProcessEvent() error = %v", err)
	}
	if stagedPath != event.Path {
		t.Fatalf("stage path = %q, want %q", stagedPath, event.Path)
	}
	if !status.Staged {
		t.Fatal("expected status to mark the event as staged")
	}
	if got := status.Line(); got != "modified parts/gearbox.sldprt staged" {
		t.Fatalf("Line() = %q", got)
	}
}

func TestProcessEventWithoutAutoStage(t *testing.T) {
	t.Parallel()

	status, err := ProcessEvent(Event{Path: "parts/gearbox.sldprt", Kind: ChangeCreated}, false, nil)
	if err != nil {
		t.Fatalf("ProcessEvent() error = %v", err)
	}
	if status.Staged {
		t.Fatal("did not expect staging when auto-stage is disabled")
	}
	if got := status.Line(); got != "created parts/gearbox.sldprt" {
		t.Fatalf("Line() = %q", got)
	}
}

func TestProcessEventStageError(t *testing.T) {
	t.Parallel()

	want := errors.New("boom")
	_, err := ProcessEvent(Event{Path: "parts/gearbox.sldprt", Kind: ChangeRemoved}, true, func(string) error {
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("ProcessEvent() error = %v, want %v", err, want)
	}
}
