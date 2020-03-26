package base

import "github.com/prometheus/client_golang/prometheus"

const (
	StatusOK     = 1 // StatusOK 正常连接状态
	StatusFailed = 0 // StatusFailed 正常连接状态
)

// Base 监控
type Base struct {
	Name   string
	metric prometheus.Gauge
}

// Up set 1
func (m *Base) Up() {
	m.metric.Set(StatusOK)
}

// Down set 0
func (m *Base) Down() {
	m.metric.Set(StatusFailed)
}

func New(name string, gauge prometheus.Gauge) Base {
	return Base{
		Name:   name,
		metric: gauge,
	}
}
