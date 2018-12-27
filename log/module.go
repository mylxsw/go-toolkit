package log

import (
	"fmt"
	"sync"
	"time"
)

// Logger 日志对象
type Logger struct {
	moduleName   string
	level        int
	formatter    Formatter
	writer       Writer
	timeLocation func() *time.Location
}

var loggers = make(map[string]*Logger)
var moduleLock sync.Mutex

// defaultConfig 默认配置对象
type defaultConfig struct {
	logLevel     int
	formatter    Formatter
	writer       Writer
	timeLocation *time.Location
}

// 默认配置信息
var defaultLogConfig = defaultConfig{
	logLevel:     LevelDebug,
	formatter:    &DefaultFormatter{},
	writer:       &DefaultWriter{},
	timeLocation: time.UTC,
}

// SetDefaultLocation set default time location
func SetDefaultLocation(loc *time.Location) {
	defaultLogConfig.timeLocation = loc
}

// SetDefaultLevel 设置全局默认日志输出级别
func SetDefaultLevel(level int) {
	defaultLogConfig.logLevel = level
}

// SetDefaultFormatter 设置全局默认的日志输出格式化器
func SetDefaultFormatter(formatter Formatter) {
	defaultLogConfig.formatter = formatter
}

// SetDefaultWriter 设置全局默认的日志输出器
func SetDefaultWriter(writer Writer) {
	defaultLogConfig.writer = writer
}

// Module 获取指定模块的日志输出对象
func Module(moduleName string) *Logger {
	moduleLock.Lock()
	defer moduleLock.Unlock()

	if logger, ok := loggers[moduleName]; ok {
		return logger
	}

	logger := &Logger{
		moduleName: moduleName,
		level:      defaultLogConfig.logLevel,
		timeLocation: func() *time.Location {
			return defaultLogConfig.timeLocation
		},
	}

	loggers[moduleName] = logger

	return logger
}

// SetTimeLocation set time location for module
func (module *Logger) SetTimeLocation(loc *time.Location) {
	module.timeLocation = func() *time.Location {
		return loc
	}
}

func (module *Logger) output(level int, context map[string]interface{}, v ...interface{}) string {
	message := module.getFormatter().Format(time.Now().In(module.timeLocation()), module.moduleName, level, context, v...)
	// 低于设定日志级别的日志不会输出
	if level >= module.level {
		if err := module.getWriter().Write(message); err != nil {
			fmt.Printf("can not write to output: %s", err)
		}
	}

	return message
}

// GetDefaultModule 获取默认的模块日志
func GetDefaultModule() *Logger {
	return Module("default")
}

// SetLevel 设置日志输出级别
func (module *Logger) SetLevel(level int) *Logger {
	module.level = level

	return module
}

// SetFormatter 设置日志格式化器
func (module *Logger) SetFormatter(formatter Formatter) *Logger {
	module.formatter = formatter
	return module
}

func (module *Logger) getFormatter() Formatter {
	moduleLock.Lock()
	defer moduleLock.Unlock()

	if module.formatter == nil {
		module.SetFormatter(defaultLogConfig.formatter)
	}

	return module.formatter
}

// SetWriter 设置日志输出器
func (module *Logger) SetWriter(writer Writer) *Logger {
	module.writer = writer
	return module
}

func (module *Logger) getWriter() Writer {
	moduleLock.Lock()
	defer moduleLock.Unlock()

	if module.writer == nil {
		module.SetWriter(defaultLogConfig.writer)
	}

	return module.writer
}

// WithContext 带有上下文信息的日志输出
func (module *Logger) WithContext(context map[string]interface{}) *ContextLogger {
	return &ContextLogger{
		logger:  module,
		context: context,
	}
}

// Emergency 记录emergency日志
func (module *Logger) Emergency(v ...interface{}) string {
	return module.output(LevelEmergency, nil, v...)
}

// Alert 记录Alert日志
func (module *Logger) Alert(v ...interface{}) string {
	return module.output(LevelAlert, nil, v...)
}

// Critical 记录Critical日志
func (module *Logger) Critical(v ...interface{}) string {
	return module.output(LevelCritical, nil, v...)
}

// Error 记录Error日志
func (module *Logger) Error(v ...interface{}) string {
	return module.output(LevelError, nil, v...)
}

// Warning 记录Warning日志
func (module *Logger) Warning(v ...interface{}) string {
	return module.output(LevelWarning, nil, v...)
}

// Notice 记录Notice日志
func (module *Logger) Notice(v ...interface{}) string {
	return module.output(LevelNotice, nil, v...)
}

// Info 记录Info日志
func (module *Logger) Info(v ...interface{}) string {
	return module.output(LevelInfo, nil, v...)
}

// Debug 记录Debug日志
func (module *Logger) Debug(v ...interface{}) string {
	return module.output(LevelDebug, nil, v...)
}

// Emergencyf 记录emergency日志
func (module *Logger) Emergencyf(format string, v ...interface{}) string {
	return module.Emergency(fmt.Sprintf(format, v...))
}

// Alertf 记录Alert日志
func (module *Logger) Alertf(format string, v ...interface{}) string {
	return module.Alert(fmt.Sprintf(format, v...))
}

// Criticalf 记录critical日志
func (module *Logger) Criticalf(format string, v ...interface{}) string {
	return module.Critical(fmt.Sprintf(format, v...))
}

// Errorf 记录error日志
func (module *Logger) Errorf(format string, v ...interface{}) string {
	return module.Error(fmt.Sprintf(format, v...))
}

// Warningf 记录warning日志
func (module *Logger) Warningf(format string, v ...interface{}) string {
	return module.Warning(fmt.Sprintf(format, v...))
}

// Noticef 记录notice日志
func (module *Logger) Noticef(format string, v ...interface{}) string {
	return module.Notice(fmt.Sprintf(format, v...))
}

// Infof 记录info日志
func (module *Logger) Infof(format string, v ...interface{}) string {
	return module.Info(fmt.Sprintf(format, v...))
}

// Debugf 记录debug日志
func (module *Logger) Debugf(format string, v ...interface{}) string {
	return module.Debug(fmt.Sprintf(format, v...))
}
