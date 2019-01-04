// +build !windows

package process

import (
	"os/exec"
	"strconv"
	"syscall"
)

func (process *Process) createCmd() *exec.Cmd {
	cmd := exec.Command(process.GetCommand(), process.GetArgs()...)
	// 这里就不重新创建进程组了，创建新的进程组，在当前程序意外退出后，启动的外部进程时无法自动关闭的
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Setpgid: true,
	// }

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
