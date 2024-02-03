package process

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"process-manager/internal/logger"
)

// Process represents a managed process.
type Process struct {
	Command          string
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	logger           *logger.Logger
	processID        int
	processName      string
	processStatus    string
	processPID       int
	restartCount     int
	lastUptimeStart  time.Time
	memoryUsageBytes int64
}

// NewProcess creates a new instance of the Process.
func NewProcess(processName string, command string) *Process {
	ctx, cancel := context.WithCancel(context.Background())
	return &Process{
		processName:     processName,
		processStatus:   "Running",
		Command:         command,
		ctx:             ctx,
		cancel:          cancel,
		lastUptimeStart: time.Now(),
	}
}

// SetProcessID sets the process ID.
func (p *Process) SetProcessID(processID int) {
	p.processID = processID
}

// Start starts the process and initializes the logger.
func (p *Process) Start() {

	// Skip starting if the process is already stopped
	if p.processStatus == "Stopped" {
		return
	}

	p.initLogger() // Initialize logger before starting the command

	cmd := exec.CommandContext(p.ctx, "bash", "-c", p.Command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error stdout:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		// Handle error and increment restart count
		p.restartCount++
		p.lastUptimeStart = time.Now()
		return
	}

	// Store the actual PID after starting the process
	p.processPID = cmd.Process.Pid

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		// Forward stdout of the command to the logger
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			p.Log(scanner.Text())
		}

		// Wait for the process to finish
		cmd.Wait()

		// If the process terminated unexpectedly and it's not already stopped, restart it
		if p.processStatus != "Stopped" {
			p.restartCount++
			p.lastUptimeStart = time.Now()
			p.Start()
		}

		// Close the logger only when not restarting
		p.logger.Close()
	}()
}

// initLogger initializes the logger for the process.
func (p *Process) initLogger() {
	p.logger, _ = logger.NewLogger(p.GetProcessID(), p.GetProcessName())
}

// Stop stops the process gracefully.
func (p *Process) Stop() {
	p.cancel()
	p.processStatus = "Stopped"
	p.processPID = 0
	p.wg.Wait()
}

// Log writes a log entry to the process's log file.
func (p *Process) Log(entry string) {
	if p.logger != nil {
		if err := p.logger.Log(entry); err != nil {
			// Handle the error (print to console or log it in another way)
			fmt.Println("Error writing log entry:", err)
		}
	}
}

// GetProcessID returns the process ID.
func (p *Process) GetProcessID() int {
	return p.processID
}

// ProcessID returns a unique identifier for the process.
func (p *Process) GetProcessName() string {
	return p.processName
}

// ProcessID returns a unique identifier for the process.
func (p *Process) GetProcessStatus() string {
	return p.processStatus
}

// ActualPID returns the actual PID of the process.
func (p *Process) GetProcessPID() int {
	return p.processPID
}

// Uptime returns the duration for which the process has been running.
func (p *Process) GetUptime() string {
	if p.processStatus != "Stopped" {
		uptime := time.Since(p.lastUptimeStart)
		return fmt.Sprintf("%.2f", uptime.Seconds())
	}

	return fmt.Sprintf("-")
}

// RestartCount returns the number of times the process has been restarted.
func (p *Process) GetRestartCount() int {
	return p.restartCount
}

// MemoryUsageBytes returns the current memory usage of the process.
func (p *Process) GetMemoryUsageBytes() int64 {
	// Implement logic to get memory usage (placeholder)
	return p.memoryUsageBytes
}
