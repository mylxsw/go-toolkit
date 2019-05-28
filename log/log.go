package log

// SetLevel 设置日志输出级别
func SetLevel(level int) *Logger {
	return GetDefaultModule().SetLevel(level)
}

// SetFormatter 设置日志格式化器
func SetFormatter(formatter Formatter) *Logger {
	return GetDefaultModule().SetFormatter(formatter)
}

// SetWriter 设置日志输出器
func SetWriter(writer Writer) *Logger {
	return GetDefaultModule().SetWriter(writer)
}

func WithContext(context C) *ContextLogger {
	return GetDefaultModule().WithContext(context)
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
