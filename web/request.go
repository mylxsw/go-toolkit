package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// Request 请求对象封装
type Request struct {
	r *http.Request
}

// Raw get the underlying http.Request
func (req *Request) Raw() *http.Request {
	return req.r
}

// UnmarshalJSON unmarshal request body as json object
// result must be reference to a variable
func (req *Request) UnmarshalJSON(v interface{}) error {
	data, err := ioutil.ReadAll(req.r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Set 设置一个变量，存储到当前请求
func (req *Request) Set(key string, value interface{}) {
	context.Set(req.r, key, value)
}

// Get 从当前请求提取设置的变量
func (req *Request) Get(key string) interface{} {
	return context.Get(req.r, key)
}

// Clear 清理掉请求中设置的变量
func (req *Request) Clear() {
	context.Clear(req.r)
}

// HTTPRequest 返回http.Request对象
func (req *Request) HTTPRequest() *http.Request {
	return req.r
}

// PathVar 获取路径中的变量
func (req *Request) PathVar(key string) string {
	if res, ok := mux.Vars(req.r)[key]; ok {
		return res
	}

	return ""
}

// PathVars 获取所有的路径变量
func (req *Request) PathVars() map[string]string {
	return mux.Vars(req.r)
}

// Input 获取表单输入
func (req *Request) Input(key string) string {
	return req.r.FormValue(key)
}

// File Retrieving Uploaded Files
func (req *Request) File(key string) (*UploadedFile, error) {
	file, header, err := req.r.FormFile(key)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	tempFile, err := ioutil.TempFile("/tmp", "yunsom-go-tools-")
	if err != nil {
		return nil, fmt.Errorf("无法创建临时文件 %s", err.Error())
	}
	defer tempFile.Close()

	io.Copy(tempFile, file)

	return &UploadedFile{
		Header:   header,
		SavePath: tempFile.Name(),
	}, nil
}

// IsXMLHTTPRequest 判断是否是XMLHTTPRequest
func (req *Request) IsXMLHTTPRequest() bool {
	return req.r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// AJAX 判断是否是AJAX请求
func (req *Request) AJAX() bool {
	return req.IsXMLHTTPRequest()
}

// Is 判断请求方法
func (req *Request) Is(method string) bool {
	return req.Method() == method
}

// IsGet 判断是否是Get请求
func (req *Request) IsGet() bool {
	return req.Is("GET")
}

// IsPost 判断是否是Post请求
func (req *Request) IsPost() bool {
	return req.Is("POST")
}

// IsHead 判断是否是HEAD请求
func (req *Request) IsHead() bool {
	return req.Is("HEAD")
}

// IsDelete 判断是是否是Delete请求
func (req *Request) IsDelete() bool {
	return req.Is("DELETE")
}

// IsPut 判断是否是Put请求
func (req *Request) IsPut() bool {
	return req.Is("PUT")
}

// IsPatch 判断是否是Patch请求
func (req *Request) IsPatch() bool {
	return req.Is("PATCH")
}

// IsOptions 判断是否是Options请求
func (req *Request) IsOptions() bool {
	return req.Is("OPTIONS")
}

// Method 获取请求方法
func (req *Request) Method() string {
	return req.r.Method
}
