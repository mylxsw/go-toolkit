package job

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	worker := New(3, func(job *Job) *JobOutput {
		result := make(map[JobOutputType]string)

		result[Stdout] = "stdout string"
		result[Stderr] = "stderr string"

		time.Sleep(20 * time.Millisecond)

		return &JobOutput{
			Job:    job,
			Result: result,
		}
	})

	taskHandle(t, worker, true)
	taskHandle(t, worker, false)
}

func taskHandle(t *testing.T, worker *Worker, async bool) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		jobName := "test-job-" + strconv.Itoa(i)
		go func(jobName string) {
			defer wg.Done()

			job := NewJobAutoID(
				jobName,
				Queue{
					Name: "default",
				},
				JobCommand{
					Name: "ping",
					Args: []interface{}{
						"-c", "baidu.com",
					},
				},
				time.Now(),
			)
			if async {
				worker.SendAsync(job, func(output *JobOutput, err error) {
					currentTime := time.Now().Format("2006-01-02 15:04:05.999")
					if err != nil {
						t.Errorf("%s -> async job [%s] handle error: %s", currentTime, output.Job.Name, err.Error())
					}

					t.Logf("%s -> async job [%s] result: %v", currentTime, output.Job.Name, output.Result)
				})

			} else {
				res, err := worker.Send(job)

				currentTime := time.Now().Format("2006-01-02 15:04:05.999")
				if err != nil {
					t.Errorf("%s -> job [%s] handle error: %s", currentTime, jobName, err.Error())
				}

				t.Logf("%s -> job [%s] result: %v", currentTime, jobName, res.Result)
			}
		}(jobName)
	}

	wg.Wait()

}
