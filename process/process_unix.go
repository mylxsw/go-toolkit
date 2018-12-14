// +build !windows

package process

import (
	"os/exec"
	"strconv"
	"syscall"
)

func (process *Process) createCmd() *exec.Cmd {
	cmd := exec.Command(process.GetCommand(), process.GetArgs()...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if process.uid != "" {
		cmd.SysProcAttr.Credential = createCredential(process.uid)
	}

	return cmd
}

func createCredential(uid string) *syscall.Credential {
	credential := syscall.Credential{}
	if uid != "" {
		uidVal, _ := strconv.Atoi(uid)
		credential.Uid = uint32(uidVal)
	}

	return &credential
}
