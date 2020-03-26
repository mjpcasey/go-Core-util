package gmiddleware

import "net/http"

// todo 等待实现
// 响应时间统计中间件
type timeStatics struct {
	name string
}

func (w *timeStatics) ProcessRequest(resp http.ResponseWriter, req *http.Request) (proceed bool) {
	panic("implement me")
}

func (w *timeStatics) ProcessResponse(resp http.ResponseWriter, req *http.Request) (proceed bool) {
	panic("implement me")
}
