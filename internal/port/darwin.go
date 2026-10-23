//go:build darwin

package port

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"net"
)

type DarwinPortInspector struct{}


func (d DarwinPortInspector) Find(port int) ([]PortInfo, error) {
	if err := ValidatePort(port); err != nil {
		return nil, err
	}

	return d.inspect(fmt.Sprintf("-iTCP:%d", port))
}

func (d DarwinPortInspector) List() ([]PortInfo, error) {
	return d.inspect("-iTCP")
}

func (d DarwinPortInspector) inspect(filter string) ([]PortInfo, error) {
	cmd := exec.Command(
		"lsof",
		"-nP",
		"-FpcLPTn",
		filter,
		"-sTCP:LISTEN",
	)

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError

		if errors.As(err, &exitErr) &&
			exitErr.ExitCode() == 1 &&
			len(output) == 0 &&
			len(exitErr.Stderr) == 0 {
			return []PortInfo{}, nil
		}

		return nil, fmt.Errorf("failed to inspect ports: %w", err)
	}

	var results []PortInfo
	var pid int
	var process string
	var user string
	var protocol string

	for _, line := range strings.Split(string(output), "\n") {
		if len(line) < 2 {
			continue
		}

		field := line[0]
		value := line[1:]

		switch field {
		case 'p':
			pid, _ = strconv.Atoi(value)
			process = ""
			user = ""
			protocol = ""

		case 'c':
			process = value

		case 'L':
			user = value

		case 'P':
			protocol = value

		case 'n':
			if pid == 0 {
				continue
			}

			results = append(results, PortInfo{
				Port:     extractPort(value),
				Protocol: protocol,
				Address:  value,
				PID:      pid,
				Process:  process,
				User:     user,
			})
		}
	}

	return results, nil
}


func extractPort(address string) int {
	_, portString, err := net.SplitHostPort(address)
	if err != nil {
		return 0
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return 0
	}

	return port
}
