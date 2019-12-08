package executor

import (
	"context"
	"sync"
	"testing"
)

func TestExecutor(t *testing.T) {
	testCommand(t, "ifconfig")
	testCommand(t, "ifconfig", "-a")

	testCommandWithOutputChan(t, "ifconfig")
	testCommandWithOutputChan(t, "ifconfig", "-a")

	//testCommand(t, "ping", "-c", "4", "baidu.com")
	//testCommandWithOutputChan(t, "ping", "-c", "2", "yunsom.com")
}

func testCommandWithOutputChan(t *testing.T, executable string, args ...string) {
	cmd := New(executable, args...)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(outputChan <-chan Output) {
		defer wg.Done()
		for out := range outputChan {
			t.Logf("%s -> %s", out.Type.String(), out.Content)
		}
	}(cmd.OpenOutputChan())

	if _, err := cmd.Run(context.TODO()); err != nil {
		t.Errorf("command execute failed: %s", err.Error())
	} else {
		t.Log("test ok")
	}

	wg.Wait()
}

func testCommand(t *testing.T, executable string, args ...string) {
	cmd := New(executable, args...)

	if _, err := cmd.Run(context.TODO()); err != nil {
		t.Errorf("command %s execute failed: %s", executable, err.Error())
	} else {
		t.Log("test ok")
	}
}
