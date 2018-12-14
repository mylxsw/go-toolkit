// +build windows

package process

import (
	"os/exec"
)

func (process *Process) createCmd() *exec.Cmd {
	cmd := exec.Command(process.GetCommand(), process.GetArgs()...)
	return cmd
}
