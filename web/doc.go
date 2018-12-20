/*
Package web 实现了HTTP请求路由，中间件，请求响应处理，配合container实现依赖注入。

例如

	func ExampleHandler(ctx *web.WebContext) web.HTTPResponse {
		return ctx.Resolve(func(messageModel *storage.MessageModel) web.HTTPResponse {
			// your codes

			return ctx.NewAPIResponse("000000", "ok", map[string]interface{} {
				"aa": "bbb",
				"ccc": "ddd",
			})
		})
	}

*/
package web
