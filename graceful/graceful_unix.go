// +build !windows

package graceful

import (
	"os"
	"syscall"
)

func NewWithDefault() *Graceful {
	return New([]os.Signal{syscall.SIGUSR2}, []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP})
}
