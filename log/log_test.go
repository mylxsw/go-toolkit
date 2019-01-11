package log_test

import (
	"testing"
	"time"

	"github.com/mylxsw/go-toolkit/log"
)

func TestModule(t *testing.T) {
	// SetDefaultLevel(LevelCritical)

	loc, _ := time.LoadLocation("Asia/Chongqing")
	log.SetDefaultLocation(loc)
	// log.SetDefaultColorful(false)


	log.GetDefaultModule().SetLevel(log.LevelDebug)
	log.Debug("xxxx")

	log.Module("order").Noticef("order %s created", "1234592")
	log.Module("order").Infof("order %s created", "1234592")
	log.Module("order").Debugf("order %s created", "1234592")
	log.Module("order").Errorf("order %s created", "1234592")
	log.Module("order").Emergencyf("order %s created", "1234592")
	log.Module("order").Warningf("order %s created", "1234592")
	log.Module("order").Alertf("order %s created", "1234592")
	log.Module("order").Criticalf("order %s created", "1234592")

	log.Module("user").SetFormatter(log.NewJSONFormatter()).Error("user create failed")

	log.WithContext(nil).Debug("error occur")
	log.Module("purchase").SetFormatter(log.NewJSONFormatter()).WithContext(map[string]interface{}{}).Infof("用户 %s 已创建", "mylxsw")
}
