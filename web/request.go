package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// Request 请求对象封装
type Request struct {
	r    *http.Request
	body []byte
}

// Raw get the underlying http.Request
func (req *Request) Raw() *http.Request {
	return req.r
}

// UnmarshalJSON unmarshal request body as json object
// result must be reference to a variable
func (req *Request) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(req.body, v)
}

// Set 设置一个变量，存储到当前请求
func (req *Request) Set(key string, value interface{}) {
	context.Set(req.r, key, value)
}

// Get 从当前请求提取设置的变量
func (req *Request) Get(key string) interface{} {
	return context.Get(req.r, key)
}

// Clear clear all variables in request
func (req *Request) Clear() {
	context.Clear(req.r)
}

// HTTPRequest return a http.Request
func (req *Request) HTTPRequest() *http.Request {
	return req.r
}

// PathVar return a path parameter
func (req *Request) PathVar(key string) string {
	if res, ok := mux.Vars(req.r)[key]; ok {
		return res
	}

	return ""
}

// PathVars return all path parameters
func (req *Request) PathVars() map[string]string {
	return mux.Vars(req.r)
}

// Input return form parameter from request
func (req *Request) Input(key string) string {
	if req.IsJSON() {
		val := req.JSONGet(key)
		if val != "" {
			return val
		}
	}

	return req.r.FormValue(key)
}

func (req *Request) JSONGet(keys ...string) string {
	value, dataType, _, err := jsonparser.Get(req.body, keys...)
	if err != nil {
		return ""
	}

	switch dataType {
	case jsonparser.String:
		if res, err := jsonparser.ParseString(value); err == nil {
			return res
		}
	case jsonparser.Number:
		if res, err := jsonparser.ParseFloat(value); err == nil {
			return strconv.FormatFloat(res, 'f', -1, 32)
		}
		if res, err := jsonparser.ParseInt(value); err == nil {
			return fmt.Sprintf("%d", res)
		}
	case jsonparser.Object:
		fallthrough
	case jsonparser.Array:
		return fmt.Sprintf("%x", value)
	case jsonparser.Boolean:
		if res, err := jsonparser.ParseBoolean(value); err == nil {
			if res {
				return "true"
			} else {
				return "false"
			}
		}
	case jsonparser.NotExist:
		fallthrough
	case jsonparser.Null:
		fallthrough
	case jsonparser.Unknown:
		return ""
	}

	return ""
}

// InputWithDefault return a form parameter with a default value
func (req *Request) InputWithDefault(key string, defaultVal string) string {
	val := req.Input(key)
	if val == "" {
		return defaultVal
	}

	return val
}

func (req *Request) ToInt(val string, defaultVal int) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return res
}

func (req *Request) ToInt64(val string, defaultVal int64) int64 {
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}

	return res
}

func (req *Request) ToFloat32(val string, defaultVal float32) float32 {
	res, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return defaultVal
	}

	return float32(res)
}

func (req *Request) ToFloat64(val string, defaultVal float64) float64 {
	res, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}

	return res
}

// IntInput return a integer form parameter
func (req *Request) IntInput(key string, defaultVal int) int {
	return req.ToInt(req.Input(key), defaultVal)
}

// Int64Input return a integer form parameter
func (req *Request) Int64Input(key string, defaultVal int64) int64 {
	return req.ToInt64(req.Input(key), defaultVal)
}

// Float32Input return a float32 form parameter
func (req *Request) Float32Input(key string, defaultVal float32) float32 {
	return req.ToFloat32(req.Input(key), defaultVal)
}

// Float64Input return a float64 form parameter
func (req *Request) Float64Input(key string, defaultVal float64) float64 {
	return req.ToFloat64(req.Input(key), defaultVal)
}

// File Retrieving Uploaded Files
func (req *Request) File(key string) (*UploadedFile, error) {
	file, header, err := req.r.FormFile(key)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	tempFile, err := ioutil.TempFile("/tmp", "yunsom-go-tools-")
	if err != nil {
		return nil, fmt.Errorf("无法创建临时文件 %s", err.Error())
	}
	defer func() {
		_ = tempFile.Close()
	}()

	if _, err := io.Copy(tempFile, file); err != nil {
		return nil, err
	}

	return &UploadedFile{
		Header:   header,
		SavePath: tempFile.Name(),
	}, nil
}

// IsXMLHTTPRequest return whether the request is a ajax request
func (req *Request) IsXMLHTTPRequest() bool {
	return req.r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// AJAX return whether the request is a ajax request
func (req *Request) AJAX() bool {
	return req.IsXMLHTTPRequest()
}

// IsJSON return whether the request is a json request
func (req *Request) IsJSON() bool {
	return req.ContentType() == "application/json"
}

// ContentType return content type for request
func (req *Request) ContentType() string {
	t := req.r.Header.Get("Content-Type")
	if t == "" {
		return "text/html"
	}

	return strings.ToLower(strings.Split(t, ";")[0])
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
