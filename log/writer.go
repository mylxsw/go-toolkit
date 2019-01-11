package log

import "fmt"

// Writer 日志输出接口
type Writer interface {
	Write(message string) error
}

// DefaultWriter 默认日志输出器
type DefaultWriter struct{}

// NewDefaultWriter create a new default writer
func NewDefaultWriter() *DefaultWriter {
	return &DefaultWriter{}
}

// Write 日志输出
func (writer DefaultWriter) Write(message string) error {
	fmt.Println(message)
	return nil
}
