package process

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Config hold the manager config
type Config struct {
	MonitCheckInterval time.Duration
}

// Manager is process manager
type Manager struct {
	sync.Mutex

	programs map[string]Program
	config   *Config
	stopped  chan struct{}
}

// NewManager create a new process manager
func NewManager(config *Config) *Manager {
	manager := Manager{
		programs: make(map[string]Program),
		config:   config,
		stopped:  make(chan struct{}),
	}

	return &manager
}

// AddProgram add a new program to process manager
func (manager *Manager) AddProgram(program Program) error {
	program.prepare()

	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.programs[program.Name]; ok {
		return fmt.Errorf("the program with name %s already exist", program.Name)
	}

	manager.programs[program.Name] = program

	return nil
}

// RemoveProgram remove a program from process manager
func (manager *Manager) RemoveProgram(name string) {
	manager.Lock()
	defer manager.Unlock()

	if program, ok := manager.programs[name]; ok {
		for _, process := range program.processes {
			process.Stop()
		}

		delete(manager.programs, name)
	}

	return
}

// RestartProgram restart all instances for a program
func (manager *Manager) RestartProgram(name string) error {
	manager.Lock()
	defer manager.Unlock()

	if program, ok := manager.programs[name]; ok {
		for _, process := range program.processes {
			process.Stop()
			// monit will start the process automatically
		}

		return nil
	}

	return fmt.Errorf("the program with name %s does not exist", name)
}

// Programs return all programs
func (manager *Manager) Programs() map[string]Program {
	return manager.programs
}

// Monit start process manager
func (manager *Manager) Monit(ctx context.Context) <-chan struct{} {
	go func() {
		ticker := time.NewTicker(manager.config.MonitCheckInterval)
		defer func() {
			ticker.Stop()
			manager.stopped <- struct{}{}
		}()

		for {
			select {
			case <-ctx.Done():
				manager.stopAllProcesses()
				return
			case <-ticker.C:
				manager.retryAllProcesses()
			}
		}
	}()

	return manager.stopped
}

// retryAllProcesses retry all processes
func (manager *Manager) retryAllProcesses() {
	manager.Lock()
	defer manager.Unlock()

	for _, program := range manager.programs {
		for _, process := range program.processes {
			if !process.IsRunning() {
				process.Start()
			}
		}
	}
}

// stopAllProcessess stop all processes
func (manager *Manager) stopAllProcesses() {
	manager.Lock()
	defer manager.Unlock()

	stoppeds := make([]<-chan struct{}, 0)
	for _, program := range manager.programs {
		for _, process := range program.processes {
			stoppeds = append(stoppeds, process.Stop())
		}
	}

	for _, stop := range stoppeds {
		<-stop
	}
}
