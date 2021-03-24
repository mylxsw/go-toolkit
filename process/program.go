package process

import (
	"fmt"
	"strings"
)

// Program is the program we want to execute
type Program struct {
	Name      string `json:"name,omitempty"`
	Command   string `json:"command,omitempty"`
	User      string `json:"user,omitempty"`
	ProcNum   int    `json:"proc_num,omitempty"`
	processes []*Process
}

// NewProgram create a new Program
func NewProgram(name, command, username string, procNum int) *Program {
	return &Program{
		Name:      name,
		Command:   command,
		User:      username,
		ProcNum:   procNum,
		processes: make([]*Process, 0),
	}
}

func (program *Program) initProcesses(outputFunc OutputHandler) *Program {
	snips := strings.Split(program.Command, " ")
	command, args := snips[0], snips[1:]

	for i := 0; i < program.ProcNum; i++ {
		program.processes = append(program.processes, NewProcess(
			fmt.Sprintf("%s/%d", program.Name, i),
			command,
			args,
			program.User,
		).setOutputFunc(outputFunc))
	}

	return program
}

// Processes get all processes for the program
func (program *Program) Processes() []*Process {
	return program.processes
}
