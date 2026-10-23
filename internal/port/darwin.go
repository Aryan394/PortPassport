//go:build darwin

package port

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type DarwinPortInspector struct{}

func (d DarwinPortInspector) Find(port int) ([]PortInfo, error) {
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", port)
	}

	cmd := exec.Command(
		"lsof",
		"-nP",
		"-FpcLPTn",
		fmt.Sprintf("-iTCP:%d", port),
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

		return nil, fmt.Errorf("failed to inspect port: %w", err)
	}

	var results []PortInfo

	var pid int
	var process string
	var user string
	var protocol string
	var address string

	addResult := func() {
		if pid == 0 || address == "" {
			return
		}

		results = append(results, PortInfo{
			Port:     port,
			Protocol: protocol,
			Address:  address,
			PID:      pid,
			Process:  process,
			User:     user,
		})
	}

	for _, line := range strings.Split(string(output), "\n") {
		if len(line) < 2 {
			continue
		}

		field := line[0]
		value := line[1:]

		switch field {
		case 'p':
			addResult()

			pid, _ = strconv.Atoi(value)
			process = ""
			user = ""
			protocol = ""
			address = ""

		case 'c':
			process = value

		case 'L':
			user = value

		case 'P':
			protocol = value

		case 'n':
			address = value
			addResult()
			address = ""
		}
	}

	addResult()

	return results, nil
}
