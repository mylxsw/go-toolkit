// +build windows

package fpm

import "fmt"

// Reload reload php-fpm process
func (proc *Process) Reload() error {
	return fmt.Errorf("windows not suppert reload signal(USR2)")
}
