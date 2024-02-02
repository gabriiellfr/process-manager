package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"process-manager/internal/manager"
	"syscall"
)

func main() {
	pm := manager.NewProcessManager()

	// Start the manager in the background
	go func() {
		// Add your actual processes here
		// pm.AddProcess("node your-node-app.js")
		// pm.AddProcess("python your-python-script.py")
	}()

	// Handle interrupt signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh // Wait for an interrupt signal

		// Stop all processes on interrupt
		pm.StopAllProcesses()
		os.Exit(0)
	}()

	// Listen for user commands
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a command: ")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "start":
			// Add a new process to the manager and start it (customize this part)
			fmt.Print("Enter the command for the new process: ")
			scanner.Scan()
			newCommand := scanner.Text()
			pm.AddProcess(newCommand)
		case "stop":
			// Stop a specific process (customize this part)
			fmt.Print("Enter the Process ID to stop: ")
			var processID int
			fmt.Scan(&processID)
			pm.StopProcess(processID)
		case "list":
			// List running processes
			runningProcesses := pm.ListProcesses()
			if len(runningProcesses) == 0 {
				fmt.Println("No running processes.")
			} else {
				fmt.Println("================================================================================================")
				fmt.Printf("%-30s | %-6s | %-8s | %-8s | %-15s | %-10s | %-20s\n",
					"Process Name", "ID", "PID", "Status", "Restart Count", "Uptime", "Log Path")
				fmt.Println("================================================================================================")
				for _, p := range runningProcesses {
					fmt.Printf("%-30s | %-6d | %-8d | %-8s | %-13d | %-15s | %-20s\n",
						p.Command, p.ProcessID(), os.Getpid(), "Running", p.RestartCount(), p.Uptime(), fmt.Sprintf("logs/process_%d.log", p.ProcessID()))
				}
				fmt.Println("================================================================================================")
			}
		case "exit":
			// Stop all processes and exit
			pm.StopAllProcesses()
			os.Exit(0)
		default:
			fmt.Println("Unknown command. Valid commands: start, stop, list, exit")
		}
	}
}
