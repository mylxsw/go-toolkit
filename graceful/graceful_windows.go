// +build windows

package graceful

import (
	"os"
)

func NewWithDefault() *Graceful {
	return New([]os.Signal{}, []os.Signal{os.Interrupt})
}
