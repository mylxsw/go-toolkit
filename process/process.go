package process

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mylxsw/asteria/log"
)

// OutputType command output type
type OutputType string

const (
	// LogTypeStderr stderr output
	LogTypeStderr = OutputType("stderr")
	// LogTypeStdout stdout output
	LogTypeStdout = OutputType("stdout")
)

// OutputHandler process output handler
type OutputHandler func(logType OutputType, line string, process *Process)

// Process is a program instance
type Process struct {
	name    string   // process name
	command string   // the command to execute
	args    []string // the arguments for command
	user    string   // the user to run the command
	uid     string   // the user id to run the command
	pid     int

	*exec.Cmd
	stat             chan *Process
	lastAliveTime    time.Duration
	timer            *time.Timer
	lock             sync.Mutex
	outputHandler    OutputHandler
	lastErrorMessage string // last error message
}

// GetPID get process pid
func (process *Process) GetPID() int {
	process.lock.Lock()
	defer process.lock.Unlock()

	return process.pid
}

// SetPID update process pid
func (process *Process) SetPID(pid int) {
	process.lock.Lock()
	defer process.lock.Unlock()

	process.pid = pid
}

// GetName get process name
func (process *Process) GetName() string {
	process.lock.Lock()
	defer process.lock.Unlock()

	return process.name
}

// GetUser get process user
func (process *Process) GetUser() string {
	process.lock.Lock()
	defer process.lock.Unlock()

	return process.user
}

// GetCommand get command
func (process *Process) GetCommand() string {
	process.lock.Lock()
	defer process.lock.Unlock()

	return process.command
}

// GetArgs get command args
func (process *Process) GetArgs() []string {
	process.lock.Lock()
	defer process.lock.Unlock()

	return process.args
}

// GetLastErrorMessage get last error message
func (process *Process) GetLastErrorMessage() string {
	process.lock.Lock()
	defer process.lock.Unlock()

	return process.lastErrorMessage
}

// SetLastErrorMessage update last error message
func (process *Process) SetLastErrorMessage(msg string) {
	process.lock.Lock()
	defer process.lock.Unlock()

	process.lastErrorMessage = msg
}

// NewProcess create a new process
func NewProcess(name string, command string, args []string, username string) *Process {
	process := Process{
		name:    name,
		command: command,
		args:    args,
		user:    username,
		stat:    make(chan *Process),
	}

	// need root privilege to set user or group, because setuid and setgid are privileged calls
	if username != "" {
		sysUser, err := user.Lookup(username)
		if err != nil {
			log.Warningf("lookup user %s failed: %s", username, err.Error())
		} else {
			process.uid = sysUser.Uid
		}
	}

	return &process
}

// setOutputFunc set a function to receive process output
func (process *Process) setOutputFunc(f OutputHandler) *Process {
	process.outputHandler = f

	return process
}

// Start start the process
func (process *Process) start() <-chan *Process {
	go func() {
		startTime := time.Now()

		defer func() {
			process.SetPID(0)
			process.lastAliveTime = time.Now().Sub(startTime)

			log.Warningf("process %s finished", process.name)
			process.stat <- process
		}()

		cmd := process.createCmd()

		if process.outputHandler != nil {
			stdoutPipe, _ := cmd.StdoutPipe()
			stderrPipe, _ := cmd.StderrPipe()

			go process.consoleLog(LogTypeStdout, &stdoutPipe)
			go process.consoleLog(LogTypeStderr, &stderrPipe)
		}

		if err := cmd.Start(); err != nil {
			log.Errorf("process %s start failed: %s", process.name, err.Error())
			process.SetLastErrorMessage(err.Error())
			return
		}

		process.SetPID(cmd.Process.Pid)

		if err := cmd.Wait(); err != nil {
			log.Warningf("process %s stopped with error : %s", process.name, err.Error())
			process.SetLastErrorMessage(err.Error())
		}
	}()

	return process.stat
}

// retryDelayTime get the next retry delay
func (process *Process) retryDelayTime() time.Duration {
	if process.lastAliveTime.Seconds() < 5 {
		return 5 * time.Second
	}

	return 0
}

// RemoveTimer remove the timer
func (process *Process) removeTimer() {
	process.lock.Lock()
	defer process.lock.Unlock()

	if process.timer != nil {
		process.timer.Stop()
		process.timer = nil
	}
}

func (process *Process) stop(timeout time.Duration) {
	process.removeTimer()

	pid := process.GetPID()
	name := process.GetName()

	process.lock.Lock()
	defer process.lock.Unlock()

	if pid <= 0 {
		log.Warningf("process %s with pid=%d is invalid", name, pid)
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Warningf("process %s with pid=%d doesn't exist", name, pid)
		return
	}

	stopped := make(chan interface{})
	go func() {
		proc.Signal(syscall.SIGTERM)
		close(stopped)

		log.Debugf("process %s gracefully stopped", name)
	}()

	select {
	case <-stopped:
		return
	case <-time.After(timeout):
		proc.Signal(syscall.SIGKILL)
		log.Debugf("process %s forced stopped", name)
	}
}

func (process *Process) consoleLog(logType OutputType, input *io.ReadCloser) error {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				return fmt.Errorf("process %s output failed: %s", process.GetName(), err.Error())
			}
			break
		}

		if process.outputHandler != nil {
			process.outputHandler(logType, strings.Trim(line, "\n"), process)
		}
	}

	return nil
}
