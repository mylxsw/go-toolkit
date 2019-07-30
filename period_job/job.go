package period_job

import (
	"context"
	"sync"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-toolkit/container"
)

// Job is a interface for a job
type Job interface {
	Handle()
}

// Manager 周期性任务管理器
type Manager struct {
	container *container.Container
	ctx       context.Context

	wg sync.WaitGroup
}

// NewManager 创建一个Manager
func NewManager(ctx context.Context, cc *container.Container) *Manager {
	return &Manager{
		container: cc,
		ctx:       ctx,
	}
}

// Run 启动周期性任务循环
func (jm *Manager) Run(name string, job Job, interval time.Duration) {
	log.Debugf("Job %s 运行中...", name)
	jm.wg.Add(1)

	go func() {
		globalTicker := time.NewTicker(interval)
		defer func() {
			globalTicker.Stop()
			jm.wg.Done()
		}()

		for {
			select {
			case <-globalTicker.C:
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Errorf("Job %s 发生异常：%s", name, err)
						}
					}()
					if err := jm.container.Resolve(job.Handle); err != nil {
						log.Errorf("Job %s 执行失败: %s", name, err)
					}
				}()
			case <-jm.ctx.Done():
				log.Debugf("Job %s 已停止", name)
				return
			}
		}
	}()
}

// Wait 等待所有任务结束
func (jm *Manager) Wait() {
	jm.wg.Wait()
}
