package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mylxsw/go-toolkit/log"
	process_util "github.com/shirou/gopsutil/process"
)

// Process is a program instance
type Process struct {
	Name       string                // process name
	Command    string                // the command to execute
	Args       []string              // the arguments for command
	User       string                // the user to run the command
	Group      string                // the group to run the command
	uid        string                // the user id to run the commmand
	gid        string                // the group id to run the command
	cancel     context.CancelFunc    // cancel function using to cancel the process
	running    bool                  // running is the state of the process
	startTime  time.Time             // the process start time, should be used with running=true
	triedCount int                   // triedCount is the process failed count
	pid        int                   // the process pid
	proc       *process_util.Process // the gopsutils process instance
}

// NewProcess create a new process
func NewProcess(process Process) *Process {
	process.running = false
	process.triedCount = 0

	if process.User != "" {
		sysUser, err := user.Lookup(process.User)
		if err != nil {
			log.Module("process").Warningf("lookup user %s failed: %s", process.User, err.Error())
		} else {
			process.uid = sysUser.Uid
			process.gid = sysUser.Gid
		}
	}

	if process.Group != "" {
		sysGroup, err := user.LookupGroup(process.Group)
		if err != nil {
			log.Module("process").Warningf("lookup group %s failed: %s", process.Group, err.Error())
		} else {
			process.gid = sysGroup.Gid
		}
	}

	return &process
}

// IsRunning return the process status
func (process *Process) IsRunning() bool {
	return process.running
}

// AliveTime return the process alive time
func (process *Process) AliveTime() float64 {
	if !process.running {
		return 0
	}

	return time.Since(process.startTime).Seconds()
}

// Start start the process
func (process *Process) Start() <-chan struct{} {
	stopped := make(chan struct{})

	go func() {
		process.run()
		stopped <- struct{}{}
	}()

	return stopped
}

func (process *Process) run() {
	if process.running {
		log.Module("process").Warningf("process %s is running, request abort", process.Name)
		return
	}

	process.running = true
	process.startTime = time.Now()
	process.triedCount++

	log.Module("process").Debugf("process %s start...", process.Name)

	defer func() {
		process.running = false
		log.Module("process").Debugf("process %s stoped, alive %.4fs...", process.Name, process.AliveTime())
	}()

	ctx, cancel := context.WithCancel(context.Background())
	process.cancel = cancel

	cmd := exec.CommandContext(ctx, process.Command, process.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if process.uid != "" || process.gid != "" {
		cmd.SysProcAttr.Credential = createCredential(process.uid, process.gid)
	}

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	go consoleLog(process.Name, "debug", &stdoutPipe)
	go consoleLog(process.Name, "error", &stderrPipe)

	if err := cmd.Start(); err != nil {
		log.Module("process").Errorf("command %s start failed: %s", process.Name, err.Error())
		return
	}

	process.pid = cmd.Process.Pid
	proc, err := process_util.NewProcess(int32(process.pid))
	process.proc = proc
	if err != nil {
		log.Module("process").Errorf("can not get process info for %s: %s", process.Name, err.Error())
	}

	if err := cmd.Wait(); err != nil {
		log.Module("process").Errorf("command %s execute failed: %s", process.Name, err.Error())
	}
}

// Stop stop the process
func (process *Process) Stop() <-chan struct{} {
	stopped := make(chan struct{})
	process.cancel()

	go func() {
		for {
			if !process.IsRunning() {
				stopped <- struct{}{}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return stopped
}

// Uptime return the uptime for the process
func (process *Process) Uptime() (uptime time.Time, err error) {
	if process.proc == nil {
		err = fmt.Errorf("can not get process info for process %s", process.Name)
		return
	}

	createTime, err := process.proc.CreateTime()
	if err != nil {
		return
	}

	uptime = time.Unix(0, createTime*int64(time.Millisecond))

	return
}

// Info return the underlying process instance
func (process *Process) Info() (*process_util.Process, error) {
	if process.proc == nil {
		return nil, fmt.Errorf("can not get process info for process %s", process.Name)
	}

	return process.proc, nil
}

func consoleLog(name, logType string, input *io.ReadCloser) error {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				return fmt.Errorf("comand %s output failed: %s", name, err.Error())
			}
			break
		}

		if logType == "error" {
			log.Module("process." + name).Error(strings.TrimRight(line, "\n"))
		} else {
			log.Module("process." + name).Debug(strings.TrimRight(line, "\n"))
		}
	}

	return nil
}

func createCredential(uid, gid string) *syscall.Credential {
	credential := syscall.Credential{}
	if uid != "" {
		uidVal, _ := strconv.Atoi(uid)
		credential.Uid = uint32(uidVal)
	}

	if gid != "" {
		gidVal, _ := strconv.Atoi(gid)
		credential.Gid = uint32(gidVal)
	}

	return &credential
}
