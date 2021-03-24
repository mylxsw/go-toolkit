package fpm

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/mylxsw/asteria/log"
)

// Fpm FPM进程管理对象
type Fpm struct {
	process *Process // phpfpm进程管理器
	started bool     // 标识进程是否已经启动（停止）
	config  *Config  // 进程配置
	lock    sync.Mutex
}

// Config Meta Fpm配置
type Config struct {
	FpmBin        string // phpfpm 可执行文件路径
	PhpConfigFile string // php.ini配置文件
	FpmConfigDir  string // php-fpm.conf配置文件目录
	WorkDir       string // 工作目录
	User          string // 执行用户
	Group         string // 执行用户组
	ErrorLogFile  string // fpm错误日志
	SlowLogFile   string // fpm 慢查询日志
	PIDFile       string // fpm pid file
	SocketFile    string // fpm socket file
	FpmConfigFile string // fpm config file

	PM              string                        // fpm进程管理方式
	MaxChildren     string                        // fpm最大子进程数目
	StartServers    string                        // fpm启动时进程数目
	MinSpareServers string                        // fpm最小空闲进程数数目
	MaxSpareServers string                        // fpm最大空闲进程数目
	SlowlogTimeout  string                        // fpm慢请求日志超时时间
	OutputHandler   func(typ LogType, msg string) // php-fpm 标准输出，标准错误输出处理器
}

// NewFpm 创建一个PFM实例
func NewFpm(config *Config) *Fpm {

	errorLog := config.ErrorLogFile
	if errorLog == "" {
		errorLog = filepath.Join(config.WorkDir, config.FpmConfigDir, "php-fpm.error.log")
	}

	slowLog := config.SlowLogFile
	if slowLog == "" {
		slowLog = filepath.Join(config.WorkDir, config.FpmConfigDir, "php-fpm.slow.log")
	}

	pidFile := config.PIDFile
	if pidFile == "" {
		pidFile = filepath.Join(config.WorkDir, config.FpmConfigDir, "php-fpm.pid")
	}

	socketFile := config.SocketFile
	if socketFile == "" {
		socketFile = filepath.Join(config.WorkDir, config.FpmConfigDir, "php-fpm.sock")
	}

	fpmConfigFile := config.FpmConfigFile
	if fpmConfigFile == "" {
		fpmConfigFile = filepath.Join(config.WorkDir, config.FpmConfigDir, "php-fpm.conf")
	}

	process := NewProcess(Meta{
		FpmBin:          config.FpmBin,
		PidFile:         pidFile,
		ErrorLog:        errorLog,
		Listen:          socketFile,
		FpmConfigFile:   fpmConfigFile,
		SlowLog:         slowLog,
		PhpConfigFile:   config.PhpConfigFile,
		User:            config.User,
		Group:           config.Group,
		PM:              config.PM,
		MaxChildren:     config.MaxChildren,
		StartServers:    config.StartServers,
		MinSpareServers: config.MinSpareServers,
		MaxSpareServers: config.MaxSpareServers,
		SlowlogTimeout:  config.SlowlogTimeout,
		OutputHandler:   config.OutputHandler,
	})

	// 更新/创建配置文件
	process.UpdateConfigFile(process.Meta.FpmConfigFile)

	return &Fpm{
		config:  config,
		process: process,
	}
}

// start 启动fpm master进程
func (fpm *Fpm) start() error {
	// 先尝试关闭已经存在的fpm（当前项目相关的）
	_ = CloseExistProcess(fpm.process.GetProcessMeta().PidFile)

	fpm.lock.Lock()
	defer fpm.lock.Unlock()

	err := fpm.process.Start()
	fpm.started = true
	if err != nil {
		return fmt.Errorf("php-fpm process start failed: %s", err.Error())
	}

	return nil
}

// Loop 循环检测fpm master是否挂掉，挂掉后自动重新启动
func (fpm *Fpm) Loop(ok chan struct{}) error {
	if err := fpm.start(); err != nil {
		return err
	}

	log.Debugf("master process for php-fpm has started with pid=%d", fpm.process.PID())
	ok <- struct{}{}

	for {
		if err := fpm.process.Wait(); err != nil {
			log.Errorf("php-fpm process has stopped: %v", err)
		}

		// 如果进程未启动（已经手动关闭），则退出循环
		if func() bool {
			fpm.lock.Lock()
			defer fpm.lock.Unlock()

			return !fpm.started
		}() {
			break
		}

		// 进程启动状态下，异常退出后重启进程
		if err := fpm.start(); err != nil {
			log.Errorf("php-fpm process start failed: %v", err)
			continue
		}

		log.Debugf("master process for php-fpm has restarted with pid=%d", fpm.process.PID())
	}

	return nil
}

// Reload reload php-fpm process
func (fpm *Fpm) Reload() error {
	log.Debug("reload php-fpm process")
	return fpm.process.Reload()
}

// Kill 停止fpm进程
func (fpm *Fpm) Kill() error {
	fpm.lock.Lock()
	defer fpm.lock.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("kill fpm process failed: %v", err)
		}
	}()

	err := fpm.process.Kill()
	if err != nil {
		log.Warningf("kill fpm process failed: %s", err.Error())
		return err
	}

	fpm.started = false

	return nil
}

// GetNetworkAddress 获取监听的网络类型和地址
func (fpm *Fpm) GetNetworkAddress() (network, address string) {
	fpm.lock.Lock()
	defer fpm.lock.Unlock()

	if !fpm.started {
		return "", ""
	}
	return fpm.process.Address()
}
