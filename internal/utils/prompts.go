package utils

import (
	"bufio"
	"fmt"
	"process-manager/internal/process"
)

// PromptAndAddProcess prompts the user for process information and adds it to the manager.
func PromptAndAddProcess(scanner *bufio.Scanner) (string, string) {
	fmt.Print("Enter the process name: ")
	scanner.Scan()
	newName := scanner.Text()

	fmt.Print("Enter the command for the new process: ")
	scanner.Scan()
	newCommand := scanner.Text()

	return newName, newCommand
}

// PromptAndStopProcess prompts the user for the process ID to stop.
func PromptAndStopProcess(scanner *bufio.Scanner) int {
	var processID int
	fmt.Print("Enter the Process ID to stop: ")
	_, err := fmt.Scan(&processID)

	if err != nil {
		fmt.Println("Error scanning process: ", err)
	}

	return processID
}

// ListProcesses displays the list of running processes.
func ListProcesses(runningProcesses []*process.Process) {
	if len(runningProcesses) == 0 {
		fmt.Println("No running processes.")
	} else {
		fmt.Println("=======================================================================================================")
		fmt.Printf("%-15s | %-6s | %-8s | %-8s | %-15s | %-15s | %-20s\n",
			"Process Name", "ID", "PID", "Status", "Restart Count", "Uptime", "Log Path")
		fmt.Println("=======================================================================================================")
		for _, p := range runningProcesses {
			fmt.Printf("%-15s | %-6d | %-8d | %-8s | %-15d | %-15s | %-20s\n",
				p.GetProcessName(), p.GetProcessID(), p.GetProcessPID(), p.GetProcessStatus(), p.GetRestartCount(), p.GetUptime(), fmt.Sprintf("logs/process_%s_%d.log", p.GetProcessName(), p.GetProcessID()))
		}
		fmt.Println("=======================================================================================================")
	}
}
