//go:build darwin

package port

import (
	"fmt"
	"os"
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
