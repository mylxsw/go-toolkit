package log

import (
	"fmt"
	"log"
	"strings"
)

var DebugEnabled = false

const (
	LevelEmergency = "EMERGENCY"
	LevelAlert     = "ALERT"
	LevelCritical  = "CRITICAL"
	LevelError     = "ERROR"
	LevelWarning   = "WARNING"
	LevelNotice    = "NOTICE"
	LevelInfo      = "INFO"
	LevelDebug     = "DEBUG"
)

func output(level string, format string, v ...interface{}) {
	message := fmt.Sprintf(fmt.Sprintf("[%s] %s", strings.ToUpper(level), format), v...)
	log.Print(strings.Trim(strings.Replace(message, "\n", "\n	", -1), "\n	"))
}

func Emergency(v ...interface{}) {
	output(LevelEmergency, fmt.Sprint(v...))
}

func Alert(v ...interface{}) {
	output(LevelAlert, fmt.Sprint(v...))
}

func Critical(v ...interface{}) {
	output(LevelCritical, fmt.Sprint(v...))
}

func Error(v ...interface{}) {
	output(LevelError, fmt.Sprint(v...))
}

func Warning(v ...interface{}) {
	output(LevelWarning, fmt.Sprint(v...))
}

func Notice(v ...interface{}) {
	output(LevelNotice, fmt.Sprint(v...))
}

func Info(v ...interface{}) {
	output(LevelInfo, fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	if !DebugEnabled {
		return
	}
	output(LevelDebug, fmt.Sprint(v...))
}

func Emergencyf(format string, v ...interface{}) {
	output(LevelEmergency, format, v...)
}

func Alertf(format string, v ...interface{}) {
	output(LevelAlert, format, v...)
}

func Criticalf(format string, v ...interface{}) {
	output(LevelCritical, format, v...)
}

func Errorf(format string, v ...interface{}) {
	output(LevelError, format, v...)
}

func Warningf(format string, v ...interface{}) {
	output(LevelWarning, format, v...)
}

func Noticef(format string, v ...interface{}) {
	output(LevelNotice, format, v...)
}

func Infof(format string, v ...interface{}) {
	output(LevelInfo, format, v...)
}

func Debugf(format string, v ...interface{}) {
	if !DebugEnabled {
		return
	}
	output(LevelDebug, format, v...)
}
