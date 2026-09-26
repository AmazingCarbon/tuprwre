package report

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/c4rb0nx1/tuprwre/internal/event"
)

// errReader fails every Read with a fixed error, simulating a log file that
// becomes unreadable mid-stream.
type errReader struct{ err error }

func (e errReader) Read([]byte) (int, error) { return 0, e.err }

// TestLoadMidStreamReadError proves a reader failure after valid records is
// returned, not swallowed as end-of-stream.
func TestLoadMidStreamReadError(t *testing.T) {
	sentinel := errors.New("disk error")
	line := `{"version":"0","id":"x","time":"2026-09-21T09:00:00Z","session_id":"s","source":"gateway","kind":"tool_call_intent","tool_call_id":"c","tool_name":"Bash","arguments":{"command":"ls"},"complete":true}` + "\n"
	_, _, err := Load(io.MultiReader(strings.NewReader(line), errReader{sentinel}))
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want %v", err, sentinel)
	}
}

// TestBuildSkipsNilProcess proves Build tolerates an unfiltered sensor event
// with no process: it skips and counts it instead of dereferencing nil.
func TestBuildSkipsNilProcess(t *testing.T) {
	rep := Build([]event.Event{{
		Version: event.SchemaVersion, ID: "s1",
		Source: event.SourceSensor, Kind: event.KindExec, Sensor: "fixture",
	}}, LoadStats{}, Options{})
	if rep.Load.InvalidEffects != 1 {
		t.Errorf("invalid effects = %d, want 1", rep.Load.InvalidEffects)
	}
	if len(rep.Sessions) != 0 {
		t.Errorf("sessions = %d, want 0", len(rep.Sessions))
	}
}
