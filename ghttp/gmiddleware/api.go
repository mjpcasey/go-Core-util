package gmiddleware

import (
	"gcore/gmonitor"
	"github.com/juju/ratelimit"
	"time"
)

// 构建 QPS 限制中间件实现
//
// @param name 标识名
//
// @param limit 每秒QPS限制
//
// @return Middleware 实现
func NewQPSLimiter(name string, limit int) Middleware {
	lim := &qpsLimiter{
		name:  name,
		limit: limit,
	}

	if limit > 0 {
		lim.bucket = ratelimit.NewBucket(time.Second/time.Duration(limit), int64(limit))
	}

	return lim
}

// 构建 QPS 统计中间件实现
//
// @param name 标识名
//
// @return Middleware 实现
func NewQPSMonitor(name string) Middleware {
	return &qpsMonitor{
		name:    name,
		monitor: gmonitor.NewRequestMonitor(name),
	}
}

// 构建响应时间统计中间件实现
//
// @param name 标识名
//
// @return Middleware 实现
func NewTimeStatics(name string) Middleware {
	return &timeStatics{
		name: name,
	}
}
