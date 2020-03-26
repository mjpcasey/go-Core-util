package ghttp

import (
	"gcore/ghttp/gmiddleware"
	"net/http"
)

// 获取 Handler 实现
//
// @param mainHandler 实际业务 mainHandler
//
// @param middles 中间件
//
// @return Handler Handler的实现
func NewHandler(mainHandler http.Handler, middles... gmiddleware.Middleware) Handler {
	return &baseHandler{
		handler: mainHandler,
		middles: middles,
	}
}