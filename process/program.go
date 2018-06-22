package process

import (
	"fmt"
	"strings"
)

// Program is the program we want to execute
type Program struct {
	Name      string
	Command   string
	User      string
	Group     string
	ProcNum   int
	processes []*Process
}

// prepare prepare a new program
func (program *Program) prepare() {
	snips := strings.Split(program.Command, " ")
	command, args := snips[0], snips[1:]

	program.processes = make([]*Process, program.ProcNum)
	for i := 0; i < program.ProcNum; i++ {
		program.processes[i] = NewProcess(Process{
			Name:    fmt.Sprintf("%s-%d", program.Name, i),
			Command: command,
			Args:    args,
			User:    program.User,
			Group:   program.Group,
		})
	}
}

func (program *Program) inspections() []Inspection {
	inspections := make([]Inspection, len(program.processes))
	for i, process := range program.processes {
		inspections[i] = *NewInspection(process)
	}

	return inspections
}
