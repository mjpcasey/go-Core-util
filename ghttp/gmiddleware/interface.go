package gmiddleware

import "net/http"

// Middleware HTTP 中间件
type Middleware interface {
	ProcessRequest(resp http.ResponseWriter, req *http.Request) (proceed bool)
	ProcessResponse(resp http.ResponseWriter, req *http.Request) (proceed bool)
}