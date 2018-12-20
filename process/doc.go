/*
Package process 是一个进程管理器，类似于supervisor，能够监控进程的运行状态，自动重启失败进程。

	manager := NewManager(3*time.Second, func(logType OutputType, line string, process *Process) {
		color.Green("%s[%s] => %s\n", process.GetName(), logType, line)
	})
	manager.AddProgram("test", "/bin/sleep 1", 5, "")
	manager.AddProgram("prometheus", "/bin/echo Hello", 1, "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	manager.Watch(ctx)
*/
package process
