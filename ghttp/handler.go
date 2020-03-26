package ghttp

import (
	"gcore/ghttp/gmiddleware"
	"net/http"
)

// 基本实现
type baseHandler struct {
	handler http.Handler

	middles []gmiddleware.Middleware
}

// 对 net/http 接口的实现
func (h *baseHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	for _, mid := range h.middles {
		if !mid.ProcessRequest(resp, req) {
			return
		}
	}

	h.handler.ServeHTTP(resp, req)

	// 这一段需要对 h.baseHandler 定接口以统一输出
	// 但目前各项目内 http 接口代码不好改动，所以暂时只能用来指标统计
	for i := len(h.middles) - 1; i >= 0; i-- {
		if !h.middles[i].ProcessResponse(resp, req) {
			return
		}
	}
}

// 添加中间件
//
// @param middles 中间件
func (h *baseHandler) AddMiddleware(middles ...gmiddleware.Middleware) {
	h.middles = append(h.middles, middles...)
}
