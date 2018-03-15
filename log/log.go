package log

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

func Emergency(v ...interface{}) string {
	return GetDefaultModule().Emergency(v...)
}

func Alert(v ...interface{}) string {
	return GetDefaultModule().Alert(v...)
}

func Critical(v ...interface{}) string {
	return GetDefaultModule().Critical(v...)
}

func Error(v ...interface{}) string {
	return GetDefaultModule().Error(v...)
}

func Warning(v ...interface{}) string {
	return GetDefaultModule().Warning(v...)
}

func Notice(v ...interface{}) string {
	return GetDefaultModule().Notice(v...)
}

func Info(v ...interface{}) string {
	return GetDefaultModule().Info(v...)
}

func Debug(v ...interface{}) string {
	return GetDefaultModule().Debug(v...)
}

func Emergencyf(format string, v ...interface{}) string {
	return GetDefaultModule().Emergencyf(format, v...)
}

func Alertf(format string, v ...interface{}) string {
	return GetDefaultModule().Alertf(format, v...)
}

func Criticalf(format string, v ...interface{}) string {
	return GetDefaultModule().Criticalf(format, v...)
}

func Errorf(format string, v ...interface{}) string {
	return GetDefaultModule().Errorf(format, v...)
}

func Warningf(format string, v ...interface{}) string {
	return GetDefaultModule().Warningf(format, v...)
}

func Noticef(format string, v ...interface{}) string {
	return GetDefaultModule().Noticef(format, v...)
}

func Infof(format string, v ...interface{}) string {
	return GetDefaultModule().Infof(format, v...)
}

func Debugf(format string, v ...interface{}) string {
	return GetDefaultModule().Debugf(format, v...)
}
