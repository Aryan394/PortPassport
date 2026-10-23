//go:build darwin

package port

import (
	"fmt"
	"os"
	"time"
)

func KillProcess(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("invalid process ID: %d", pid)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}
	if err := process.Kill(); err != nil {
		return fmt.Errorf("terminate process %d: %w", pid, err)
	}
	return nil
}

func WaitForProcessExit(pid int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		process, err := os.FindProcess(pid)
		if err != nil {
			return nil
		}

		err = process.Signal(os.Signal(nil))
		if err != nil {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("process %d may still be running", pid)
}