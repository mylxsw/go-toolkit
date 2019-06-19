package file

import (
	"os"
	"testing"
)

func TestCreateZipArchiveFile(t *testing.T) {

	files := []ZipFile{
		{Name: "file_test.go", Path: "file_test.go"},
		{Name: "file.go", Path: "file.go"},
		{Name: "zip.go", Path: "zip.go"},
	}

	zipfile := "xxx.zip"

	if err := CreateZipArchiveFile(zipfile, files, false); err != nil {
		t.Errorf("test failed: %s", err.Error())
	}

	os.Remove(zipfile)
}
