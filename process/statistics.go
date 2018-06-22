package process

import (
	"strings"
	"time"
)

// Inspection contain a process info for human
type Inspection struct {
	Name       string    `json:"name,omitempty"`
	Command    string    `json:"command,omitempty"`
	Args       string    `json:"args,omitempty"`
	Uptime     time.Time `json:"uptime,omitempty"`
	AliveTime  float64   `json:"alive_time,omitempty"`
	Status     string    `json:"status,omitempty"`
	IsRunning  bool      `json:"is_running,omitempty"`
	PID        int       `json:"pid,omitempty"`
	TriedCount int       `json:"tried_count,omitempty"`
	User       string    `json:"user,omitempty"`
}

// NewInspection create a new inspection for a process
func NewInspection(process *Process) *Inspection {
	inspection := Inspection{}
	inspection.PID = process.pid
	inspection.IsRunning = process.IsRunning()
	inspection.AliveTime = process.AliveTime()
	inspection.TriedCount = process.triedCount
	inspection.User = process.User
	inspection.Name = process.Name
	inspection.Command = process.Command
	inspection.Args = strings.Join(process.Args, " ")

	if uptime, err := process.Uptime(); err == nil {
		inspection.Uptime = uptime
	}

	if proc, err := process.Info(); err == nil {
		inspection.Status, _ = proc.Status()
		if username, err := proc.Username(); err == nil {
			inspection.User = username
		}
	}

	return &inspection
}
