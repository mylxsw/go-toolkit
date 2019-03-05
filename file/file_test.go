package file

import (
	"fmt"
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

func TestFileSize(t *testing.T) {
	fileSize := Size("zip.go")
	fmt.Printf("file size is %.2f K", float64(fileSize) / 1024)

	if fileSize <= 0 {
		t.Errorf("check file size failed")
	}
}