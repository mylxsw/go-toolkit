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

	"github.com/mylxsw/go-toolkit/log"
)

// OutputType command output type
type OutputType string

const (
	// LogTypeStderr stderr output
	LogTypeStderr = OutputType("stderr")
	// LogTypeStdout stdout output
	LogTypeStdout = OutputType("stdout")
)

// OutputFunc process output handler
type OutputFunc func(logType OutputType, line string, process *Process)

// Process is a program instance
type Process struct {
	Name    string   // process name
	Command string   // the command to execute
	Args    []string // the arguments for command
	User    string   // the user to run the command
	uid     string   // the user id to run the command
	PID     int

	*exec.Cmd
	stat          chan *Process
	lastAliveTime time.Duration
	timer         *time.Timer
	lock          sync.Mutex
	logFunc       OutputFunc
	LastErrorMsg  string // last error message
}

// NewProcess create a new process
func NewProcess(name string, command string, args []string, username string) *Process {
	process := Process{
		Name:    name,
		Command: command,
		Args:    args,
		User:    username,
		stat:    make(chan *Process),
	}

	// need root privilege to set user or group, because setuid and setgid are privileged calls
	if username != "" {
		sysUser, err := user.Lookup(username)
		if err != nil {
			log.Module("process").Warningf("lookup user %s failed: %s", username, err.Error())
		} else {
			process.uid = sysUser.Uid
		}
	}

	return &process
}

// setOutputFunc set a function to receive process output
func (process *Process) setOutputFunc(f OutputFunc) *Process {
	process.logFunc = f

	return process
}

// Start start the process
func (process *Process) start() <-chan *Process {
	go func() {
		startTime := time.Now()

		defer func() {
			process.PID = 0
			process.lastAliveTime = time.Now().Sub(startTime)
			log.Module("process").Warningf("process %s finished", process.Name)
			process.stat <- process
		}()

		cmd := process.createCmd()

		stdoutPipe, _ := cmd.StdoutPipe()
		stderrPipe, _ := cmd.StderrPipe()

		go process.consoleLog(LogTypeStdout, &stdoutPipe)
		go process.consoleLog(LogTypeStderr, &stderrPipe)

		if err := cmd.Start(); err != nil {
			log.Module("process").Errorf("process %s start failed: %s", process.Name, err.Error())
			process.LastErrorMsg = err.Error()
			return
		}

		process.lock.Lock()
		process.PID = cmd.Process.Pid
		process.lock.Unlock()

		if err := cmd.Wait(); err != nil {
			log.Module("process").Warningf("process %s stopped with error : %s", process.Name, err.Error())
			process.LastErrorMsg = err.Error()
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

	process.lock.Lock()
	defer process.lock.Unlock()

	if process.PID <= 0 {
		return
	}

	proc, err := os.FindProcess(process.PID)
	if err != nil {
		log.Module("process").Warningf("process %s with pid=%d doesn't exist", process.Name, process.PID)
		return
	}

	stopped := make(chan interface{})
	go func() {
		proc.Signal(syscall.SIGTERM)
		close(stopped)

		log.Module("process").Debugf("process %s gracefully stopped", process.Name)
	}()

	select {
	case <-stopped:
		return
	case <-time.After(timeout):
		proc.Signal(syscall.SIGKILL)
		log.Module("process").Debugf("process %s forced stopped", process.Name)
	}
}

func (process *Process) consoleLog(logType OutputType, input *io.ReadCloser) error {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				return fmt.Errorf("process %s output failed: %s", process.Name, err.Error())
			}
			break
		}

		if process.logFunc != nil {
			process.logFunc(logType, strings.Trim(line, "\n"), process)
		}
	}

	return nil
}
