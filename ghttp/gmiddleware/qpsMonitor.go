package gmiddleware

import (
	"gcore/gmonitor"
	"net/http"
	"strconv"
	"time"
)

// QPS 统计中间件
type qpsMonitor struct {
	name    string
	monitor gmonitor.RequestMonitor
}

func (mon *qpsMonitor) ProcessRequest(resp http.ResponseWriter, req *http.Request) (proceed bool) {
	reqTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	req.Header.Set("req-start-time", reqTime)
	return true
}

func (mon *qpsMonitor) ProcessResponse(resp http.ResponseWriter, req *http.Request) (proceed bool) {
	reqTimeNano := req.Header.Get("req-start-time")
	start, err := strconv.ParseInt(reqTimeNano, 10, 64)
	if err != nil {
		return true // 出错时不影响正常请求处理
	}

	now := time.Now()
	startTime := time.Unix(0, start)
	duration := now.Sub(startTime)
	mon.monitor.RecordRequest("", duration)

	mon.monitor.AddRequest("")

	return true
}
