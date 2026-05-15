package runner_test

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/yourorg/procfile-run/runner"
)

func TestSupervisor_MultipleProcesses_AllExit(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)

	p1 := runner.NewProcess("svc1", "echo svc1", logger)
	p2 := runner.NewProcess("svc2", "echo svc2", logger)

	sv := runner.NewSupervisor([]*runner.Process{p1, p2}, logger)

	done := make(chan struct{})
	go func() {
		sv.Start()
		close(done)
	}()

	select {
	case <-done:
		// both fast processes exited, supervisor cleaned up
	case <-time.After(5 * time.Second):
		t.Fatal("supervisor did not finish within timeout")
	}
}
