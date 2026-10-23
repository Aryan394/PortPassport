package main

import (
	"bufio"
	"fmt"
	"github.com/Aryan394/PortPassport/internal/port"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		printUsage()
		os.Exit(1)
	}
	killMode := false
	portArg := ""
	if len(os.Args) == 2 {
		portArg = os.Args[1]
	} else {
		if os.Args[1] != "--kill" {
			printUsage()
			os.Exit(1)
		}
		killMode = true
		portArg = os.Args[2]
	}

	portNumber, err := strconv.Atoi(portArg)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		fmt.Println("Error: port must be between 1 and 65535")
		os.Exit(1)
	}

	inspector := port.NewInspector()
	results, err := inspector.Find(portNumber)
	if err != nil {
		fmt.Println("Inspection failed:", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Printf("No TCP listener found on port %d\n", portNumber)
		return
	}

	fmt.Printf("Port %d: IN USE\n\n", portNumber)
	fmt.Printf("%-8s %-18s %-12s %-8s %s\n",
		"PID", "PROCESS", "USER", "PROTOCOL", "ADDRESS")

	for _, result := range results {
		fmt.Printf("%-8d %-18s %-12s %-8s %s\n",
			result.PID,
			result.Process,
			result.User,
			result.Protocol,
			result.Address,
		)
	}
	if killMode {
		terminateProcess(results)
	}
}

func terminateProcess(results []port.PortInfo) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the PID you want to terminate (0 to cancel): ")
	input, err := reader.ReadString('\n')
	if err != nil && len(input) == 0 {
		fmt.Println("Could not read input. Cancelled.")
		return
	}

	pid, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || pid <= 0 {
		fmt.Println("Cancelled.")
		return
	}

	var selected *port.PortInfo

	for i := range results {
		if results[i].PID == pid {
			selected = &results[i]
			break
		}
	}
	if selected == nil {
		fmt.Println("That PID was not found in the port results. Cancelled.")
		return
	}
	if pid == os.Getpid() {
		fmt.Println("Refusing to terminate PortPassport itself.")
		return
	}
	fmt.Printf("Terminate %s (PID %d)? Type 'yes' to confirm: ",
		selected.Process, pid)

	confirmation, err := reader.ReadString('\n')
	if err != nil && len(confirmation) == 0 {
		fmt.Println("Could not read confirmation. Cancelled.")
		return
	}
	if strings.TrimSpace(confirmation) != "yes" {
		fmt.Println("Cancelled.")
		return
	}

	inspector := port.NewInspector()

	currentResults, err := inspector.Find(selected.Port)
	if err != nil {
		fmt.Println("Could not recheck the port:", err)
		return
	}
	stillListening := false
	for _, result := range currentResults {
		if result.PID == pid {
			stillListening = true
			break
		}
	}
	if !stillListening {
		fmt.Println("The process is no longer listed on this port. Cancelled.")
		return
	}
	if err := port.KillProcess(pid); err != nil {
		fmt.Println("Failed to terminate process:", err)
		return
	}

	fmt.Printf("Termination request sent to PID %d.\n", pid)
	fmt.Println("Run PortPassport again to check whether the port is free.")
	fmt.Printf("Termination signal sent to PID %d.\n", pid)
	fmt.Println("Run PortPassport again to verify the port is free.")
	fmt.Printf("Termination signal sent to PID %d.\n", pid)
	fmt.Println("Run PortPassport again to check whether the port is free.")
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  portpassport <port>")
	fmt.Println("  portpassport --kill <port>")
}
