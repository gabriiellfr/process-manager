package manager

import (
	"process-manager/internal/process"
	"sync"
)

// ProcessManager manages multiple processes.
type ProcessManager struct {
	processes     []*process.Process
	mu            sync.Mutex
	nextProcessID int
}

// NewProcessManager creates a new instance of the ProcessManager.
func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

// AddProcess adds a new process to the manager and starts it.
func (pm *ProcessManager) AddProcess(command string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p := process.NewProcess(command)
	processID := pm.getNextProcessID()
	p.SetProcessID(processID)

	pm.processes = append(pm.processes, p)

	// Start the added process
	go p.Start()
}

// getNextProcessID gets the next available process ID.
func (pm *ProcessManager) getNextProcessID() int {
	id := pm.nextProcessID
	pm.nextProcessID++
	return id
}

// StopProcess stops a specific process gracefully.
func (pm *ProcessManager) StopProcess(processID int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, p := range pm.processes {
		if p.GetProcessID() == processID {
			p.Stop()
			break
		}
	}
}

// StopAllProcesses stops all managed processes gracefully.
func (pm *ProcessManager) StopAllProcesses() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, p := range pm.processes {
		p.Stop()
	}
}

// ListProcesses returns information about the currently running processes.
func (pm *ProcessManager) ListProcesses() []*process.Process {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Return a copy of the processes slice to avoid data race
	return append([]*process.Process{}, pm.processes...)
}
