package web

import (
	"fmt"
	"net/http"

	"github.com/mylxsw/go-toolkit/container"
)

// WebContext 作为一个web请求的上下文信息
type WebContext struct {
	Response  *Response
	Request   *Request
	Container *container.Container
}

type webHandler struct {
	handle    WebHandler
	container *container.Container
}

// WebHandler 控制器方法
type WebHandler func(context *WebContext) HTTPResponse

// NewWebHandler 创建一个WebHandler，用于传递给Router
func NewWebHandler(c *container.Container, handler WebHandler, decors ...HandlerDecorator) webHandler {
	for i := range decors {
		d := decors[len(decors)-i-1]
		handler = d(handler)
	}

	return webHandler{
		handle:    handler,
		container: c,
	}
}

// ServeHTTP 实现http.HandlerFunc接口
func (h webHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &WebContext{
		Response: &Response{
			w:       w,
			headers: make(map[string]string),
		},
		Request:   &Request{r: r},
		Container: h.container,
	}

	resp := h.handle(context)
	if resp != nil {
		_ = resp.CreateResponse()
	}
}

// JSON is a shortcut for NewJSONResponse func
func (ctx *WebContext) JSON(res interface{}) *JSONResponse {
	return ctx.NewJSONResponse(res)
}

// NewJSONResponse create a new JSONResponse with the http status code equal 200
func (ctx *WebContext) NewJSONResponse(res interface{}) *JSONResponse {
	return NewJSONResponse(ctx.Response, http.StatusOK, res)
}

// JSONWithCode create a json response with a http status code
func (ctx *WebContext) JSONWithCode(res interface{}, code int) *JSONResponse {
	return NewJSONResponse(ctx.Response, code, res)
}

// API is a shortcut for NewAPIResponse func
func (ctx *WebContext) API(businessCode string, message string, data interface{}) *JSONResponse {
	return ctx.NewAPIResponse(businessCode, message, data)
}

// NewAPIResponse create a new APIResponse
func (ctx *WebContext) NewAPIResponse(businessCode string, message string, data interface{}) *JSONResponse {
	return ctx.NewJSONResponse(struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    businessCode,
		Message: message,
		Data:    data,
	})
}

// NewRawResponse create a new RawResponse
func (ctx *WebContext) NewRawResponse() *RawResponse {
	return NewRawResponse(ctx.Response)
}

// NewHTMLResponse create a new HTMLResponse
func (ctx *WebContext) NewHTMLResponse(res string) *HTMLResponse {
	return NewHTMLResponse(ctx.Response, http.StatusOK, res)
}

// HTML is a shortcut for NewHTMLResponse func
func (ctx *WebContext) HTML(res string) *HTMLResponse {
	return ctx.NewHTMLResponse(res)
}

// HTMLWithCode create a HTMLResponse with http status code
func (ctx *WebContext) HTMLWithCode(res string, code int) *HTMLResponse {
	return NewHTMLResponse(ctx.Response, code, res)
}

// Error is a shortcut for NewErrorResponse func
func (ctx *WebContext) Error(res string, code int) *ErrorResponse {
	return ctx.NewErrorResponse(res, code)
}

// NewErrorResponse create a error response
func (ctx *WebContext) NewErrorResponse(res string, code int) *ErrorResponse {
	return NewErrorResponse(ctx.Response, res, code)
}

// Redirect replies to the request with a redirect to url
func (ctx *WebContext) Redirect(location string, code int) *RedirectResponse {
	return NewRedirectResponse(ctx.Response, ctx.Request, location, code)
}

// Resolve resolve implements dependency injection for http handler
func (ctx *WebContext) Resolve(callback interface{}) HTTPResponse {
	results, err := ctx.Container.Call(callback)
	if err != nil {
		return ctx.NewErrorResponse(fmt.Sprintf("resolve dependency error: %s", err.Error()), 500)
	}

	if len(results) == 0 {
		return ctx.NewHTMLResponse("")
	}

	resp, ok := results[0].(HTTPResponse)
	if ok {
		return resp
	}

	return ctx.NewJSONResponse(results)
}
