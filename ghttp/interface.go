package ghttp

import (
	"gcore/ghttp/gmiddleware"
	"net/http"
)

// ghttp.Handler 定义
type Handler interface {
	http.Handler
	AddMiddleware(middles ...gmiddleware.Middleware) // 添加中间件接口
}
