package log

import (
	"testing"
)

func TestModule(t *testing.T) {
	GetDefaultModule().Debug("xxxx")

	Module("order").Noticef("order %s created", "1234592")
	Module("user").SetFormatter(JSONFormatter{}).Error("user create failed")
}
