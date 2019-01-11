package log_test

import (
	"fmt"
	"testing"

	"github.com/mylxsw/go-toolkit/log"
)

func TestColorText(t *testing.T) {
	fmt.Println(log.ColorTextWrap(log.TextLightBlue, "Hello, world"))
	fmt.Println(log.ColorBackgroundWrap(log.TextLightCyan, log.TextLightBlue, "中文"))
}
