package log

import (
	"testing"
	"time"
)

func TestModule(t *testing.T) {
	// SetDefaultLevel(LevelCritical)

	loc, _ := time.LoadLocation("Asia/Chongqing")
	SetDefaultLocation(loc)

	GetDefaultModule().SetLevel(LevelDebug)
	Debug("xxxx")

	Module("order").Noticef("order %s created", "1234592")
	Module("user").SetFormatter(JSONFormatter{}).Error("user create failed")

	WithContext(nil).Debug("error occur")
	Module("purchase").SetFormatter(JSONFormatter{}).WithContext(map[string]interface{}{}).Infof("用户 %s 已创建", "mylxsw")
}
