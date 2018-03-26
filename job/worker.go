package job

import (
	"fmt"

	"github.com/jeffail/tunny"
)

// Worker 任务处理worker
type Worker struct {
	pool       *tunny.WorkPool
	numWorkers int
}

// JobHandler 任务处理函数
type JobHandler func(job *Job) *JobOutput

// JobAfterHandler 异步任务处理回调函数
type JobAfterHandler func(output *JobOutput, err error)

// New 创建一个worker
func New(numWorkers int, handler JobHandler) *Worker {

	pool, _ := tunny.CreatePool(numWorkers, func(data interface{}) interface{} {
		return handler(data.(*Job))
	}).Open()

	return &Worker{
		pool:       pool,
		numWorkers: numWorkers,
	}
}

// Send 同步发送一个job到worker执行，并获取执行结果
func (worker *Worker) Send(job *Job) (*JobOutput, error) {
	res, err := worker.pool.SendWork(job)
	if err != nil {
		return nil, fmt.Errorf("worker execute failed: %s", err.Error())
	}

	return res.(*JobOutput), nil
}

// SendAsync 异步发送任务到worker执行
func (worker *Worker) SendAsync(job *Job, after JobAfterHandler) {
	worker.pool.SendWorkAsync(job, func(res interface{}, err error) {
		after(res.(*JobOutput), err)
	})
}
