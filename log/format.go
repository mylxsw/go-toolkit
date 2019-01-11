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
	Format(colorful bool, currentTime time.Time, moduleName string, level int, context map[string]interface{}, v ...interface{}) string
}

// DefaultFormatter 默认日志格式化
type DefaultFormatter struct{}

// NewDefaultFormatter create a new default formatter
func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{}
}

// Format 日志格式化
func (formatter DefaultFormatter) Format(colorful bool, currentTime time.Time, moduleName string, level int, context map[string]interface{}, v ...interface{}) string {

	var message string
	if colorful {

		if len(moduleName) > 20 {
			moduleName = "..." + moduleName[len(moduleName) - 17:]
		}


		message = fmt.Sprintf(
			"[%s] %-20s %s %s %s",
			ColorTextWrap(TextLightWhite, currentTime.Format(time.RFC3339)),
			moduleName,
			colorfulLevelName(level),
			strings.Trim(fmt.Sprint(v...), "\n"),
			ColorTextWrap(TextLightGrey, formatContext(context)),
		)
	} else {
		message = fmt.Sprintf(
			"[%s] %s.%s: %s %s",
			currentTime.Format(time.RFC3339),
			moduleName,
			GetLevelName(level),
			strings.Trim(fmt.Sprint(v...), "\n"),
			formatContext(context),
		)
	}

	// 将多行内容增加前缀tab，与第一行内容分开
	return strings.Trim(strings.Replace(message, "\n", "\n	", -1), "\n	")
}

// JSONFormatter json输格式化
type JSONFormatter struct{}

// NewJSONFormatter create a new json formatter
func NewJSONFormatter() *JSONFormatter  {
	return &JSONFormatter{}
}

type jsonOutput struct {
	ModuleName string                 `json:"module"`
	LevelName  string                 `json:"level_name"`
	Level      int                    `json:"level"`
	Context    map[string]interface{} `json:"context"`
	Message    string                 `json:"message"`
	DateTime   string                 `json:"datetime"`
}

// Format 日志格式化
func (formatter JSONFormatter) Format(colorful bool, currentTime time.Time, moduleName string, level int, context map[string]interface{}, v ...interface{}) string {
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

	message := string(res)
	if colorful {
		datetime = ColorTextWrap(TextLightWhite, datetime)
		message = ColorTextWrap(TextLightGrey, message)
	}

	return fmt.Sprintf("[%s] %s", datetime, message)
}

func formatContext(context map[string]interface{}) string {
	if context == nil {
		context = make(map[string]interface{})
	}

	contextJSON, _ := json.Marshal(context)
	return string(contextJSON)
}

func colorfulLevelName(level int) string {
	levelName := fmt.Sprintf("[%s]", GetLevelNameAbbreviation(level))

	switch level {
	case LevelDebug:
		levelName = ColorBackgroundWrap(TextLightWhite, TextLightBlue, levelName)
	case LevelInfo:
		levelName = ColorBackgroundWrap(TextWhite, TextLightCyan, levelName)
	case LevelNotice:
		levelName = ColorBackgroundWrap(TextLightWhite, TextYellow, levelName)
	case LevelWarning:
		levelName = ColorBackgroundWrap(TextRed, TextYellow, levelName)
	case LevelError:
		levelName = ColorBackgroundWrap(TextLightWhite, TextRed, levelName)
	case LevelCritical:
		levelName = ColorBackgroundWrap(TextLightWhite, TextLightRed, levelName)
	case LevelAlert:
		levelName = ColorBackgroundWrap(TextLightWhite, TextLightRed, levelName)
	case LevelEmergency:
		levelName = ColorBackgroundWrap(TextLightWhite, TextLightRed, levelName)
	}

	return levelName
}
