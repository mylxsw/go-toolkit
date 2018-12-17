package fpm

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/go-ini/ini"
	"github.com/mylxsw/go-toolkit/file"
	"github.com/mylxsw/go-toolkit/pidfile"
)

// Process PHP-FPM进程管理器
type Process struct {
	Meta Meta
	cmd  *exec.Cmd
}

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

	f.Section("global").NewKey("pid", proc.Meta.PidFile)
	f.Section("global").NewKey("error_log", proc.Meta.ErrorLog)
	f.Section("www").NewKey("listen", proc.Meta.Listen)
	f.Section("www").NewKey("slowlog", proc.Meta.SlowLog)
	f.Section("www").NewKey("pm", proc.Meta.PM)
	f.Section("www").NewKey("pm.max_children", proc.Meta.MaxChildren)
	f.Section("www").NewKey("pm.start_servers", proc.Meta.StartServers)
	f.Section("www").NewKey("pm.min_spare_servers", proc.Meta.MinSpareServers)
	f.Section("www").NewKey("pm.max_spare_servers", proc.Meta.MaxSpareServers)
	f.Section("www").NewKey("request_slowlog_timeout", proc.Meta.SlowlogTimeout)

	if proc.Meta.User != "" {
		f.Section("www").NewKey("user", proc.Meta.User)
	}
	if proc.Meta.Group != "" {
		f.Section("www").NewKey("group", proc.Meta.Group)
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

// CloseExistProcess close php-fpm process already exist for this project
func CloseExistProcess(pidfileName string) error {
	pid, _ := pidfile.ReadPIDFile(pidfileName)
	if pid > 0 {
		process, _ := os.FindProcess(pid)
		process.Signal(os.Interrupt)
	}

	return nil
}

func loadIniFromFile(configFile string) *ini.File {
	var f *ini.File
	if !file.Exist(configFile) {
		f = ini.Empty()
		f.NewSection("global")
		f.NewSection("www")
	} else {
		f, err := ini.Load(configFile)
		if err != nil {
			panic(err)
		}
	}

	return f
}
