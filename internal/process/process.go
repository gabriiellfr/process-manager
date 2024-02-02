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
	restartCount     int
	lastRestartTime  time.Time
	lastUptimeStart  time.Time
	memoryUsageBytes int64
}

// NewProcess creates a new instance of the Process.
func NewProcess(command string) *Process {
	ctx, cancel := context.WithCancel(context.Background())
	return &Process{
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

// GetProcessID returns the process ID.
func (p *Process) GetProcessID() int {
	return p.processID
}

// Start starts the process and initializes the logger.
// Start starts the process and initializes the logger.
func (p *Process) Start() {
	p.initLogger() // Initialize logger before starting the command

	cmd := exec.CommandContext(p.ctx, "bash", "-c", p.Command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		// Handle error
		return
	}

	if err := cmd.Start(); err != nil {
		// Handle error and increment restart count
		p.restartCount++
		p.lastRestartTime = time.Now()
		return
	}

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

		// If the process terminated unexpectedly, restart it
		p.restartCount++
		p.lastRestartTime = time.Now()
		p.Start()

		// Close the logger only when not restarting
		p.logger.Close()
	}()
}

// initLogger initializes the logger for the process.
func (p *Process) initLogger() {
	p.logger, _ = logger.NewLogger(p.ProcessID())
}

// Stop stops the process gracefully.
func (p *Process) Stop() {
	p.cancel()
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

// ProcessID returns a unique identifier for the process.
func (p *Process) ProcessID() int {
	return p.processID
}

// Uptime returns the duration for which the process has been running.
func (p *Process) Uptime() string {
	uptime := time.Since(p.lastUptimeStart)
	return fmt.Sprintf("%.2f", uptime.Seconds())
}

// RestartCount returns the number of times the process has been restarted.
func (p *Process) RestartCount() int {
	return p.restartCount
}

// MemoryUsageBytes returns the current memory usage of the process.
func (p *Process) MemoryUsageBytes() int64 {
	// Implement logic to get memory usage (placeholder)
	return p.memoryUsageBytes
}
