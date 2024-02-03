package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"process-manager/internal/manager"
	"process-manager/internal/utils"
	"syscall"
)

func main() {
	pm := manager.NewProcessManager()

	// Start the manager in the background
	go func() {
		pm.AddProcess("node server", "NODE_ENV=development node /Users/gabrielsantos/Documents/GitHub/react-dashboard-server/src/index.js")
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
			newName, newCommand := utils.PromptAndAddProcess(scanner)

			pm.AddProcess(newName, newCommand)
		case "stop":
			processID := utils.PromptAndStopProcess(scanner)

			pm.StopProcess(processID)
		case "list":
			runningProcesses := pm.ListProcesses()

			utils.ListProcesses(runningProcesses)
		case "exit":
			// Stop all processes and exit
			pm.StopAllProcesses()
			os.Exit(0)
		default:
			fmt.Println("Unknown command. Valid commands: start, stop, list, exit")
		}
	}
}
