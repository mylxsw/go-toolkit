package graceful

import (
	"os"
	"os/signal"
	"sync"

	"github.com/mylxsw/go-toolkit/log"
)

var logger = log.Module("toolkit.graceful")

type Graceful struct {
	lock sync.Mutex

	reloadSignals   []os.Signal
	shutdownSignals []os.Signal

	reloadHandlers   []func()
	shutdownHandlers []func()
}

func New(reloadSignals []os.Signal, shutdownSignals []os.Signal) *Graceful {
	return &Graceful{
		reloadSignals:    reloadSignals,
		shutdownSignals:  shutdownSignals,
		reloadHandlers:   make([]func(), 0),
		shutdownHandlers: make([]func(), 0),
	}
}

func (gf *Graceful) RegisterReloadHandler(h func()) {
	gf.lock.Lock()
	defer gf.lock.Unlock()

	gf.reloadHandlers = append(gf.reloadHandlers, h)
}

func (gf *Graceful) RegisterShutdownHandler(h func()) {
	gf.lock.Lock()
	defer gf.lock.Unlock()

	gf.shutdownHandlers = append(gf.shutdownHandlers, h)
}

func (gf *Graceful) Reload() {
	logger.Debug("execute reload...")
	go gf.reload()
}

func (gf *Graceful) Shutdown() {
	logger.Debug("shutdown...")

	if err := gf.SignalSelf(os.Interrupt); err != nil {
		logger.Errorf("shutdown failed: %s", err)
	}
}

func (gf *Graceful) SignalSelf(sig os.Signal) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}

	return p.Signal(sig)
}

func (gf *Graceful) shutdown() {
	gf.lock.Lock()
	defer gf.lock.Unlock()

	for i := len(gf.shutdownHandlers) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("execute shutdown handler failed: %s", err)
				}
			}()
			gf.shutdownHandlers[i]()
		}()
	}
}

func (gf *Graceful) reload() {
	gf.lock.Lock()
	defer gf.lock.Unlock()

	for i := len(gf.reloadHandlers) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("execute reload handler failed: %s", err)
				}
			}()
			gf.reloadHandlers[i]()
		}()
	}
}

func (gf *Graceful) Start() error {
	// 平滑退出
	signalChan := make(chan os.Signal)

	signals := make([]os.Signal, 0)
	signals = append(signals, gf.reloadSignals...)
	signals = append(signals, gf.shutdownSignals...)

	signal.Notify(signalChan, signals...)
	for {
		sig := <-signalChan

		for _, s := range gf.shutdownSignals {
			if s == sig {
				goto FINAL
			}
		}

		for _, s := range gf.reloadSignals {
			if s == sig {
				logger.Debugf("received a reload signal %s", sig.String())
				gf.reload()
				break
			}
		}
	}
FINAL:

	logger.Debug("received a shutdown signal")

	gf.shutdown()

	return nil
}
