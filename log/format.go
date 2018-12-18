package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Formatter 日志格式化接口
type Formatter interface {
	// Format 日志格式化
	Format(currentTime time.Time, moduleName string, level int, context map[string]interface{}, v ...interface{}) string
}

// DefaultFormatter 默认日志格式化
type DefaultFormatter struct{}

// Format 日志格式化
func (formatter DefaultFormatter) Format(currentTime time.Time, moduleName string, level int, context map[string]interface{}, v ...interface{}) string {
	message := fmt.Sprintf("[%s] %s.%s: %s %s", currentTime.Format(time.RFC3339), moduleName, GetLevelName(level), fmt.Sprint(v...), formatContext(context))

	// 将多行内容增加前缀tab，与第一行内容分开
	return strings.Trim(strings.Replace(message, "\n", "\n	", -1), "\n	")
}

// JSONFormatter json输格式化
type JSONFormatter struct{}

type jsonOutput struct {
	ModuleName string                 `json:"module"`
	LevelName  string                 `json:"level_name"`
	Level      int                    `json:"level"`
	Context    map[string]interface{} `json:"context"`
	Message    string                 `json:"message"`
	DateTime   string                 `json:"datetime"`
}

// Format 日志格式化
func (formatter JSONFormatter) Format(currentTime time.Time, moduleName string, level int, context map[string]interface{}, v ...interface{}) string {
	datetime := currentTime.Format(time.RFC3339)

	if context == nil {
		context = make(map[string]interface{})
	}

	res, _ := json.Marshal(jsonOutput{
		DateTime:   datetime,
		Message:    fmt.Sprint(v...),
		Level:      level,
		ModuleName: moduleName,
		LevelName:  GetLevelName(level),
		Context:    context,
	})

	return fmt.Sprintf("[%s] %s", datetime, string(res))
}

func formatContext(context map[string]interface{}) string {
	if context == nil {
		context = make(map[string]interface{})
	}

	contextJSON, _ := json.Marshal(context)
	return string(contextJSON)
}
