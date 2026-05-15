package runner

import (
	"sync"
	"testing"
	"time"
)

func newTestLogger(t *testing.T) *Logger {
	t.Helper()
	palette := NewPalette()
	return NewLogger(palette)
}

func TestNewProcess(t *testing.T) {
	logger := newTestLogger(t)
	p := NewProcess("web", "echo hello", logger)

	if p.Name != "web" {
		t.Errorf("expected name %q, got %q", "web", p.Name)
	}
	if p.Command != "echo hello" {
		t.Errorf("expected command %q, got %q", "echo hello", p.Command)
	}
	if p.Pid() != -1 {
		t.Errorf("expected Pid -1 before start, got %d", p.Pid())
	}
}

func TestProcess_Start_Success(t *testing.T) {
	logger := newTestLogger(t)
	p := NewProcess("echo", "echo hello", logger)

	var wg sync.WaitGroup
	done := make(chan string, 1)

	if err := p.Start(&wg, done); err != nil {
		t.Fatalf("unexpected error starting process: %v", err)
	}

	if p.Pid() <= 0 {
		t.Errorf("expected positive PID after start, got %d", p.Pid())
	}

	select {
	case name := <-done:
		if name != "echo" {
			t.Errorf("expected done signal for %q, got %q", "echo", name)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("process did not finish in time")
	}

	wg.Wait()
}

func TestProcess_Start_EmptyCommand(t *testing.T) {
	logger := newTestLogger(t)
	p := NewProcess("bad", "", logger)

	var wg sync.WaitGroup
	done := make(chan string, 1)

	err := p.Start(&wg, done)
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
}

func TestProcess_Stop_NotStarted(t *testing.T) {
	logger := newTestLogger(t)
	p := NewProcess("noop", "sleep 100", logger)

	if err := p.Stop(); err != nil {
		t.Errorf("expected no error stopping unstarted process, got: %v", err)
	}
}
