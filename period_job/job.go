package period_job

import (
	"context"
	"sync"
	"time"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
)

// Job is a interface for a job
type Job interface {
	Handle()
}

// Manager 周期性任务管理器
type Manager struct {
	container container.Container
	ctx       context.Context
	pauseJobs map[string]bool
	lock      sync.RWMutex

	wg sync.WaitGroup
}

// NewManager 创建一个Manager
func NewManager(ctx context.Context, cc container.Container) *Manager {
	return &Manager{
		container: cc,
		ctx:       ctx,
		pauseJobs: make(map[string]bool),
	}
}

// Paused return whether the named job has been paused
func (jm *Manager) Paused(name string) bool {
	jm.lock.RLock()
	defer jm.lock.RUnlock()

	paused, _ := jm.pauseJobs[name]
	return paused
}

func (jm *Manager) Pause(name string, pause bool) {
	jm.lock.Lock()
	defer jm.lock.Unlock()

	jm.pauseJobs[name] = pause
}

// Run 启动周期性任务循环
func (jm *Manager) Run(name string, job Job, interval time.Duration) {
	log.Debugf("Job %s running...", name)

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
				if jm.Paused(name) {
					continue
				}

				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Errorf("Job %s has some error：%s", name, err)
						}
					}()
					if err := jm.container.Resolve(job.Handle); err != nil {
						log.Errorf("Job %s failed: %s", name, err)
					}
				}()
			case <-jm.ctx.Done():
				log.Debugf("Job %s stopped", name)
				return
			}
		}
	}()
}

// Wait 等待所有任务结束
func (jm *Manager) Wait() {
	jm.wg.Wait()
}
