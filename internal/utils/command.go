package utils

import (
	"os/exec"
)

// ExecuteCommand executes a shell command and returns the output
func ExecuteCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
