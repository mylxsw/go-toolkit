package fpm

import (
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	process := NewProcess(Meta{
		FpmBin:          "/Users/mylxsw/.phpbrew/php/php-7.2.4/sbin/php-fpm",
		PidFile:         "/tmp/fpm.pid",
		ErrorLog:        "/tmp/fpm-error.log",
		SlowLog:         "/tmp/fpm-slow.log",
		Listen:          "/tmp/fpm.sock",
		FpmConfigFile:   "/tmp/fpm.ini",
		PhpConfigFile:   "/Users/mylxsw/.phpbrew/php/php-7.2.4/etc/php.ini",
		PM:              "dynamic",
		MaxChildren:     "10",
		StartServers:    "4",  // fpm启动时进程数目
		MinSpareServers: "4",  // fpm最小空闲进程数数目
		MaxSpareServers: "8",  // fpm最大空闲进程数目
		SlowlogTimeout:  "3s", // fpm慢请求日志超时时间
	})

	process.UpdateConfigFile("/tmp/fpm.ini")

	if err := process.Start(); err != nil {
		t.Fatal("无法启动fpm进程")
	}

	go func() {
		for {
			if err := process.Wait(); err != nil {
				t.Logf("process exited: %s", err.Error())
			}

			t.Logf("process exited without error")

			if err := process.Start(); err != nil {
				t.Fatal("无法启动fpm进程")
			}

			t.Logf("new process started with pid=%d", process.PID())
		}
	}()

	time.Sleep(3 * time.Second)
	process.Kill()

	time.Sleep(3 * time.Second)
	process.Kill()

	time.Sleep(3 * time.Second)
	process.Kill()

	time.Sleep(10 * time.Second)
	process.Kill()
}

func TestCloseExistProcess(t *testing.T) {
	CloseExistProcess("/Users/mylxsw/codes/golang/src/git.yunsom.cn/golang/php-server/tmp/s1/.fpm/php-fpm.pid")
}
