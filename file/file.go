package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Exist 判断文件是否存在
func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// Size return file size for path
func Size(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return stat.Size()
}

// InsertSuffix insert a suffix to filepath
func InsertSuffix(src string, suffix string) string {
	ext := path.Ext(src)

	return fmt.Sprintf("%s%s%s", src[:len(src)-len(ext)], suffix, ext)
}

// ReplaceExt replace ext for src
func ReplaceExt(src string, ext string) string {
	ext1 := path.Ext(src)

	return fmt.Sprintf("%s%s", src[:len(src)-len(ext1)], ext)
}

// FileGetContents reads entire file into a string
func FileGetContents(filename string) (string, error) {
	res, err := ioutil.ReadFile(filename)
	return string(res), err
}
