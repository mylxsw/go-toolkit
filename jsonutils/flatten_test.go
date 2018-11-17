package jsonutils

import (
	"testing"
)

var message1 = `{
    "message": "ack_confirm",
    "context": {
        "msg": "层级消息内容",
        "sms": {
            "id": 44444,
            "app_id": 1,
            "template_params": {
                "username": "李逍遥",
                "password": "lixiaoyao",
                "gender": "boy",
                "created_at": "2018-11-12 13:47:55"
            },
            "status": 1,
            "created_at": "2018-11-12 13:47:55",
            "updated_at": "2018-11-14 13:49:04"
        },
        "ack": {
            "msg": "😄",
            "code": "6460"
        },
        "file": "/webroot/your/project/Test.php:322"
    },
    "level": 400,
    "level_name": "ERROR",
    "channel": "custom_cmd",
    "datetime": "2018-11-16 13:51:01",
    "extra": {
        "ref": "5bee5ac564a71bbb33cai2jkk"
    }
}`

func TestToKvPairs(t *testing.T) {
	ju, err := New([]byte(message1))
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	kvPairs := ju.ToKvPairs()
	if len(kvPairs) == 0 {
		t.Error("convert to kv pairs failed")
	}

	if len(kvPairs) != 19 {
		t.Errorf("kv pairs not matched")
	}
}

func TestToKvPairsArray(t *testing.T) {
	ju, err := New([]byte(message1))
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	kvPairs := ju.ToKvPairsArray()
	if len(kvPairs) != 19 {
		t.Errorf("kv pairs not matched")
	}
}
