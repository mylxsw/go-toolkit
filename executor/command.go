package executor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

const outputChanSize = 1000

// Command 命令行命令
type Command struct {
	Executable string
	Args       []string

	init func(cmd *exec.Cmd) error

	output chan Output
	stdout string
	stderr string

	lock sync.Mutex
}

// Output 命令式输出
type Output struct {
	Type    OutputType
	Content string
}

// OutputType Job输出类型
type OutputType int

const (
	// Stdout 标准输出
	Stdout OutputType = iota
	// Stderr 标准错误输出
	Stderr
)

// String 输出类型的字符串表示
func (outputType OutputType) String() string {
	switch outputType {
	case Stdout:
		return "stdout"
	case Stderr:
		return "stderr"
	}

	return ""
}

// New 创建一个新的命令
func New(executable string, args ...string) *Command {
	return &Command{
		Executable: executable,
		Args:       args,
	}
}

// Init initialize the command
// you can set cmd properties in init callback
// such as
//     cmd.SysProcAttr = &syscall.SysProcAttr{
//	       Setpgid: true,
//     }
func (command *Command) Init(init func(cmd *exec.Cmd) error) {
	command.init = init
}

// Run 执行命令
func (command *Command) Run(ctx context.Context) (bool, error) {
	defer command.close()

	cmd := exec.CommandContext(ctx, command.Executable, command.Args...)
	if command.init != nil {
		if err := command.init(cmd); err != nil {
			return false, err
		}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("can not open stdout pipe: %s", err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return false, fmt.Errorf("can not open stderr pipe: %s", err.Error())
	}

	if err = cmd.Start(); err != nil {
		return false, fmt.Errorf("can not start command: %s", err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_ = command.bindOutputChan(&stdout, Stdout)
	}()

	go func() {
		defer wg.Done()
		_ = command.bindOutputChan(&stderr, Stderr)
	}()

	if err = cmd.Wait(); err != nil {
		return false, fmt.Errorf("wait for command failed: %s", err.Error())
	}

	wg.Wait()

	if cmd.ProcessState.Success() {
		return true, nil
	}

	return false, errors.New("command execute finished, but an error occured")
}

// StdoutString 命令执行后标准输出
func (command *Command) StdoutString() string {
	return command.stdout
}

// StderrString 命令执行后标准错误输出
func (command *Command) StderrString() string {
	return command.stderr
}

// OpenOutputChan 打开输出channel
func (command *Command) OpenOutputChan() <-chan Output {
	command.lock.Lock()
	defer command.lock.Unlock()

	if command.output == nil {
		command.output = make(chan Output, outputChanSize)
	}

	return command.output
}

func (command *Command) close() {
	command.lock.Lock()
	defer command.lock.Unlock()

	if command.output != nil {
		close(command.output)
		command.output = nil
	}
}

func (command *Command) bindOutputChan(input *io.ReadCloser, outputType OutputType) error {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				return fmt.Errorf("read output failed: %s", err.Error())
			}
			break
		}

		if outputType == Stdout {
			command.stdout += line
		} else {
			command.stderr += line
		}

		command.lock.Lock()
		outputIsNotNil := command.output != nil
		command.lock.Unlock()

		if outputIsNotNil {
			command.output <- Output{
				Type:    outputType,
				Content: strings.TrimRight(line, "\n"),
			}
		}
	}

	return nil
}
