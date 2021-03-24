package fpm

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/mylxsw/go-toolkit/file"
	"github.com/mylxsw/go-toolkit/pidfile"
	"gopkg.in/ini.v1"
)

// Process PHP-FPM进程管理器
type Process struct {
	Meta Meta
	cmd  *exec.Cmd
}

type LogType string

const (
	LogTypeStdout LogType = "stdout"
	LogTypeStderr LogType = "stderr"
)

// Meta 进程运行配置信息
type Meta struct {
	FpmBin          string
	PidFile         string
	ErrorLog        string
	SlowLog         string
	Listen          string
	FpmConfigFile   string
	PhpConfigFile   string
	User            string
	Group           string
	PM              string // fpm进程管理方式
	MaxChildren     string // fpm最大子进程数目
	StartServers    string // fpm启动时进程数目
	MinSpareServers string // fpm最小空闲进程数数目
	MaxSpareServers string // fpm最大空闲进程数目
	SlowlogTimeout  string // fpm慢请求日志超时时间
	OutputHandler   func(typ LogType, msg string)
}

// NewProcess creates a new process descriptor
func NewProcess(meta Meta) *Process {
	return &Process{
		Meta: meta,
	}
}

// GetProcessMeta get process configuration info
func (proc *Process) GetProcessMeta() Meta {
	return proc.Meta
}

// UpdateConfigFile 更新或创建fpm配置文件
func (proc *Process) UpdateConfigFile(configFile string) {
	f := loadIniFromFile(configFile)

	_, _ = f.Section("global").NewKey("pid", proc.Meta.PidFile)
	_, _ = f.Section("global").NewKey("error_log", proc.Meta.ErrorLog)
	_, _ = f.Section("www").NewKey("listen", proc.Meta.Listen)
	_, _ = f.Section("www").NewKey("slowlog", proc.Meta.SlowLog)
	_, _ = f.Section("www").NewKey("pm", proc.Meta.PM)
	_, _ = f.Section("www").NewKey("pm.max_children", proc.Meta.MaxChildren)
	_, _ = f.Section("www").NewKey("pm.start_servers", proc.Meta.StartServers)
	_, _ = f.Section("www").NewKey("pm.min_spare_servers", proc.Meta.MinSpareServers)
	_, _ = f.Section("www").NewKey("pm.max_spare_servers", proc.Meta.MaxSpareServers)
	_, _ = f.Section("www").NewKey("request_slowlog_timeout", proc.Meta.SlowlogTimeout)

	if proc.Meta.User != "" {
		_, _ = f.Section("www").NewKey("user", proc.Meta.User)
	}
	if proc.Meta.Group != "" {
		_, _ = f.Section("www").NewKey("group", proc.Meta.Group)
	}

	if err := f.SaveTo(configFile); err != nil {
		panic(err)
	}
}

// Start 启动php-fpm主进程
func (proc *Process) Start() (err error) {
	args := []string{
		proc.Meta.FpmBin,
		"--fpm-config",
		proc.Meta.FpmConfigFile,
		"--nodaemonize",
		"--allow-to-run-as-root",
	}
	// look for php.ini file
	if proc.Meta.PhpConfigFile == "" {
		args = append(args, "-n")
	} else {
		args = append(args, "-c", proc.Meta.PhpConfigFile)
	}

	// generate extended information for debugger/profiler
	// args = append(args, "-e")

	proc.cmd = &exec.Cmd{
		Path: proc.Meta.FpmBin,
		Args: args,
	}

	if proc.Meta.OutputHandler != nil {
		stdoutPipe, _ := proc.cmd.StdoutPipe()
		stderrPipe, _ := proc.cmd.StderrPipe()

		go proc.consoleLog(LogTypeStdout, &stdoutPipe)
		go proc.consoleLog(LogTypeStderr, &stderrPipe)
	}

	if err := proc.cmd.Start(); err != nil {
		return err
	}

	// wait until the service is connectable
	// or time out
	select {
	case <-proc.waitConn():
		// do nothing
	case <-time.After(time.Second * 10):
		// wait 10 seconds or timeout
		err = fmt.Errorf("time out")
	}

	return
}

// PID get process pid
func (proc *Process) PID() int {
	return proc.cmd.Process.Pid
}

// waitConn test whether process is ok
func (proc *Process) waitConn() <-chan net.Conn {
	chanConn := make(chan net.Conn)
	go func() {
		for {
			if conn, err := net.Dial(proc.Address()); err != nil {
				time.Sleep(time.Millisecond * 2)
			} else {
				chanConn <- conn
				break
			}
		}
	}()
	return chanConn
}

// Address 返回监听方式和地址
func (proc *Process) Address() (network, address string) {
	reIP := regexp.MustCompile("^(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})\\:(\\d{2,5}$)")
	rePort := regexp.MustCompile("^(\\d+)$")
	switch {
	case reIP.MatchString(proc.Meta.Listen):
		network = "tcp"
		address = proc.Meta.Listen
	case rePort.MatchString(proc.Meta.Listen):
		network = "tcp"
		address = ":" + proc.Meta.Listen
	default:
		network = "unix"
		address = proc.Meta.Listen
	}
	return
}

// Kill kill the process
func (proc *Process) Kill() error {
	if proc.cmd == nil || proc.cmd.Process == nil {
		return fmt.Errorf("fpm process is not initialized")
	}

	return proc.cmd.Process.Signal(os.Interrupt)
}

// Wait wait for process to exit
func (proc *Process) Wait() (err error) {
	return proc.cmd.Wait()
}

func (proc *Process) consoleLog(typ LogType, input *io.ReadCloser) error {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				return fmt.Errorf("process fpm output failed: %s", err.Error())
			}
			break
		}

		proc.Meta.OutputHandler(typ, strings.Trim(line, "\n"))
	}

	return nil
}

// CloseExistProcess close php-fpm process already exist for this project
func CloseExistProcess(pidfileName string) error {
	pid, _ := pidfile.ReadPIDFile(pidfileName)
	if pid > 0 {
		process, _ := os.FindProcess(pid)
		_ = process.Signal(os.Interrupt)
	}

	return nil
}

func loadIniFromFile(configFile string) *ini.File {
	var f *ini.File
	if !file.Exist(configFile) {
		f = ini.Empty()
		_, _ = f.NewSection("global")
		_, _ = f.NewSection("www")
	} else {
		f2, err := ini.Load(configFile)
		if err != nil {
			panic(err)
		}

		f = f2
	}

	return f
}
