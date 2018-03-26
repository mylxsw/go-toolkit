package job

import (
	"strconv"
	"strings"
	"time"

	"github.com/docker/distribution/uuid"
)

type Queue struct {
	Name string `json:"name"`
}

// Job 任务对象
type Job struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Queue      Queue      `json:"queue"`
	Command    JobCommand `json:"command"`
	ExecTime   time.Time  `json:"exec_time"`
	RetryCount int        `json:"retry_count"`
	FailedTime time.Time  `json:"failed_time"`
	Status     JobStatus  `json:"status"`
}

// NewJob 创建一个新的Job
func NewJob(ID string, name string, queue Queue, command JobCommand, execTime time.Time) *Job {
	return &Job{
		ID:       ID,
		Name:     name,
		Queue:    queue,
		Command:  command,
		ExecTime: execTime,
		Status:   None,
	}
}

// NewJobAutoID 创建一个新的Job，自动生成ID
func NewJobAutoID(name string, queue Queue, command JobCommand, execTime time.Time) *Job {
	return NewJob(
		uuid.Generate().String(),
		name,
		queue,
		command,
		execTime,
	)
}

// JobStatus Job执行状态
type JobStatus int

const (
	// None 任务刚创建，无状态
	None JobStatus = iota
	// Queuing 任务排队中
	Queuing
	// Ready 任务已就绪
	Ready
	// Running 任务运行中
	Running
	// Pending 任务等待中
	Pending
	// Interrupt 任务已中断
	Interrupt
)

// String 将任务状态转换为字符串表示
func (status JobStatus) String() string {
	switch status {
	case None:
		return "none"
	case Queuing:
		return "queuing"
	case Ready:
		return "ready"
	case Running:
		return "running"
	case Pending:
		return "pending"
	case Interrupt:
		return "interrupt"
	}

	return ""
}

// JobOutput Job输出
type JobOutput struct {
	Job    *Job
	Result map[JobOutputType]string
}

// JobOutputType Job输出类型
type JobOutputType int

const (
	// Stdout 标准输出
	Stdout JobOutputType = iota
	// Stderr 标准错误输出
	Stderr
)

// JobCommand Job对应的要执行的命令
type JobCommand struct {
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}

// String 命令格式化为可执行的字符串下形式
func (cmd JobCommand) String() string {
	return cmd.Name + " " + strings.Join(cmd.GetArgsString(), " ")
}

// GetArgsString 以字符串数组的形式返回参数集合
func (cmd JobCommand) GetArgsString() []string {
	res := make([]string, len(cmd.Args))
	for i, s := range cmd.Args {
		switch s.(type) {
		case string:
			res[i] = s.(string)
		case float64:
			res[i] = strconv.FormatFloat(s.(float64), 'f', -1, 64)
		}
	}

	return res
}
