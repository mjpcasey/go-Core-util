package gmiddleware

import (
	"github.com/juju/ratelimit"
	"net/http"
)

// QPS 限制中间件
type qpsLimiter struct {
	name string
	// qps设置
	limit int
	// 限流桶
	bucket *ratelimit.Bucket
}

func (lim *qpsLimiter) ProcessRequest(resp http.ResponseWriter, req *http.Request) (proceed bool) {
	if lim.limit > 0 && lim.bucket != nil {
		if lim.bucket.Available() > 0 { // 令牌数满足，取一个，继续执行后续请求
			lim.bucket.Take(1)

			return true
		} else { // 服务繁忙，拒绝请求
			http.Error(resp, "Server Busy", http.StatusServiceUnavailable)

			return false
		}
	}

	return true
}

func (lim *qpsLimiter) ProcessResponse(resp http.ResponseWriter, req *http.Request) (proceed bool) {
	return true
}

