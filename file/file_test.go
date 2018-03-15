package file

import (
	"testing"
)

func TestFileExist(t *testing.T) {
	if !Exist("file.go") {
		t.Errorf("check file existence failed")
	}

	if Exist("not_exist_file") {
		t.Errorf("check file existence failed")
	}
}
