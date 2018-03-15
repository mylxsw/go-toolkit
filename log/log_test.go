package log

import (
	"testing"
)

func TestModule(t *testing.T) {
	// SetDefaultLevel(LevelCritical)

	GetDefaultModule().SetLevel(LevelDebug)
	Debug("xxxx")

	Module("order").Noticef("order %s created", "1234592")
	Module("user").SetFormatter(JSONFormatter{}).Error("user create failed")

	WithContext(nil).Debug("error occur")
	Module("purchase").SetFormatter(JSONFormatter{}).WithContext(map[string]interface{}{}).Infof("用户 %s 已创建", "mylxsw")
}
