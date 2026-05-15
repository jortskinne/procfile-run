package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Process represents a single managed process from the Procfile.
type Process struct {
	Name    string
	Command string
	cmd     *exec.Cmd
	logger  *Logger
}

// NewProcess creates a new Process with the given name, command, and logger.
func NewProcess(name, command string, logger *Logger) *Process {
	return &Process{
		Name:    name,
		Command: command,
		logger:  logger,
	}
}

// Start launches the process and streams its output through the logger.
func (p *Process) Start(wg *sync.WaitGroup, done chan<- string) error {
	parts := strings.Fields(p.Command)
	if len(parts) == 0 {
		return fmt.Errorf("process %q has empty command", p.Name)
	}

	p.cmd = exec.Command(parts[0], parts[1:]...)
	p.cmd.Env = os.Environ()

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe for %q: %w", p.Name, err)
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe for %q: %w", p.Name, err)
	}

	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("start %q: %w", p.Name, err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.logger.Stream(p.Name, stdout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.logger.Stream(p.Name, stderr)
	}()

	go func() {
		p.cmd.Wait()
		done <- p.Name
	}()

	return nil
}

// Stop sends an interrupt signal to the process.
func (p *Process) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	return p.cmd.Process.Signal(os.Interrupt)
}

// Pid returns the OS process ID, or -1 if not started.
func (p *Process) Pid() int {
	if p.cmd == nil || p.cmd.Process == nil {
		return -1
	}
	return p.cmd.Process.Pid
}
