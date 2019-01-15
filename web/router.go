package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mylxsw/go-toolkit/container"
)

// Router 定制的路由
type Router struct {
	router      *mux.Router
	container   *container.Container
	routes      []route
	middlewares []HandlerDecorator
	prefix      string
}

// route 路由规则
type route struct {
	method      string
	path        string
	webHandler  WebHandler
	middlewares []HandlerDecorator
}

// NewRouter 创建一个路由器
func NewRouter(middlewares ...HandlerDecorator) *Router {
	return NewRouterWithContainer(container.New(), middlewares...)
}

// NewRouterWithContainer 创建一个路由器，带有依赖注入容器支持
func NewRouterWithContainer(c *container.Container, middlewares ...HandlerDecorator) *Router {
	return create(c, mux.NewRouter(), middlewares...)
}

// create 创建定制路由器
func create(c *container.Container, router *mux.Router, middlewares ...HandlerDecorator) *Router {
	return &Router{
		router:      router,
		routes:      []route{},
		middlewares: middlewares,
		prefix:      "",
		container:   c,
	}
}

// Group 创建路由组
func (router *Router) Group(prefix string, f func(rou *Router), decors ...HandlerDecorator) {
	r := create(router.container, router.router, decors...)
	r.prefix = prefix

	f(r)
	r.parse()

	for _, route := range r.getRoutes() {
		router.addWebHandler(route.method, route.path, route.webHandler, route.middlewares...)
	}
}

// Perform 将路由规则添加到路由器
func (router *Router) Perform() *mux.Router {
	for _, r := range router.routes {
		var handler http.Handler
		handler = NewWebHandler(router.container, r.webHandler, r.middlewares...)
		route := router.router.Handle(r.path, handler)
		if r.method != "" {
			route.Methods(r.method)
		}
	}

	return router.router
}

// GetRoutes 获取所有路由规则
func (router *Router) getRoutes() []route {
	return router.routes
}

func (router *Router) addWebHandler(method string, path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.routes = append(router.routes, route{
		method:      method,
		path:        path,
		webHandler:  handler,
		middlewares: middlewares,
	})
}

// Parse 解析路由规则，将中间件信息同步到路由规则
func (router *Router) parse() {
	for i := range router.routes {
		router.routes[i].path = fmt.Sprintf("%s/%s", strings.TrimRight(router.prefix, "/"), strings.TrimLeft(router.routes[i].path, "/"))
		router.routes[i].middlewares = append(router.routes[i].middlewares, router.middlewares...)
	}
}

func (router *Router) addHandler(method string, path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addWebHandler(method, path, func(ctx *WebContext) HTTPResponse {
		return ctx.Resolve(handler)
	}, middlewares...)
}

// Any 指定所有请求方式的路由规则
func (router *Router) Any(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("", path, handler, middlewares...)
}

// Get 指定所有GET方式的路由规则
func (router *Router) Get(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("GET", path, handler, middlewares...)
}

// Post 指定所有Post方式的路由规则
func (router *Router) Post(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("POST", path, handler, middlewares...)
}

// Delete 指定所有DELETE方式的路由规则
func (router *Router) Delete(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("DELETE", path, handler, middlewares...)
}

// Put 指定所有Put方式的路由规则
func (router *Router) Put(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("PUT", path, handler, middlewares...)
}

// Patch 指定所有Patch方式的路由规则
func (router *Router) Patch(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("PATCH", path, handler, middlewares...)
}

// Head 指定所有Head方式的路由规则
func (router *Router) Head(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("HEAD", path, handler, middlewares...)
}

// Options 指定所有OPTIONS方式的路由规则
func (router *Router) Options(path string, handler interface{}, middlewares ...HandlerDecorator) {
	router.addHandler("OPTIONS", path, handler, middlewares...)
}

// WebAny 指定所有请求方式的路由规则，WebHandler方式
func (router *Router) WebAny(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("", path, handler, middlewares...)
}

// WebGet 指定GET请求方式的路由规则，WebHandler方式
func (router *Router) WebGet(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("GET", path, handler, middlewares...)
}

// WebPost 指定POST请求方式的路由规则，WebHandler方式
func (router *Router) WebPost(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("POST", path, handler, middlewares...)
}

// WebPut 指定所有Put方式的路由规则，WebHandler方式
func (router *Router) WebPut(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("PUT", path, handler, middlewares...)
}

// WebDelete 指定所有DELETE方式的路由规则，WebHandler方式
func (router *Router) WebDelete(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("DELETE", path, handler, middlewares...)
}

// WebPatch 指定所有PATCH方式的路由规则，WebHandler方式
func (router *Router) WebPatch(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("PATCH", path, handler, middlewares...)
}

// WebHead 指定所有HEAD方式的路由规则，WebHandler方式
func (router *Router) WebHead(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("HEAD", path, handler, middlewares...)
}

// WebOptions 指定所有OPTIONS方式的路由规则，WebHandler方式
func (router *Router) WebOptions(path string, handler WebHandler, middlewares ...HandlerDecorator) {
	router.addWebHandler("OPTIONS", path, handler, middlewares...)
}
