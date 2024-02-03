package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger represents a logger for a process.
type Logger struct {
	ProcessID int
	File      *os.File
	mutex     sync.Mutex // Add a mutex for synchronization
}

// NewLogger creates a new instance of the Logger.
func NewLogger(processID int, processName string) (*Logger, error) {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating logs directory: %w", err)
	}

	fileName := fmt.Sprintf("logs/process_%s_%d.log", processName, processID)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}

	return &Logger{
		ProcessID: processID,
		File:      file,
	}, nil
}

// Log writes a log entry to the log file.
func (l *Logger) Log(entry string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.File == nil {
		return fmt.Errorf("log file is not open")
	}

	logLine := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), entry)
	_, err := l.File.WriteString(logLine)
	return err
}

// Close closes the log file.
func (l *Logger) Close() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.File != nil {
		l.File.Close()
		l.File = nil // Set File to nil after closing
	}
}
