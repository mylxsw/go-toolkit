package web

import (
	"net/http"

	"github.com/mylxsw/go-toolkit/container"
)

// HandlerDecorator 该函数是http handler的装饰器
type HandlerDecorator func(WebHandler) WebHandler

type handleFunc struct {
	callback  http.Handler
	decors    []HandlerDecorator
	container *container.Container
}

// Middleware 用于包装http handler，对其进行装饰
func Middleware(c *container.Container, h http.Handler, decors ...HandlerDecorator) http.Handler {
	return handleFunc{callback: h, decors: decors, container: c}
}

// ServeHTTP 实现http.HandlerFunc接口
func (f handleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var callback = func(context *WebContext) HTTPResponse {
		f.callback.ServeHTTP(context.Response.w, context.Request.r)
		return nil
	}

	for i := range f.decors {
		d := f.decors[len(f.decors)-i-1]
		callback = d(callback)
	}

	context := &WebContext{
		Response: &Response{
			w:       w,
			headers: make(map[string]string),
		},
		Request:   &Request{r: r},
		Container: f.container,
	}

	callback(context)
}
