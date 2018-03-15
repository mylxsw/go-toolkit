package file_test

import (
	"testing"

	"git.yunsom.cn/golang/broadcast/utils/file"
)

func TestFileExist(t *testing.T) {
	if !file.Exist("file.go") {
		t.Errorf("check file existence failed")
	}

	if file.Exist("not_exist_file") {
		t.Errorf("check file existence failed")
	}
}
