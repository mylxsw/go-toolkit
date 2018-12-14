package process

import (
	"context"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestManager(t *testing.T) {

	manager := NewManager(3*time.Second, func(logType OutputType, line string, process *Process) {
		color.Green("%s[%s] => %s\n", process.GetName(), logType, line)
	})
	manager.AddProgram("test", "/bin/sleep 1", 5, "")
	manager.AddProgram("prometheus", "/bin/echo Hello", 1, "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	manager.Watch(ctx)
}

//func printInspections(manager *Manager) {
//	var template = "%-8s %-10s %-10s %-10s %-10s %-10s %-10s %-20s %-20s\n"
//	fmt.Println()
//	fmt.Printf(color.New(color.BgBlue).Sprint(template), "pid", "name", "running", "alive", "status", "user", "tried", "uptime", "command")
//	for _, program := range manager.Programs() {
//		for _, insp := range program.inspections() {
//			pid := "-"
//			if insp.PID > 0 {
//				pid = strconv.Itoa(insp.PID)
//			}
//
//			uptime := "-"
//			if !insp.Uptime.IsZero() {
//				uptime = insp.Uptime.Format("2006-01-02 15:04:05")
//			}
//
//			runningState := strconv.FormatBool(insp.IsRunning)
//			if insp.IsRunning {
//				runningState = color.GreenString("%-10s", "ok")
//			} else {
//				runningState = color.RedString("%-10s", "failed")
//			}
//
//			fmt.Printf(
//				template,
//				pid,
//				strWithDefault(insp.Name, "-"),
//				runningState,
//				strWithDefault(fmt.Sprintf("%.4fs", insp.AliveTime), "-"),
//				strWithDefault(insp.Status, "-"),
//				strWithDefault(insp.User, "-"),
//				strWithDefault(strconv.Itoa(insp.TriedCount), "-"),
//				uptime,
//				strWithDefault(fmt.Sprintf("%s %s", insp.Command, insp.Args), "-"),
//			)
//		}
//	}
//
//	fmt.Println()
//}
//
//func strWithDefault(str string, defaultVal string) string {
//	if str == "" {
//		return defaultVal
//	}
//
//	return str
//}
