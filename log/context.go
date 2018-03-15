package log

import "fmt"

// ContextLogger 带有上下文信息的日志输出
type ContextLogger struct {
	logger  *Logger
	context map[string]interface{}
}

// Emergency 记录emergency日志
func (logger *ContextLogger) Emergency(v ...interface{}) string {
	return logger.logger.output(LevelEmergency, logger.context, v...)
}

// Alert 记录Alert日志
func (logger *ContextLogger) Alert(v ...interface{}) string {
	return logger.logger.output(LevelAlert, logger.context, v...)
}

// Critical 记录Critical日志
func (logger *ContextLogger) Critical(v ...interface{}) string {
	return logger.logger.output(LevelCritical, logger.context, v...)
}

// Error 记录Error日志
func (logger *ContextLogger) Error(v ...interface{}) string {
	return logger.logger.output(LevelError, logger.context, v...)
}

// Warning 记录Warning日志
func (logger *ContextLogger) Warning(v ...interface{}) string {
	return logger.logger.output(LevelWarning, logger.context, v...)
}

// Notice 记录Notice日志
func (logger *ContextLogger) Notice(v ...interface{}) string {
	return logger.logger.output(LevelNotice, logger.context, v...)
}

// Info 记录Info日志
func (logger *ContextLogger) Info(v ...interface{}) string {
	return logger.logger.output(LevelInfo, logger.context, v...)
}

// Debug 记录Debug日志
func (logger *ContextLogger) Debug(v ...interface{}) string {
	return logger.logger.output(LevelDebug, logger.context, v...)
}

// Emergencyf 记录emergency日志
func (logger *ContextLogger) Emergencyf(format string, v ...interface{}) string {
	return logger.Emergency(fmt.Sprintf(format, v...))
}

// Alertf 记录Alert日志
func (logger *ContextLogger) Alertf(format string, v ...interface{}) string {
	return logger.Alert(fmt.Sprintf(format, v...))
}

// Criticalf 记录critical日志
func (logger *ContextLogger) Criticalf(format string, v ...interface{}) string {
	return logger.Critical(fmt.Sprintf(format, v...))
}

// Errorf 记录error日志
func (logger *ContextLogger) Errorf(format string, v ...interface{}) string {
	return logger.Error(fmt.Sprintf(format, v...))
}

// Warningf 记录warning日志
func (logger *ContextLogger) Warningf(format string, v ...interface{}) string {
	return logger.Warning(fmt.Sprintf(format, v...))
}

// Noticef 记录notice日志
func (logger *ContextLogger) Noticef(format string, v ...interface{}) string {
	return logger.Notice(fmt.Sprintf(format, v...))
}

// Infof 记录info日志
func (logger *ContextLogger) Infof(format string, v ...interface{}) string {
	return logger.Info(fmt.Sprintf(format, v...))
}

// Debugf 记录debug日志
func (logger *ContextLogger) Debugf(format string, v ...interface{}) string {
	return logger.Debug(fmt.Sprintf(format, v...))
}
