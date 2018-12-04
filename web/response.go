package web

import "net/http"

// HTTPResponse HTTP响应接口
type HTTPResponse interface {
	CreateResponse() error
}

// Response 响应对象封装
type Response struct {
	w        http.ResponseWriter
	headers  map[string]string
	cookie   *http.Cookie
	original string
}

// ResponseWriter 获取http.ResponseWriter对象
func (resp *Response) ResponseWriter() http.ResponseWriter {
	return resp.w
}

// SetContent 设置响应内容
func (resp *Response) SetContent(content string) {
	resp.original = content
}

// Header 设置要发送的header
func (resp *Response) Header(key, value string) {
	resp.headers[key] = value
}

// Cookie 设置cookie
func (resp *Response) Cookie(cookie *http.Cookie) {
	// http.SetCookie(resp.w, cookie)
	resp.cookie = cookie
}

// Flush 发送响应内容给客户端
func (resp *Response) Flush() {
	// 发送响应头
	for key, value := range resp.headers {
		resp.w.Header().Set(key, value)
	}

	// 发送Cookie
	if resp.cookie != nil {
		http.SetCookie(resp.w, resp.cookie)
	}

	// 发送body
	resp.w.Write([]byte(resp.original))
}

// SendError 发送错误响应
func (resp *Response) SendError(code int, message string) {
	resp.w.WriteHeader(code)
	resp.w.Write([]byte(message))
}
