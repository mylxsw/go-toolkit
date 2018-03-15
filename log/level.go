package log

import "strings"

// 日志输出级别
const (
	LevelEmergency = 600
	LevelAlert     = 550
	LevelCritical  = 500
	LevelError     = 400
	LevelWarning   = 300
	LevelNotice    = 250
	LevelInfo      = 200
	LevelDebug     = 100
)

// GetLevelName 获取日志级别名称
func GetLevelName(level int) string {
	switch level {
	case LevelEmergency:
		return "EMERGENCY"
	case LevelAlert:
		return "ALERT"
	case LevelCritical:
		return "CRITICAL"
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARNING"
	case LevelNotice:
		return "NOTICE"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	}

	return "UNKNOWN"
}

// GetLevelByName 使用名称获取Level真实的数值
func GetLevelByName(levelName string) int {
	switch strings.ToUpper(levelName) {
	case "EMERGENCY":
		return LevelEmergency
	case "ALERT":
		return LevelAlert
	case "CRITICAL":
		return LevelCritical
	case "ERROR":
		return LevelError
	case "WARNING":
		return LevelWarning
	case "NOTICE":
		return LevelNotice
	case "INFO":
		return LevelInfo
	case "DEBUG":
		return LevelDebug
	}

	return 0
}
