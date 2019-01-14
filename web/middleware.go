package web

import (
	"net/http"
	"time"

	"github.com/mylxsw/go-toolkit/container"
	"github.com/mylxsw/go-toolkit/log"
)

var logger = log.Module("toolkit.web")

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

// RequestMiddleware is a middleware collections
type RequestMiddleware struct{}

// NewRequestMiddleware create a new RequestMiddleware
func NewRequestMiddleware() RequestMiddleware {
	return RequestMiddleware{}
}

// AccessLog create a access log middleware
func (rm RequestMiddleware) AccessLog() HandlerDecorator {
	return func(handler WebHandler) WebHandler {
		return func(ctx *WebContext) HTTPResponse {
			defer func(startTime time.Time) {
				logger.Debugf(
					"%-8s %s [%.4fms]",
					ctx.Request.Method(),
					ctx.Request.HTTPRequest().URL.String(),
					time.Now().Sub(startTime).Seconds()*1000,
				)
			}(time.Now())

			return handler(ctx)
		}
	}
}

// CORS create a CORS middleware
func (rm RequestMiddleware) CORS(origin string) HandlerDecorator {
	return func(handler WebHandler) WebHandler {
		return func(ctx *WebContext) HTTPResponse {
			ctx.Response.Header("Access-Control-Allow-Origin", origin)
			ctx.Response.Header("Access-Control-Allow-Headers", "X-Requested-With")
			ctx.Response.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS,HEAD,PUT,PATCH,DELETE")

			return handler(ctx)
		}
	}
}

// ExceptionHandler create a Exception handler middleware
func (rm RequestMiddleware) ExceptionHandler() HandlerDecorator {
	return func(handler WebHandler) WebHandler {
		return func(ctx *WebContext) (resp HTTPResponse) {
			if err := recover(); err != nil {
				logger.Errorf("request failed %s", err)
				resp = ctx.Error("Internal Server Error", http.StatusInternalServerError)
			}

			return handler(ctx)
		}
	}
}
