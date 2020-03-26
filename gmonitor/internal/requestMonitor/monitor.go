package requestMonitor

import (
	"fmt"
	"gcore/gmonitor/internal/util"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// monitor 记录请求数和请求分位数
type monitor struct {
	Name string //

	counter  *prometheus.CounterVec
	duration *prometheus.SummaryVec
}

// Unregister 注销指标
func (req *monitor) Unregister() {
	prometheus.Unregister(req.counter)
	prometheus.Unregister(req.duration)
}

// AddRequest 请求计数器+1
//
// @status: 请求结果状态(e.q. http status 200)
func (req *monitor) AddRequest(status string) {
	req.counter.WithLabelValues(status).Inc()
}

// RecordRequest 记录请求耗时
//
// @status: 请求结果状态
func (req *monitor) RecordRequest(status string, duration time.Duration) {
	req.duration.WithLabelValues(status).Observe(float64(duration / time.Millisecond))
}

func New(name string) (m *monitor) {
	util.Check(name + "_request_total")
	util.Check(name + "_duration_milliseconds")

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_request_total",
			Help: fmt.Sprintf("The counter of %s request", name),
		},
		[]string{"status"},
	)
	durationSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: name + "_duration_milliseconds",
			Help: fmt.Sprintf("Summary the duration time of %s request", name),
			// key 0.5:    TP50 50%的请求时间小于某个值
			// value 0.05: 精度范围
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"status"},
	)

	prometheus.MustRegister(counter)
	prometheus.MustRegister(durationSummary)

	m = &monitor{
		Name:     name,
		counter:  counter,
		duration: durationSummary,
	}

	return
}
