// +build !windows

package fpm

import "syscall"

// Reload reload php-fpm process
func (proc *Process) Reload() error {
	return proc.cmd.Process.Signal(syscall.SIGUSR2)
}
