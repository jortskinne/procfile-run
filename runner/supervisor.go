package runner

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Supervisor manages a set of processes and handles graceful shutdown.
type Supervisor struct {
	processes []*Process
	logger    *log.Logger
	wg        sync.WaitGroup
}

// NewSupervisor creates a new Supervisor with the given processes.
func NewSupervisor(processes []*Process, logger *log.Logger) *Supervisor {
	return &Supervisor{
		processes: processes,
		logger:    logger,
	}
}

// Start launches all processes and blocks until they all exit or a signal is received.
func (s *Supervisor) Start() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for _, p := range s.processes {
		s.wg.Add(1)
		go func(proc *Process) {
			defer s.wg.Done()
			if err := proc.Start(); err != nil {
				s.logger.Printf("[supervisor] process %q exited with error: %v", proc.Name, err)
			}
			// When any process exits, trigger shutdown of others.
			quit <- syscall.SIGTERM
		}(p)
	}

	<-quit
	s.logger.Println("[supervisor] shutting down all processes...")
	s.stopAll()
	s.wg.Wait()
	s.logger.Println("[supervisor] all processes stopped")
}

// stopAll sends stop signals to all running processes.
func (s *Supervisor) stopAll() {
	for _, p := range s.processes {
		if err := p.Stop(); err != nil {
			s.logger.Printf("[supervisor] error stopping %q: %v", p.Name, err)
		}
	}
}
