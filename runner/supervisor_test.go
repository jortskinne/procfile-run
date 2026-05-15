package runner

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func newSupervisorLogger() *log.Logger {
	return log.New(&bytes.Buffer{}, "", 0)
}

func TestNewSupervisor(t *testing.T) {
	logger := newSupervisorLogger()
	p := NewProcess("web", "echo hello", logger)
	sv := NewSupervisor([]*Process{p}, logger)

	if sv == nil {
		t.Fatal("expected non-nil supervisor")
	}
	if len(sv.processes) != 1 {
		t.Fatalf("expected 1 process, got %d", len(sv.processes))
	}
}

func TestSupervisor_StopAll_NotStarted(t *testing.T) {
	logger := newSupervisorLogger()
	p := NewProcess("web", "echo hello", logger)
	sv := NewSupervisor([]*Process{p}, logger)

	// stopAll on processes that were never started should not panic or error.
	svDone := make(chan struct{})
	go func() {
		sv.stopAll()
		close(svDone)
	}()

	select {
	case <-svDone:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("stopAll timed out")
	}
}

func TestSupervisor_Start_ExitsCleanly(t *testing.T) {
	logger := newSupervisorLogger()
	// Use a fast-exiting command so the supervisor shuts down quickly.
	p := NewProcess("fast", "echo done", logger)
	sv := NewSupervisor([]*Process{p}, logger)

	done := make(chan struct{})
	go func() {
		sv.Start()
		close(done)
	}()

	select {
	case <-done:
		// supervisor exited after process finished
	case <-time.After(5 * time.Second):
		t.Fatal("supervisor did not exit in time")
	}
}
