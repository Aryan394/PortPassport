package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/Aryan394/PortPassport/internal/port"
)

func main() {
	fmt.Println("PortPassport!")
	if(len(os.Args) < 2) {
		fmt.Println("Usage: PortPassport <port>")
		return
	}

	portNumber, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println("Invalid port:", os.Args[1])
		return
	}
	var inspector port.PortInspector

	inspector = port.MockPortInspector{}

	results, err := inspector.Find(portNumber)

	if err != nil {
		fmt.Println("Error finding port:", err)
		return
	}

	for _, result := range results {
		fmt.Println("Port:", result.Port)
		fmt.Println("Protocol:", result.Protocol)
		fmt.Println("Address:", result.Address)
		fmt.Println("PID:", result.PId)
		fmt.Println("Process:", result.Process)
		fmt.Println("User:", result.User)
	}
}