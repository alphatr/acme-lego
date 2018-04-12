package misc

import (
	"fmt"
	"os"
	"os/exec"
)

func NewCmd(cmdStr string) *exec.Cmd {
	return exec.Command("/bin/bash", "-c", cmdStr)
}

func RunCmd(cmdStr string) (string, error) {
	cmd := NewCmd(cmdStr)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("run-command-error: %s")
	}

	return string(output), nil
}

func RunCmdBindTerminal(cmdStr string) {
	cmd := NewCmd(cmdStr)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
