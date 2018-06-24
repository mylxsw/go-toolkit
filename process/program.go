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
	ProcNum   int
	processes []*Process
}

func NewProgram(name, command, username string, procNum int) *Program {
	return &Program{
		Name:      name,
		Command:   command,
		User:      username,
		ProcNum:   procNum,
		processes: make([]*Process, 0),
	}
}

func (program *Program) initProcesses(outputFunc OutputFunc) *Program {
	snips := strings.Split(program.Command, " ")
	command, args := snips[0], snips[1:]

	for i := 0; i < program.ProcNum; i++ {
		program.processes = append(program.processes, NewProcess(
			fmt.Sprintf("%s-%d", program.Name, i),
			command,
			args,
			program.User,
		).setOutputFunc(outputFunc))
	}

	return program
}

func (program *Program) Processes() []*Process {
	return program.processes
}
