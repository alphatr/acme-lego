package common

import (
	"os"
	"os/exec"

	"github.com/alphatr/acme-lego/common/errors"
)

// NewBash 创建新的 Bash Shell
func NewBash(command string) *exec.Cmd {
	return exec.Command("/bin/bash", "-c", command)
}

// RunCommand 执行命令
func RunCommand(command string) (string, *errors.Error) {
	cmd := NewBash(command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), errors.NewError(errors.CommonCommandRunErrno, err, command)
	}

	return string(output), nil
}

// RunCommandBindTerminal 在终端中执行命令
func RunCommandBindTerminal(command string) {
	cmd := NewBash(command)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
